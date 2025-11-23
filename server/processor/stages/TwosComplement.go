package stages

import (
	"context"
	"csrspServer/pipeline"
	"fmt"
)

func NewTwosComplementStage(config ComplementConfig, errChan chan<- error) pipeline.StageOneToOne {

	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			bytesPerWord := config.BitsPerWord / 8
			if bytesPerWord == 0 {
				errChan <- fmt.Errorf("TwosComplement: BitsPerWord must be >= 8")
				return
			}
			mask := makeMask(bytesPerWord, config.BitsToComplement)

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in TwosComplement on frame %d: %v", frame.ID, r)
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					payload := frame.Payload
					if len(payload)%bytesPerWord != 0 {
						errChan <- fmt.Errorf("TwosComplement: frame %d length (%d) is not a multiple of BitsPerWord (%d)", frame.ID, len(payload), config.BitsPerWord)
						return
					}

					newPayload := make([]byte, len(payload))

					for i := 0; i < len(payload); i += bytesPerWord {
						word := payload[i : i+bytesPerWord]
						complementedWord := applyTwosComplement(word, mask)
						copy(newPayload[i:], complementedWord)
					}

					output <- pipeline.Frame{ID: frame.ID, Payload: newPayload}
				}()
			}
		}()

		return output
	}
}

func applyTwosComplement(input []byte, mask []byte) []byte {
	out := make([]byte, len(input))
	carry := byte(1)
	for i := len(input) - 1; i >= 0; i-- {
		xordByte := input[i] ^ mask[i]
		sum := uint16(xordByte) + uint16(carry)
		out[i] = byte(sum)
		carry = byte(sum >> 8)
	}
	return out
}
