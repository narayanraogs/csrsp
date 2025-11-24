package stages

import (
	"context"
	"csrsp/server/pipeline"
	"fmt"
)

type TransposeConfig struct {
	NoOfRows      int
	NoOfColumns   int
	BytesPerPixel int
}

func NewTransposeStage(config TransposeConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			if config.NoOfRows <= 0 || config.NoOfColumns <= 0 || config.BytesPerPixel <= 0 {
				errChan <- fmt.Errorf("transpose: configuration values must be positive")
				return
			}

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in Transpose on frame %d: %v", frame.ID, r)
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					payload := frame.Payload
					expectedLen := config.NoOfRows * config.NoOfColumns * config.BytesPerPixel
					if len(payload) != expectedLen {
						errChan <- fmt.Errorf("transpose: frame %d length (%d) does not match expected length (%d)", frame.ID, len(payload), expectedLen)
						return
					}

					newPayload := make([]byte, expectedLen)
					inputOffset := 0

					for row := 0; row < config.NoOfRows; row++ {
						for col := 0; col < config.NoOfColumns; col++ {
							outputOffset := (col*config.NoOfRows + row) * config.BytesPerPixel
							copy(newPayload[outputOffset:outputOffset+config.BytesPerPixel], payload[inputOffset:inputOffset+config.BytesPerPixel])
							inputOffset += config.BytesPerPixel
						}
					}
					output <- pipeline.Frame{ID: frame.ID, Payload: newPayload}
				}()
			}
		}()
		return output
	}
}
