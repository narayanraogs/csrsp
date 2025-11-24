package dwt

import (
	"csrsp/server/utils/binary"
	"fmt"
	"log/slog"
	"math"
	"sync"
)

// Decompress takes a slice of raw segment payloads and returns the decompressed image data.
func Decompress(segmentsData [][]byte) ([][]float64, error) {
	parser := NewSegmentParser()
	segments := make([]Segment, len(segmentsData))
	for i, segData := range segmentsData {
		seg, err := parser.Parse(segData)
		if err != nil {
			return nil, fmt.Errorf("error parsing segment %d: %w", i, err)
		}
		segments[i] = seg
	}

	var wg sync.WaitGroup
	var segWg sync.WaitGroup
	var pixelDepth, imageWidth, segSizeInBlocks int

	assemble, getRecoveredImage := getAssembler()
	outputChannel := make(chan chan dcac, len(segments))

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for dcACChan := range outputChannel {
			temp := <-dcACChan
			pixelDepth = temp.pixelDepth
			imageWidth = temp.imageWidth
			segSizeInBlocks = temp.segSizeInBlocks
			assemble(temp)
		}
	}(&wg)

	segWg.Add(len(segments))
	for i := range segments {
		tempChan := make(chan dcac, 1)
		outputChannel <- tempChan
		go segments[i].decompress(tempChan, &segWg)
	}

	segWg.Wait()
	close(outputChannel)
	wg.Wait()

	recoveredImage := getRecoveredImage()
	if len(recoveredImage) == 0 {
		return nil, fmt.Errorf("decompression resulted in an empty image")
	}

	recoveredImage = inverseWeightAfterTransform(recoveredImage, false, make([]int, 10), segSizeInBlocks, imageWidth)
	recoveredImage = transform(recoveredImage, pixelDepth)
	recoveredImage = normalize(recoveredImage, pixelDepth)

	return recoveredImage, nil
}

// --- Internal Decompression Logic ---

type decompressor struct {
	seg    *Segment
	dc2C   []int // Add this field to hold state for refineDC
	wordNo int
	bitNo  int
}

// dcac is a data transfer object used to pass decoded data to the assembler.
type dcac struct {
	dc2C            []int
	acs             map[int]map[ACKey]int
	segSizeInBlocks int
	pixelDepth      int
	imageWidth      int
	dcStop          bool
}

func (seg *Segment) decompress(output chan dcac, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in DWT segment decompression", "error", r)
		}
	}()

	d := &decompressor{seg: seg}

	segmentStatus := make(map[ACKey]int)
	significantPyramid := make(map[ACKey]int)
	acs := make(map[ACKey]int)
	ds := make(map[int]int)
	dMap := make(map[int]map[int]int)

	for i := 0; i < seg.Three.SegmentSizeInBlocks; i++ {
		dMap[i] = make(map[int]int)
	}

	bitPlane := d.getBitPlane()
	quantization := d.getQuantizationFactor(bitPlane)
	dc2C := d.initialDecoding(gaggleDCSize, quantization, true)
	d.dc2C = dc2C // Store dc2C in the decompressor state
	dc := putDCs(d.dc2C, quantization, seg.OneA.BitDepthDC, bitPlane[0])

	for bp := quantization - 1; bp > seg.OneA.BitDepthAC-1 && bp >= bitPlane[0]; bp-- {
		d.refineDC(dc, bp, bitPlane[0])
	}

	bitDepthACBlock := d.initialDecoding(gaggleACSize, quantization, false)
	numberOfGaggles := seg.Three.SegmentSizeInBlocks/gaggleACSize + (seg.Three.SegmentSizeInBlocks%gaggleACSize+gaggleACSize-1)/gaggleACSize
	entropy, entReset := getEntropyWord(numberOfGaggles, d.provider)

	for bp := seg.OneA.BitDepthAC - 1; bp >= seg.Two.BitPlaneStop; bp-- {
		if bp < quantization && bp >= bitPlane[0] {
			d.refineDC(dc, bp, bitPlane[0])
		}
		entReset()
		for rLevel := 0; rLevel < resolutionLevels; rLevel++ {
			if bp <= seg.Two.BitPlaneStop && rLevel >= seg.Two.StageStop {
				d.sendResult(output, dc2C, acs)
				return
			}
			for block := 0; block < seg.Three.SegmentSizeInBlocks; block++ {
				if bitDepthACBlock[block] > bp {
					decode(block, bitPlane, rLevel, bp, ds, segmentStatus, significantPyramid, acs, dMap[block], entropy)
				}
			}
		}
		for block := 0; block < seg.Three.SegmentSizeInBlocks; block++ {
			for rLevel := 0; rLevel < resolutionLevels; rLevel++ {
				if bitDepthACBlock[block] > bp {
					refineACValues(bp, rLevel, block/gaggleACSize, segmentStatus, acs, bitPlane, entropy)
				}
			}
		}
	}
	d.sendResult(output, dc2C, acs)
}

