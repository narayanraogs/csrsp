package stages

import (
	"context"
	"csrsp/server/pipeline"
	"csrsp/server/utils/binary"
	"fmt"
	"runtime"
)

// BayersDemosaicingConfig holds the configuration for a BayersDemosaicing stage.
type BayersDemosaicingConfig struct {
	NoOfRows      int
	NoOfColumns   int
	BytesPerPixel int
}

// --- Internal structs for the worker pool ---
type bayersJob struct {
	sequenceID uint64
	frame      pipeline.Frame
}

type bayersResult struct {
	sequenceID  uint64
	newPayload  []byte
	decodeError error
}

// NewBayersDemosaicingStage creates a new one-to-one stage that performs Bayers Demosaicing.
func NewBayersDemosaicingStage(config BayersDemosaicingConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			numWorkers := runtime.NumCPU()
			jobs := make(chan bayersJob, numWorkers)
			results := make(chan bayersResult, numWorkers)

			for i := 0; i < numWorkers; i++ {
				go bayersWorker(ctx, config, jobs, results)
			}

			sequencerDone := make(chan struct{})
			go func() {
				bayersSequenceResults(ctx, output, results, errChan) // Can reuse the same sequencer
				close(sequencerDone)
			}()

			var currentSequenceID uint64 = 0
		outerFor:
			for frame := range input {
				select {
				case jobs <- bayersJob{sequenceID: currentSequenceID, frame: frame}:
					currentSequenceID++
				case <-ctx.Done():
					break outerFor
				}
			}

			close(jobs)
			<-sequencerDone
		}()

		return output
	}
}

func bayersWorker(ctx context.Context, config BayersDemosaicingConfig, jobs <-chan bayersJob, results chan<- bayersResult) {
	for job := range jobs {
		newPayload, err := handleBayersDemosaicing(job.frame.Payload, config)
		select {
		case results <- bayersResult{sequenceID: job.sequenceID, newPayload: newPayload, decodeError: err}:
		case <-ctx.Done():
			return
		}
	}
}

func handleBayersDemosaicing(payload []byte, config BayersDemosaicingConfig) ([]byte, error) {
	expectedLen := config.NoOfRows * config.NoOfColumns * config.BytesPerPixel
	if len(payload) != expectedLen {
		return nil, fmt.Errorf("payload length %d does not match expected length %d", len(payload), expectedLen)
	}

	outputPayload := make([]byte, config.NoOfRows*config.NoOfColumns*config.BytesPerPixel*3) // 3 channels for RGB

	getPixelValue := func(row, col int) (uint64, bool) {
		if row < 0 || row >= config.NoOfRows || col < 0 || col >= config.NoOfColumns {
			return 0, false
		}
		index := (row*config.NoOfColumns + col) * config.BytesPerPixel
		value, err := binary.BytesToUint64BE(payload[index : index+config.BytesPerPixel])
		return value, err == nil
	}

	getInterpolatedValue := func(rows, cols []int) uint64 {
		var sum float64
		var count int
		for i := 0; i < len(rows); i++ {
			if temp, ok := getPixelValue(rows[i], cols[i]); ok {
				sum += float64(temp)
				count++
			}
		}
		if count == 0 {
			return 0
		}
		return uint64(sum / float64(count))
	}

	for row := 0; row < config.NoOfRows; row++ {
		for col := 0; col < config.NoOfColumns; col++ {
			var r, g, b uint64
			// Determine the color of the current pixel based on the Bayer pattern (RGGB is assumed)
			if row%2 == 0 {
				if col%2 == 0 { // Red pixel
					r, _ = getPixelValue(row, col)
					g = getInterpolatedValue([]int{row - 1, row + 1, row, row}, []int{col, col, col - 1, col + 1})
					b = getInterpolatedValue([]int{row - 1, row + 1, row - 1, row + 1}, []int{col - 1, col - 1, col + 1, col + 1})
				} else { // Green pixel (in a red row)
					g, _ = getPixelValue(row, col)
					r = getInterpolatedValue([]int{row, row}, []int{col - 1, col + 1})
					b = getInterpolatedValue([]int{row - 1, row + 1}, []int{col, col})
				}
			} else {
				if col%2 == 0 { // Green pixel (in a blue row)
					g, _ = getPixelValue(row, col)
					r = getInterpolatedValue([]int{row - 1, row + 1}, []int{col, col})
					b = getInterpolatedValue([]int{row, row}, []int{col - 1, col + 1})
				} else { // Blue pixel
					b, _ = getPixelValue(row, col)
					r = getInterpolatedValue([]int{row - 1, row + 1, row - 1, row + 1}, []int{col - 1, col - 1, col + 1, col + 1})
					g = getInterpolatedValue([]int{row - 1, row + 1, row, row}, []int{col, col, col - 1, col + 1})
				}
			}

			// Place the RGB values into the output payload
			outputIndex := (row*config.NoOfColumns + col) * config.BytesPerPixel * 3

			copy(outputPayload[outputIndex:], binary.Uint64ToBytesBE(r)[8-config.BytesPerPixel:])
			copy(outputPayload[outputIndex+config.BytesPerPixel:], binary.Uint64ToBytesBE(g)[8-config.BytesPerPixel:])
			copy(outputPayload[outputIndex+config.BytesPerPixel*2:], binary.Uint64ToBytesBE(b)[8-config.BytesPerPixel:])
		}
	}

	return outputPayload, nil
}

func bayersSequenceResults(ctx context.Context, output chan<- pipeline.Frame, results <-chan bayersResult, errChan chan<- error) {
	buffer := make(map[uint64]bayersResult)
	var nextSequenceID uint64 = 0
	for {
		result, ok := buffer[nextSequenceID]
		if ok {
			if result.decodeError != nil {
				errChan <- fmt.Errorf("DWTDecompression failed for image sequence %d: %w", result.sequenceID, result.decodeError)
			} else {
				output <- pipeline.Frame{ID: int(result.sequenceID), Payload: result.newPayload}
			}
			delete(buffer, nextSequenceID)
			nextSequenceID++
			continue
		}
		select {
		case result, chanOpen := <-results:
			if !chanOpen {
				if len(buffer) > 0 {
					errChan <- fmt.Errorf("DWTDecompression: missing %d images at end of stream (e.g., %d)", len(buffer), nextSequenceID)
				}
				return
			}
			buffer[result.sequenceID] = result
		case <-ctx.Done():
			return
		}
	}
}
