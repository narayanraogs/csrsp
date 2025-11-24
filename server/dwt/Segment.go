package dwt

// Header1A represents the first part of the DWT segment header.
type Header1A struct {
	StartImg     bool
	EndImg       bool
	SegmentCount int
	BitDepthDC   int
	BitDepthAC   int
	Part2Flag    bool
	Part3Flag    bool
	Part4Flag    bool
}

// Header1B represents the second part of the DWT segment header.
type Header1B struct {
	PadRows int
}

// Header2 represents the third part of the DWT segment header.
type Header2 struct {
	SegByteLimit int
	DcStop       bool
	BitPlaneStop int
	StageStop    int
	UseFill      bool
}

// Header3 represents the fourth part of the DWT segment header.
type Header3 struct {
	SegmentSizeInBlocks int
	OptDCSelect         bool
	OptACSelect         bool
}

// Header4 represents the fifth part of the DWT segment header.
type Header4 struct {
	DWTtype               bool
	ExtendedPixelBitDepth bool
	SignedPixels          bool
	PixelBitDepth         int
	ImageWidth            int
	TransposeImg          bool
	CodeWordLength        int
	CustomWtFlag          bool
	CustomWt              []int
}

// Segment holds all the parsed information and data for a single DWT segment.
type Segment struct {
	OneA  Header1A
	OneB  Header1B
	Two   Header2
	Three Header3
	Four  Header4
	Data  []byte
}
