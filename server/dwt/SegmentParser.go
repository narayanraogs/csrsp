package dwt

import (
	"csrsp/server/utils/binary"
	"fmt"
)

// SegmentParser holds the state required to parse a stream of DWT segments.
// Storing the previous headers in a struct makes the parsing process stateful but thread-safe.
type SegmentParser struct {
	prevTwo   Header2
	prevThree Header3
	prevFour  Header4
}

// NewSegmentParser creates a new stateful parser.
func NewSegmentParser() *SegmentParser {
	return &SegmentParser{}
}

// Parse takes a raw byte payload and parses it into a Segment struct.
// It uses the stored state from previous segments to fill in missing headers.
func (p *SegmentParser) Parse(data []byte) (Segment, error) {
	var seg Segment
	var err error

	headerLength := 3
	provider := getDataProvider(data)

	header1ABytes, ok := provider(3)
	if !ok {
		return seg, fmt.Errorf("incomplete segment data: cannot read Header1A")
	}
	seg.OneA, err = p.parseHeader1A(header1ABytes)
	if err != nil {
		return seg, fmt.Errorf("failed to parse Header1A: %w", err)
	}

	if seg.OneA.EndImg {
		header1BBytes, ok := provider(1)
		if !ok {
			return seg, fmt.Errorf("incomplete segment data: cannot read Header1B")
		}
		headerLength++
		seg.OneB.PadRows = int((header1BBytes[0] >> 5) & 0x07)
	}

	if seg.OneA.Part2Flag {
		header2Bytes, ok := provider(5)
		if !ok {
			return seg, fmt.Errorf("incomplete segment data: cannot read Header2")
		}
		headerLength += 5
		seg.Two, err = p.parseHeader2(header2Bytes)
		if err != nil {
			return seg, fmt.Errorf("failed to parse Header2: %w", err)
		}
		p.prevTwo = seg.Two
	} else {
		seg.Two = p.prevTwo
	}

	if seg.OneA.Part3Flag {
		header3Bytes, ok := provider(3)
		if !ok {
			return seg, fmt.Errorf("incomplete segment data: cannot read Header3")
		}
		headerLength += 3
		seg.Three, err = p.parseHeader3(header3Bytes)
		if err != nil {
			return seg, fmt.Errorf("failed to parse Header3: %w", err)
		}
		p.prevThree = seg.Three
	} else {
		seg.Three = p.prevThree
	}

	if seg.OneA.Part4Flag {
		header4Bytes, ok := provider(8)
		if !ok {
			return seg, fmt.Errorf("incomplete segment data: cannot read Header4")
		}
		headerLength += 8
		seg.Four, err = p.parseHeader4(header4Bytes)
		if err != nil {
			return seg, fmt.Errorf("failed to parse Header4: %w", err)
		}
		p.prevFour = seg.Four
	} else {
		seg.Four = p.prevFour
	}

	seg.Data, _ = provider(seg.Two.SegByteLimit - headerLength)
	return seg, nil
}

// --- Private header parsing methods ---

func (p *SegmentParser) parseHeader1A(h []byte) (Header1A, error) {
	var r Header1A
	mask, _ := binary.NewContinuousMask(0, 0, 1)
	temp, _ := mask.ExtractUint64(h)
	r.StartImg = temp == 1
	mask, _ = binary.NewContinuousMask(0, 1, 1)
	temp, _ = mask.ExtractUint64(h)
	r.EndImg = temp == 1
	mask, _ = binary.NewContinuousMask(0, 2, 8)
	temp, _ = mask.ExtractUint64(h)
	r.SegmentCount = int(temp)
	mask, _ = binary.NewContinuousMask(1, 2, 5)
	temp, _ = mask.ExtractUint64(h)
	r.BitDepthDC = int(temp)
	mask, _ = binary.NewContinuousMask(1, 7, 5)
	temp, _ = mask.ExtractUint64(h)
	r.BitDepthAC = int(temp)
	mask, _ = binary.NewContinuousMask(2, 5, 1)
	temp, _ = mask.ExtractUint64(h)
	r.Part2Flag = temp == 1
	mask, _ = binary.NewContinuousMask(2, 6, 1)
	temp, _ = mask.ExtractUint64(h)
	r.Part3Flag = temp == 1
	mask, _ = binary.NewContinuousMask(2, 7, 1)
	temp, _ = mask.ExtractUint64(h)
	r.Part4Flag = temp == 1
	return r, nil
}

