// Package util provides utility functions and types for common operations.
package util

// IsASCIIDigits returns true if s is non-empty and every rune is an ASCII digit ('0'..'9').
func IsASCIIDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !isASCIIDigit(r) {
			return false
		}
	}
	return true
}

// isASCIIDigit reports whether r is in the ASCII range '0'..'9'.
func isASCIIDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
