// Package util provides utility functions and types for common operations.
//
// This file contains simple helpers for computing arithmetic means: a running
// average over all observed values (Average) and a smoothed moving average
// (MovingAverage). These types are lightweight and not concurrency-safe; if
// you need to use them from multiple goroutines, add your own synchronization.
package util

// Average represents a simple arithmetic mean calculator over an entire set of
// values. It maintains a running sum (stored in the field "last") and a count
// of observed values to compute the mean on demand.
//
// Zero value and initialization:
//   - The zero value of Average is not ready for use. Calling Get() on a zero
//     value will cause a division-by-zero panic because count is 0.
//   - Use NewAverage(initial) to construct a valid instance or call Reset(initial)
//     before the first use.
//
// Concurrency: Average is not safe for concurrent use by multiple goroutines.
//
// Note: The field name "last" actually stores the running sum of all values
// seen so far, not the last sample value.
type Average struct {
	last  float64 // Running sum of all values
	count int64   // Number of values included in the average
}

// NewAverage returns a new Average initialized with the provided first value.
//
// Parameters:
//   - initial: The first value to include in the average.
//
// Returns: A pointer to a new Average instance initialized with the given value.
func NewAverage(initial float64) *Average {
	return &Average{initial, 1}
}

// Reset discards all previously accumulated values and initializes the average
// with the supplied initial value.
//
// Parameters:
//   - initial: The new initial value to set for the average.
func (m *Average) Reset(initial float64) {
	m.count = 1
	m.last = initial
}

// Update adds a new sample to the running average.
//
// Parameters:
//   - v: The new value to include in the average.
func (m *Average) Update(v float64) {
	m.last += v
	m.count++
}

// UpdateGet adds a new sample and returns the current average in a single call.
func (m *Average) UpdateGet(v float64) float64 {
	m.Update(v)
	return m.Get()
}

// Get returns the current arithmetic mean of all values observed so far.
//
// Edge cases:
//   - If the Average was not initialized via NewAverage or Reset, and count is 0,
//     this will panic due to division by zero.
func (m *Average) Get() float64 {
	return m.last / float64(m.count)
}

// MovingAverage maintains state for a smoothed moving average.
//
// Semantics:
//   - NewMovingAverage(n, initial) configures smoothing based on n: larger n
//     produces a slower, more stable response; smaller n reacts faster to new
//     samples. The update rule applies a fraction of the delta between the new
//     value and the current value. While not a textbook EMA, it provides a
//     simple, stable smoothing behavior for many use cases.
//
// Concurrency: MovingAverage is not safe for concurrent use by multiple goroutines.
type MovingAverage struct {
	last   float64
	factor float64
}

// NewMovingAverage returns a MovingAverage configured with a smoothing window
// parameter n and an initial value. 'n' represents an intended window size: a
// higher 'n' yields a more stable (less volatile) output.
//
// Panics if n <= 0.
func NewMovingAverage(n uint, initial float64) *MovingAverage {
	if n <= 0 {
		panic("n must be > 0")
	}
	return &MovingAverage{initial, 1.0 / float64(n)}
}

// Reset sets the moving average back to the provided initial value without
// changing the smoothing configuration.
func (m *MovingAverage) Reset(initial float64) {
	m.last = initial
}

// Update incorporates another data point into the moving average.
// The internal smoothing factor is derived from n passed to NewMovingAverage.
func (m *MovingAverage) Update(v float64) {
	m.last = (((v-m.last)*m.factor + m.last) + m.last) / 2.0
}

// UpdateAndGet adds a new data point and returns the current moving average in
// one call.
func (m *MovingAverage) UpdateAndGet(v float64) float64 {
	m.Update(v)
	return m.Get()
}

// Get returns the current moving average value.
func (m *MovingAverage) Get() float64 {
	return m.last
}
