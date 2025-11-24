package dwt

import (
	"math"
)

func transform(imageSamples [][]float64, pixelDepth int) [][]float64 {
	ySize := len(imageSamples)
	xSize := len(imageSamples[0])
	if wtType != 0 && wtLevels > 0 {
		xSubBandSizes := make([]int, wtLevels)
		ySubBandSizes := make([]int, wtLevels)

		xSubBandSizes[wtLevels-1] = xSize
		ySubBandSizes[wtLevels-1] = ySize

		for k := wtLevels - 2; k >= 0; k-- {
			xSubBandSizes[k] = xSubBandSizes[k+1]/2 + xSubBandSizes[k+1]%2
			ySubBandSizes[k] = ySubBandSizes[k+1]/2 + ySubBandSizes[k+1]%2
		}
		for currentLevel := 0; currentLevel < wtLevels; currentLevel++ {
			xSubBandSize := xSubBandSizes[currentLevel]
			ySubBandSize := ySubBandSizes[currentLevel]
			if wtOrder == 2 {
				ySubBandSize = ySize
			}
			if wtOrder == 0 {
				for x := 0; x < xSubBandSize; x++ {
					currentColumn := make([]float64, ySubBandSize)
					for y := 0; y < ySubBandSize; y++ {
						currentColumn[y] = float64(imageSamples[y][x])
					}
					currentColumn = filtering(currentColumn)
					for y := 0; y < ySubBandSize; y++ {
						imageSamples[y][x] = currentColumn[y]
					}
				}
			}

			for y := 0; y < ySubBandSize; y++ {
				currentRow := make([]float64, xSubBandSize)
				for x := 0; x < xSubBandSize; x++ {
					currentRow[x] = imageSamples[y][x]
				}
				currentRow = filtering(currentRow)
				for x := 0; x < xSubBandSize; x++ {
					imageSamples[y][x] = currentRow[x]
				}
			}
			if wtOrder == 1 && wtOrder != 2 {
				for x := 0; x < xSubBandSize; x++ {
					currentColumn := make([]float64, 0)
					for y := 0; y < ySubBandSize; y++ {
						currentColumn[y] = imageSamples[y][x]
					}
					currentColumn = filtering(currentColumn)
					for y := 0; y < ySubBandSize; y++ {
						imageSamples[y][x] = currentColumn[y]
					}
				}
			}

		}
	}
	return imageSamples
}

func filtering(src []float64) []float64 {
	if len(src) == 1 {
		return src
	}
	if len(src)%2 == 0 {
		return evenFiltering(src)
	} else {
		return oddFiltering(src)
	}
}

func evenFiltering(src []float64) []float64 {
	subbandSize := len(src)
	half := subbandSize / 2
	dst := make([]float64, subbandSize)
	for k := 0; k < half; k++ {
		dst[2*k] = src[k]
		dst[2*k+1] = src[half+k]
	}
	if wtType == 1 {
		dst[0] = dst[0] - math.Floor(float64((dst[1]+dst[1]+2)/4))
		for k := 2; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] - math.Floor(float64((dst[k-1]+dst[k+1]+2)/4))
		}
		for k := 1; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] + math.Floor(float64((dst[k-1]+dst[k+1])/2))
		}
		dst[subbandSize-1] = dst[subbandSize-1] + math.Floor(float64(dst[subbandSize-2]+dst[subbandSize-2])/2)
	} else if wtType == 2 || wtType == 3 {
		var nh_97, nl_97 float64
		if wtType == 2 {
			nh_97 = 1.23017410491400
			nl_97 = 1 / nh_97
		} else {
			nl_97 = 1.14960430535816
			nh_97 = -1 / nl_97
		}

		for k := 0; k < subbandSize; k += 2 {
			dst[k] = (dst[k]) / nl_97
			dst[k+1] = dst[k+1] / nh_97
		}
		dst[0] = dst[0] - delta97*(dst[1]+dst[1])
		for k := 2; k < subbandSize; k += 2 {
			dst[k] = dst[k] - delta97*(dst[k-1]+dst[k+1])
		}
		for k := 1; k < subbandSize-2; k += 2 {
			dst[k] = dst[k] - gamma97*(dst[k-1]+dst[k+1])
		}
		dst[subbandSize-1] = dst[subbandSize-1] - gamma97*(dst[subbandSize-2]+dst[subbandSize-2])
		dst[0] = dst[0] - beta97*(dst[1]+dst[1])
		for k := 2; k < subbandSize; k += 2 {
			dst[k] = dst[k] - beta97*(dst[k-1]+dst[k+1])
		}
		for k := 1; k < subbandSize-2; k += 2 {
			dst[k] = dst[k] - alpha97*(dst[k-1]+dst[k+1])
		}
		dst[subbandSize-1] = dst[subbandSize-1] - alpha97*(dst[subbandSize-2]+dst[subbandSize-2])
	} else if wtType == 4 {
		if subbandSize >= 6 {
			var alpha1 float64 = 0.5625
			var alpha2 float64 = 0.0625
			var beta float64 = 0.25

			dst[0] = dst[0] + math.Floor(-beta*(dst[1]+dst[1])+0.5)
			for k := 2; k < subbandSize; k += 2 {
				dst[k] = dst[k] + math.Floor(-beta*(dst[k-1]+dst[k+1])+0.5)
			}
			dst[1] = dst[1] + math.Floor(alpha1*(dst[0]+dst[2])-alpha2*(dst[2]+dst[4])+0.5)
			for k := 3; k < subbandSize-3; k += 2 {
				dst[k] = dst[k] + math.Floor(alpha1*(dst[k-1]+dst[k+1])-alpha2*(dst[k-3]+dst[k+3])+0.5)
			}
			dst[subbandSize-3] = dst[subbandSize-3] + math.Floor(alpha1*(dst[subbandSize-4]+dst[subbandSize-2])-alpha2*(dst[subbandSize-6]+dst[subbandSize-2])+0.5)
			dst[subbandSize-1] = dst[subbandSize-1] + math.Floor(alpha1*(dst[subbandSize-2]+dst[subbandSize-2])-alpha2*(dst[subbandSize-4]+dst[subbandSize-4])+0.5)
		}
	} else if wtType == 5 || wtType == 6 {
		var alpha, beta, gamma, delta float64
		if wtType == 6 {
			alpha = -1.58615986717275
			beta = -0.05297864003258
			gamma = 0.88293362717904
			delta = 0.44350482244527
		} else {
			alpha = -0.5
			beta = 0.25
			gamma = 0
			delta = 0
		}

		if wtType == 6 {
			dst[0] = dst[0] - math.Floor(delta*(dst[1]+dst[1])+0.5)
			for k := 2; k < subbandSize; k += 2 {
				dst[k] = dst[k] - math.Floor(delta*(dst[k-1]+dst[k+1])+0.5)
			}
			for k := 1; k < subbandSize-2; k += 2 {
				dst[k] = dst[k] - math.Floor(gamma*(dst[k-1]+dst[k+1])+0.5)
			}
			dst[subbandSize-1] = dst[subbandSize-1] - math.Floor(gamma*(dst[subbandSize-2]+dst[subbandSize-2])+0.5)
		}
		dst[0] = dst[0] - math.Floor(beta*(dst[1]+dst[1])+0.5)
		for k := 2; k < subbandSize; k += 2 {
			dst[k] = dst[k] - math.Floor(beta*(dst[k-1]+dst[k+1])+0.5)
		}

		for k := 1; k < subbandSize-2; k += 2 {
			dst[k] = dst[k] - math.Floor(alpha*(dst[k-1]+dst[k+1])+0.5)
		}
		dst[subbandSize-1] = dst[subbandSize-1] - math.Floor(alpha*(dst[subbandSize-2]+dst[subbandSize-2])+0.5)

	}
	return dst
}

