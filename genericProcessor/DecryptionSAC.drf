package stages

import (
	"context"
	"csrspServer/cgo/decryptionsac"
	"csrspServer/pipeline"
	"fmt"
)

// DecryptionSACConfig holds the configuration for the SAC Decryption stage.
type DecryptionSACConfig struct {
	KeyFilePath        string
	BatchSize          int
	DefaultFrameLength int

	// Bypass configuration
	BypassProvider func(payload []byte) (int, bool)
	BypassValue    int
}

// NewDecryptionSACStage creates a new one-to-one stage that decrypts frames using the SAC library.
func NewDecryptionSACStage(config DecryptionSACConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic in DecryptionSAC stage: %v", r)
				}
				close(output)
			}()

			// --- Decryptor Initialization ---
			decryptor, err := decryptionsac.New(config.KeyFilePath)
			if err != nil {
				errChan <- fmt.Errorf("failed to initialize SAC decryptor: %w", err)
				return
			}
			defer decryptor.Close()
			// --- End Initialization ---

			var buffer []byte
			var frameLength = config.DefaultFrameLength
			var isSynced = false

			decryptAndSend := func(payload []byte, frameLen int) {
				if len(payload) == 0 || frameLen == 0 {
					return
				}
				decryptedData, err := decryptor.Decrypt(payload)
				if err != nil {
					errChan <- fmt.Errorf("DecryptionSAC: decrypt call failed: %w", err)
					return // Stop processing if decryption fails
				}

				// Re-chunk the decrypted data back into frames
				for len(decryptedData) >= frameLen {
					output <- pipeline.Frame{Payload: decryptedData[:frameLen]}
					decryptedData = decryptedData[frameLen:]
				}
			}

			for frame := range input {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if frameLength == 0 {
					frameLength = len(frame.Payload)
				}

				bypass, ok := config.BypassProvider(frame.Payload)
				if ok && bypass == config.BypassValue {
					output <- frame
					continue
				}

				// Phase 1: Wait for a sync frame.
				if !isSynced {
					if decryptor.IsSynced(frame.Payload) {
						isSynced = true
					}
					continue // Drop frames until synced.
				}

				// Phase 2: Once synced, buffer frames for batching.
				buffer = append(buffer, frame.Payload...)
				if len(buffer)/frameLength >= config.BatchSize {
					// Decrypt a full batch
					batchSizeInBytes := (len(buffer) / frameLength) * frameLength
					decryptAndSend(buffer[:batchSizeInBytes], frameLength)
					buffer = buffer[batchSizeInBytes:] // Keep the remainder
				}
			}

			// Process the final partial buffer, which the original code ignored.
			decryptAndSend(buffer, frameLength)
		}()

		return output
	}
}
