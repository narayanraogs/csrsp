// Package binary provides utilities for bit and byte-level manipulations.
package binary

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"
)

// --- Bitmasking --- //

// MaskPart represents a single component of a larger, possibly non-contiguous, bit mask.
type MaskPart struct {
	ByteIndex     int  // Index of the byte in the source slice
	Mask          byte // The byte mask (e.g., 0b00111000)
	Shift         uint // How many bits to right-shift the masked value to align it to the right
	ValidBitCount int  // The number of bits this part contributes
}

// Mask represents a complete bit mask to extract a value from a byte slice.
// It can be composed of multiple parts to support non-contiguous data extraction.
type Mask struct {
	Parts           []MaskPart
	TotalValidBits  int
	isNegativeIndex bool // Internal flag to handle index resolution from the end of a slice
}

// NewContinuousMask creates a Mask for a continuous sequence of bits.
// It spans from a start byte and start bit for a given number of bits.
func NewContinuousMask(startWord int, startBit int, numberOfBits int) (*Mask, error) {
	if startBit < 0 || startBit > 7 {
		return nil, fmt.Errorf("startBit must be between 0 and 7, got %d", startBit)
	}
	if numberOfBits <= 0 {
		return nil, fmt.Errorf("numberOfBits must be positive, got %d", numberOfBits)
	}

	m := &Mask{
		TotalValidBits: numberOfBits,
	}

	currentWord := startWord
	bitsToMask := numberOfBits
	currentStartBit := startBit

	for bitsToMask > 0 {
		if currentWord < 0 {
			m.isNegativeIndex = true
		}

		bitsInThisWord := 8 - currentStartBit
		if bitsInThisWord > bitsToMask {
			bitsInThisWord = bitsToMask
		}

		mask, shift := generateMask(currentStartBit, bitsInThisWord)

		m.Parts = append(m.Parts, MaskPart{
			ByteIndex:     currentWord,
			Mask:          mask,
			Shift:         shift,
			ValidBitCount: bitsInThisWord,
		})

		bitsToMask -= bitsInThisWord
		currentWord++
		currentStartBit = 0 // Subsequent parts always start at bit 0
	}

	return m, nil
}

// generateMask creates a byte mask and right-shift amount for a given start bit and number of bits.
func generateMask(startBit int, numBits int) (mask byte, shift uint) {
	if numBits <= 0 {
		return 0, 0
	}
	// Create a block of 1s of length numBits.
	mask = byte((1 << numBits) - 1)
	// Shift the block to the correct starting position (from the left).
	mask <<= (8 - (startBit + numBits))
	// The shift value is the position of the rightmost bit.
	shift = uint(8 - (startBit + numBits))
	return mask, shift
}

// NewBitwiseMask creates a mask from specific word numbers and bit positions.
// `wordNo` contains the byte indices.
// `validBits` contains strings where each character represents a bit position (e.g., "012" for the first 3 bits).
func NewBitwiseMask(wordNo []int, validBits []string) (*Mask, error) {
	if len(wordNo) != len(validBits) {
		return nil, fmt.Errorf("wordNo and validBits slices must have the same length")
	}

	m := &Mask{}
	for i, word := range wordNo {
		if word < 0 {
			m.isNegativeIndex = true
		}

		bitStr := validBits[i]
		mask, shift := generateBitwiseMask(bitStr)
		bitCount := len(bitStr)

		m.Parts = append(m.Parts, MaskPart{
			ByteIndex:     word,
			Mask:          mask,
			Shift:         shift,
			ValidBitCount: bitCount,
		})
		m.TotalValidBits += bitCount
	}

	return m, nil
}

