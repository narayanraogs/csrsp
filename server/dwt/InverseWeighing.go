package dwt

import "math"

func weightingRequiredT(customWtFlag bool) bool {
	if !customWtFlag && wtType != 4 {
		return false
	}
	return true
}
func inverseWeightAfterTransform(imageSamples [][]float64, customWtFlag bool, customWeight []int, segSizeInBlocks int, imageWidth int) [][]float64 {

	if !weightingRequiredT(customWtFlag) {
		return imageSamples
	}

	ySize := len(imageSamples)
	xSize := len(imageSamples[0])

	xSubBandSize := make([]int, wtLevels+1)
	ySubBandSize := make([]int, wtLevels+1)

	xSubBandSize[wtLevels] = xSize
	ySubBandSize[wtLevels] = ySize

	for k := wtLevels - 1; k >= 0; k-- {
		xSubBandSize[k] = xSubBandSize[k+1]/2 + xSubBandSize[k+1]%2
		ySubBandSize[k] = ySubBandSize[k+1]/2 + ySubBandSize[k+1]%2
	}
	for currentLevel := 0; currentLevel < wtLevels; currentLevel++ {
		var weight float64

		if currentLevel == 0 {
			weight = getWeight(currentLevel*3+0, customWtFlag, customWeight)
			for y := 0; y < ySubBandSize[currentLevel]; y++ {
				for x := 0; x < xSubBandSize[currentLevel]; x++ {
					imageSamples[y][x] /= weight
				}
			}
		}

		weight = getWeight(currentLevel*3+1, customWtFlag, customWeight)
		for y := 0; y < ySubBandSize[currentLevel]; y++ {
			for x := xSubBandSize[currentLevel]; x < xSubBandSize[currentLevel+1]; x++ {
				imageSamples[y][x] /= weight
			}
		}

		weight = getWeight(currentLevel*3+2, customWtFlag, customWeight)
		for y := ySubBandSize[currentLevel]; y < ySubBandSize[currentLevel+1]; y++ {
			for x := 0; x < xSubBandSize[currentLevel]; x++ {
				imageSamples[y][x] /= weight
			}
		}

		weight = getWeight(currentLevel*3+3, customWtFlag, customWeight)
		for y := ySubBandSize[currentLevel]; y < ySubBandSize[currentLevel+1]; y++ {
			for x := xSubBandSize[currentLevel]; x < xSubBandSize[currentLevel+1]; x++ {
				imageSamples[y][x] /= weight
			}
		}
	}
	return imageSamples
}

func getWeight(subband int, customWtFlag bool, customWeight []int) float64 {
	var weight float64 = 1

	if !customWtFlag {
		if wtType == 4 {
			exponent := wtLevels - subband/3
			weight = math.Pow(2, float64(exponent))
		} else {
			weight = 1
		}
	} else if customWtFlag {
		weight = math.Pow(2, float64(customWeight[subband]))
	}
	return weight
}
