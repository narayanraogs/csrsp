package stages

// ComplementConfig holds the shared configuration for Ones and Twos Complement stages.
type ComplementConfig struct {
	// The size of the word in bits to operate on (e.g., 16, 32). Must be a multiple of 8.
	BitsPerWord int
	// The number of least significant bits to apply the complement to.
	BitsToComplement int
}

// makeMask creates the bitmask used for the complement operation.
func makeMask(noOfBytes int, bitsToComplement int) []byte {
	var mask = make([]byte, noOfBytes)
	for i := noOfBytes - 1; i >= 0; i-- {
		mask[i] = getByteMask(bitsToComplement)
		bitsToComplement -= 8
	}
	return mask
}

// getByteMask calculates the mask for a single byte.
func getByteMask(noOfBits int) byte {
	if noOfBits >= 8 {
		return 0xFF
	}
	if noOfBits <= 0 {
		return 0x00
	}
	// Equivalent to (1 << noOfBits) - 1
	return byte((1 << noOfBits) - 1)
}
