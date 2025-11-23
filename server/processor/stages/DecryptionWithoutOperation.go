package stages

import (
	"context"
	"csrspServer/cgo/decryption"
	"csrspServer/pipeline"
	"fmt"
	"runtime"
)

// DecryptionWithoutOperationConfig holds the configuration for the stage.
type DecryptionWithoutOperationConfig struct {
	KeyFile  []byte
	Mode     int
	InfoFile string
}

// --- Internal structs for the worker pool ---
type decryJob struct {
	sequenceID uint64
	frame      pipeline.Frame
}

type decryResult struct {
	sequenceID  uint64
	newPayload  []byte
	decodeError error
}

// NewDecryptionWithoutOperationStage creates a stage that decrypts frames using the standard library.
func NewDecryptionWithoutOperationStage(config DecryptionWithoutOperationConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			numWorkers := runtime.NumCPU()
			jobs := make(chan decryJob, numWorkers)
			results := make(chan decryResult, numWorkers)

			for i := 0; i < numWorkers; i++ {
				go decryWorker(ctx, config, jobs, results)
			}

			sequencerDone := make(chan struct{})
			go func() {
				// This simple sequencer can be defined locally.
				buffer := make(map[uint64]decryResult)
				var nextSequenceID uint64 = 0
			outerFor:
				for {
					result, ok := buffer[nextSequenceID]
					if ok {
						if result.decodeError != nil {
							errChan <- fmt.Errorf("decryption failed for frame sequence %d: %w", result.sequenceID, result.decodeError)
						} else {
							output <- pipeline.Frame{ID: int(result.sequenceID), Payload: result.newPayload}
						}
						delete(buffer, nextSequenceID)
						nextSequenceID++
						continue
					}
					select {
					case res, chanOpen := <-results:
						if !chanOpen {
							break outerFor // All workers are done.
						}
						buffer[res.sequenceID] = res
					case <-ctx.Done():
						break outerFor
					}
				}
				close(sequencerDone)
			}()

			var currentSequenceID uint64 = 0
		sequencerFor:
			for frame := range input {
				select {
				case jobs <- decryJob{sequenceID: currentSequenceID, frame: frame}:
					currentSequenceID++
				case <-ctx.Done():
					break sequencerFor
				}
			}

			close(jobs)
			<-sequencerDone
		}()

		return output
	}
}

func decryWorker(ctx context.Context, config DecryptionWithoutOperationConfig, jobs <-chan decryJob, results chan<- decryResult) {
	for job := range jobs {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					results <- decryResult{sequenceID: job.sequenceID, decodeError: fmt.Errorf("panic in decryWorker: %v", r)}
				}
			}()
			// The number of frames is assumed to be 1 for a single frame processing stage.
			decryptedPayload, err := decryption.Decrypt(1, job.frame.Payload, config.KeyFile, config.Mode, config.InfoFile)

			select {
			case results <- decryResult{sequenceID: job.sequenceID, newPayload: decryptedPayload, decodeError: err}:
			case <-ctx.Done():
				return
			}
		}()
	}
}
