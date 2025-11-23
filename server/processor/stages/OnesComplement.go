package stages

import (
	"context"
	"csrspServer/pipeline"
	"fmt"
)

func NewOnesComplementStage(config ComplementConfig, errChan chan<- error) pipeline.StageOneToOne {

	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			bytesPerWord := config.BitsPerWord / 8
			if bytesPerWord == 0 {
				errChan <- fmt.Errorf("OnesComplement: BitsPerWord must be >= 8")
				return
			}
			mask := makeMask(bytesPerWord, config.BitsToComplement)

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in OnesComplement on frame %d: %v", frame.ID, r)
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					payload := frame.Payload
					if len(payload)%bytesPerWord != 0 {
						errChan <- fmt.Errorf("OnesComplement: frame %d length (%d) is not a multiple of BitsPerWord (%d)", frame.ID, len(payload), config.BitsPerWord)
						return
					}

					newPayload := make([]byte, len(payload))

					for i := 0; i < len(payload); i += bytesPerWord {
						word := payload[i : i+bytesPerWord]
						complementedWord := applyOnesComplement(word, mask)
						copy(newPayload[i:], complementedWord)
					}

					output <- pipeline.Frame{ID: frame.ID, Payload: newPayload}
				}()
			}
		}()

		return output
	}
}

func applyOnesComplement(input []byte, mask []byte) []byte {
	out := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		out[i] = input[i] ^ mask[i]
	}
	return out
}
