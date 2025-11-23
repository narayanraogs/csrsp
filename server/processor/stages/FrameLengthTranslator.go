package stages

import (
	"context"
	"csrspServer/pipeline"
	"fmt"
)

// FrameLengthTranslatorConfig holds the configuration for a FrameLengthTranslator stage.
type FrameLengthTranslatorConfig struct {
	// The desired length in bytes for all outgoing frames.
	OutputFrameLength int
}

// NewFrameLengthTranslatorStage creates a new one-to-one stage that rebuffers an
// incoming stream of frames into a new stream of frames with a fixed length.
func NewFrameLengthTranslatorStage(config FrameLengthTranslatorConfig, errChan chan<- error) pipeline.StageOneToOne {

	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			// --- State Variable ---
			// This buffer holds data that has been received but not yet emitted.
			// Its lifecycle is the same as the stage itself.
			var buffer []byte
			// --- End State Variable ---

			// If OutputFrameLength is not set, it would cause an infinite loop or panic.
			if config.OutputFrameLength <= 0 {
				errChan <- fmt.Errorf("FrameLengthTranslator: OutputFrameLength must be a positive integer")
				return
			}

			for frame := range input {
				// No need for a separate panic-recovery closure here, as this stage's
				// core logic does not have high-risk operations like complex type assertions
				// or unsafe memory access. Appending and slicing are safe.

				select {
				case <-ctx.Done():
					return
				default:
				}

				// Append the new frame's data to the internal buffer.
				buffer = append(buffer, frame.Payload...)

				// Emit as many full frames as possible from the buffer.
				for len(buffer) >= config.OutputFrameLength {
					// Create the new output frame from the start of the buffer.
					newPayload := buffer[:config.OutputFrameLength]

					// Send the new frame downstream.
					output <- pipeline.Frame{ID: frame.ID, Payload: newPayload} // Note: ID is from the last input frame.

					// Reslice the buffer to remove the data that was just sent.
					buffer = buffer[config.OutputFrameLength:]
				}
			}

			// After the input channel closes, check if any data is left in the buffer.
			if len(buffer) > 0 {
				errChan <- fmt.Errorf("frame length translator: %d bytes of data were left in the buffer after the input channel closed. This may indicate a data alignment issue", len(buffer))
			}
		}()

		return output
	}
}
