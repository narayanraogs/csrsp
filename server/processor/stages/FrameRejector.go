package stages

import (
	"context"
	"csrspServer/pipeline"
	"fmt"
)

// FrameRejectorConfig holds the configuration for a FrameRejector stage.
type FrameRejectorConfig struct {
	// The number of initial frames to reject from the stream.
	FramesToReject int
}

// NewFrameRejectorStage creates a new one-to-one stage that rejects the first N frames.
func NewFrameRejectorStage(config FrameRejectorConfig, errChan chan<- error) pipeline.StageOneToOne {

	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			// --- State Variable ---
			// This counter maintains the state for the stage.
			var framesRejected = 0
			// --- End State Variable ---

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in FrameRejector on frame %d: %v", frame.ID, r)
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					// --- Start of Business Logic ---
					if framesRejected < config.FramesToReject {
						framesRejected++
						return // Reject this frame.
					}

					// The rejection count has been met, so pass the frame through.
					output <- frame
					// --- End of Business Logic ---
				}()
			}
		}()

		return output
	}
}
