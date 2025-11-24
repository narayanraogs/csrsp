package stages

import (
	"context"
	"csrsp/server/pipeline"
	"fmt"
)

// DemultiplexerConfig holds the configuration for a Demultiplexer stage.
type DemultiplexerConfig struct {
	// The size of the chunk in bytes to deal to each output channel.
	ChunkSize int
	// The number of output channels to create.
	NumOutputs int
}

// NewDemultiplexerStage creates a new one-to-many stage that demultiplexes a single
// input stream into multiple output streams.
func NewDemultiplexerStage(config DemultiplexerConfig, errChan chan<- error) pipeline.StageOneToMany {
	return func(ctx context.Context, input <-chan pipeline.Frame) []<-chan pipeline.Frame {
		outputs := make([]chan pipeline.Frame, config.NumOutputs)
		readOnlyOutputs := make([]<-chan pipeline.Frame, config.NumOutputs)
		for i := 0; i < config.NumOutputs; i++ {
			outputs[i] = make(chan pipeline.Frame)
			readOnlyOutputs[i] = outputs[i]
		}

		go func() {
			defer func() {
				for _, c := range outputs {
					close(c)
				}
			}()

			if config.ChunkSize <= 0 {
				errChan <- fmt.Errorf("demultiplexer: ChunkSize must be a positive integer")
				return
			}

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in Demultiplexer on frame %d: %v", frame.ID, r)
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					payload := frame.Payload
					outputPayloads := make([][]byte, config.NumOutputs)
					for i := range outputPayloads {
						// Pre-allocate with a reasonable capacity to reduce re-allocations.
						outputPayloads[i] = make([]byte, 0, len(payload)/config.NumOutputs)
					}

					offset := 0
					for offset < len(payload) {
						for i := 0; i < config.NumOutputs; i++ {
							end := offset + config.ChunkSize
							if end > len(payload) {
								end = len(payload)
							}
							if offset >= end {
								break
							}

							outputPayloads[i] = append(outputPayloads[i], payload[offset:end]...)
							offset = end
						}
					}

					// Send the resulting frames to their respective output channels.
					for i, outPayload := range outputPayloads {
						outputs[i] <- pipeline.Frame{ID: frame.ID, Payload: outPayload}
					}
				}()
			}
		}()

		return readOnlyOutputs
	}
}
