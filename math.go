// Package util provides utility functions and types for common operations.
package util

import "math"

// AlmostEqual compares two floating-point numbers and determines if they are approximately equal.
// It uses a fixed epsilon value of 1e-6 for the comparison, which is suitable for most general-purpose
// floating-point comparisons where absolute precision is not required.
//
// Parameters:
//   - a, b: The two floating-point values to compare (can be either float32 or float64)
//
// Returns:
//   - true if the absolute difference between a and b is less than 1e-6
//   - false otherwise
func AlmostEqual[T ~float64 | ~float32](a, b T) bool {
	delta := math.Abs(float64(a - b))
	return delta < 1e-6
}
