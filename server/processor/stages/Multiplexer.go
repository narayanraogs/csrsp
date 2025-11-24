package stages

import (
	"context"
	"csrsp/server/pipeline"
	"fmt"
)

type MultiplexerConfig struct {
	ChunkSize         int
	OutputFrameLength int
}

func NewMultiplexerStage(config MultiplexerConfig, errChan chan<- error) pipeline.StageManyToOne {
	return func(ctx context.Context, inputs []<-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			if config.ChunkSize <= 0 {
				errChan <- fmt.Errorf("multiplexer: ChunkSize must be a positive integer")
				return
			}

			for {
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
					break
				}

				outputPayload := make([]byte, 0, config.OutputFrameLength)
				index := 0

				inputLength := len(inputFrames[0].Payload)

				for index < inputLength {
					for _, frame := range inputFrames {
						end := index + config.ChunkSize
						if end > len(frame.Payload) {
							end = len(frame.Payload)
						}
						outputPayload = append(outputPayload, frame.Payload[index:end]...)
					}
					index += config.ChunkSize
				}

				if len(outputPayload) != config.OutputFrameLength {
					errChan <- fmt.Errorf("multiplexer: final payload length %d does not match expected length %d", len(outputPayload), config.OutputFrameLength)
				}
				output <- pipeline.Frame{ID: firstFrameID, Payload: outputPayload}
			}
		}()
		return output
	}
}
