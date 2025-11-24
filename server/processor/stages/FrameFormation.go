package stages

import (
	"context"
	"csrsp/server/pipeline"
	"fmt"
)

// FrameFormationConfig holds the configuration for a FrameFormation stage.
type FrameFormationConfig struct {
	// Defines the order and source for each chunk of the output frame.
	Chunks []FrameChunk
}

// FrameChunk defines one piece of the final assembled frame.
type FrameChunk struct {
	// The index of the input channel from which to take the data.
	InputIndex int
	// The starting byte in the input frame's payload.
	StartByte int
	// The number of bytes to take.
	NoOfBytes int
}

// NewFrameFormationStage creates a new many-to-one stage that assembles a single
// output frame from chunks of multiple input frames.
func NewFrameFormationStage(config FrameFormationConfig, errChan chan<- error) pipeline.StageManyToOne {
	return func(ctx context.Context, inputs []<-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			for {
				// In each iteration, read one frame from every input channel.
				inputFrames := make([]*pipeline.Frame, len(inputs))
				var firstFrameID int
				var allInputsClosed bool

				for i, inputChan := range inputs {
					select {
					case frame, ok := <-inputChan:
						if !ok {
							allInputsClosed = true
							break
						}
						if i == 0 {
							firstFrameID = frame.ID
						}
						inputFrames[i] = &frame
					case <-ctx.Done():
						return
					}
				}
				if allInputsClosed {
					break // Exit the main loop.
				}

				// --- Start of Business Logic ---
				var outputPayload []byte
				for _, chunk := range config.Chunks {
					if chunk.InputIndex >= len(inputFrames) || inputFrames[chunk.InputIndex] == nil {
						errChan <- fmt.Errorf("FrameFormation: invalid InputIndex %d in chunk config", chunk.InputIndex)
						continue // Skip this iteration
					}
					inputPayload := inputFrames[chunk.InputIndex].Payload
					end := chunk.StartByte + chunk.NoOfBytes
					if end > len(inputPayload) {
						errChan <- fmt.Errorf("FrameFormation: chunk bounds [%d:%d] out of range for input frame %d", chunk.StartByte, end, chunk.InputIndex)
						continue // Skip this iteration
					}
					outputPayload = append(outputPayload, inputPayload[chunk.StartByte:end]...)
				}

				output <- pipeline.Frame{ID: firstFrameID, Payload: outputPayload}
				// --- End of Business Logic ---
			}
		}()

		return output
	}
}
