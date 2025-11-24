package microwave

import (
	"csrsp/server/session"
	"csrsp/server/utils/binary"
	"log/slog"
	"math"
	"math/cmplx"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/mjibson/go-dsp/fft"
	"github.com/montanaflynn/stats"
)

// iqSampleProvider performs the reverse BAQ decompression to extract I and Q samples.
// It is a direct port of the original logic.
func iqSampleProvider(sarModeValue string, baqValue uint64, videoDataLen uint64, getBitsToUint64 []func(uint64, uint64) uint64) ([]float64, []float64) {
	var m uint64 = 0
	var baqBlockSize uint64 = 128
	var baqScaleFactorI1 float64 = 0
	var baqScaleFactorQ1 float64 = 0
	iSamples := make([]float64, 0)
	qSamples := make([]float64, 0)
	if baqValue == 8 {
		baqBlockSize = 128
		if strings.EqualFold(sarModeValue, "stretch") {
			baqBlockSize = 64
		}
	} else if baqValue >= 2 && baqValue <= 6 {
		baqBlockSize = baqValue * 16
	} else if baqValue == 1 || baqValue == 7 {
		// Log error in the future if needed
		baqBlockSize = 128
		baqValue = 8
	}
	baqBlockCount := (videoDataLen) / (baqBlockSize + baqHeaderSize)
	twosComplement := getTwosComplementProcessor(baqValue)
	for blkIndex := 0; uint64(blkIndex) < baqBlockCount; blkIndex++ {
		blkOffset := (uint64(blkIndex) * baqBlockSize) + (baqHeaderSize * (uint64(blkIndex) + 1))
		blockAvgI1 := getBitsToUint64[0]((blkOffset-1)*8, 7)
		blockAvgQ1 := getBitsToUint64[1]((blkOffset-1)*8, 7)
		if baqValue != 8 {
			baqScaleFactorI1 = baqScaleFactorTable[blockAvgI1][baqValue-2]
			baqScaleFactorQ1 = baqScaleFactorTable[blockAvgQ1][baqValue-2]
		}
		mergedDataI := make([]byte, 0)
		mergedDataQ := make([]byte, 0)
		if strings.EqualFold(sarModeValue, "stretch") && len(getBitsToUint64) == 4 {
			blockAvgQ1 = getBitsToUint64[1]((blkOffset-1)*8, 7)
			if baqValue != 8 {
				baqScaleFactorQ1 = baqScaleFactorTable[blockAvgQ1][baqValue-2]
			}
			for m = 0; m < baqBlockSize; m++ {
				mergedDataI = append(mergedDataI, byte(getBitsToUint64[0]((blkOffset+m)*8, 8)))
				mergedDataI = append(mergedDataI, byte(getBitsToUint64[1]((blkOffset+m)*8, 8)))
				mergedDataQ = append(mergedDataQ, byte(getBitsToUint64[2]((blkOffset+m)*8, 8)))
				mergedDataQ = append(mergedDataQ, byte(getBitsToUint64[3]((blkOffset+m)*8, 8)))
			}
		} else if strings.EqualFold(sarModeValue, "rt") && len(getBitsToUint64) == 2 {
			for m = 0; m < baqBlockSize; m++ {
				mergedDataI = append(mergedDataI, byte(getBitsToUint64[0]((blkOffset+m)*8, 8)))
				mergedDataQ = append(mergedDataQ, byte(getBitsToUint64[1]((blkOffset+m)*8, 8)))
			}
		} else if strings.EqualFold(sarModeValue, "nrt") && len(getBitsToUint64) == 2 {
			for m = 0; m < baqBlockSize; m += 2 {
				mergedDataI = append(mergedDataI, byte(getBitsToUint64[0]((blkOffset+m+1)*8, 8)))
				mergedDataI = append(mergedDataI, byte(getBitsToUint64[0]((blkOffset+m)*8, 8)))
				mergedDataQ = append(mergedDataQ, byte(getBitsToUint64[1]((blkOffset+m+1)*8, 8)))
				mergedDataQ = append(mergedDataQ, byte(getBitsToUint64[1]((blkOffset+m)*8, 8)))
			}
		}
		getBitsToUint64I := getDataProviderForContinuousData(mergedDataI, baqValue)
		getBitsToUint64Q := getDataProviderForContinuousData(mergedDataQ, baqValue)

		for sample := 0; sample < numberOfSamplesPerBlock; sample++ {
			iVal := twosComplement(getBitsToUint64I(baqValue))
			qVal := twosComplement(getBitsToUint64Q(baqValue))
			if baqValue == 8 {
				iSamples = append(iSamples, float64(iVal))
				qSamples = append(qSamples, float64(qVal))
			} else {
				iSamples = append(iSamples, (float64(iVal)+0.5)/baqScaleFactorI1)
				qSamples = append(qSamples, (float64(qVal)+0.5)/baqScaleFactorQ1)
			}
		}
	}
	return iSamples, qSamples
}

