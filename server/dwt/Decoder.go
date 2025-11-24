package dwt

import "math"

// ACKey is a struct used as a map key for accessing AC coefficients,
// replacing the inefficient string concatenation from the original implementation.
type ACKey struct {
	RLevel  int
	Subband int
	Y       int
	X       int
}

func decode(blockNo int, bitPlane []int, rLevel int, bp int, Ds map[int]int,
	segmentStatus map[ACKey]int, significantPyramid map[ACKey]int,
	ac map[ACKey]int, d map[int]int, entropy func(int, int, bool, int) int) {
	gaggleNo := blockNo / gaggleACSize

	updateLeastSignificantBit(segmentStatus, bitPlane, bp)
	if rLevel == 0 {
		decodeParents(segmentStatus, gaggleNo, bitPlane, bp, ac, entropy)
		return
	}
	if rLevel == 1 {
		setPreviousSignificance(segmentStatus, significantPyramid)
		Ds[blockNo] = decodeChildren(Ds[blockNo], gaggleNo, significantPyramid, segmentStatus, bitPlane, bp, d, ac, entropy)
		return
	}
	decodeGenerationFunction(rLevel, bitPlane, bp, gaggleNo, d, significantPyramid, segmentStatus, ac, entropy)

}

func updateLeastSignificantBit(segmentStatus map[ACKey]int, bitPlane []int, bp int) {
	for rLevel := 0; rLevel < resolutionLevels; rLevel++ {
		sizeResolution := 1 << rLevel
		for subBand := 0; subBand < 3; subBand++ {
			subBandNo := 3*rLevel + subBand + 1
			if bp >= bitPlane[subBandNo] {
				continue
			}
			setLeastSignificantBit(segmentStatus, rLevel, sizeResolution, subBand)
		}
	}
}

func setLeastSignificantBit(segmentStatus map[ACKey]int, rLevel int, sizeResolution int, subBand int) {
	for y := 0; y < sizeResolution; y++ {
		for x := 0; x < sizeResolution; x++ {
			key := ACKey{RLevel: rLevel, Subband: subBand, Y: y, X: x}
			segmentStatus[key] = -1
		}
	}
}

func decodeParents(segmentStatus map[ACKey]int, gaggleNo int, bitPlane []int, bp int, ac map[ACKey]int, entropy func(int, int, bool, int) int) {
	wordLength := 0
	for subband := 0; subband < families; subband++ {
		key := ACKey{RLevel: 0, Subband: subband, Y: 0, X: 0}
		if segmentStatus[key] == 0 {
			wordLength++
		}
	}
	if wordLength != 0 {
		word := entropy(0, wordLength, false, gaggleNo)
		wordLength = 0
		for subband := families - 1; subband >= 0; subband-- {
			key := ACKey{RLevel: 0, Subband: subband, Y: 0, X: 0}
			if segmentStatus[key] == 0 {
				if word%2 == 1 {
					segmentStatus[key] = 1
					wordLength++
				}
				word = word >> 1
			}
		}
	}
	if wordLength != 0 {
		word := entropy(1, wordLength, false, gaggleNo)
		for subband := families - 1; subband >= 0; subband-- {
			key := ACKey{RLevel: 0, Subband: subband, Y: 0, X: 0}
			if segmentStatus[key] == 1 {
				if word%2 == 1 {
					ac[key] = -getRecoveryValue(bitPlane, bp, 0, subband)
				} else {
					ac[key] = getRecoveryValue(bitPlane, bp, 0, subband)
				}
				word = word >> 1
			}
		}
	}
}

func getRecoveryValue(bitPlane []int, bp int, rLevel int, subBand int) int {
	subbandNumber := (3 * rLevel) + subBand + 1
	recoveryValue := math.Pow(2, float64(bp))
	if bp > bitPlane[subbandNumber] {
		gamma := math.Round(gammaValue * recoveryValue)
		recoveryValue = recoveryValue + gamma
	}
	return int(recoveryValue)
}