func oddFiltering(src []float64) []float64 {
	subbandSize := len(src)
	half := subbandSize / 2
	dst := make([]float64, subbandSize)
	for k := 0; k < half; k++ {
		dst[2*k] = src[k]
		dst[2*k+1] = src[half+k+1]
	}
	dst[subbandSize-1] = src[half]
	if wtType == 1 {
		dst[0] = dst[0] - math.Floor(((dst[1] + dst[1] + 2) / 4))
		for k := 2; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] - math.Floor(((dst[k-1] + dst[k+1] + 2) / 4))
		}
		dst[subbandSize-1] = dst[subbandSize-1] - math.Floor((dst[subbandSize-2]+dst[subbandSize-2]+2)/4)
		for k := 1; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] + math.Floor(((dst[k-1] + dst[k+1]) / 2))
		}
	} else if wtType == 2 || wtType == 3 {
		var nh_97, nl_97 float64
		if wtType == 2 {
			nh_97 = 1.230174104914001
			nl_97 = 1 / nh_97
		} else {
			nh_97 = 1.14960430535816
			nl_97 = 1 / nh_97
		}

		for k := 0; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] / nl_97
			dst[k+1] = dst[k+1] / nh_97
		}
		dst[subbandSize-1] = dst[subbandSize-1] / nl_97
		dst[0] = dst[0] - delta97*(dst[1]+dst[1])
		for k := 2; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] - delta97*(dst[k-1]+dst[k+1])
		}
		dst[subbandSize-1] = dst[subbandSize-1] - delta97*(dst[subbandSize-2]+dst[subbandSize-2])

		for k := 1; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] - gamma97*(dst[k-1]+dst[k+1])
		}

		dst[0] = dst[0] - beta97*(dst[1]+dst[1])
		for k := 2; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] - beta97*(dst[k-1]+dst[k+1])
		}
		dst[subbandSize-1] = dst[subbandSize-1] - beta97*(dst[subbandSize-2]+dst[subbandSize-2])

		for k := 1; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] - alpha97*(dst[k-1]+dst[k+1])
		}

	} else if wtType == 5 || wtType == 6 {
		var alpha, beta, gamma, delta float64
		if wtType == 6 {
			alpha = -1.58615986717275
			beta = -0.05297864003258
			gamma = 0.88293362717904
			delta = 0.44350482244527
		} else {
			alpha = -0.5
			beta = 0.25
			gamma = 0
			delta = 0
		}
		if wtType == 6 {
			dst[0] = dst[0] - math.Floor(delta*(dst[1]+dst[1])+0.5)
			for k := 2; k < subbandSize-1; k += 2 {
				dst[k] = dst[k] - math.Floor(delta*(dst[k-1]+dst[k+1])+0.5)
			}
			dst[subbandSize-1] = dst[subbandSize-1] - math.Floor(delta*(dst[subbandSize-2]+dst[subbandSize-2])+0.5)

			for k := 1; k < subbandSize-1; k += 2 {
				dst[k] = dst[k] - math.Floor(gamma*(dst[k-1]+dst[k+1])+0.5)
			}
		}
		dst[0] = dst[0] - (math.Floor(beta*(dst[1]+dst[1]) + 0.5))
		for k := 2; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] - math.Floor(beta*(dst[k-1]+dst[k+1])+0.5)
		}
		dst[subbandSize-1] = dst[subbandSize-1] - math.Floor(beta*(dst[subbandSize-2]+dst[subbandSize-2])+0.5)
		for k := 1; k < subbandSize-1; k += 2 {
			dst[k] = dst[k] - math.Floor(alpha*(dst[k-1]+dst[k+1])+0.5)
		}
	}
	return (dst)
}
