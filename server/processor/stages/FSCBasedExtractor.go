package stages

import (
	"context"
	"csrspServer/pipeline"
	"fmt"
)

type FSCBasedExtractorConfig struct {
	FSC                    []byte
	FSCMask                []byte
	FixedFrameLength       int
	VariableLengthProvider func(payload []byte) (int, bool)
}

func NewFSCBasedExtractorStage(config FSCBasedExtractorConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			var buffer []byte

			if len(config.FSC) == 0 || len(config.FSC) != len(config.FSCMask) {
				errChan <- fmt.Errorf("FSCBasedExtractor: FSC and FSCMask must have the same non-zero length")
				return
			}

			for frame := range input {
				select {
				case <-ctx.Done():
					return
				default:
				}

				buffer = append(buffer, frame.Payload...)

				for {
					index := getIndexWithMask(buffer, config.FSC, config.FSCMask)
					if index == -1 {
						if len(buffer) > len(config.FSC) {
							buffer = buffer[len(buffer)-len(config.FSC):]
						}
						break
					}

					buffer = buffer[index:]

					frameLength := config.FixedFrameLength
					if frameLength == 0 {
						if config.VariableLengthProvider == nil {
							errChan <- fmt.Errorf("FSCBasedExtractor: FixedFrameLength is 0 but VariableLengthProvider is not set")
							return
						}
						var ok bool
						frameLength, ok = config.VariableLengthProvider(buffer)
						if !ok {
							errChan <- fmt.Errorf("FSCBasedExtractor: VariableLengthProvider failed on frame %d", frame.ID)
							buffer = buffer[len(config.FSC):]
							continue
						}
					}

					if len(buffer) < frameLength {
						break
					}

					outputFramePayload := buffer[:frameLength]
					output <- pipeline.Frame{ID: frame.ID, Payload: outputFramePayload}

					buffer = buffer[frameLength:]
				}
			}
		}()

		return output
	}
}

func getIndexWithMask(input []byte, fsc []byte, mask []byte) int {
	for i := 0; i <= len(input)-len(fsc); i++ {
		match := true
		for j := 0; j < len(fsc); j++ {
			if (input[i+j] & mask[j]) != fsc[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
