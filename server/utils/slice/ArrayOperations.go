// Package slice provides common utility functions for slice operations.
package slice

import (
	"fmt"
	"strconv"
	"strings"
)

// IndexOfInt returns the index of the first occurrence of element in the slice.
// It returns -1 if the element is not found.
func IndexOfInt(slice []int, element int) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}

// IndexOfStringFold returns the index of the first occurrence of element in the slice,
// performing a case-insensitive comparison after trimming whitespace.
// It returns -1 if the element is not found.
func IndexOfStringFold(slice []string, element string) int {
	trimmedElement := strings.TrimSpace(element)
	for i, v := range slice {
		if strings.EqualFold(trimmedElement, strings.TrimSpace(v)) {
			return i
		}
	}
	return -1
}

// HexStringsToBytes converts a slice of hexadecimal strings to a byte slice.
// It returns an error if any string is not a valid hex representation.
func HexStringsToBytes(hexStrings []string) ([]byte, error) {
	bytes := make([]byte, len(hexStrings))
	for i, hexStr := range hexStrings {
		val, err := strconv.ParseUint(hexStr, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("failed to parse hex string '%s' at index %d: %w", hexStr, i, err)
		}
		bytes[i] = byte(val)
	}
	return bytes, nil
}

// Unique returns a new slice containing only the unique elements from the input slice.
// It works with any comparable type (string, int, float64, etc.).
// The order of elements in the output slice is the order of their first appearance.
func Unique[T comparable](slice []T) []T {
	if slice == nil {
		return nil
	}
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// IndicesGreaterThan returns a slice of indices where the element in the input slice
// is greater than the given value.
func IndicesGreaterThan(slice []float64, value float64) []int {
	var indices []int
	for i, v := range slice {
		if v > value {
			indices = append(indices, i)
		}
	}
	return indices
}

// IndicesEqualTo returns a slice of indices where the element in the input slice
// is equal to the given value.
func IndicesEqualTo(slice []float64, value float64) []int {
	var indices []int
	for i, v := range slice {
		if v == value {
			indices = append(indices, i)
		}
	}
	return indices
}
