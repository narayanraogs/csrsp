package stages

import (
	"context"
	"csrsp/server/pipeline"
	"fmt"
)

type FrameStripperConfig struct {
	StartByte int
	NoOfBytes int
	FullData  bool
}

// NewFrameStripper creates a new one-to-one FrameStripper stage.
// It returns a StageOneToOne function that can be added to a pipeline.
func NewFrameStripperStage(config FrameStripperConfig, errChan chan<- error) pipeline.StageOneToOne {

	// The returned function is the actual pipeline stage.
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		// Launch a single, long-lived goroutine to process all frames.
		go func() {
			defer close(output)

			for frame := range input {
				// Process each frame in a panic-safe closure.
				func() {
					defer func() {
						if r := recover(); r != nil {
							// If a panic occurs (e.g., out-of-bounds slice), send the error
							// to the supervisor and continue to the next frame.
							errChan <- fmt.Errorf("panic in FrameStripper on frame %d: %v", frame.ID, r)
						}
					}()

					// Check for context cancellation between frames.
					select {
					case <-ctx.Done():
						return // Exit the loop if the context is cancelled.
					default:
					}

					// --- Start of Business Logic ---
					payload := frame.Payload
					payloadLen := len(payload)
					sb := config.StartByte
					nb := config.NoOfBytes

					if sb < 0 {
						sb = payloadLen + sb
					}

					// Basic bounds checking to prevent panics.
					if sb < 0 || sb > payloadLen {
						errChan <- fmt.Errorf("FrameStripper: StartByte [%d] is out of bounds for frame %d with length %d", sb, frame.ID, payloadLen)
						return // Skip this frame.
					}

					var strippedPayload []byte
					if config.FullData {
						strippedPayload = payload[sb:]
					} else {
						endByte := sb + nb
						if nb < 0 {
							endByte = payloadLen + nb
						}
						if endByte < sb || endByte > payloadLen {
							errChan <- fmt.Errorf("FrameStripper: EndByte [%d] is out of bounds for frame %d with length %d", endByte, frame.ID, payloadLen)
							return // Skip this frame.
						}
						strippedPayload = payload[sb:endByte]
					}
					// --- End of Business Logic ---

					// Send the new frame to the output channel.
					output <- pipeline.Frame{ID: frame.ID, Payload: strippedPayload}
				}()
			}
		}()

		return output
	}
}
