// Package fs provides file system utility functions.
package fs

import "os"

// Exists checks if a file or directory exists at the given path.
// It returns false for an empty path. For non-empty paths, it returns true
// if and only if os.Stat returns a nil error.
func Exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}