func setPreviousSignificance(segmentStatus map[ACKey]int, significantPyramid map[ACKey]int) {
	for rLevel := 0; rLevel < resolutionLevels; rLevel++ {
		sizeResolution := 1 << rLevel
		for subband := 0; subband < families; subband++ {
			for y := 0; y < sizeResolution; y++ {
				for x := 0; x < sizeResolution; x++ {
					key := ACKey{RLevel: rLevel, Subband: subband, Y: y, X: x}
					pyramidKey := ACKey{RLevel: rLevel, Subband: subband}
					if segmentStatus[key] > significantPyramid[pyramidKey] {
						significantPyramid[pyramidKey] = segmentStatus[key]
					}
				}
			}
		}
	}
	for rLevel := resolutionLevels - 2; rLevel >= 1; rLevel-- {
		for subband := 0; subband < families; subband++ {
			key1 := ACKey{RLevel: rLevel, Subband: subband}
			key2 := ACKey{RLevel: rLevel + 1, Subband: subband}
			if significantPyramid[key1] < significantPyramid[key2] {
				significantPyramid[key1] = significantPyramid[key2]
			}
		}
	}
}

func decodeChildren(Ds int, gaggleNo int, significantPyramid map[ACKey]int, segmentStatus map[ACKey]int, bitPlane []int, bp int,
	d map[int]int, ac map[ACKey]int, entropy func(int, int, bool, int) int) int {
	rLevel := 1
	if Ds == 0 {
		temp := entropy(0, 1, false, gaggleNo)
		if temp == 1 {
			Ds = 1
		}
	} else {
		Ds = 2
	}
	if Ds <= 0 {
		return Ds
	}
	wordLength := 0
	for subband := 0; subband < families; subband++ {
		key := ACKey{RLevel: 1, Subband: subband}
		if significantPyramid[key] == 0 && !bitPlaneOfZeros(resolutionLevels-1, subband, bp, bitPlane) {
			wordLength++
		}
	}
	word := 0
	if wordLength != 0 {
		word = entropy(0, wordLength, true, gaggleNo)
	} else {
		word = 0
	}
	for subband := families - 1; subband >= 0; subband-- {
		key := ACKey{RLevel: rLevel, Subband: subband}
		if significantPyramid[key] == 0 && !bitPlaneOfZeros(resolutionLevels-1, subband, bp, bitPlane) {
			if word%2 == 1 {
				significantPyramid[key] = 1
				d[subband] = 2
			}
			word = word >> 1
		}
	}
	for subband := 0; subband < 3; subband++ {
		key := ACKey{RLevel: rLevel, Subband: subband}
		recoveryValue := getRecoveryValue(bitPlane, bp, rLevel, subband)
		if significantPyramid[key] != 0 && !bitPlaneOfZeros(rLevel, subband, bp, bitPlane) {
			wordLength = 0
			for k := 0; k < 4; k++ {
				childKey := ACKey{RLevel: rLevel, Subband: subband, Y: k / 2, X: k % 2}
				if segmentStatus[childKey] == 0 {
					wordLength++
				}
			}
			if wordLength != 0 {
				word = entropy(0, wordLength, false, gaggleNo)
				signLength := 0
				for k := 3; k >= 0; k-- {
					childKey := ACKey{RLevel: rLevel, Subband: subband, Y: k / 2, X: k % 2}
					if segmentStatus[childKey] == 0 {
						bit := word % 2
						if bit == 1 {
							segmentStatus[childKey] = 1
							signLength++
						}
						word = word >> 1
					}
				}
				if signLength != 0 {
					word = entropy(1, signLength, false, gaggleNo)
					for k := 3; k >= 0; k-- {
						childKey := ACKey{RLevel: rLevel, Subband: subband, Y: k / 2, X: k % 2}
						if segmentStatus[childKey] == 1 {
							bit := word % 2
							if bit == 1 {
								ac[childKey] = -recoveryValue
							} else {
								ac[childKey] = recoveryValue
							}
							word = word >> 1
						}
					}
				}
			}
		}
	}
	return Ds
}

func bitPlaneOfZeros(rLevel int, subBand int, bp int, bitPlane []int) bool {
	var zeros bool
	subbandNumber := (3 * rLevel) + subBand + 1
	if bitPlane[subbandNumber] > bp {
		zeros = true
	} else {
		zeros = false
	}
	return zeros
}

