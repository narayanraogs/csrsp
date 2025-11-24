package stages

import (
	"context"
	"csrspServer/cgo/decryptionwithop"
	"csrspServer/pipeline"
	"fmt"
	"runtime"
)

// DecryptionWithOperationConfig holds the configuration for the stage.
type DecryptionWithOperationConfig struct {
	OpCode   int
	KeyFile  []byte
	Mode     int
	InfoFile string
}

// --- Internal structs for the worker pool ---
type decryWithOpJob struct {
	sequenceID uint64
	frame      pipeline.Frame
}

type decryWithOpResult struct {
	sequenceID  uint64
	newPayload  []byte
	decodeError error
}

// NewDecryptionWithOperationStage creates a stage that decrypts frames using the with-operation library.
func NewDecryptionWithOperationStage(config DecryptionWithOperationConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			numWorkers := runtime.NumCPU()
			jobs := make(chan decryWithOpJob, numWorkers)
			results := make(chan decryWithOpResult, numWorkers)

			for i := 0; i < numWorkers; i++ {
				go decryWithOpWorker(ctx, config, jobs, results)
			}

			// This simple sequencer can be defined locally.
			buffer := make(map[uint64]decryWithOpResult)
			var nextSequenceID uint64 = 0

			// Start the input processing loop.
			var currentSequenceID uint64 = 0
			for frame := range input {
				select {
				case jobs <- decryWithOpJob{sequenceID: currentSequenceID, frame: frame}:
					currentSequenceID++
				case <-ctx.Done():
					close(jobs)
					return
				}
			}
			close(jobs)

			// Process all results from the workers.
			for i := 0; i < int(currentSequenceID); i++ {
				select {
				case res := <-results:
					buffer[res.sequenceID] = res
					for {
						result, ok := buffer[nextSequenceID]
						if !ok {
							break // Not the next in sequence yet.
						}
						if result.decodeError != nil {
							errChan <- fmt.Errorf("DecryptionWithOp failed for frame %d: %w", result.sequenceID, result.decodeError)
						} else {
							output <- pipeline.Frame{ID: int(result.sequenceID), Payload: result.newPayload}
						}
						delete(buffer, nextSequenceID)
						nextSequenceID++
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		return output
	}
}

func decryWithOpWorker(ctx context.Context, config DecryptionWithOperationConfig, jobs <-chan decryWithOpJob, results chan<- decryWithOpResult) {
	for job := range jobs {
		// Use a closure to capture the job for panic recovery.
		func(currentJob decryWithOpJob) {
			defer func() {
				if r := recover(); r != nil {
					results <- decryWithOpResult{sequenceID: currentJob.sequenceID, decodeError: fmt.Errorf("panic in decryWithOpWorker: %v", r)}
				}
			}()

			decryptedPayload, err := decryptionwithop.DecryptWithOp(1, currentJob.frame.Payload, config.KeyFile, config.Mode, config.OpCode, config.InfoFile)

			select {
			case results <- decryWithOpResult{sequenceID: currentJob.sequenceID, newPayload: decryptedPayload, decodeError: err}:
			case <-ctx.Done():
				return
			}
		}(job)
	}
}
