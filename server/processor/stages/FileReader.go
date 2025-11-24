package stages

import (
	"context"
	"csrsp/server/pipeline"
	"fmt"
	"io"
	"os"
)

// FileReaderConfig holds the configuration for the FileReader source stage.
type FileReaderConfig struct {
	FilePath  string
	FrameSize int // If 0, the entire file is read as a single frame.
}

// NewFileReaderSource creates a new pipeline Source that reads data from a file.
func NewFileReaderSource(config FileReaderConfig, errChan chan<- error) pipeline.Source {
	return func(ctx context.Context) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in FileReader: %v", r)
				}
			}()

			file, err := os.Open(config.FilePath)
			if err != nil {
				errChan <- fmt.Errorf("FileReader: failed to open file %s: %w", config.FilePath, err)
				return
			}
			defer file.Close()

			frameID := 0
			buffer := make([]byte, config.FrameSize)

			// Handle case where the whole file is one frame.
			if config.FrameSize == 0 {
				payload, err := io.ReadAll(file)
				if err != nil {
					errChan <- fmt.Errorf("FileReader: failed to read entire file: %w", err)
					return
				}
				select {
				case output <- pipeline.Frame{ID: 0, Payload: payload}:
				case <-ctx.Done():
				}
				return
			}

			// Handle chunked reading.
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				n, err := io.ReadFull(file, buffer)
				if err == io.EOF {
					break // End of file, normal exit.
				} else if err == io.ErrUnexpectedEOF {
					// The last chunk is smaller than FrameSize, send it and finish.
					output <- pipeline.Frame{ID: frameID, Payload: buffer[:n]}
					break
				} else if err != nil {
					errChan <- fmt.Errorf("FileReader: error reading file chunk: %w", err)
					return
				}

				// Send a copy of the buffer, as it will be reused.
				payloadCopy := make([]byte, n)
				copy(payloadCopy, buffer[:n])
				output <- pipeline.Frame{ID: frameID, Payload: payloadCopy}
				frameID++
			}
		}()

		return output
	}
}
