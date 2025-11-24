package dwt

func getEntropyWord(noOfGaggles int, provider func(int) int) (func(int, int, bool, int) int, func()) {
	var ids = make([][]int, noOfGaggles)
	var reset = func() {
		for i := 0; i < noOfGaggles; i++ {
			ids[i] = make([]int, 3)
			for j := 0; j < 3; j++ {
				ids[i][j] = -1
			}
		}
	}
	codeIDLength := []int{1, 2, 2}

	var two = func(gaggle int) int {
		if ids[gaggle][0] == 1 {
			symbol := provider(2)
			return symbol
		}
		temp, _ := getNext1ModLimit(provider, 3)
		return temp
	}

	var three = func(gaggle int) int {
		var symbol int
		var found bool
		if ids[gaggle][1] == 3 {
			symbol = provider(3)
			return symbol
		}
		if ids[gaggle][1] == 0 {
			symbol, found = getNext1ModLimit(provider, 3)
			if !found {
				bits := provider(2)
				symbol = 3 + bits
				if bits == 3 {
					temp := provider(1)
					symbol = symbol + temp
				}
			}
		} else if ids[gaggle][1] == 1 {
			bits := provider(2)
			symbol = bits - 2
			switch bits {
			case 1:
				temp := provider(1)
				symbol = 2 + temp
			case 0:
				temp := provider(2)
				symbol = temp + 2
				switch temp {
				case 0:
					symbol = 6
				case 1:
					symbol = 7
				}
			}
		}
		return symbol
	}

	var four = func(gaggle int) int {
		var symbol = 0
		var found bool
		if ids[gaggle][2] == 3 {
			symbol = provider(4)
			return symbol
		}
		switch ids[gaggle][2] {
		case 0:
			symbol, found = getNext1ModLimit(provider, 4)
			if !found {
				if temp := provider(1); temp == 0 {
					bits := provider(2)
					symbol = 4 + bits
				} else {
					bits := provider(3)
					symbol = 8 + bits
				}
			}
		case 1:
			bits := provider(2)
			symbol = bits - 2
			switch bits {
			case 1:
				temp := provider(1)
				symbol = 2 + temp
			case 0:
				temp := provider(2)
				symbol = 2 + temp
				switch temp {
				case 0:
					temp1 := provider(2)
					symbol = temp1 + 6
				case 1:
					temp2 := provider(2)
					symbol = temp2 + 10
					if temp2 == 2 {
						temp3 := provider(1)
						symbol = temp3 + 12
					}
					if temp2 == 3 {
						temp3 := provider(1)
						symbol = temp3 + 14
					}
				}
			}
		case 2:
			bits := provider(3)
			symbol = bits - 4
			if bits > 1 && bits < 4 {
				temp := provider(1)
				symbol = 2*bits + temp
			} else if bits < 2 {
				temp := provider(2)
				symbol = 12 - 4*bits + temp
			}
		}

		return symbol
	}

	var word = func(stage int, length int, impossible bool, gaggle int) int {
		symbol := 0
		if length >= 2 && stage == 0 {
			if ids[gaggle][length-2] == -1 {
				idLength := codeIDLength[length-2]
				ids[gaggle][length-2] = provider(idLength)
			}
			switch length {
			case 2:
				symbol = two(gaggle)
			case 3:
				symbol = three(gaggle)
			case 4:
				symbol = four(gaggle)
			}
		} else {
			symbol = provider(length)
		}

		context := getRecommendedEntropyContext(stage, length, impossible)
		word := unMapSymbol(symbol, context)
		return word
	}
	return word, reset

}

func getNext1ModLimit(provider func(int) int, limit int) (int, bool) {
	var symbol int
	var found bool
	for symbol = 0; symbol < limit; symbol++ {
		if val := provider(1); val != 0 {
			found = true
			break
		}
	}
	return symbol, found
}

func getRecommendedEntropyContext(stage int, length int, impossiblePattern bool) int {
	if stage <= 0 {
		return getStageZeroContext(length, impossiblePattern)
	}
	if stage == 1 {
		return getStageOneContext(length)
	}
	return 74
}

func getStageZeroContext(length int, impossiblePattern bool) int {
	var context int
	switch length {
	case 1:
		context = 64
	case 2:
		context = 65
	case 3:
		if impossiblePattern {
			context = 66
		} else {
			context = 67
		}
	case 4:
		if impossiblePattern {
			context = 68
		} else {
			context = 69
		}
	}
	return context
}

func getStageOneContext(length int) int {
	var context int
	switch length {
	case 1:
		context = 70
	case 2:
		context = 71
	case 3:
		context = 72
	case 4:
		context = 73
	}
	return context
}

func unMapSymbol(symbol int, context int) int {
	wordMap := map[int][]int{
		65: {0, 2, 1, 3},
		67: {2, 0, 4, 6, 1, 3, 5, 7},
		66: {2, 4, 6, 1, 3, 5, 7, 0},
		69: {8, 1, 4, 2, 12, 5, 3, 10, 9, 6, 0, 14, 7, 11, 13, 15},
		68: {8, 1, 4, 2, 12, 5, 3, 10, 9, 6, 14, 7, 11, 13, 15, 0},
	}
	wordArray, ok := wordMap[context]
	if !ok {
		return symbol
	}
	return wordArray[symbol]
}