func getDataProviderForContinuousData(frame []byte, baqValue uint64) func(noOfBits uint64) uint64 {
	var startBit uint64 = 0
	var provider = func(noOfBits uint64) uint64 {
		startWord := startBit / 8
		mask, _ := binary.NewContinuousMask(int(startWord), int(startBit%8), int(noOfBits))
		value, _ := mask.ExtractUint64(frame)
		startBit += noOfBits
		if baqValue != 8 {
			value = uint64(baqMap[int(baqValue)][int(value)])
		}
		return value
	}
	return provider
}

func getDataProviderForNonContinuousData(frame []byte) func(startBit uint64, noOfBits uint64) uint64 {
	var provider = func(startBit uint64, noOfBits uint64) uint64 {
		startWord := startBit / 8
		startBit = startBit % 8
		mask, _ := binary.NewContinuousMask(int(startWord), int(startBit), int(noOfBits))
		value, _ := mask.ExtractUint64(frame)
		return value
	}
	return provider
}

func getTwosComplementProcessor(baqBits uint64) func(uint64) int {
	maxValueWithSign := uint64(math.Pow(2, float64(baqBits-1)) - 1)
	maxValueWithoutSign := uint64(math.Pow(2, float64(baqBits)))

	var twosComplement = func(value uint64) int {
		if value > maxValueWithSign {
			return int(value) - int(maxValueWithoutSign)
		}
		return int(value)

	}
	return twosComplement
}

// iqStatisticsCalculator calculates basic statistics for a slice of float64 values.
func iqStatisticsCalculator(values []float64) session.IQStats {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in Microwave - IQ Calculator", "panic", r, "stack", string(debug.Stack()))
		}
	}()
	var rmsSqr float64
	minValue := slices.Min(values)
	maxValue := slices.Max(values)
	for _, value := range values {
		rmsSqr += value * value
	}

	mean, _ := stats.Mean(values)
	sd, _ := stats.StandardDeviation(values)
	rms := math.Sqrt(rmsSqr / float64(len(values)))

	return session.IQStats{Min: minValue, Max: maxValue, Mean: mean, SD: sd, RMS: rms}
}

