package stages

import (
	"context"
	"csrspServer/pipeline"
	"fmt"
)

type NibbleShifterConfig struct {
	ShiftRight bool
	StartByte  int
	NoOfBytes  int
}

func NewNibbleShifterStage(config NibbleShifterConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			var shifter func([]byte) []byte
			if config.ShiftRight {
				shifter = rightShiftNibbles
			} else {
				shifter = leftShiftNibbles
			}

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in NibbleShifter on frame %d: %v", frame.ID, r)
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					payload := frame.Payload
					start := config.StartByte
					end := start + config.NoOfBytes

					if start < 0 || start > len(payload) || end > len(payload) {
						errChan <- fmt.Errorf("NibbleShifter: bounds [%d:%d] are out of range for frame %d with length %d", start, end, frame.ID, len(payload))
						return
					}

					newPayload := make([]byte, 0, len(payload))
					newPayload = append(newPayload, payload[:start]...)
					newPayload = append(newPayload, shifter(payload[start:end])...)
					newPayload = append(newPayload, payload[end:]...)

					output <- pipeline.Frame{ID: frame.ID, Payload: newPayload}
				}()
			}
		}()

		return output
	}
}

func rightShiftNibbles(payload []byte) []byte {
	shiftedData := make([]byte, len(payload))
	var previousNibble byte // The low nibble from the previous byte.

	for i, currentByte := range payload {
		shiftedData[i] = (previousNibble << 4) | (currentByte >> 4)
		previousNibble = currentByte & 0x0F
	}
	return shiftedData
}

func leftShiftNibbles(payload []byte) []byte {
	shiftedData := make([]byte, len(payload))
	var nextNibble byte // The high nibble from the next byte.

	for i := len(payload) - 1; i >= 0; i-- {
		currentByte := payload[i]
		shiftedData[i] = (currentByte << 4) | (nextNibble >> 4)
		nextNibble = currentByte & 0xF0
	}
	return shiftedData
}
