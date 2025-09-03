// Package util provides utility functions and types for common operations.
package util

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Symbol constants used for indicating value trends in TrackedValue.
const (
	upArrow   = "↑" // Symbol indicating an upward trend
	downArrow = "↓" // Symbol indicating a downward trend
	balance   = "☯" // Symbol indicating a balanced or unchanged value
)

// ValueIF is an interface that represents a value which can be tracked and updated.
// It extends `fmt.Stringer` interface and requires implementations to have the following methods:
// - Update(v T) bool: updates the value with the provided value `v`, returns `true` if the update was successful.
// - Value() T: retrieves the current value of the object.
type ValueIF[T constraints.Float] interface {
	fmt.Stringer
	Update(v T) bool
	Value() T
}

// TrackedValue represents a value that can be tracked and updated.
type TrackedValue[T constraints.Float] struct {
	value       T
	stringValue string
	symbol      string
	isCurrency  bool
	mvAvg       *MovingAverage
}

// NewTrackedValue initializes a new TrackedValue with the given parameters.
// The TrackedValue tracks the value and calculates a string representation including a trending symbol. The calculation of
// value is based on a moving average. The trend indicator is adjusted based on the last update value compared to the average.
//
// Parameters:
//   - v: the initial value of the TrackedValue
//   - isCurrency: determines whether the TrackedValue represents a currency or not
//   - window: specifies the number of values to consider when calculating the moving average
//
// Returns:
//   - *TrackedValue[T]: a pointer to the newly created TrackedValue instance
func NewTrackedValue[T constraints.Float](v T, isCurrency bool, window uint) *TrackedValue[T] {
	tv := &TrackedValue[T]{symbol: balance, value: v, isCurrency: isCurrency, mvAvg: NewMovingAverage(window, float64(v))}
	tv.stringValue = tv.calcString()
	return tv
}

// String returns the string representation of the TrackedValue.
// It retrieves the pre-calculated stringValue field.
// Returns:
//   - string: the string representation of the TrackedValue.
func (t *TrackedValue[T]) String() string {
	return t.stringValue
}

// Update updates the TrackedValue. It calculates the average value, compares it to the new value, and updates
// the symbol accordingly. Finally, it calculates a new string representation based on the value and symbol.
// If the string representation changes, the returned value is true, else false.
//
// Parameters:
// - v: the new value to update the TrackedValue with
//
// Returns:
// - changed: a boolean indicating if the value has changed compared to the previous value
func (t *TrackedValue[T]) Update(v T) (changed bool) {
	avg := int64(t.mvAvg.UpdateAndGet(float64(v)) * 100.00)
	datum := int64(v * 100.00)

	if avg == datum {
		t.symbol = balance
	} else if avg > datum {
		t.symbol = downArrow
	} else {
		t.symbol = upArrow
	}

	t.value = v

	newString := t.calcString()

	changed = newString != t.stringValue

	t.stringValue = newString

	return
}

// Value returns the current value stored in the TrackedValue.
// It retrieves the value field.
// Returns:
//   - T: the current value stored in the TrackedValue.
func (t *TrackedValue[T]) Value() T {
	return t.value
}

// calcString returns the string representation of the TrackedValue based on its value, symbol, and isCurrency flag.
// If isCurrency is true, it formats the value as dollars and appends the symbol.
// If isCurrency is false, it formats the value as a number and appends the symbol.
// Returns:
//   - string: the formatted string representation of the TrackedValue.
func (t *TrackedValue[T]) calcString() string {
	if t.isCurrency {
		return DollarsFormat.FormatMoney(t.value) + t.symbol
	} else {
		return NumberFormat.FormatMoney(t.value) + t.symbol
	}
}
