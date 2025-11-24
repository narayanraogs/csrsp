package stages

import (
	"context"
	"csrsp/server/pipeline"
	"fmt"
)

// ByteSwapperConfig holds the configuration for a ByteSwapper stage.
type ByteSwapperConfig struct {
	// The size of the byte chunks to be swapped (e.g., 2 for 16-bit, 4 for 32-bit).
	ChunkSize int
}

// NewByteSwapperStage creates a new one-to-one stage that reverses the order of bytes
// within fixed-size chunks of the frame payload.
func NewByteSwapperStage(config ByteSwapperConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			if config.ChunkSize <= 0 {
				errChan <- fmt.Errorf("ByteSwapper: ChunkSize must be a positive integer")
				return
			}

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in ByteSwapper on frame %d: %v", frame.ID, r)
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					payload := frame.Payload
					if len(payload)%config.ChunkSize != 0 {
						errChan <- fmt.Errorf("ByteSwapper: frame %d length (%d) is not a multiple of ChunkSize (%d)", frame.ID, len(payload), config.ChunkSize)
						return // Skip this invalid frame.
					}

					newPayload := make([]byte, len(payload))
					for i := 0; i < len(payload); i += config.ChunkSize {
						chunk := payload[i : i+config.ChunkSize]
						reversedChunk := reverseBytes(chunk)
						copy(newPayload[i:], reversedChunk)
					}

					output <- pipeline.Frame{ID: frame.ID, Payload: newPayload}
				}()
			}
		}()

		return output
	}
}

// reverseBytes creates a new slice with the byte order reversed.
func reverseBytes(input []byte) []byte {
	output := make([]byte, len(input))
	for i, j := 0, len(input)-1; i < len(input); i, j = i+1, j-1 {
		output[i] = input[j]
	}
	return output
}
