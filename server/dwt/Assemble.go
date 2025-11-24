package dwt

import "log/slog"

func getAssembler() (func(dcac), func() [][]float64) {
	recoveredImage := make([][]float64, 0)
	numberOfBlocks := 0
	var locationOfBlocks []int

	var assemble = func(value dcac) {
		blocksPerSegment := value.segSizeInBlocks
		locationOfBlocks = initLocationOfBlocks(blocksPerSegment, numberOfBlocks, locationOfBlocks)
		recoveredImage = resizeImage(recoveredImage, value.imageWidth, value.dcStop, blocksPerSegment, numberOfBlocks)
		for block := 0; block < blocksPerSegment; block++ {
			putBlockInSegment(block, locationOfBlocks[block], value.dcStop, recoveredImage, value.dc2C, value.acs)
		}
		numberOfBlocks += blocksPerSegment
	}

	var getRecoveredImage = func() [][]float64 {
		return recoveredImage
	}
	return assemble, getRecoveredImage
}

func resizeImage(recoveredImage [][]float64, xSize int, DCStop bool, blocksPerSegment int, numberOfBlocks int) [][]float64 {
	yInitialSize := 0

	if recoveredImage != nil {
		yInitialSize = len(recoveredImage)
	}

	pixelsPerBlock := 1 << (wtLevels * 2)
	blocksAvailable := (xSize * yInitialSize) / pixelsPerBlock
	emptyBlocks := blocksAvailable - numberOfBlocks
	if emptyBlocks < blocksPerSegment {
		blocksToAdd := blocksPerSegment - emptyBlocks
		sizeSideBlock := 1 << (wtLevels)

		rowsToAdd := (blocksToAdd * sizeSideBlock) / xSize

		if (blocksToAdd*sizeSideBlock)%xSize != 0 {
			rowsToAdd++
		}
		rowsToAdd *= sizeSideBlock

		yExtendedSize := yInitialSize + rowsToAdd

		extendedChannel := make([][]float64, yExtendedSize)
		for i := 0; i < len(extendedChannel); i++ {
			extendedChannel[i] = make([]float64, xSize)
		}

		xSubbandSize := xSize >> (wtLevels)
		ySubbandSizeInitial := (yInitialSize >> (wtLevels))
		for y := 0; y < ySubbandSizeInitial; y++ {
			for x := 0; x < xSubbandSize; x++ {
				extendedChannel[y][x] = recoveredImage[y][x]
			}
		}
		if DCStop == false {
			for rLevel := 0; rLevel < resolutionLevels; rLevel++ {

				xSubbandSize := (xSize >> (wtLevels - rLevel))
				ySubbandSizeInitial := (yInitialSize >> (wtLevels - rLevel))
				ySubbandSizeExtended := (yExtendedSize >> (wtLevels - rLevel))

				x0 := xSubbandSize
				y0Initial := 0
				y0Extended := 0
				for y := 0; y < ySubbandSizeInitial; y++ {
					for x := x0; x < xSubbandSize+x0; x++ {
						extendedChannel[y+y0Extended][x] = recoveredImage[y+y0Initial][x]
					}
				}

				x0 = 0
				y0Initial = ySubbandSizeInitial
				y0Extended = ySubbandSizeExtended
				for y := 0; y < ySubbandSizeInitial; y++ {
					for x := x0; x < xSubbandSize+x0; x++ {
						extendedChannel[y+y0Extended][x] = recoveredImage[y+y0Initial][x]
					}
				}

				x0 = xSubbandSize
				y0Initial = ySubbandSizeInitial
				y0Extended = ySubbandSizeExtended
				for y := 0; y < ySubbandSizeInitial; y++ {
					for x := x0; x < xSubbandSize+x0; x++ {
						extendedChannel[y+y0Extended][x] = recoveredImage[y+y0Initial][x]
					}
				}
			}

		}
		recoveredImage = nil
		recoveredImage = extendedChannel
	}
	return recoveredImage
}

func putBlockInSegment(block int, blockNumber int, DCStop bool, recoveredImage [][]float64, DCs []int, acs map[int]map[ACKey]int) {
	xSize := len(recoveredImage[0])
	xResidualSubBandSize := (xSize >> wtLevels)
	x0 := blockNumber % xResidualSubBandSize
	y0 := blockNumber / xResidualSubBandSize
	recoveredImage[y0][x0] = float64(DCs[block])

	if wtLevels != 0 && DCStop == false {

		ySize := len(recoveredImage)
		squaredBlockSize := 1 << resolutionLevels

		if xSize%squaredBlockSize != 0 || ySize%squaredBlockSize != 0 {
			slog.Error("Bit Plane decoder cannot run with this image dimensions at the channel")
		}

		xSubBandSize := (xSize >> wtLevels)
		ySubBandSize := (ySize >> wtLevels)

		x0 = blockNumber % xSubBandSize
		y0 = blockNumber / ySubBandSize

		for rLevel := 0; rLevel < resolutionLevels; rLevel++ {
			sizeResolutionLevel := (1 << rLevel)
			xSubBandSize_r := (xSize >> (wtLevels - rLevel))
			ySubBandSize_r := (ySize >> (wtLevels - rLevel))

			xInit := x0*sizeResolutionLevel + xSubBandSize_r
			yInit := y0 * sizeResolutionLevel
			for y := 0; y < sizeResolutionLevel; y++ {
				for x := 0; x < sizeResolutionLevel; x++ {
					key := ACKey{RLevel: rLevel, Subband: 0, Y: y, X: x}
					recoveredImage[yInit+y][xInit+x] = float64(acs[block][key])
				}
			}

			xInit = x0 * sizeResolutionLevel
			yInit = y0*sizeResolutionLevel + ySubBandSize_r
			for y := 0; y < sizeResolutionLevel; y++ {
				for x := 0; x < sizeResolutionLevel; x++ {
					key := ACKey{RLevel: rLevel, Subband: 1, Y: y, X: x}
					recoveredImage[yInit+y][xInit+x] = float64(acs[block][key])
				}
			}

			xInit = x0*sizeResolutionLevel + xSubBandSize_r
			yInit = y0*sizeResolutionLevel + ySubBandSize_r
			for y := 0; y < sizeResolutionLevel; y++ {
				for x := 0; x < sizeResolutionLevel; x++ {
					key := ACKey{RLevel: rLevel, Subband: 2, Y: y, X: x}
					recoveredImage[yInit+y][xInit+x] = float64(acs[block][key])
				}
			}
		}
	}
}

func initLocationOfBlocks(blocksPerSegment int, numberOfBlocks int, locationOfBlocks []int) []int {
	toBeReturned := make([]int, blocksPerSegment)
	for block := 0; block < blocksPerSegment; block++ {
		toBeReturned[block] = numberOfBlocks + block
	}
	return toBeReturned
}
