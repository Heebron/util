package util

import (
	"sync/atomic"
)

// WorkerPool is a generic worker pool that processes work items concurrently and stores results.
// It uses channels to dispatch work, collect results, and retains a function to process each work item.
type WorkerPool[W any, R any] struct {
	work   chan W
	result chan R
	f      func(W) R
	count  atomic.Int32
}

// NewWorkerPool creates and initializes a new WorkerPool with the given worker function, backlog size, and number of workers.
// It panics if backlog is negative or numWorkers is less than 1.
func NewWorkerPool[W any, R any](f func(W) R, backlog int, numWorkers int) *WorkerPool[W, R] {
	if backlog < 0 {
		panic("backlog must be greater than -1")
	}
	if numWorkers < 1 {
		panic("numWorkers must be greater than zero")
	}
	pool := &WorkerPool[W, R]{result: make(chan R, backlog), work: make(chan W, backlog), f: f}

	// start n workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for w := range pool.work { // until closed
				pool.result <- pool.f(w)
			}
		}()
	}
	return pool
}

// Close terminates the work channel, signaling that no more work items will be submitted to the WorkerPool.
func (wp *WorkerPool[W, R]) Close() {
	close(wp.work)
}

// Post submits a work item to the WorkerPool for processing and increments the active work count. This will block when the channel is full.
func (wp *WorkerPool[W, R]) Post(w W) {
	wp.count.Add(1)
	wp.work <- w
}

// Result retrieves and returns the next available result from the worker pool, decrementing the active work count.
func (wp *WorkerPool[W, R]) Result() R {
	v := <-wp.result
	wp.count.Add(-1)
	return v
}

// Len returns the current number of active work items in the WorkerPool.
func (wp *WorkerPool[W, R]) Len() int32 {
	return wp.count.Load()
}

// IsActive checks if the WorkerPool has active work items by verifying if the active work count is greater than zero.
func (wp *WorkerPool[W, R]) IsActive() bool {
	return wp.count.Load() > 0
}
