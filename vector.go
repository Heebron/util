// Package util provides utility functions and types for common operations.
package util

import (
	"fmt"
	"math"
)

// Signed is a constraint interface that allows for numeric types that can be
// used in vector calculations. This includes integers and floating point numbers.
type Signed interface {
	~int | ~int64 | ~int32 | ~float64 | ~float32
}

// Vector2D represents a 2D vector with x and y components.
// It is a generic type that can work with any numeric type that satisfies the Signed interface.
type Vector2D[T Signed] struct {
	dx, dy T // dx represents the x-component, dy represents the y-component
}

// NewVector2D creates and returns a new Vector2D with the specified x and y components.
// It uses generics to work with any numeric type that satisfies the Signed interface.
//
// Parameters:
//   - dx: The x-component of the vector
//   - dy: The y-component of the vector
//
// Returns:
//   - A pointer to a new Vector2D instance with the specified components
func NewVector2D[T Signed](dx, dy T) *Vector2D[T] {
	return &Vector2D[T]{dx, dy}
}

// Magnitude computes and returns the length (magnitude) of the vector as a float64.
// It uses the Pythagorean theorem to calculate the distance from the origin.
//
// Returns:
//   - The magnitude of the vector as a float64 value
func (v *Vector2D[T]) Magnitude() float64 {
	return math.Sqrt(float64(v.dx*v.dx) + float64(v.dy*v.dy))
}

// ThetaRadians returns the angle this vector represents in radians.
// Angle 0 is on the positive x-axis (horizontal line) and increases counter-clockwise
// around to 2*PI. Negative angles are converted to their positive equivalent.
//
// Returns:
//   - The angle of the vector in radians (range 0 to 2*PI)
func (v *Vector2D[T]) ThetaRadians() float64 {
	r := math.Atan2(float64(v.dy), float64(v.dx))
	if r < 0 {
		return (2.0 * math.Pi) + r
	}
	return r
}

// ThetaDegrees returns the angle this vector represents in degrees.
// Angle 0 is on the positive x-axis (horizontal line) and increases counter-clockwise
// around to 360 degrees. This is a convenience method that converts the result of
// ThetaRadians to degrees.
//
// Returns:
//   - The angle of the vector in degrees (range 0 to 360)
func (v *Vector2D[T]) ThetaDegrees() float64 {
	return v.ThetaRadians() * 180.0 / math.Pi
}

// Add updates this vector to be the sum of this vector and the other vector.
// The operation is performed component-wise, adding the dx and dy values separately.
// This method modifies the vector in place.
//
// Parameters:
//   - other: The vector to add to this vector
//
// Returns:
//   - A pointer to this vector after the addition operation, allowing for method chaining
func (v *Vector2D[T]) Add(other *Vector2D[T]) *Vector2D[T] {
	v.dx = v.dx + other.dx
	v.dy = v.dy + other.dy
	return v
}

// Scale multiplies both components of this vector by the given factor.
// This method modifies the vector in place, scaling its magnitude while
// preserving its direction (unless factor is negative).
//
// Parameters:
//   - factor: The value to multiply both components by
//
// Returns:
//   - A pointer to this vector after the scaling operation, allowing for method chaining
func (v *Vector2D[T]) Scale(factor float64) *Vector2D[T] {
	v.dx = T(float64(v.dx) * factor)
	v.dy = T(float64(v.dy) * factor)
	return v
}

// ScaleToMaxComponent modifies this vector so that the largest component becomes 1.0,
// and the other component is scaled proportionally. This is not a true
// normalization to unit length (see NormalizeToUnit for that), but rather
// a simplification of the vector's representation.
//
// Returns:
//   - A pointer to this vector after the scaling operation, allowing for method chaining
func (v *Vector2D[T]) ScaleToMaxComponent() *Vector2D[T] {
	if v.dy > v.dx {
		v.dx = v.dx / v.dy
		v.dy = 1.0
	} else {
		v.dy = v.dy / v.dx
		v.dx = 1.0
	}
	return v
}

// NormalizeToUnit modifies this vector to have a magnitude (length) of 1.0 while
// preserving its direction. This is a true normalization that creates a unit vector.
// If the vector has zero magnitude, it remains unchanged to avoid division by zero.
//
// Returns:
//   - A pointer to this vector after the normalization operation, allowing for method chaining
func (v *Vector2D[T]) NormalizeToUnit() *Vector2D[T] {
	mag := v.Magnitude()
	if mag > 0 {
		v.dx = T(float64(v.dx) / mag)
		v.dy = T(float64(v.dy) / mag)
	}
	return v
}

// String returns a string representation of the vector in the format {dx, dy}.
// This method implements the fmt.Stringer interface, allowing the vector to be
// used directly with fmt.Print functions.
//
// Returns:
//   - A string representation of the vector with components formatted to 4 significant digits
func (v *Vector2D[T]) String() string {
	return fmt.Sprintf("{%.4v, %.4v}", v.dx, v.dy)
}

// Copy creates and returns a new Vector2D with the same component values as this vector.
// This method does not modify the original vector and returns a pointer to the new copy.
//
// Returns:
//   - A pointer to a new Vector2D instance with the same component values as this vector
func (v *Vector2D[T]) Copy() *Vector2D[T] {
	return NewVector2D(v.dx, v.dy)
}

// Translate moves this vector by adding the specified dx and dy values to its components.
// This method modifies the vector in place.
//
// Parameters:
//   - dx: The amount to add to the x-component
//   - dy: The amount to add to the y-component
//
// Returns:
//   - A pointer to this vector after the translation operation, allowing for method chaining
func (v *Vector2D[T]) Translate(dx, dy T) *Vector2D[T] {
	v.dx += dx
	v.dy += dy
	return v
}

// X returns the x-component of the vector.
//
// Returns:
//   - The x-component (dx) of the vector
func (v *Vector2D[T]) X() T {
	return v.dx
}

// Y returns the y-component of the vector.
//
// Returns:
//   - The y-component (dy) of the vector
func (v *Vector2D[T]) Y() T {
	return v.dy
}

// ScaleXY multiplies the x-component by factorX and the y-component by factorY.
// This method modifies the vector in place, allowing for non-uniform scaling
// that can change both the magnitude and direction of the vector.
//
// Parameters:
//   - factorX: The value to multiply the x-component by
//   - factorY: The value to multiply the y-component by
//
// Returns:
//   - A pointer to this vector after the scaling operation, allowing for method chaining
func (v *Vector2D[T]) ScaleXY(factorX, factorY float64) *Vector2D[T] {
	v.dx = T(float64(v.dx) * factorX)
	v.dy = T(float64(v.dy) * factorY)
	return v
}

// Rotate rotates this vector counter-clockwise by the given angle in radians.
// The rotation is performed in place using the standard 2D rotation matrix:
//   x' = x*cos(theta) - y*sin(theta)
//   y' = x*sin(theta) + y*cos(theta)
//
// Parameters:
//   - theta: rotation angle in radians (positive is counter-clockwise)
//
// Returns:
//   - A pointer to this vector after rotation, allowing method chaining
func (v *Vector2D[T]) Rotate(theta float64) *Vector2D[T] {
	cosT := math.Cos(theta)
	sinT := math.Sin(theta)
	x := float64(v.dx)
	y := float64(v.dy)
	vx := x*cosT - y*sinT
	vy := x*sinT + y*cosT
	v.dx = T(vx)
	v.dy = T(vy)
	return v
}