func (d *decompressor) provider(noOfBits int) int {
	temp := d.wordNo*8 + d.bitNo + noOfBits
	endWord := temp / 8
	if temp%8 != 0 {
		endWord++
	}
	if endWord > len(d.seg.Data) {
		panic("bitstream provider read past end of data")
	}

	mask, _ := binary.NewContinuousMask(d.wordNo, d.bitNo, noOfBits)
	value, _ := mask.ExtractUint64(d.seg.Data)

	d.bitNo += noOfBits
	d.wordNo += d.bitNo / 8
	d.bitNo %= 8

	return int(value)
}

func (d *decompressor) sendResult(output chan dcac, dc2C []int, acs map[ACKey]int) {
	acsForAssembler := make(map[int]map[ACKey]int)
	// The V1 code processed one block at a time, so we replicate that by assigning all ACs to block 0.
	acsForAssembler[0] = acs

	var result dcac
	result.dc2C = dc2C
	result.acs = acsForAssembler
	result.dcStop = d.seg.Two.DcStop
	result.imageWidth = d.seg.Four.ImageWidth
	result.pixelDepth = d.seg.Four.PixelBitDepth
	result.segSizeInBlocks = d.seg.Three.SegmentSizeInBlocks
	output <- result
}

// ... other helper functions from V1 like initialDecoding, refineDC, etc. would be methods on *decompressor ...

func (d *decompressor) getBitPlane() []int {
	limit := 3*wtLevels + 1
	var bitPlane []int
	if d.seg.Four.CustomWtFlag {
		bitPlane = make([]int, 0)
		bitPlane = append(bitPlane, d.seg.Four.CustomWt...)
	} else if d.seg.Four.DWTtype {
		bitPlane = make([]int, limit)
		for i := 0; i < limit; i++ {
			bitPlane[i] = wtLevels - i/3
		}
	} else {
		bitPlane = make([]int, limit)
	}
	return bitPlane
}

func (d *decompressor) getQuantizationFactor(bp []int) int {
	var quant int
	quant = 1 + d.seg.OneA.BitDepthAC/2
	if d.seg.OneA.BitDepthAC == 0 || d.seg.OneA.BitDepthDC <= 3 {
		quant = 0
	} else if d.seg.OneA.BitDepthDC-quant <= 1 {
		quant = d.seg.OneA.BitDepthDC - 3
	}

	if quant < bp[0] {
		quant = bp[0]
	}
	if d.seg.OneA.BitDepthDC-quant > 10 {
		quant = d.seg.OneA.BitDepthDC - 10
	}
	return quant
}

func (d *decompressor) initialDecoding(gaggleSize int, quantization int, isDC bool) []int {
	var N int
	var min int
	var max int
	if isDC {
		N = 1
		if d.seg.OneA.BitDepthDC-quantization > N {
			N = d.seg.OneA.BitDepthDC - quantization
		}
		max = (1 << (N - 1)) - 1
		min = -(1 << (N - 1))

	} else {
		N = int(math.Log2(float64(1 + d.seg.OneA.BitDepthAC)))
		min = 0
		max = (1 << N) - 1
	}
	var quantizedDecodedValues = make([]int, 0)
	if N > 1 {
		mappedDPCM, _ := d.entropyDecode(0, N)
		quantizedDecodedValues = getDPCMDecodedValues(mappedDPCM, min, max)
	} else if N == 1 {
		quantizedDecodedValues = []int{}
		decodedValues := make([]int, d.seg.Three.SegmentSizeInBlocks)
		numberOfGaggles := d.seg.Three.SegmentSizeInBlocks / gaggleSize
		if d.seg.Three.SegmentSizeInBlocks%gaggleSize != 0 {
			numberOfGaggles++
		}
		for gaggle := 0; gaggle < numberOfGaggles; gaggle++ {
			for block := gaggle * gaggleDCSize; block < (gaggle+1)*gaggleSize && block < d.seg.Three.SegmentSizeInBlocks; block++ {
				decodedValues[block] = d.provider(N)
			}
		}
	}
	if isDC {
		quantizedDecodedValues = unQuantize(quantizedDecodedValues, N, quantization)
	}
	return quantizedDecodedValues
}

