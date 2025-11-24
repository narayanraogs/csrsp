package stages

import (
	"context"
	"csrsp/server/pipeline"
	"fmt"
	"strings"
)

// FrameSkipConfig holds the configuration for a FrameSkip stage.
type FrameSkipConfig struct {
	// ValueProvider is a function that extracts the specific string value from a frame's
	// payload that will be monitored for changes.
	ValueProvider func(payload []byte) (string, bool)
}

// NewFrameSkipStage creates a new one-to-one FrameSkip stage.
// This stage will drop all incoming frames until the value provided by the
// ValueProvider function changes. After that first change, all subsequent
// frames are passed through.
func NewFrameSkipStage(config FrameSkipConfig, errChan chan<- error) pipeline.StageOneToOne {

	// The returned function is the actual pipeline stage.
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		// Launch a single, long-lived goroutine to process all frames.
		go func() {
			defer close(output)

			// --- State Variables ---
			// These variables live for the entire duration of the stage, maintaining
			// state across multiple frames.
			var isFirstFrame = true
			var isTransitionDetected = false
			var previousValue string
			// --- End State Variables ---

			for frame := range input {
				// Process each frame in a panic-safe closure
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in FrameSkip on frame %d: %v", frame.ID, r)
						}
					}()

					// Check for context cancellation between frames.
					select {
					case <-ctx.Done():
						return
					default:
					}

					// --- Start of Business Logic ---

					// If the transition has already been found, just pass the frame through.
					if isTransitionDetected {
						output <- frame
						return
					}

					// Use the injected ValueProvider to get the current value.
					currentValue, ok := config.ValueProvider(frame.Payload)
					if !ok {
						errChan <- fmt.Errorf("FrameSkip: ValueProvider failed for frame %d", frame.ID)
						return // Skip this frame.
					}

					// The first frame is used only to establish the initial `previousValue`.
					if isFirstFrame {
						previousValue = currentValue
						isFirstFrame = false
						return // Drop the first frame.
					}

					// Compare the current value with the previous one.
					if strings.Compare(currentValue, previousValue) != 0 {
						// A transition has been detected.
						isTransitionDetected = true
						output <- frame // Pass this frame through.
					} else {
						// No transition yet, update previous value and drop the frame.
						previousValue = currentValue
					}
					// --- End of Business Logic ---
				}()
			}
		}()

		return output
	}
}