// generateBitwiseMask creates a mask from a string of bit positions (e.g., "017").
func generateBitwiseMask(validBits string) (byte, uint) {
	var mask byte
	// Initialize shift to a value larger than any possible shift.
	shift := uint(8)
	for _, char := range validBits {
		var currentShift uint
		switch char {
		case '0':
			mask |= 0x80
			currentShift = 7
		case '1':
			mask |= 0x40
			currentShift = 6
		case '2':
			mask |= 0x20
			currentShift = 5
		case '3':
			mask |= 0x10
			currentShift = 4
		case '4':
			mask |= 0x08
			currentShift = 3
		case '5':
			mask |= 0x04
			currentShift = 2
		case '6':
			mask |= 0x02
			currentShift = 1
		case '7':
			mask |= 0x01
			currentShift = 0
		default:
			continue // Ignore invalid characters
		}
		// Keep the smallest shift value, which corresponds to the rightmost bit.
		if currentShift < shift {
			shift = currentShift
		}
	}

	if shift == 8 { // No valid bits were found
		return 0, 0
	}
	return mask, shift
}

// ExtractUint64 applies the mask to a byte slice and returns the extracted value as a uint64.
func (m *Mask) ExtractUint64(frame []byte) (uint64, error) {
	if m.TotalValidBits > 64 {
		return 0, fmt.Errorf("mask extracts %d bits, which is too wide for uint64", m.TotalValidBits)
	}

	var result uint64
	remainingBits := m.TotalValidBits

	for _, part := range m.Parts {
		byteIndex := part.ByteIndex
		if m.isNegativeIndex {
			byteIndex = len(frame) + part.ByteIndex
		}

		if byteIndex < 0 || byteIndex >= len(frame) {
			return 0, fmt.Errorf("mask index %d is out of bounds for frame of length %d", part.ByteIndex, len(frame))
		}

		// Extract the relevant bits from the byte.
		val := (frame[byteIndex] & part.Mask) >> part.Shift

		// Shift the extracted value to its final position in the result and OR it.
		shiftAmount := remainingBits - part.ValidBitCount
		result |= (uint64(val) << shiftAmount)

		remainingBits -= part.ValidBitCount
	}
	return result, nil
}

// ExtractString applies the mask to a byte slice and returns the extracted value as a binary string.
func (m *Mask) ExtractString(frame []byte) (string, error) {
	var builder strings.Builder

	for _, part := range m.Parts {
		byteIndex := part.ByteIndex
		if m.isNegativeIndex {
			byteIndex = len(frame) + part.ByteIndex
		}

		if byteIndex < 0 || byteIndex >= len(frame) {
			return "", fmt.Errorf("mask index %d is out of bounds for frame of length %d", part.ByteIndex, len(frame))
		}

		val := (frame[byteIndex] & part.Mask) >> part.Shift

		// Use fmt.Sprintf with the %0*b verb to get a zero-padded binary string.
		builder.WriteString(fmt.Sprintf("%0*b", part.ValidBitCount, val))
	}

	return builder.String(), nil
}

// --- Integer Conversions (Big Endian) --- //

// Uint32ToBytesBE converts a uint32 to a 4-byte big-endian slice.
func Uint32ToBytesBE(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}

// BytesToUint32BE converts a big-endian byte slice to a uint32.
// Returns an error if the slice is smaller than 4 bytes.
func BytesToUint32BE(b []byte) (uint32, error) {
	if len(b) < 4 {
		return 0, fmt.Errorf("input slice too short: expected 4 bytes, got %d", len(b))
	}
	return binary.BigEndian.Uint32(b), nil
}

