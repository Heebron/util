package util

import (
	"runtime"
)

// Count counts the elements in the slice `d` satisfying the predicate function `f` and returns the count as int64.
func Count[T any](d []T, f func(ele T) bool) int64 {
	cpuCount := runtime.NumCPU()
	var total int64

	block := len(d) / cpuCount
	blockRemainder := len(d) % cpuCount
	start := 0

	if block > 0 {
		results := make(chan int64)
		for i := 0; i < cpuCount; i++ {
			go func(s, e int) {
				var count int64
				for _, ele := range d[s:e] {
					if f(ele) {
						count++
					}
				}
				results <- count
			}(start, start+block)
			start += block
		}

		for ; cpuCount > 0; cpuCount-- {
			total += <-results
		}
	}

	if blockRemainder > 0 {
		for _, ele := range d[len(d)-blockRemainder:] {
			if f(ele) {
				total++
			}
		}
	}

	return total
}
