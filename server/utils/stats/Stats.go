// Package stats provides basic statistical functions.
package stats

import (
	"golang.org/x/exp/constraints"
	"math"
)

// Number is a constraint that permits any integer or floating-point type.
type Number interface {
	constraints.Integer | constraints.Float
}

// Summary holds basic descriptive statistics for a dataset.
type Summary[T Number] struct {
	Mean T
	Max  T
	Min  T
}

// Describe calculates the mean, max, and min for a slice of numbers.
// It returns a zero-value Summary if the input slice is empty.
func Describe[T Number](data []T) Summary[T] {
	if len(data) == 0 {
		return Summary[T]{}
	}

	var sum T
	max, min := data[0], data[0]

	for _, value := range data {
		sum += value
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}

	return Summary[T]{
		Mean: sum / T(len(data)),
		Max:  max,
		Min:  min,
	}
}

// StdDev calculates the population standard deviation for a slice of float64s.
// It returns 0 if the slice is empty.
func StdDev(data []float64, mean float64) float64 {
	if len(data) == 0 {
		return 0.0
	}

	var sumOfSquares float64
	for _, value := range data {
		diff := value - mean
		sumOfSquares += diff * diff
	}

	variance := sumOfSquares / float64(len(data))
	return math.Sqrt(variance)
}

// Log2Ceil calculates the smallest integer i such that 2^i >= n.
// It returns 0 for n <= 0.
func Log2Ceil(n int) int {
	if n <= 0 {
		return 0
	}
	return int(math.Ceil(math.Log2(float64(n))))
}