func decodeGenerationFunction(rLevel int, bitPlane []int, bp int, gaggleNo int, d map[int]int, significantPyramid map[ACKey]int,
	segmentStatus map[ACKey]int, ac map[ACKey]int, entropy func(int, int, bool, int) int) {
	wordLength := 0
	for subband := 0; subband < families; subband++ {
		if d[subband] == 2 {
			key := ACKey{RLevel: rLevel, Subband: subband}
			if significantPyramid[key] == 0 && !bitPlaneOfZeros(rLevel, subband, bp, bitPlane) {
				wordLength++
			}
		}
	}
	if wordLength != 0 {
		var impossiblePattern = false
		word := entropy(0, wordLength, impossiblePattern, gaggleNo)
		for subband := families - 1; subband >= 0; subband-- {
			if d[subband] == 2 {
				key := ACKey{RLevel: rLevel, Subband: subband}
				if significantPyramid[key] == 0 && !bitPlaneOfZeros(rLevel, subband, bp, bitPlane) {
					if word%2 == 1 {
						significantPyramid[key] = 1
					}
					word = word >> 1
				}
			}
		}
	}

	key0 := ACKey{RLevel: rLevel, Subband: 0}
	key1 := ACKey{RLevel: rLevel, Subband: 1}
	key2 := ACKey{RLevel: rLevel, Subband: 2}

	if significantPyramid[key0] <= 0 && significantPyramid[key1] <= 0 && significantPyramid[key2] <= 0 {
		return
	}

	xSize := 1 << rLevel
	var decodeSquare = initDecodeSquare(rLevel)
	for pass := 0; pass <= rLevel; pass++ {
		numberOfWords := 1 << (2 * pass)
		for subband := 0; subband < 3; subband++ {
			key := ACKey{RLevel: rLevel, Subband: subband}
			if significantPyramid[key] > 0 && !bitPlaneOfZeros(rLevel, subband, bp, bitPlane) {
				decodeSquare[pass][subband] = make([]int, numberOfWords)
			}
		}
	}

	for subband := 0; subband < 3; subband++ {
		key := ACKey{RLevel: rLevel, Subband: subband}
		if significantPyramid[key] <= 0 || bitPlaneOfZeros(rLevel, subband, bp, bitPlane) {
			continue
		}
		numberOfSquares := 1 << (2 * rLevel)
		for sq := 0; sq < numberOfSquares; sq += 4 {
			quadrantxSize := xSize / 2
			x := 0
			y := 0
			positionInQuadrant := sq
			componentsInQuadrant := numberOfSquares / 4
			for componentsInQuadrant > 1 {
				whichQuadrant := positionInQuadrant / componentsInQuadrant
				positionInQuadrant = positionInQuadrant % componentsInQuadrant
				componentsInQuadrant = componentsInQuadrant / 4
				if whichQuadrant == 1 {
					x += quadrantxSize
				} else if whichQuadrant == 2 {
					y += quadrantxSize
				} else if whichQuadrant == 3 {
					x += quadrantxSize
					y += quadrantxSize
				}
				quadrantxSize = quadrantxSize >> 1
			}
			for i := 0; i < 4; i++ {
				statusKey := ACKey{RLevel: rLevel, Subband: subband, Y: y + i/2, X: x + i%2}
				decodeSquare[rLevel][subband][sq+i] = segmentStatus[statusKey]
			}

		}
		for pass := rLevel; pass > 0; pass-- {
			numberOfWords := 1 << (2 * pass)
			for k := 0; k < numberOfWords; k++ {
				if decodeSquare[pass-1][subband][k/4] < decodeSquare[pass][subband][k] {
					decodeSquare[pass-1][subband][k/4] = decodeSquare[pass][subband][k]
				}
			}
		}
		pyramidKey := ACKey{RLevel: rLevel, Subband: subband}
		decodeSquare[0][subband][0] = significantPyramid[pyramidKey]
	}
	for pass := 1; pass < rLevel; pass++ {
		for subband := 0; subband < 3; subband++ {
			key := ACKey{RLevel: rLevel, Subband: subband}
			if significantPyramid[key] <= 0 || bitPlaneOfZeros(rLevel, subband, bp, bitPlane) {
				continue
			}
			numberOfSquares := 1 << (2 * pass)
			wordLength = 0
			for sq := 0; sq < numberOfSquares; sq++ {
				if decodeSquare[pass-1][subband][sq/4] > 0 && decodeSquare[pass][subband][sq] == 0 {
					wordLength++
				}
				if sq%4 != 3 || wordLength == 0 {
					continue
				}
				word := 0
				if wordLength == 4 {
					word = entropy(0, wordLength, true, gaggleNo)
				} else if wordLength != 0 {
					word = entropy(0, wordLength, false, gaggleNo)
				}
				for i := 0; i < 4; i++ {
					if decodeSquare[pass-1][subband][sq/4] > 0 && decodeSquare[pass][subband][sq-i] == 0 {
						if word%2 == 1 {
							decodeSquare[pass][subband][sq-i] = 1
						}
						word = word >> 1
					}
				}
				wordLength = 0
			}
		}
	}
	for subband := 0; subband < 3; subband++ {
		recoveryValue := getRecoveryValue(bitPlane, bp, rLevel, subband)
		key := ACKey{RLevel: rLevel, Subband: subband}
		if significantPyramid[key] <= 0 || bitPlaneOfZeros(rLevel, subband, bp, bitPlane) {
			continue
		}
		numberOfSquares := 1 << (2 * rLevel)
		wordLength = 0
		for sq := 0; sq < numberOfSquares; sq++ {
			if decodeSquare[rLevel-1][subband][sq/4] > 0 && decodeSquare[rLevel][subband][sq] == 0 {
				wordLength++
			}
			if sq%4 != 3 || wordLength == 0 {
				continue
			}
			word := 0
			signLength := 0
			if wordLength == 4 {
				word = entropy(0, wordLength, true, gaggleNo)
			} else if wordLength != 0 {
				word = entropy(0, wordLength, false, gaggleNo)
			}
			for i := 0; i < 4; i++ {
				if decodeSquare[rLevel-1][subband][sq/4] > 0 && decodeSquare[rLevel][subband][sq-i] == 0 {
					if word%2 == 1 {
						decodeSquare[rLevel][subband][sq-i] = 1
						signLength++
					}
					word = word >> 1
				}
			}
			if signLength == 0 {
				wordLength = 0
				continue
			}
			word = entropy(1, signLength, false, gaggleNo)
			quadrantxSize := xSize / 2
			x := 0
			y := 0
			positionInQuadrant := sq - 3
			componentsInQuadrant := numberOfSquares / 4
			for componentsInQuadrant > 1 {
				whichQuadrant := positionInQuadrant / componentsInQuadrant
				positionInQuadrant = positionInQuadrant % componentsInQuadrant
				componentsInQuadrant = componentsInQuadrant / 4
				if whichQuadrant == 1 {
					x += quadrantxSize
				} else if whichQuadrant == 2 {
					y += quadrantxSize
				} else if whichQuadrant == 3 {
					x += quadrantxSize
					y += quadrantxSize
				}
				quadrantxSize = quadrantxSize >> 1
			}
			for i := 0; i < 4; i++ {
				if decodeSquare[rLevel][subband][sq-i] == 1 {
					k := 3 - i
					statusKey := ACKey{RLevel: rLevel, Subband: subband, Y: y + k/2, X: x + k%2}
					if word%2 == 1 {
						ac[statusKey] = -recoveryValue
						segmentStatus[statusKey] = 1
					} else {
						ac[statusKey] = recoveryValue
						segmentStatus[statusKey] = 1
					}
					word = word >> 1
				}
			}
			wordLength = 0
		}
	}
	return
}

func initDecodeSquare(rLevel int) [][][]int {
	var decodeSquare = make([][][]int, rLevel+1)

	for i := 0; i < rLevel+1; i++ {
		decodeSquare[i] = make([][]int, 3)
	}
	return decodeSquare
}
