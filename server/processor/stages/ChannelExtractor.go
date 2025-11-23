package stages

import (
	"context"
	"csrspServer/pipeline"
	"fmt"
)

// ChannelExtractorConfig holds the configuration for a ChannelExtractor stage.
type ChannelExtractorConfig struct {
	// ValueProvider extracts the routing value from a frame's payload.
	ValueProvider func(payload []byte) (int, bool)
	// ValueToOutputIndex maps the extracted value to a specific output channel index.
	ValueToOutputIndex map[int]int
	// The total number of output channels this stage will produce.
	NumOutputs int
}

// NewChannelExtractorStage creates a new one-to-many stage that routes frames to different
// output channels based on a value extracted from the frame's content.
func NewChannelExtractorStage(config ChannelExtractorConfig, errChan chan<- error) pipeline.StageOneToMany {
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

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in ChannelExtractor on frame %d: %v", frame.ID, r)
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					// Extract the routing value from the frame.
					channelValue, ok := config.ValueProvider(frame.Payload)
					if !ok {
						errChan <- fmt.Errorf("ChannelExtractor: ValueProvider failed for frame %d", frame.ID)
						return // Skip this frame.
					}

					// Find the corresponding output channel index.
					outputIndex, exists := config.ValueToOutputIndex[channelValue]
					if !exists {
						// Optional: Log or send an error for unroutable frames.
						// errChan <- fmt.Errorf("ChannelExtractor: no route for value %d in frame %d", channelValue, frame.ID)
						return // Drop the frame.
					}

					if outputIndex < 0 || outputIndex >= len(outputs) {
						errChan <- fmt.Errorf("ChannelExtractor: mapped output index %d is out of bounds", outputIndex)
						return
					}

					// Send the frame to the correct output channel.
					outputs[outputIndex] <- frame
				}()
			}
		}()

		return readOnlyOutputs
	}
}
