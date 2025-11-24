package dwt

const wtLevels = 3
const wtType = 4
const wtOrder = 0
const gaggleDCSize = 16
const gaggleACSize = 16
const resolutionLevels = 3
const gammaValue = 0.375
const families = 3
const alpha97 = -1.586134342059924
const beta97 = -0.052980118572961
const gamma97 = 0.882911075530934
const delta97 = 0.443506852043971

func getDataProvider(array []byte) func(int) ([]byte, bool) {
	var byteProvider = func(noOfBytes int) ([]byte, bool) {
		if len(array) < noOfBytes {
			return []byte{}, false
		}
		var tbr = make([]byte, 0)
		tbr = append(tbr, array[:noOfBytes]...)
		array = array[noOfBytes:]
		return tbr, true

	}
	return byteProvider
}
