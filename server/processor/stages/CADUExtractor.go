package stages

import (
	"bytes"
	"context"
	"csrsp/server/pipeline"
	"fmt"
)

// CADUExtractorConfig holds the configuration for a CADUExtractor stage.
type CADUExtractorConfig struct {
	// The Frame Synchronization Code to search for.
	FSC []byte
	// The fixed length of the CADU frame to extract.
	CADUFrameLength int

	// --- Routing Configuration ---
	// Provider to get the VCID (Virtual Channel ID) from a potential frame.
	VCIDProvider func(payload []byte) (int, bool)
	// Maps the extracted VCID to a specific output channel index.
	VCIDToOutputIndex map[int]int

	// --- Idle Frame Configuration ---
	// Provider to check if a frame is an idle frame.
	IdleValueProvider func(payload []byte) (string, bool)
	// The specific value that indicates a valid, non-idle frame.
	ValidIdleValue string

	// --- HSTM Configuration ---
	HasHSTM           bool
	HSTMCheckProvider func(payload []byte) (bool, bool)
	HSTMOutputIndex   int
	HSTMNoOfWords     int

	// --- General Configuration ---
	// The number of output channels this stage will produce.
	NumOutputs int
	StartWord  int
	DataSize   int
}

// NewCADUExtractorStage creates a new one-to-many stage that extracts and routes CADU frames.
func NewCADUExtractorStage(config CADUExtractorConfig, errChan chan<- error) pipeline.StageOneToMany {
	return func(ctx context.Context, input <-chan pipeline.Frame) []<-chan pipeline.Frame {
		outputs := make([]chan pipeline.Frame, config.NumOutputs)
		readOnlyOutputs := make([]<-chan pipeline.Frame, config.NumOutputs)
		for i := 0; i < config.NumOutputs; i++ {
			outputs[i] = make(chan pipeline.Frame)
			readOnlyOutputs[i] = outputs[i]
		}

		go func() {
			defer func() {
				for _, c := range outputs {
					close(c)
				}
			}()

			var buffer []byte

			for frame := range input {
				select {
				case <-ctx.Done():
					return
				default:
				}

				// --- Start of Business Logic ---
				idleValue, ok := config.IdleValueProvider(frame.Payload)
				if !ok || idleValue != config.ValidIdleValue {
					continue // Skip idle or invalid frames.
				}

				sw, ds := config.StartWord, config.DataSize
				if config.HasHSTM {
					isHSTM, ok := config.HSTMCheckProvider(frame.Payload)
					if ok && isHSTM {
						sw += config.HSTMNoOfWords
						ds -= config.HSTMNoOfWords
						if config.HSTMOutputIndex < len(outputs) {
							outputs[config.HSTMOutputIndex] <- pipeline.Frame{ID: frame.ID, Payload: frame.Payload[config.StartWord : config.StartWord+config.HSTMNoOfWords]}
						}
					}
				}

				buffer = append(buffer, frame.Payload[sw:sw+ds]...)

				for {
					index := bytes.Index(buffer, config.FSC)
					if index == -1 {
						if len(buffer) > len(config.FSC) {
							buffer = buffer[len(buffer)-len(config.FSC):]
						}
						break
					}

					// A potential frame starts at the index. We need enough data for a full frame.
					if len(buffer) < index+config.CADUFrameLength {
						// Not enough data. Keep the buffer from the sync code onwards.
						buffer = buffer[index:]
						break
					}

					// Extract the full CADU frame.
					caduPayload := buffer[index : index+config.CADUFrameLength]

					// Route the frame based on its VCID.
					vcid, ok := config.VCIDProvider(caduPayload)
					if !ok {
						errChan <- fmt.Errorf("CADUExtractor: VCIDProvider failed for frame %d", frame.ID)
					} else {
						if outputIndex, exists := config.VCIDToOutputIndex[vcid]; exists {
							outputs[outputIndex] <- pipeline.Frame{ID: frame.ID, Payload: caduPayload}
						} else {
							// Optional: Log or send error for unroutable VCIDs
						}
					}

					// Move buffer past the extracted frame.
					buffer = buffer[index+config.CADUFrameLength:]
				}
				// --- End of Business Logic ---
			}
		}()

		return readOnlyOutputs
	}
}
