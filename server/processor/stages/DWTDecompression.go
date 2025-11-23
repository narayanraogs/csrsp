package stages

import (
	"context"
	"csrspServer/dwt"
	"csrspServer/pipeline"
	"csrspServer/utils/binary"
	"fmt"
)

// DWTDecompressionConfig holds the configuration for the DWT stage.
type DWTDecompressionConfig struct {
	MaxSegmentsPerImage int
	PixelDepth          int
}

// NewDWTDecompressionStage creates a new stage that performs DWT decompression.
func NewDWTDecompressionStage(config DWTDecompressionConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in DWTDecompression stage: %v", r)
				}
			}()

			segmentBuffer := make([][]byte, 0, config.MaxSegmentsPerImage)
			var currentImageID int

			flushBuffer := func() {
				if len(segmentBuffer) == 0 {
					return
				}

				decompressedImage, err := dwt.Decompress(segmentBuffer)
				if err != nil {
					errChan <- fmt.Errorf("DWT decompression failed for image starting with frame %d: %w", currentImageID, err)
					segmentBuffer = nil // Clear buffer on error
					return
				}

				payload := imageToPayload(decompressedImage, config.PixelDepth)
				output <- pipeline.Frame{ID: currentImageID, Payload: payload}
				segmentBuffer = nil
			}

			for frame := range input {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if len(segmentBuffer) == 0 {
					currentImageID = frame.ID
				}

				segmentBuffer = append(segmentBuffer, frame.Payload)

				// Check for end of image flag (assuming it's in the first 3 bytes of the payload)
				// This is a simplified check based on the V1 logic.
				endImgMask, _ := binary.NewContinuousMask(0, 1, 1)
				endImgFlag, _ := endImgMask.ExtractUint64(frame.Payload)

				if endImgFlag == 1 || len(segmentBuffer) >= config.MaxSegmentsPerImage {
					flushBuffer()
				}
			}

			// Flush any remaining segments in the buffer at the end of the stream.
			flushBuffer()
		}()

		return output
	}
}

// imageToPayload converts the float64 image data back into a byte payload.
func imageToPayload(frame [][]float64, pixelDepth int) []byte {
	var array []byte
	noOfBytes := pixelDepth/8 + (pixelDepth%8+7)/8
	index := 8 - noOfBytes
	for i := range frame {
		for j := range frame[i] {
			array = append(array, binary.Uint64ToBytesBE(uint64(frame[i][j]))[index:]...)
		}
	}
	return array
}