func (d *decompressor) entropyDecode(id int, N int) ([]int, int) {
	var returnValue = make([]int, d.seg.Three.SegmentSizeInBlocks)
	numberOfGaggles := d.seg.Three.SegmentSizeInBlocks / gaggleDCSize
	remain := d.seg.Three.SegmentSizeInBlocks % gaggleDCSize
	if remain != 0 {
		numberOfGaggles++
	}
	if id == 0 {
		id = d.seg.Three.SegmentSizeInBlocks
	}
	codeOptionLength := int(math.Log2(float64(N)))

	var gaggleNo = 0
	var prevGaggleNo = -1
	var codeOption int
	var uncodedOption int
	var readFirst = false

	for block := 0; block < d.seg.Three.SegmentSizeInBlocks; block++ {
		gaggleNo = block / gaggleDCSize
		if gaggleNo != prevGaggleNo {
			if readFirst {
				readFirst = false
				d.readSecond(prevGaggleNo, id, codeOption, returnValue)
			}
			prevGaggleNo = gaggleNo
			codeOption = d.provider(codeOptionLength)
			uncodedOption = 1<<codeOptionLength - 1
		}
		if codeOption == uncodedOption {
			returnValue[block] = d.provider(codeOption)
		} else {
			if block%id == 0 {
				returnValue[block] = d.provider(N)
			} else {
				readFirst = true
				firstPart := 0
				for {
					val := d.provider(1)
					if val == 1 {
						break
					}
					firstPart++
				}
				returnValue[block] = firstPart << codeOption
			}
		}
	}
	if readFirst {
		d.readSecond(prevGaggleNo, id, codeOption, returnValue)
	}
	return returnValue, id
}

func (d *decompressor) readSecond(gaggleNo int, id int, codeOption int, values []int) {
	start := gaggleNo * gaggleDCSize
	end := start + gaggleDCSize
	if end > len(values) {
		end = len(values)
	}
	for block := start; block < end; block++ {
		if block%id != 0 {
			secondPart := d.provider(codeOption)
			values[block] = values[block] + secondPart
		}
	}
}

func getDPCMDecodedValues(delta []int, min int, max int) []int {
	returnValue := make([]int, len(delta))
	for block := 0; block < len(delta); block++ {
		if block == 0 {
			if delta[block] >= (max + 1) {
				returnValue[block] = delta[block] - 2*(max+1)
			} else {
				returnValue[block] = delta[block]
			}
		} else {
			theta := max - returnValue[block-1]
			if theta > (returnValue[block-1] - min) {
				theta = returnValue[block-1] - min
			}

			if delta[block] <= 2*theta && delta[block] >= 0 {
				if (delta[block] % 2) == 0 {
					returnValue[block] = delta[block] / 2
				} else {
					returnValue[block] = -(delta[block] + 1) / 2
				}
			} else {
				if theta == returnValue[block-1]-min {
					returnValue[block] = delta[block] - theta
				} else {
					returnValue[block] = theta - delta[block]
				}
			}
			returnValue[block] = returnValue[block] + returnValue[block-1]
		}
	}
	return returnValue
}

func unQuantize(decodedValues []int, N int, quantizationFactor int) []int {
	maxValue := 1 << N
	returnValue := make([]int, len(decodedValues))
	for block := 0; block < len(decodedValues); block++ {
		if decodedValues[block] < 0 {
			returnValue[block] = decodedValues[block] + int(maxValue)
		}
		returnValue[block] = decodedValues[block] << quantizationFactor
	}
	return returnValue
}

func putDCs(dc2C []int, quantization int, bitDepthDC int, bitPlane int) []int {
	var returnValue = make([]int, len(dc2C))
	for block := 0; block < len(dc2C); block++ {
		returnValue[block] = putDCBlock(quantization, bitDepthDC, dc2C[block], bitPlane)
	}
	return returnValue
}

func putDCBlock(quantization int, bitDepthDC int, dc2C int, bitPlane int) int {
	threshold := math.Pow(2, float64(quantization))
	var dc int
	N := 1 << bitDepthDC
	var gamma float64

	if dc2C > bitPlane {
		gamma = math.Round(gammaValue * threshold)
	}
	if dc2C < N/2 {
		dc = dc2C + int(gamma)
	} else {
		dc = dc2C - N - int(gamma)
	}
	return dc
}

func (d *decompressor) refineDC(dc []int, bitPlane int, bp int) {
	threshold := 1 << bitPlane
	for block := 0; block < len(dc); block++ {
		temp := d.provider(1)
		if temp == 1 {
			d.dc2C[block] += int(threshold)
		}
		dc[block] = putDCBlock(bitPlane, d.seg.OneA.BitDepthDC, d.dc2C[block], bp)
	}
}

func normalize(data [][]float64, pixelDepth int) [][]float64 {
	max := float64((int(1) << uint(pixelDepth)) - 1)
	for i := range data {
		for j := range data[i] {
			if data[i][j] > max {
				data[i][j] = max
			}
			if data[i][j] < 0 {
				data[i][j] = 0.0
			}
			data[i][j] = math.Round(data[i][j])
		}
	}
	return data
}
