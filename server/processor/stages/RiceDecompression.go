package stages

import (
	"context"
	"csrspServer/pipeline"
	"csrspServer/utils/binary"
	"fmt"
)

import (
	"runtime"
)

type RiceDecompressionConfig struct {
	CompressionStatusProvider func(payload []byte) (uint64, bool)
	RefCountProvider          func(payload []byte) (uint64, bool)
	NoOfBlocksProvider        func(payload []byte) (uint64, bool)
	BypassValue               int
	IsCompressionBypassPacked int
	StartByte                 int
	PixelDepth                int
	RefCountAvailSts          int
}

type riceJob struct {
	sequenceID uint64
	frame      pipeline.Frame
}

type riceResult struct {
	sequenceID uint64
	newPayload []byte
	decomError error
}

func NewRiceDecompressionStage(config RiceDecompressionConfig, errChan chan<- error) pipeline.StageOneToOne {
	return func(ctx context.Context, input <-chan pipeline.Frame) <-chan pipeline.Frame {
		output := make(chan pipeline.Frame)

		go func() {
			defer close(output)

			numWorkers := runtime.NumCPU()
			jobs := make(chan riceJob, numWorkers)
			results := make(chan riceResult, numWorkers)

			for i := 0; i < numWorkers; i++ {
				go riceWorker(ctx, config, jobs, results)
			}

			sequencerDone := make(chan struct{})
			go func() {
				riceSequenceResults(ctx, output, results, errChan)
				close(sequencerDone)
			}()

			var currentSequenceID uint64 = 0
		outerFor:
			for frame := range input {
				select {
				case jobs <- riceJob{sequenceID: currentSequenceID, frame: frame}:
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

func riceWorker(ctx context.Context, config RiceDecompressionConfig, jobs <-chan riceJob, results chan<- riceResult) {
	for job := range jobs {
		compressionStatus, ok := config.CompressionStatusProvider(job.frame.Payload)
		var newPayload []byte
		var err error
		if !ok {
			err = fmt.Errorf("CompressionStatusProvider failed")
		} else {
			if compressionStatus == uint64(config.BypassValue) {
				newPayload = handleRiceBypass(job.frame.Payload, config)
			} else {
				newPayload, err = handleRiceDecompression(job.frame.Payload, config)
			}
		}

		select {
		case results <- riceResult{sequenceID: job.sequenceID, newPayload: newPayload, decomError: err}:
		case <-ctx.Done():
			return
		}
	}
}

func handleRiceBypass(payload []byte, config RiceDecompressionConfig) []byte {
	if config.IsCompressionBypassPacked == 1 {
		return unpackPayload(payload[config.StartByte:], config.PixelDepth)
	}
	return payload
}

func handleRiceDecompression(payload []byte, config RiceDecompressionConfig) ([]byte, error) {
	var compdataBitIndex uint64 = 0
	var samplesPerBlock uint64 = 16
	var packedPixels []uint64
	var prevPixel uint64 = 0

	compDat := payload[config.StartByte:]
	pixelDepth := uint64(config.PixelDepth)

	refCount, ok := config.RefCountProvider(payload)
	if !ok {
		return nil, fmt.Errorf("RefCountProvider failed")
	}
	noOfBlocks, ok := config.NoOfBlocksProvider(payload)
	if !ok {
		return nil, fmt.Errorf("NoOfBlocksProvider failed")
	}

	optIDLen, optByPassVal := getOptIDLength(pixelDepth)
	if config.RefCountAvailSts == 1 {
		noOfBlocks *= refCount
	}

	getBitsToUint64 := getDataProvider(compDat)

	for blockIndex := uint64(0); blockIndex < noOfBlocks; blockIndex++ {
		var refAvail bool
		if (config.RefCountAvailSts == 1 && blockIndex%refCount == 0) || (config.RefCountAvailSts == 0 && blockIndex == 0) {
			refAvail = true
		} else {
			refAvail = false
		}

		optIDVal, ok := getBitsToUint64(compdataBitIndex, optIDLen)
		if !ok {
			return nil, fmt.Errorf("unable to extract option ID value")
		}
		compdataBitIndex += optIDLen
		optIDVal &= 0x1F

		if (optIDVal == 0) || (pixelDepth <= 8 && optIDVal > 7) || (pixelDepth > 8 && pixelDepth <= 16 && optIDVal > 15) || (pixelDepth > 16 && optIDVal <= 32 && optIDVal > 31) {
			return nil, fmt.Errorf("invalid option ID value %d for pixel depth %d", optIDVal, pixelDepth)
		}

		var mapValues, pixels []uint64
		var startMapIndex uint64 = 0

		if optIDVal == optByPassVal {
			for j := 0; j < int(samplesPerBlock); j++ {
				sampleVal, ok := getBitsToUint64(compdataBitIndex, pixelDepth)
				if !ok {
					return nil, fmt.Errorf("unable to extract sample pixel value")
				}
				mapValues = append(mapValues, sampleVal)
				compdataBitIndex += pixelDepth
			}
		} else {
			if refAvail {
				refPixelVal, ok := getBitsToUint64(compdataBitIndex, pixelDepth)
				if !ok {
					return nil, fmt.Errorf("unable to extract reference pixel value")
				}
				mapValues = append(mapValues, refPixelVal)
				compdataBitIndex += pixelDepth
				pixels = append(pixels, refPixelVal)
				prevPixel = refPixelVal
				startMapIndex = 1
			}
			values, index, err := getMapValues(getBitsToUint64, compdataBitIndex, samplesPerBlock-startMapIndex, optIDVal-1)
			if err != nil {
				return nil, err
			}
			compdataBitIndex = index
			mapValues = append(mapValues, values...)
		}

		for i := startMapIndex; i < samplesPerBlock; i++ {
			val := getPackedPixelValues(int32(prevPixel), int32(mapValues[i]), int32(pixelDepth))
			pixels = append(pixels, val)
			prevPixel = val
		}
		packedPixels = append(packedPixels, pixels...)
	}

	var outputPixels []byte
	noOfBytes := pixelDepth/8 + (pixelDepth%8+7)/8
	byteIndex := 4 - noOfBytes
	for _, pixel := range packedPixels {
		tempBytes := binary.Uint32ToBytesBE(uint32(pixel))
		outputPixels = append(outputPixels, tempBytes[byteIndex:]...)
	}

	finalPayload := make([]byte, 0, len(payload))
	finalPayload = append(finalPayload, payload[:config.StartByte]...)
	finalPayload = append(finalPayload, outputPixels...)

	if len(finalPayload) > 7 {
		outputPacketLength := len(finalPayload) - 7
		finalPayload[5] = byte(outputPacketLength & 0xFF)
		finalPayload[4] = byte((outputPacketLength >> 8) & 0xFF)
	}
	if len(finalPayload) > 25 {
		finalPayload[25] &= 0xCF
	}

	return finalPayload, nil
}

func getOptIDLength(pixelDepth uint64) (uint64, uint64) {
	if pixelDepth <= 8 {
		return 3, 7
	}
	if pixelDepth <= 16 {
		return 4, 15
	}
	return 5, 31
}

func getDataProvider(frame []byte) func(uint64, uint64) (uint64, bool) {
	bitWise := make([]bool, len(frame)*8)
	for i, d := range frame {
		for j := 0; j < 8; j++ {
			bitWise[i*8+j] = (d>>(7-j))&1 == 1
		}
	}
	return func(startBit, noOfBits uint64) (uint64, bool) {
		if startBit+noOfBits > uint64(len(bitWise)) {
			return 0, false
		}
		var value uint64 = 0
		for i := uint64(0); i < noOfBits; i++ {
			value <<= 1
			if bitWise[startBit+i] {
				value |= 1
			}
		}
		return value, true
	}
}

func getMapValues(getBits func(uint64, uint64) (uint64, bool), bitIndex, samples, klen uint64) ([]uint64, uint64, error) {
	var fsPositions []uint64
	var zerosCount uint64 = 1
	for uint64(len(fsPositions)) < samples {
		val, ok := getBits(bitIndex, 1)
		if !ok {
			return nil, 0, fmt.Errorf("failed to read FS bit")
		}
		bitIndex++
		if val == 0 {
			zerosCount++
		} else {
			fsPositions = append(fsPositions, zerosCount-1)
			zerosCount = 1
		}
	}
	var mapValues []uint64
	for i := uint64(0); i < samples; i++ {
		kVal, ok := getBits(bitIndex, klen)
		if !ok {
			return nil, 0, fmt.Errorf("failed to read k-value")
		}
		mapValues = append(mapValues, (fsPositions[i]<<klen)+kVal)
		bitIndex += klen
	}
	return mapValues, bitIndex, nil
}

func getPackedPixelValues(prevPixel, nextPixel, pixelDepth int32) uint64 {
	var theta, delta int32
	pixelMaxVal := (1 << pixelDepth) - 1
	maxValue := int32(pixelMaxVal) - prevPixel
	minValue := prevPixel
	if minValue <= maxValue {
		theta = minValue
	} else {
		theta = maxValue
	}
	if (nextPixel%2 == 0) && (nextPixel <= 2*theta) {
		delta = nextPixel / 2
	} else if (nextPixel%2 != 0) && (nextPixel <= 2*theta-1) {
		delta = -(nextPixel + 1) / 2
	} else if theta == minValue {
		delta = nextPixel - theta
	} else {
		delta = -(nextPixel - theta)
	}
	return uint64(prevPixel + delta)
}

func riceSequenceResults(ctx context.Context, output chan<- pipeline.Frame, results <-chan riceResult, errChan chan<- error) {
	buffer := make(map[uint64]riceResult)
	var nextSequenceID uint64 = 0
	for {
		result, ok := buffer[nextSequenceID]
		if ok {
			if result.decomError != nil {
				errChan <- fmt.Errorf("RiceCecompression failed for image sequence %d: %w", result.sequenceID, result.decomError)
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
