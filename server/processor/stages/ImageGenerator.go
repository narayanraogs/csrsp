package stages

import (
	"bytes"
	"context"
	"csrspServer/pipeline"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png" // Import to register PNG decoder
	"math"
	"os"
	"path/filepath"
)

type ImageType string

const (
	// ImageTypeRGB processes raw, interleaved RGB pixel data into a JPEG.
	ImageTypeRGB ImageType = "RGB"
	// ImageTypeGray processes raw, grayscale pixel data into a JPEG.
	ImageTypeGray ImageType = "GRAY"
	// ImageTypeDirectDecode decodes a standard image format (like PNG, GIF) and re-encodes it as a JPEG.
	ImageTypeDirectDecode ImageType = "DIRECT_DECODE"
	// ImageTypePassthroughJPEG writes the payload directly to a .jpg file, assuming it's already a valid JPEG.
	ImageTypePassthroughJPEG ImageType = "PASSTHROUGH_JPEG"
)

// ImageGeneratorConfig holds the configuration for the unified ImageGenerator sink.
type ImageGeneratorConfig struct {
	ImageType       ImageType
	OutputPath      string
	FrameType       string
	FrameIdentifier string

	// The following fields are only required for RGB and GRAY processing.
	NoOfRows         int
	NoOfColumns      int
	PackedPixelDepth int
	BytesPerPixel    int
}

func NewImageGeneratorSink(config ImageGeneratorConfig, errChan chan<- error) pipeline.Sink {
	return func(ctx context.Context, input <-chan pipeline.Frame) {
		go func() {
			var frameCount = 0

			if err := os.MkdirAll(config.OutputPath, 0755); err != nil {
				errChan <- fmt.Errorf("ImageGenerator: could not create output path %s: %w", config.OutputPath, err)
				return
			}

			for frame := range input {
				func() {
					defer func() {
						if r := recover(); r != nil {
							errChan <- fmt.Errorf("panic in ImageGenerator on frame %d: %v", frame.ID, r)
						}
					}()

					select {
					case <-ctx.Done():
						return
					default:
					}

					fileName := fmt.Sprintf("%s_%s_%d.jpg", config.FrameType, config.FrameIdentifier, frameCount)
					fullPath := filepath.Join(config.OutputPath, fileName)
					var err error

					switch config.ImageType {
					case ImageTypePassthroughJPEG:
						err = os.WriteFile(fullPath, frame.Payload, 0644)
					case ImageTypeDirectDecode:
						err = decodeAndSave(frame.Payload, fullPath)
					case ImageTypeRGB:
						err = processRGB(frame.Payload, fullPath, config)
					case ImageTypeGray:
						err = processGray(frame.Payload, fullPath, config)
					default:
						err = fmt.Errorf("unsupported image type: %s", config.ImageType)
					}

					if err != nil {
						errChan <- fmt.Errorf("ImageGenerator: failed to process frame %d: %w", frame.ID, err)
					}

					frameCount++
				}()
			}
		}()
	}
}

func decodeAndSave(payload []byte, path string) error {
	img, _, err := image.Decode(bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}
	return saveJPEG(img, path)
}

func saveJPEG(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
}

func processRGB(payload []byte, path string, config ImageGeneratorConfig) error {
	img := image.NewRGBA(image.Rect(0, 0, config.NoOfColumns, config.NoOfRows))
	scale := 255.0 / math.Pow(2, float64(config.PackedPixelDepth))

	for row := 0; row < config.NoOfRows; row++ {
		for col := 0; col < config.NoOfColumns; col++ {
			index := (row*config.NoOfColumns + col) * config.BytesPerPixel * 3
			if index+config.BytesPerPixel*3 > len(payload) {
				return fmt.Errorf("payload access out of bounds at row %d, col %d", row, col)
			}

			r := float64(payload[index]) * scale
			g := float64(payload[index+config.BytesPerPixel]) * scale
			b := float64(payload[index+config.BytesPerPixel*2]) * scale

			img.SetRGBA(col, row, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255})
		}
	}
	return saveJPEG(img, path)
}

func processGray(payload []byte, path string, config ImageGeneratorConfig) error {
	img := image.NewGray(image.Rect(0, 0, config.NoOfColumns, config.NoOfRows))
	scale := 255.0 / math.Pow(2, float64(config.PackedPixelDepth))

	for row := 0; row < config.NoOfRows; row++ {
		for col := 0; col < config.NoOfColumns; col++ {
			index := (row*config.NoOfColumns + col) * config.BytesPerPixel
			if index+config.BytesPerPixel > len(payload) {
				return fmt.Errorf("payload access out of bounds at row %d, col %d", row, col)
			}

			val := float64(payload[index]) * scale
			img.SetGray(col, row, color.Gray{Y: uint8(val)})
		}
	}
	return saveJPEG(img, path)
}
