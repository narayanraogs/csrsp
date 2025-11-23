package stages

import (
	"context"
	"csrspServer/pipeline"
	"csrspServer/utils/binary"
	"fmt"
	"math"
)

type UnpackConfig struct {
	StartByte     int
	NoOfBytes     int
	InputNoOfBits int
}

func NewUnpackStage(config UnpackConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			if config.NoOfBytes > 0 && config.NoOfBytes*8%config.InputNoOfBits != 0 {
				errChan <- fmt.Errorf("unpack: NoOfBytes*8 (%d) is not a multiple of InputNoOfBits (%d)", config.NoOfBytes*8, config.InputNoOfBits)
				return
			}

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in Unpack on frame %d: %v", frame.ID, r)
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
					if config.NoOfBytes <= 0 {
						end = len(payload)
					}

					if start < 0 || start > len(payload) || end > len(payload) {
						errChan <- fmt.Errorf("unpack: bounds [%d:%d] are out of range for frame %d with length %d", start, end, frame.ID, len(payload))
						return
					}

					newPayload := unpackPayload(payload[start:end], config.InputNoOfBits)
					output <- pipeline.Frame{ID: frame.ID, Payload: newPayload}
				}()
			}
		}()

		return output
	}
}

func unpackPayload(payload []byte, noOfBits int) []byte {
	if noOfBits%8 == 0 {
		return payload
	}

	output := make([]byte, 0)
	startBit := 0
	index := 0

	outputLen, _ := getPackedWordSize(noOfBits)
	mask := createUnpackMask(noOfBits)

	for index < len(payload) {
		endBit := startBit + noOfBits
		noOfWords, rightShift := getPackedWordSize(endBit)

		if index+noOfWords > len(payload) {
			break
		}

		array := payload[index : index+noOfWords]
		var tempArray = make([]byte, 8)
		copy(tempArray[8-len(array):], array)

		value, _ := binary.BytesToUint64BE(tempArray)
		value = value >> rightShift
		value = value & mask

		tempOutput := binary.Uint64ToBytesBE(value)
		output = append(output, tempOutput[8-outputLen:]...)

		index = index + endBit/8
		startBit = endBit % 8
	}
	return output
}

func getPackedWordSize(endBit int) (int, int) {
	noOfWords := endBit / 8
	if endBit%8 != 0 {
		noOfWords++
	}
	return noOfWords, noOfWords*8 - endBit
}

func createUnpackMask(numberOfBits int) uint64 {
	if numberOfBits >= 64 {
		return math.MaxUint64
	}
	return (1 << numberOfBits) - 1
}
