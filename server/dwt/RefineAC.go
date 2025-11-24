package dwt

import "math"

func refineACValues(bp int, rLevel int, gaggleNo int, segmentStatus map[ACKey]int, ac map[ACKey]int,
	bitPlane []int, entropy func(int, int, bool, int) int) {
	threshold := 1 << uint(bp)
	sizeResolutionLevel := 1 << uint(rLevel)

	for subband := 0; subband < families; subband++ {
		gamma := getGamma(bp, lastBitPlane(bp, rLevel, subband, bitPlane))
		previousGamma := getGamma(bp+1, lastBitPlane(bp+1, rLevel, subband, bitPlane))
		refinementBit := 0
		if rLevel != 2 {
			for y := 0; y < sizeResolutionLevel; y++ {
				for x := 0; x < sizeResolutionLevel; x++ {
					key := ACKey{RLevel: rLevel, Subband: subband, Y: y, X: x}
					if segmentStatus[key] == 2 {
						temp := entropy(2, 1, false, gaggleNo)
						if temp == -1 {
							return
						}
						refinementBit = 0
						if temp == 1 {
							refinementBit = 1
						}
						if ac[key] < 0 {
							ac[key] = ac[key] + previousGamma - threshold*refinementBit - gamma
						} else {
							ac[key] = ac[key] - previousGamma + threshold*refinementBit + gamma
						}
					} else if segmentStatus[key] == 1 {
						segmentStatus[key] = 2
					}
				}
			}
			continue
		}
		for y := 0; y < 2; y++ {
			for x := 0; x < 2; x++ {
				key := ACKey{RLevel: rLevel, Subband: subband, Y: y, X: x}
				if segmentStatus[key] == 2 {
					temp := entropy(2, 1, false, gaggleNo)
					if temp == -1 {
						return
					}
					refinementBit = 0
					if temp == 1 {
						refinementBit = 1
					}
					if ac[key] < 0 {
						ac[key] = ac[key] + previousGamma - threshold*refinementBit - gamma
					} else {
						ac[key] = ac[key] - previousGamma + threshold*refinementBit + gamma
					}
				} else if segmentStatus[key] == 1 {
					segmentStatus[key] = 2
				}
			}
		}
		for y := 0; y < 2; y++ {
			for x := 2; x < 4; x++ {
				key := ACKey{RLevel: rLevel, Subband: subband, Y: y, X: x}
				if segmentStatus[key] == 2 {
					temp := entropy(2, 1, false, gaggleNo)
					if temp == -1 {
						return
					}
					refinementBit = 0
					if temp == 1 {
						refinementBit = 1
					}
					if ac[key] < 0 {
						ac[key] = ac[key] + previousGamma - threshold*refinementBit - gamma
					} else {
						ac[key] = ac[key] - previousGamma + threshold*refinementBit + gamma
					}
				} else if segmentStatus[key] == 1 {
					segmentStatus[key] = 2
				}
			}
		}
		for y := 2; y < 4; y++ {
			for x := 0; x < 2; x++ {
				key := ACKey{RLevel: rLevel, Subband: subband, Y: y, X: x}
				if segmentStatus[key] == 2 {
					temp := entropy(2, 1, false, gaggleNo)
					if temp == -1 {
						return
					}
					refinementBit = 0
					if temp == 1 {
						refinementBit = 1
					}
					if ac[key] < 0 {
						ac[key] = ac[key] + previousGamma - threshold*refinementBit - gamma
					} else {
						ac[key] = ac[key] - previousGamma + threshold*refinementBit + gamma
					}
				} else if segmentStatus[key] == 1 {
					segmentStatus[key] = 2
				}
			}
		}
		for y := 2; y < 4; y++ {
			for x := 2; x < 4; x++ {
				key := ACKey{RLevel: rLevel, Subband: subband, Y: y, X: x}
				if segmentStatus[key] == 2 {
					temp := entropy(2, 1, false, gaggleNo)
					if temp == -1 {
						return
					}
					refinementBit = 0
					if temp == 1 {
						refinementBit = 1
					}
					if ac[key] < 0 {
						ac[key] = ac[key] + previousGamma - threshold*refinementBit - gamma
					} else {
						ac[key] = ac[key] - previousGamma + threshold*refinementBit + gamma
					}
				} else if segmentStatus[key] == 1 {
					segmentStatus[key] = 2
				}
			}
		}

	}
}

func getGamma(bitPlane int, lastBitPlane bool) int {
	gamma := 0
	threshold := 1 << uint(bitPlane)
	if !lastBitPlane {
		gamma = int(math.Round(gammaValue * float64(threshold)))
	} else {
		gamma = 0
	}
	return gamma
}

func lastBitPlane(bp int, rLevel int, subband int, bitPlane []int) bool {
	lastBitPlane := true
	subbandNumber := (3 * rLevel) + subband + 1
	if bp > bitPlane[subbandNumber] {
		lastBitPlane = false
	}
	return lastBitPlane
}
