// Package util provides utility functions and types for common operations.
package util

import (
	"context"
	"math"
	"sync"
	"time"
)

// Debouncer is a generic struct that caches a value with a delay to prevent excessive recomputation or fetching.
// It uses a fetcher function to retrieve the value and a timeout to manage the validity of the cached value.
type Debouncer[T any] struct {
	fetcher   func() (T, error)
	lastValue T
	timeOut   time.Time // time when value can be renewed
	delay     time.Duration
}

// NewDebouncer creates a new Debouncer with a specified delay and a fetcher function to retrieve values.
//
// Parameters:
//   - delay: The duration to wait before allowing a new value to be fetched
//   - fetcher: A function that returns a value of type T and an error
//
// Returns:
//   - A pointer to a new Debouncer instance
func NewDebouncer[T any](delay time.Duration, fetcher func() (T, error)) *Debouncer[T] {
	return &Debouncer[T]{
		fetcher: fetcher,
		delay:   delay,
	}
}

// GetValue retrieves the cached value or fetches a new one if the timeout has expired, updating the cache and timeout.
// If the timeout has expired, it calls the fetcher function to get a new value, updates the cache,
// and resets the timeout.
//
// Returns:
//   - The cached value of type T
//   - An error if one occurred during fetching (nil if using cached value)
func (d *Debouncer[T]) GetValue() (T, error) {
	if time.Now().After(d.timeOut) {
		value, err := d.fetcher()
		if err == nil {
			d.timeOut = time.Now().Add(d.delay)
			d.lastValue = value
		}
	}
	return d.lastValue, nil
}

// PubSubDebouncer is a utility for debouncing values with a pub-sub mechanism for notifying listeners.
// It holds the latest value and enforces a delay duration before updates are allowed.
// Values are fetched using a provided fetcher function, and a context is used for cancellation.
// It uses a PubSub instance to broadcast updates to registered listeners.
type PubSubDebouncer[T comparable] struct {
	lastValue   T
	listeners   *PubSub[T]
	delay       time.Duration
	timeOut     time.Time
	fetcherFunc func() (T, error)
	context     context.Context
	sync.RWMutex
}

// NewPubSubDebouncer creates a new PubSubDebouncer with a specified debounce delay and value-fetching function.
// Returns a pointer to the PubSubDebouncer and a context.CancelFunc for managing its lifecycle.
// The cancel function should be called when the debouncer is no longer needed to clean up resources.
//
// Parameters:
//   - delay: The minimum duration between value updates (must be >= 100ms)
//   - fetcher: A function that returns a value of type T and an error
//
// Returns:
//   - A pointer to a new PubSubDebouncer instance
//   - A context.CancelFunc that should be called to stop the debouncer when no longer needed
//
// Panics:
//   - If the delay is less than 100 milliseconds
func NewPubSubDebouncer[T comparable](delay time.Duration, fetcher func() (T, error)) (*PubSubDebouncer[T], context.CancelFunc) {
	if delay < 100*time.Millisecond {
		panic("delay must be greater >= 100 milliseconds")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	ret := &PubSubDebouncer[T]{
		listeners:   NewPubSub[T](),
		delay:       delay,
		fetcherFunc: fetcher,
		context:     ctx,
	}

	return ret, cancelFunc
}

// GetValue retrieves the current value, refreshing it using the fetcher function if the timeout has expired.
// This method is thread-safe and uses a read lock to protect access to the debouncer's state.
//
// Returns:
//   - The current value of type T (either cached or freshly fetched)
//   - An error if one occurred during fetching (nil if using cached value or if fetch was successful)
func (d *PubSubDebouncer[T]) GetValue() (T, error) {
	d.RLock()
	defer d.RUnlock()
	var err error

	if time.Now().After(d.timeOut) {
		d.lastValue, err = d.fetcherFunc()
		if err == nil {
			d.timeOut = time.Now().Add(d.delay)
		}
	}

	return d.lastValue, err
}

// GetValueMust retrieves the latest value from the debouncer and panics if an error occurs during fetching.
// This is a convenience method for cases where errors are not expected or should cause program termination.
//
// Returns:
//   - The current value of type T
//
// Panics:
//   - If an error occurs during value fetching
func (d *PubSubDebouncer[T]) GetValueMust() T {
	value, err := d.GetValue()
	if err != nil {
		panic(err)
	}
	return value
}

// SetValue updates the stored value if it differs from the previous value, considering type-specific precision thresholds.
// Triggers a broadcast to notify listeners if the value changes and adjusts the timeout duration.
// For floating-point types (float32, float64), it uses a small epsilon value to determine if the value has changed.
// For other types, it uses direct equality comparison.
//
// Parameters:
//   - value: The new value to store and potentially broadcast to listeners
func (d *PubSubDebouncer[T]) SetValue(value T) {
	doBroadcast := false

	d.Lock()

	switch any(value).(type) {
	case float64:
		if math.Abs(any(d.lastValue).(float64)-any(value).(float64)) > 0.00000001 {
			d.lastValue = value
			doBroadcast = true
		}
	case float32:
		if math.Abs(float64(any(d.lastValue).(float32)-any(value).(float32))) > 0.00001 {
			d.lastValue = value
			doBroadcast = true
		}
	default:
		if d.lastValue != value {
			d.lastValue = value
			doBroadcast = true
		}
	}

	d.timeOut = time.Now().Add(d.delay)
	d.Unlock()

	if doBroadcast {
		d.listeners.Broadcast(value)
		//log.Printf("changed value = %v", value)
	}
}

// Register registers a new listener channel for receiving debounced values. Starts the fetcher if no listeners are active.
// This method automatically starts the background fetcher goroutine if this is the first active listener.
//
// Returns:
//   - A receive-only channel that will receive updates when the debounced value changes
func (d *PubSubDebouncer[T]) Register() <-chan T {
	wasActive := d.listeners.IsActive()
	newListener := d.listeners.Register(1)
	if !wasActive {
		go d.fetcher() // start the fetcher if it wasn't running before
	}
	return newListener
}

// Unregister removes a specified listener channel from the debouncer's PubSub, closing the channel and freeing resources.
// If this was the last active listener, the background fetcher goroutine will automatically stop on its next iteration.
//
// Parameters:
//   - c: The channel previously returned by Register() that should be unregistered
func (d *PubSubDebouncer[T]) Unregister(c <-chan T) {
	d.listeners.Unregister(c)
}

// fetcher periodically fetches a value using fetcherFunc and updates the debouncer's state, notifying listeners if active.
// This is an internal method that runs as a goroutine, continuously fetching values at the specified delay interval.
// It automatically exits when either the context is canceled or there are no more active listeners.
func (d *PubSubDebouncer[T]) fetcher() {
	//log.Printf("starting fetcher")
	for {
		value, err := d.fetcherFunc()
		if err == nil {
			d.SetValue(value)
		}
		select {
		case <-d.context.Done():
			return
		case <-time.After(d.delay):
		}

		if !d.listeners.IsActive() {
			//		log.Printf("no listeners, exiting fetcher")
			return
		}
	}
}