// Uint64ToBytesBE converts a uint64 to an 8-byte big-endian slice.
func Uint64ToBytesBE(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// BytesToUint64BE converts a big-endian byte slice to a uint64.
// Returns an error if the slice is smaller than 8 bytes.
func BytesToUint64BE(b []byte) (uint64, error) {
	if len(b) < 8 {
		return 0, fmt.Errorf("input slice too short: expected 8 bytes, got %d", len(b))
	}
	return binary.BigEndian.Uint64(b), nil
}

// --- Specialized Floating-Point Conversions --- //

const (
	// MIL-STD-1750A 32-bit format: 24-bit 2's complement mantissa, 8-bit 2's complement exponent.
	mantissaMask32    = 0xFFFFFF00
	exponentMask32    = 0x000000FF
	mantissaSignBit32 = 0x00800000         // Sign bit of the 24-bit mantissa
	mantissaSignExt32 = 0xFFFFFFFFFF000000 // Extension for negative mantissa
	mantissaNorm32    = -23                // Normalization factor for 24-bit mantissa (2^-23)

	// MIL-STD-1750A 48-bit format: 40-bit 2's complement mantissa, 8-bit 2's complement exponent.
	mantissaMask48    = 0xFFFFFFFFFF00
	exponentMask48    = 0x0000000000FF
	mantissaSignBit48 = 0x8000000000       // Sign bit of the 40-bit mantissa
	mantissaSignExt48 = 0xFFFF000000000000 // Extension for negative mantissa
	mantissaNorm48    = -39                // Normalization factor for 40-bit mantissa (2^-39)
)

// DecodeF1750A32 decodes a uint64 representing a MIL-STD-1750A 32-bit floating point number.
func DecodeF1750A32(input uint64) float64 {
	// Exponent is the last 8 bits.
	exponent := int8(input & exponentMask32)

	// Mantissa is the first 24 bits of the 32-bit value.
	mantissa := (input & mantissaMask32) >> 8

	// Manually sign-extend the 24-bit mantissa to 64 bits.
	if (mantissa & mantissaSignBit32) != 0 {
		mantissa |= mantissaSignExt32
	}

	// Normalize the mantissa and apply the exponent.
	fmant := math.Pow(2.0, mantissaNorm32) * float64(int64(mantissa))
	return fmant * math.Pow(2.0, float64(exponent))
}

// DecodeF1750A48 decodes a uint64 representing a MIL-STD-1750A 48-bit floating point number.
func DecodeF1750A48(input uint64) float64 {
	// In the 48-bit representation, the exponent is typically in bytes 2 and 3.
	// Assuming the input uint64 has the 48-bit value in the lower bits.
	// This implementation is based on the original code's bit shifting.
	exponent := int8((input >> 16) & 0xFF)

	// The original code constructs the mantissa from non-contiguous parts.
	// This suggests a specific packing format for the 48 bits within the 64-bit uint.
	mantissa := (input & 0xFFFFFF000000) >> 8
	mantissa |= (input & 0xFFFF)

	// Manually sign-extend the 40-bit mantissa to 64 bits.
	if (mantissa & mantissaSignBit48) != 0 {
		mantissa |= mantissaSignExt48
	}

	// Normalize the mantissa and apply the exponent.
	fmant := math.Pow(2.0, mantissaNorm48) * float64(int64(mantissa))
	return fmant * math.Pow(2.0, float64(exponent))
}

// --- String & Byte Conversions --- //

// BitStringToHexString converts a string of '0's and '1's to its hexadecimal representation.
func BitStringToHexString(s string) (string, error) {
	if len(s)%8 != 0 {
		return "", fmt.Errorf("input string length must be a multiple of 8, got %d", len(s))
	}

	var builder strings.Builder
	builder.Grow(len(s) / 4) // Pre-allocate memory: 8 bits -> 2 hex chars

	for i := 0; i < len(s); i += 8 {
		byteStr := s[i : i+8]
		val, err := strconv.ParseUint(byteStr, 2, 8)
		if err != nil {
			return "", fmt.Errorf("failed to parse binary string '%s': %w", byteStr, err)
		}
		// Format as a two-digit, zero-padded, uppercase hex string.
		if _, err := fmt.Fprintf(&builder, "%02X", val); err != nil {
			// This error is not expected with a strings.Builder
			return "", err
		}
	}
	return builder.String(), nil
}

// XORCardinality calculates the bitwise XOR between two byte slices and returns the cardinality
// (number of set bits) of the result.
func XORCardinality(a, b []byte) (int, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("slices must have the same length, got %d and %d", len(a), len(b))
	}
	cardinality := 0
	for i := range a {
		xorValue := a[i] ^ b[i]
		cardinality += bits.OnesCount8(xorValue)
	}
	return cardinality, nil
}
