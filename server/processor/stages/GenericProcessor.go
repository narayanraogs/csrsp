package stages

import (
	"context"
	"csrsp/server/pipeline"
	"csrsp/server/utils/binary"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
)

type GenericProcessorConfig struct {
	ExePath         string
	WriteSocketPath string
	ReadSocketPath  string
	InstanceID      int
}

func NewGenericProcessorStage(config GenericProcessorConfig, errChan chan<- error) pipeline.StageManyToMany {
	return func(ctx context.Context, inputs []<-chan pipeline.Frame) []<-chan pipeline.Frame {
		numOutputs := len(inputs) // Assuming a symmetric number of inputs and outputs
		outputs := make([]chan pipeline.Frame, numOutputs)
		readOnlyOutputs := make([]<-chan pipeline.Frame, numOutputs)
		for i := 0; i < numOutputs; i++ {
			outputs[i] = make(chan pipeline.Frame)
			readOnlyOutputs[i] = outputs[i]
		}

		go func() {
			defer func() {
				for _, c := range outputs {
					close(c)
				}
			}()

			writeSockPath := fmt.Sprintf("%s_%d.sock", strings.TrimSuffix(config.WriteSocketPath, ".sock"), config.InstanceID)
			readSockPath := fmt.Sprintf("%s_%d.sock", strings.TrimSuffix(config.ReadSocketPath, ".sock"), config.InstanceID)
			os.Remove(readSockPath)
			os.Remove(writeSockPath)

			readListener, err := net.Listen("unix", readSockPath)
			if err != nil {
				errChan <- fmt.Errorf("GenericProcessor: failed to listen on read socket %s: %w", readSockPath, err)
				return
			}
			defer readListener.Close()

			writeListener, err := net.Listen("unix", writeSockPath)
			if err != nil {
				errChan <- fmt.Errorf("GenericProcessor: failed to listen on write socket %s: %w", writeSockPath, err)
				return
			}
			defer writeListener.Close()

			cmd := exec.CommandContext(ctx, config.ExePath, readSockPath, writeSockPath)
			cmd.Stdout = os.Stdout // Or a custom logger
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				errChan <- fmt.Errorf("GenericProcessor: failed to start executable %s: %w", config.ExePath, err)
				return
			}

			readConnChan := make(chan net.Conn)
			writeConnChan := make(chan net.Conn)
			go func() {
				conn, err := readListener.Accept()
				if err == nil {
					readConnChan <- conn
				}
			}()
			go func() {
				conn, err := writeListener.Accept()
				if err == nil {
					writeConnChan <- conn
				}
			}()

			var readConn, writeConn net.Conn
			select {
			case readConn = <-readConnChan:
			case <-ctx.Done():
				return
			}
			select {
			case writeConn = <-writeConnChan:
			case <-ctx.Done():
				return
			}
			defer readConn.Close()
			defer writeConn.Close()
			// --- End Setup ---

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				var writeWg sync.WaitGroup
				writeWg.Add(len(inputs))
				for i, inputChan := range inputs {
					go func(channelIndex int, ch <-chan pipeline.Frame) {
						defer writeWg.Done()
						primaryHeader := binary.Uint64ToBytesBE(uint64(channelIndex))
						for frame := range ch {
							secondaryHeader := binary.Uint64ToBytesBE(uint64(len(frame.Payload)))
							header := slices.Concat(primaryHeader, secondaryHeader)
							fullMessage := slices.Concat(header, frame.Payload)
							if _, err := writeConn.Write(fullMessage); err != nil {
								// Error is sent, but we don't stop the whole stage
								errChan <- fmt.Errorf("GenericProcessor: write error: %w", err)
							}
						}
					}(i, inputChan)
				}
				writeWg.Wait()
			}()

			go func() {
				defer wg.Done()
				const headerLen = 16 // 8 bytes for channel index, 8 for frame length
				for {
					headerBuf := make([]byte, headerLen)
					_, err := io.ReadFull(readConn, headerBuf)
					if err != nil {
						// EOF or closed connection will break the loop
						break
					}
					c, _ := binary.BytesToUint64BE(headerBuf[:8])
					f, _ := binary.BytesToUint64BE(headerBuf[8:])
					chanIndex := int(c)
					frameLength := int(f)

					if chanIndex < 0 || chanIndex >= len(outputs) {
						errChan <- fmt.Errorf("GenericProcessor: received invalid channel index %d", chanIndex)
						continue
					}

					payloadBuf := make([]byte, frameLength)
					_, err = io.ReadFull(readConn, payloadBuf)
					if err != nil {
						break
					}
					outputs[chanIndex] <- pipeline.Frame{Payload: payloadBuf}
				}
			}()

			wg.Wait()          // Wait for both reader and writer to finish
			cmd.Process.Kill() // Ensure the external process is terminated
		}()

		return readOnlyOutputs
	}
}
