package stages

import (
	"context"
	"csrsp/server/pipeline"
	"fmt"
)

// APIDExtractorConfig holds the configuration for an APIDExtractor stage.
type APIDExtractorConfig struct {
	// --- VCDU Configuration ---
	FirstHeaderPtrProvider func(payload []byte) (int, bool)
	VCDUDataStart          int
	VCDUDataLength         int
	IdleFrameValue         int
	NoNewPacketValue       int

	// --- Space Packet Configuration ---
	SPHeaderLength       int
	APIDProvider         func(payload []byte) (int, bool)
	PacketLengthProvider func(payload []byte) (int, bool)

	// --- Routing Configuration ---
	// Maps an APID to a specific output channel index.
	APIDToOutputIndex map[int]int
	// Maps an APID to its expected length (0 for variable length).
	APIDToExpectedLength map[int]int
	// Maps an APID to its auxiliary data length.
	APIDToAuxLength map[int]int

	NumOutputs int
}

// NewAPIDExtractorStage creates a new one-to-many stage that extracts and routes Space Packets from a VCDU stream.
func NewAPIDExtractorStage(config APIDExtractorConfig, errChan chan<- error) pipeline.StageOneToMany {
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

				firstHdrPtr, ok := config.FirstHeaderPtrProvider(frame.Payload)
				if !ok {
					errChan <- fmt.Errorf("APIDExtractor: FirstHeaderPtrProvider failed for frame %d", frame.ID)
					continue
				}

				if firstHdrPtr == config.IdleFrameValue {
					continue // Skip idle frames
				}

				vcduPayload := frame.Payload[config.VCDUDataStart : config.VCDUDataStart+config.VCDUDataLength]

				if firstHdrPtr == config.NoNewPacketValue {
					// This VCDU contains no new packet, so append its entire payload to the buffer.
					buffer = append(buffer, vcduPayload...)
					continue
				}

				// A new packet starts in this VCDU. Append the data before the pointer.
				if firstHdrPtr > len(vcduPayload) {
					errChan <- fmt.Errorf("APIDExtractor: FirstHeaderPtr %d is out of bounds for VCDU payload length %d", firstHdrPtr, len(vcduPayload))
					continue
				}
				buffer = append(buffer, vcduPayload[:firstHdrPtr]...)

				// Now, process the buffer which is guaranteed to contain at least one full packet.
				for len(buffer) >= config.SPHeaderLength {
					spHeader := buffer[:config.SPHeaderLength]
					apid, ok := config.APIDProvider(spHeader)
					if !ok {
						errChan <- fmt.Errorf("APIDProvider failed")
						break
					}

					packetLen, ok := config.PacketLengthProvider(spHeader)
					if !ok {
						errChan <- fmt.Errorf("PacketLengthProvider failed")
						break
					}

					expectedLen, exists := config.APIDToExpectedLength[apid]
					if !exists {
						// Unidentified APID, cannot proceed. Clear buffer and wait for next sync.
						errChan <- fmt.Errorf("unidentified APID %d found", apid)
						buffer = vcduPayload[firstHdrPtr:]
						break
					}

					// If length is defined in DB, it must match. Otherwise, use length from packet.
					actualContentLength := 0
					if expectedLen > 0 {
						if expectedLen != packetLen {
							errChan <- fmt.Errorf("packet length mismatch for APID %d", apid)
						}
						actualContentLength = expectedLen - config.SPHeaderLength
					} else {
						actualContentLength = packetLen - config.SPHeaderLength
					}

					if len(buffer) < config.SPHeaderLength+actualContentLength {
						// Not enough data in the buffer for the full packet. Wait for more.
						buffer = append(buffer, vcduPayload[firstHdrPtr:]...)
						break
					}

					// Extract the full space packet.
					fullPacket := buffer[:config.SPHeaderLength+actualContentLength]

					// Route the packet.
					if outputIndex, exists := config.APIDToOutputIndex[apid]; exists {
						outputs[outputIndex] <- pipeline.Frame{ID: frame.ID, Payload: fullPacket}
					}

					// Move buffer past the extracted packet.
					buffer = buffer[config.SPHeaderLength+actualContentLength:]
				}
			}
		}()

		return readOnlyOutputs
	}
}