// generateSARChirpParameters performs FFT and analysis to calculate chirp parameters.
func generateSARChirpParameters(iSamples []float64, qSamples []float64, samplingFrequency float64) session.ChirpParams {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Panic recovered in Microwave - Chirp Parameters", "panic", r, "stack", string(debug.Stack()))
		}
	}()

	originalLen := len(iSamples)
	fftLen := nextPow2(originalLen)

	iqSamples := make([]complex128, fftLen)
	for i := 0; i < originalLen; i++ {
		iqSamples[i] = complex(iSamples[i], qSamples[i])
	}

	peakBias := (fftLen - originalLen) * 8

	fft.FFT(iqSamples)

	m := make([]complex128, fftLen)
	for i, sample := range iqSamples {
		temp := cmplx.Conj(sample)
		m[i] = sample * temp
	}

	fft.IFFT(m)

	// IFFTShift
	corr := make([]complex128, fftLen)
	half := fftLen / 2
	copy(corr, m[half:])
	copy(corr[fftLen-half:], m[:half])

	var interpolationFactor = 16
	numberOfInterpolateSamples := len(corr) * interpolationFactor

	fft.FFT(corr) // Re-using corr for the next FFT

	center := len(corr) / 2
	inter := make([]complex128, numberOfInterpolateSamples)
	copy(inter, corr[:center+1])
	copy(inter[numberOfInterpolateSamples-center:], corr[center+1:])

	if len(corr)%2 == 0 {
		inter[center] = inter[center] / 2
		inter[center+numberOfInterpolateSamples-len(corr)] = inter[center]
	}

	fft.IFFT(inter)

	multiplyFactor := complex(float64(interpolationFactor), 0)
	for i := 0; i < len(inter); i++ {
		inter[i] *= multiplyFactor
	}

	var maxValue float64 = -1
	var maxIndex = -1
	intersq := make([]float64, len(inter))
	for i := 0; i < len(inter); i++ {
		intersq[i] = cmplx.Abs(inter[i]) * cmplx.Abs(inter[i])
		if intersq[i] > maxValue {
			maxValue = intersq[i]
			maxIndex = i
		}
	}

	peakAmp := 10 * math.Log10(maxValue)
	peakAmpLoc := maxIndex

	_, peaksLoc, peakLocIndexIntersq := findPeaks(intersq, peakAmp)
	if peakLocIndexIntersq < 10 || peakLocIndexIntersq > len(peaksLoc)-11 {
		// Not enough side lobes to calculate, return what we have
		return session.ChirpParams{PeakAmp: peakAmp, PeakAmpLoc: uint64(peakAmpLoc - peakBias)}
	}

	chirpPeaks10 := intersq[peaksLoc[peakLocIndexIntersq-10] : peaksLoc[peakLocIndexIntersq+10]+1]
	tempPeaks, _, peakLocIndex := findPeaks(chirpPeaks10, peakAmp)
	pslrLeft := tempPeaks[peakLocIndex-1] - tempPeaks[peakLocIndex]
	pslrRight := tempPeaks[peakLocIndex+1] - tempPeaks[peakLocIndex]

	chirpPeaksWithOneLobe := intersq[peaksLoc[peakLocIndexIntersq-1] : peaksLoc[peakLocIndexIntersq+1]+1]
	leftLobe := intersq[peaksLoc[peakLocIndexIntersq-1]:peaksLoc[peakLocIndexIntersq]]
	rightLobe := intersq[peaksLoc[peakLocIndexIntersq] : peaksLoc[peakLocIndexIntersq+1]+1]

	var minLeft = math.MaxFloat64
	var minLeftIndex = -1
	for i := 0; i < len(leftLobe); i++ {
		if leftLobe[i] < minLeft {
			minLeft = leftLobe[i]
			minLeftIndex = i
		}
	}

	var minRight = math.MaxFloat64
	var minRightIndex = -1
	for i := 0; i < len(rightLobe); i++ {
		if rightLobe[i] < minRight {
			minRight = rightLobe[i]
			minRightIndex = i
		}
	}

	mainLobe := chirpPeaksWithOneLobe[minLeftIndex : len(leftLobe)+minRightIndex]
	var mainLobeSum float64
	for i := 0; i < len(mainLobe); i++ {
		mainLobeSum += mainLobe[i]
	}

	var chirpWith10LobesSum float64
	for i := 0; i < len(chirpPeaks10); i++ {
		chirpWith10LobesSum += chirpPeaks10[i]
	}
	islr := 10 * math.Log10((chirpWith10LobesSum-mainLobeSum)/mainLobeSum)

	var speed float64 = 300000000
	var fs = samplingFrequency
	pixelSpacing := speed / (2 * fs * 16)
	maxValue = math.Sqrt(maxValue)
	halfPowerPointRight := -1
	for i := maxIndex; i < len(inter); i++ {
		if (cmplx.Abs(inter[i]) / maxValue) <= (1 / math.Sqrt(2)) {
			halfPowerPointRight = i
			break
		}
	}
	halfPowerPointLeft := -1
	for i := maxIndex; i > 0; i-- {
		if (cmplx.Abs(inter[i]) / maxValue) <= (1 / math.Sqrt(2)) {
			halfPowerPointLeft = i
			break
		}
	}
	resolution := float64(halfPowerPointRight-halfPowerPointLeft) * pixelSpacing

	peakAmpLoc = peakAmpLoc - peakBias

	return session.ChirpParams{
		PeakAmp:    peakAmp,
		PeakAmpLoc: uint64(peakAmpLoc),
		PSLRLeft:   pslrLeft,
		PSLRRight:  pslrRight,
		ISLR:       islr,
		Resolution: resolution,
	}
}

func findPeaks(data []float64, peakAmp float64) ([]float64, []int, int) {
	peaks := make([]float64, 0)
	peaksLoc := make([]int, 0)
	peakLocIndex := -1
	for i := 1; i < len(data)-1; i++ {
		currNum := data[i]
		if data[i-1] < currNum && data[i+1] < currNum {
			peaks = append(peaks, 10*math.Log10(currNum))
			peaksLoc = append(peaksLoc, i)
			if peakAmp == 10*math.Log10(currNum) {
				peakLocIndex = len(peaks) - 1
			}
		}
	}
	return peaks, peaksLoc, peakLocIndex
}

func nextPow2(orig int) int {
	value := 1
	for value < orig {
		value = value * 2
	}
	return value
}
