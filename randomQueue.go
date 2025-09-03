// Package util contains small helper types and functions shared across the project.
//
// This file provides a generic, concurrency-safe RandomQueue that allows adding
// items and retrieving/removing a random item in O(1) time on average using a
// swap-with-last pop technique.
package util

import (
	"math/rand"
	"sync"
)

// RandomQueue is a generic, mutex-protected container that stores items and
// returns a uniformly random element upon removal.
//
// Concurrency:
//   - All exported methods acquire an internal mutex; callers can safely use the
//     same queue from multiple goroutines.
//
// Performance:
//   - GetAndRemove removes a random element in amortized O(1) time by swapping
//     the chosen element with the last element, then shrinking the slice.
//
// Zero value:
//   - The zero value of RandomQueue is not ready for use because the mutex is a
//     pointer. Use NewRandomQueue to construct a working instance.
//
// Edge cases:
//   - Calling GetAndRemove on an empty queue will panic due to rand.Intn(0).
//     Call Len() first if emptiness is possible.
//
// Example:
//
//	q := util.NewRandomQueue[int]()
//	q.Add(1); q.Add(2); q.Add(3)
//	x := q.GetAndRemove() // x is 1, 2, or 3 with equal probability
//	n := q.Len()          // remaining item count
//
// T can be any type.
type RandomQueue[T any] struct {
	lock  *sync.Mutex
	files []T
}

// Add inserts a single item into the queue.
func (r *RandomQueue[T]) Add(f T) {
	r.lock.Lock()
	r.files = append(r.files, f)
	r.lock.Unlock()
}

// GetAndRemove selects a uniformly random item from the queue, removes it, and
// returns it.
//
// Note: This will panic if the queue is empty. Use Len() to guard if needed.
func (r *RandomQueue[T]) GetAndRemove() T {
	r.lock.Lock()
	offset := rand.Intn(len(r.files))
	f := r.files[offset]
	r.files[offset] = r.files[len(r.files)-1] // move the last into the slice
	r.files = r.files[:len(r.files)-1]        // shorten the slice by 1
	r.lock.Unlock()
	return f
}

// Len returns the number of items currently stored in the queue.
func (r *RandomQueue[T]) Len() int {
	r.lock.Lock()
	defer r.lock.Unlock()
	return len(r.files)
}

// NewRandomQueue constructs a new, ready-to-use RandomQueue.
func NewRandomQueue[T any]() *RandomQueue[T] {
	return &RandomQueue[T]{
		lock:  &sync.Mutex{},
		files: make([]T, 0),
	}
}
