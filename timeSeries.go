// Package util provides utility functions and types for common operations.
package util

import (
	"slices"
	"time"
)

// TimeSeries represents a collection of time-stamped entries, indexed with specified time precision.
// PrecisionMask defines the time granularity used to group entries.
// The data field stores the actual time series as a map of time.Time to values of type *T.
type TimeSeries[T any] struct {
	PrecisionMask time.Duration
	data          map[time.Time]*T
}

// NewTimeSeriesMap initializes and returns a new TimeSeries instance with the specified time precision. This implementation
// is NOT concurrent safe.
func NewTimeSeriesMap[T any](precision time.Duration) *TimeSeries[T] {
	return &TimeSeries[T]{
		PrecisionMask: precision,
		data:          make(map[time.Time]*T),
	}
}

// Get retrieves the value associated with a specific time, truncated to the series' precision. Returns nil if not found.
func (ts *TimeSeries[T]) Get(t time.Time) *T {
	return ts.data[t.In(time.UTC).Truncate(ts.PrecisionMask)]
}

// Update adds or updates an entry in the time series for the specified time.
// If the time entry does not exist, it is created using the init function.
// If the time entry exists, it is modified using the update function.
func (ts *TimeSeries[T]) Update(t time.Time, init func() *T, update func(*T)) {
	key := t.In(time.UTC).Truncate(ts.PrecisionMask)
	v, ok := ts.data[key]
	if !ok {
		v = init()
		ts.data[key] = v
	} else {
		update(v)
	}
}

// Len returns the number of entries in the TimeSeries.
func (ts *TimeSeries[T]) Len() int {
	return len(ts.data)
}

// Clear removes all entries from the TimeSeries, resetting its data storage to an empty state.
func (ts *TimeSeries[T]) Clear() {
	clear(ts.data)
}

// Keys returns a sorted slice of all keys present in the TimeSeries.
func (ts *TimeSeries[T]) Keys() []time.Time {
	keys := make([]time.Time, 0, len(ts.data))
	for k := range ts.data {
		keys = append(keys, k)
	}

	slices.SortFunc(keys, func(i, j time.Time) int {
		return int(i.UnixMicro() - j.UnixMicro())
	})

	return keys
}
