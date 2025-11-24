package stages

import (
	"context"
	"csrsp/server/pipeline"
	"log/slog"
	"runtime"
	"time"
)

type RSDecodingConfig struct {
	BlockSize  int
	NoOfBlocks int
	StartByte  int
	NoOfRoots  int
	FieldSize  int
	Polynomial int

	FrameIdentifier string
	FrameType       string

	// Channel for sending operational telemetry.
	TelemetryChan chan<- pipeline.TelemetryEvent
}

type rsJob struct {
	sequenceID uint64
	frame      pipeline.Frame
}

type rsResult struct {
	sequenceID  uint64
	newPayload  []byte
	decodeError error
}

func NewRSDecodingStage(config RSDecodingConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			logTable := generateLogTable(config.FieldSize, config.Polynomial)
			exponentTable := generateExponentTable(logTable, config.FieldSize)

			numWorkers := runtime.NumCPU()
			jobs := make(chan rsJob, numWorkers)
			results := make(chan rsResult, numWorkers)

			for i := 0; i < numWorkers; i++ {
				go rsWorker(ctx, config, jobs, results, logTable, exponentTable)
			}

			// Sequencer logic can be simplified and placed here directly.
			buffer := make(map[uint64]rsResult)
			var nextSequenceID uint64 = 0

			var currentSequenceID uint64 = 0
			for frame := range input {
				select {
				case jobs <- rsJob{sequenceID: currentSequenceID, frame: frame}:
					currentSequenceID++
				case <-ctx.Done():
					close(jobs)
					return
				}
			}
			close(jobs)

			for i := 0; i < int(currentSequenceID); i++ {
				select {
				case res := <-results:
					buffer[res.sequenceID] = res
					for {
						result, ok := buffer[nextSequenceID]
						if !ok {
							break
						}
						if result.decodeError != nil {
							errChan <- result.decodeError
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

func rsWorker(ctx context.Context, config RSDecodingConfig, jobs <-chan rsJob, results chan<- rsResult, logTable, exponentTable []int) {
	for job := range jobs {
		newPayload, err := handleRSDecoding(job.frame.Payload, int(job.sequenceID), config, logTable, exponentTable)
		select {
		case results <- rsResult{sequenceID: job.sequenceID, newPayload: newPayload, decodeError: err}:
		case <-ctx.Done():
			return
		}
	}
}

func handleRSDecoding(payload []byte, frameNo int, config RSDecodingConfig, logTable, exponentTable []int) ([]byte, error) {
	blocks := deinterleave(payload, config.NoOfBlocks, config.BlockSize, config.StartByte)
	decodedBlocks := make([][]byte, config.NoOfBlocks)
	for i := 0; i < len(blocks); i++ {
		temp := decode(blocks[i], frameNo, i, config, logTable, exponentTable)
		decodedBlocks[i] = temp
	}
	temp := interleave(decodedBlocks, config.NoOfBlocks, config.BlockSize)
	newPayload := make([]byte, 0, len(payload))
	newPayload = append(newPayload, payload[:config.StartByte]...)
	newPayload = append(newPayload, temp...)
	return newPayload, nil
}

// ... (The complex 'decode' function and its helpers: deinterleave, interleave, etc.)
// ... (This part is long, so I will replace the DB call inside it separately)

func deinterleave(frame []byte, noOfBlocks int, blockSize int, startByte int) [][]byte {
	blocks := make([][]byte, noOfBlocks)
	for i := 0; i < noOfBlocks; i++ {
		blocks[i] = make([]byte, blockSize)
	}
	var index = startByte
	for j := 0; j < blockSize; j++ {
		for i := 0; i < noOfBlocks; i++ {
			blocks[i][j] = frame[index]
			index++
		}
	}
	return blocks
}

func interleave(blocks [][]byte, noOfBlocks int, blockSize int) []byte {
	frame := make([]byte, blockSize*noOfBlocks)
	k := 0
	for i := 0; i < blockSize; i++ {
		for j := 0; j < noOfBlocks; j++ {
			frame[k] = blocks[j][i]
			k++
		}
	}
	return frame
}

func getMod255(value int) int {
	return value % 255
}

var decode = func(block []byte, frameNo int, blockNo int, config RSDecodingConfig,
	logTable []int, exponentTable []int) []byte {
	blockSize := config.BlockSize
	noOfRoots := config.NoOfRoots
	decodedBlock := make([]int, blockSize)
	for i := 0; i < blockSize; i++ {
		decodedBlock[i] = logTable[block[blockSize-1-i]]
	}
	s := make([]int, noOfRoots+1)
	var syndromeError = 0
	for i := 1; i <= noOfRoots; i++ {
		s[i] = 0
		for j := 0; j < blockSize; j++ {
			if decodedBlock[j] != -1 {
				s[i] ^= exponentTable[getMod255(decodedBlock[j]+i*j)]
			}
		}
		if s[i] != 0 {
			syndromeError = 1 // set flag if non-zero syndrome => error
		}
		s[i] = logTable[s[i]] // convert syndrome from polynomial form to index form
	}
	flag := 0
	cnt := 0
	if syndromeError == 1 {
		elp := make([][]int, noOfRoots+2)
		for i := 0; i < len(elp); i++ {
			elp[i] = make([]int, noOfRoots)
		}
		d := make([]int, noOfRoots+2)
		l := make([]int, noOfRoots+2)
		uLu := make([]int, noOfRoots+2)
		root := make([]int, noOfRoots/2)
		loc := make([]int, noOfRoots/2)
		z := make([]int, (noOfRoots/2)+1)
		reg := make([]int, (noOfRoots/2)+1)
		err := make([]int, blockSize)
		count := 0
		q := 0
		u := 0
		index := 0
		flag = 1
		d[0] = 0
		d[1] = s[1]
		elp[0][0] = 0
		elp[1][0] = 1
		for i := 1; i < noOfRoots; i++ {
			elp[0][i] = -1
			elp[1][i] = 0
		}
		l[0] = 0
		l[1] = 0
		uLu[0] = -1
		uLu[1] = 0
		for {
			u++
			if d[u] == -1 {
				l[u+1] = l[u]
				for i := 0; i <= l[u]; i++ {
					elp[u+1][i] = elp[u][i]
					elp[u][i] = logTable[elp[u][i]]
				}
			} else {
				q = u - 1
				for (d[q] == -1) && (q > 0) {
					q--
				}
				if q > 0 {
					index = q
					for {
						index--
						if d[index] != -1 && uLu[q] < uLu[index] {
							q = index
						}
						if !(index > 0) {
							break
						}
					}
				}
				if l[u] > l[q]+u-q {
					l[u+1] = l[u]
				} else {
					l[u+1] = l[q] + u - q
				}
				for i := 0; i < noOfRoots; i++ {
					elp[u+1][i] = 0
				}
				for i := 0; i <= l[q]; i++ {
					if elp[q][i] != -1 {
						elp[u+1][i+u-q] = exponentTable[getMod255(d[u]+blockSize-d[q]+elp[q][i])]
					}
				}
				for i := 0; i <= l[u]; i++ {
					elp[u+1][i] ^= elp[u][i]
					elp[u][i] = logTable[elp[u][i]]
				}
			}
			uLu[u+1] = u - l[u+1]
			if u < noOfRoots {
				if s[u+1] != -1 {
					d[u+1] = exponentTable[s[u+1]]
				} else {
					d[u+1] = 0
				}
				for i := 1; i <= l[u+1]; i++ {
					if s[u+1-i] != -1 && elp[u+1][i] != 0 {
						d[u+1] ^= exponentTable[getMod255(s[u+1-i]+logTable[elp[u+1][i]])]
					}
				}
				d[u+1] = logTable[d[u+1]]
			}
			if !(u < noOfRoots && l[u+1] <= (noOfRoots/2)) {
				break
			}
		}
		u++
		if l[u] <= 4 {
			for i := 0; i <= l[u]; i++ {
				elp[u][i] = logTable[elp[u][i]]
			}
			for i := 1; i <= l[u]; i++ {
				reg[i] = elp[u][i]
			}
			count = 0
			for i := 1; i <= blockSize; i++ {
				q = 1
				for j := 1; j <= l[u]; j++ {
					if reg[j] != -1 {
						reg[j] = getMod255(reg[j] + j)
						q ^= exponentTable[reg[j]]
					}
				}
				if q == 0 {
					root[count] = i
					loc[count] = blockSize - i
					count++
					cnt++
				}
			}
			if count == l[u] {
				for i := 1; i <= l[u]; i++ {
					if s[i] != -1 && elp[u][i] != -1 {
						z[i] = exponentTable[s[i]] ^ exponentTable[elp[u][i]]
					} else if s[i] != -1 && elp[u][i] == -1 {
						z[i] = exponentTable[s[i]]
					} else if s[i] == -1 && elp[u][i] != -1 {
						z[i] = exponentTable[elp[u][i]]
					} else {
						z[i] = 0
					}
					for j := 1; j < i; j++ {
						if s[j] != -1 && elp[u][i-j] != -1 {
							z[i] ^= exponentTable[getMod255(elp[u][i-j]+s[j])]
						}
					}
					z[i] = logTable[z[i]]
				}
				for i := 0; i < blockSize; i++ {
					err[i] = 0
					if decodedBlock[i] != -1 {
						decodedBlock[i] = exponentTable[decodedBlock[i]]
					} else {
						decodedBlock[i] = 0
					}
				}
				for i := 0; i < l[u]; i++ {
					err[loc[i]] = 1
					for j := 1; j <= l[u]; j++ {
						if z[j] != -1 {
							err[loc[i]] ^= exponentTable[getMod255(z[j]+j*root[i])]
						}
					}
					if err[loc[i]] != 0 {
						err[loc[i]] = logTable[err[loc[i]]]
						q = 0
						for j := 0; j < l[u]; j++ {
							if j != i {
								q += logTable[1^exponentTable[getMod255(loc[j]+root[i])]]
							}
						}
						q = getMod255(q)
						err[loc[i]] = exponentTable[getMod255(err[loc[i]]-q+blockSize)]
						decodedBlock[loc[i]] ^= err[loc[i]]
					}
				}
			} else {
				flag = 2
			}
		} else {
			flag = 2
		}
	} else {
		flag = 0
	}
	if flag == 0 || flag == 2 {
		for i := 0; i < blockSize; i++ {
			if decodedBlock[i] != -1 {
				decodedBlock[i] = exponentTable[decodedBlock[i]]
			} else {
				decodedBlock[i] = 0
			}
		}
	}
	correctedBlock := make([]byte, blockSize)
	for i := 0; i < blockSize; i++ {
		correctedBlock[i] = byte(decodedBlock[blockSize-1-i])
	}

	if config.TelemetryChan != nil {
		switch flag {
		case 1: // Correctable
			slog.Warn("RS Correctable Errors", slog.Int("FrameNo", frameNo), slog.Int("BlockNo", blockNo), "FrameType", config.FrameType, "FrameIdentifier", config.FrameIdentifier, slog.Int("errors", cnt))
			config.TelemetryChan <- pipeline.TelemetryEvent{
				StageName: "RSDecoding",
				Timestamp: time.Now(),
				EventType: pipeline.EventTypeRSCorrection,
				Severity:  pipeline.SeverityWarning,
				Details: pipeline.RSCorrectionDetails{
					FrameIdentifier: config.FrameIdentifier,
					LineNumber:      frameNo,
					BlockNumber:     blockNo,
					ErrorsCorrected: cnt,
					IsUncorrectable: false,
				},
			}
		case 2: // Uncorrectable
			slog.Error("RS uncorrectable error", slog.Int("Line", frameNo), slog.Int("Block", blockNo), "Frame Type", config.FrameType, "FrameIdentifer", config.FrameIdentifier)
			config.TelemetryChan <- pipeline.TelemetryEvent{
				StageName: "RSDecoding",
				Timestamp: time.Now(),
				EventType: pipeline.EventTypeRSCorrection,
				Severity:  pipeline.SeverityError,
				Details: pipeline.RSCorrectionDetails{
					FrameIdentifier: config.FrameIdentifier,
					LineNumber:      frameNo,
					BlockNumber:     blockNo,
					ErrorsCorrected: 0,
					IsUncorrectable: true,
				},
			}
		}
	}
	return correctedBlock
}

func generateLogTable(fieldSize int, polynomial int) []int {
	logTable := make([]int, fieldSize)
	for i := 0; i < fieldSize; i++ {
		logTable[i] = -1
	}
	b := 1
	for l := 0; l < fieldSize-1; l++ {
		if logTable[b] != -1 {
		}
		logTable[b] = l
		b = b << 1
		if fieldSize <= b {
			b = (b - fieldSize) ^ polynomial
		}
	}
	logTable[0] = -1
	return logTable
}

func generateExponentTable(logTable []int, fieldSize int) []int {
	exponentTable := make([]int, fieldSize)
	for i := 1; i < fieldSize; i++ {
		logValue := logTable[i]
		if logValue != -1 {
			exponentTable[logValue] = i
		}
	}
	return exponentTable
}
