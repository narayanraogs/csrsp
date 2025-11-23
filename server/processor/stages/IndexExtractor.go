package stages

import (
	"context"
	"csrspServer/pipeline"
	"fmt"
	"math"
)

type IndexExtractorConfig struct {
	IndexProvider func(payload []byte) (int, bool)

	StartByte  int
	NoOfBytes  int
	StartIndex int
	EndIndex   int
}

func NewIndexExtractorStage(config IndexExtractorConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			frameMap := make(map[int][]byte)
			expectedIndex := config.StartIndex
			lastSeenIndex := math.MinInt32

			sendFrame := func(endIdx int, lastFrameID uint64) {
				if endIdx < config.StartIndex {
					return
				}
				// Determine the total size of the final payload.
				totalSize := 0
				for i := config.StartIndex; i <= endIdx; i++ {
					chunk, ok := frameMap[i]
					if !ok {
						errChan <- fmt.Errorf("IndexExtractor: missing chunk for index %d when assembling frame", i)
						return
					}
					totalSize += len(chunk)
				}

				newPayload := make([]byte, 0, totalSize)
				for i := config.StartIndex; i <= endIdx; i++ {
					newPayload = append(newPayload, frameMap[i]...)
				}

				output <- pipeline.Frame{ID: int(lastFrameID), Payload: newPayload}
			}

			resetState := func() {
				frameMap = make(map[int][]byte)
				expectedIndex = config.StartIndex
				lastSeenIndex = math.MinInt32
			}

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in IndexExtractor on frame %d: %v", frame.ID, r)
							resetState()
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					receivedIndex, ok := config.IndexProvider(frame.Payload)
					if !ok {
						errChan <- fmt.Errorf("IndexExtractor: IndexProvider failed for frame %d", frame.ID)
						return
					}

					if lastSeenIndex == receivedIndex {
						return
					}

					extractChunk := func(p []byte) []byte {
						if config.NoOfBytes == -1 {
							return p[config.StartByte:]
						}
						return p[config.StartByte : config.StartByte+config.NoOfBytes]
					}

					if lastSeenIndex > receivedIndex && receivedIndex == config.StartIndex {
						finalIndex := expectedIndex - 1
						if config.EndIndex != -1 && finalIndex == config.EndIndex {
							sendFrame(finalIndex, uint64(frame.ID))
						}
						resetState()
					}

					if expectedIndex == receivedIndex {
						lastSeenIndex = receivedIndex
						frameMap[expectedIndex] = extractChunk(frame.Payload)
						expectedIndex++
						if config.EndIndex != -1 && expectedIndex > config.EndIndex {
							sendFrame(config.EndIndex, uint64(frame.ID))
							resetState()
						}
					} else {
						errChan <- fmt.Errorf("IndexExtractor: sequence broken. Expected index %d, got %d", expectedIndex, receivedIndex)
						resetState()
					}
				}()
			}
		}()

		return output
	}
}