func (p *SegmentParser) parseHeader2(h []byte) (Header2, error) {
	var r Header2
	mask, _ := binary.NewContinuousMask(0, 0, 27)
	temp, _ := mask.ExtractUint64(h)
	r.SegByteLimit = int(temp)
	if r.SegByteLimit == 0 {
		r.SegByteLimit = 1 << 27
	}
	mask, _ = binary.NewContinuousMask(3, 3, 1)
	temp, _ = mask.ExtractUint64(h)
	r.DcStop = temp == 1
	mask, _ = binary.NewContinuousMask(3, 4, 5)
	temp, _ = mask.ExtractUint64(h)
	r.BitPlaneStop = int(temp)
	mask, _ = binary.NewContinuousMask(4, 1, 2)
	temp, _ = mask.ExtractUint64(h)
	r.StageStop = int(temp) + 1
	mask, _ = binary.NewContinuousMask(4, 3, 1)
	temp, _ = mask.ExtractUint64(h)
	r.UseFill = temp == 1
	return r, nil
}

func (p *SegmentParser) parseHeader3(h []byte) (Header3, error) {
	var r Header3
	mask, _ := binary.NewContinuousMask(0, 0, 20)
	temp, _ := mask.ExtractUint64(h)
	r.SegmentSizeInBlocks = int(temp)
	mask, _ = binary.NewContinuousMask(2, 4, 1)
	temp, _ = mask.ExtractUint64(h)
	r.OptDCSelect = temp == 1
	mask, _ = binary.NewContinuousMask(2, 5, 1)
	temp, _ = mask.ExtractUint64(h)
	r.OptACSelect = temp == 1
	return r, nil
}

func (p *SegmentParser) parseHeader4(h []byte) (Header4, error) {
	var r Header4
	mask, _ := binary.NewContinuousMask(0, 0, 1)
	temp, _ := mask.ExtractUint64(h)
	r.DWTtype = temp == 1
	mask, _ = binary.NewContinuousMask(0, 2, 1)
	temp, _ = mask.ExtractUint64(h)
	r.ExtendedPixelBitDepth = temp == 1
	mask, _ = binary.NewContinuousMask(0, 3, 1)
	temp, _ = mask.ExtractUint64(h)
	r.SignedPixels = temp == 1
	mask, _ = binary.NewContinuousMask(0, 4, 4)
	temp, _ = mask.ExtractUint64(h)
	r.PixelBitDepth = int(temp)
	if r.PixelBitDepth == 0 {
		r.PixelBitDepth = 16
	}
	mask, _ = binary.NewContinuousMask(1, 0, 20)
	temp, _ = mask.ExtractUint64(h)
	r.ImageWidth = int(temp)
	mask, _ = binary.NewContinuousMask(3, 4, 1)
	temp, _ = mask.ExtractUint64(h)
	r.TransposeImg = temp == 1
	mask, _ = binary.NewContinuousMask(3, 5, 3)
	temp, _ = mask.ExtractUint64(h)
	r.CodeWordLength = int(temp)
	mask, _ = binary.NewContinuousMask(4, 0, 1)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWtFlag = temp == 1
	r.CustomWt = make([]int, 10)
	mask, _ = binary.NewContinuousMask(4, 1, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[0] = int(temp)
	mask, _ = binary.NewContinuousMask(4, 3, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[1] = int(temp)
	mask, _ = binary.NewContinuousMask(4, 5, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[2] = int(temp)
	mask, _ = binary.NewContinuousMask(4, 7, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[3] = int(temp)
	mask, _ = binary.NewContinuousMask(5, 1, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[4] = int(temp)
	mask, _ = binary.NewContinuousMask(5, 3, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[5] = int(temp)
	mask, _ = binary.NewContinuousMask(5, 5, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[6] = int(temp)
	mask, _ = binary.NewContinuousMask(5, 7, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[7] = int(temp)
	mask, _ = binary.NewContinuousMask(6, 1, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[8] = int(temp)
	mask, _ = binary.NewContinuousMask(6, 3, 2)
	temp, _ = mask.ExtractUint64(h)
	r.CustomWt[9] = int(temp)
	return r, nil
}
