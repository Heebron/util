package util

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {

	wp := NewWorkerPool(func(a int) int { return a }, 0, runtime.NumCPU())

	sumA := 0
	sumB := 0

	go func() {
		for i := 1; i <= 10; i++ {
			wp.Post(i)
			sumA += i
		}
		wp.Close()
	}()

	for i := 1; i <= 10; i++ {
		sumB += wp.Result()
	}

	if wp.IsActive() {
		t.Fail()
	}

	if sumA != sumB {
		t.Fail()
	}
}

func TestInterruptibleWorkerPool(t *testing.T) {
	wp := NewWorkerPool(func(a context.Context) int {
		select {
		case <-a.Done():
			return 1
		case <-time.NewTicker(5000 * time.Millisecond).C:
			return 0
		}

	}, 10, runtime.NumCPU())

	ctx, cancelFunc := context.WithCancel(context.Background())

	wp.Post(ctx)

	time.Sleep(100 * time.Millisecond)
	cancelFunc()
	time.Sleep(100 * time.Millisecond)
	if wp.Result() != 1 {
		t.Fail()
	}
	if wp.IsActive() {
		t.Fail()
	}
}
