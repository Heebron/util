package main

import (
	"runtime"
	"testing"
)

func TestCount(t *testing.T) {

	var sampleData []int

	for i := 0; i < runtime.NumCPU()*3; i++ {
		sampleData = make([]int, 0)
		for j := 0; j < i; j++ {
			sampleData = append(sampleData, j)
		}
		c := Count(sampleData, func(ele int) bool { return true })

		if int(c) != i {
			t.Fatal("expected ", i, " got ", c)
		}
	}
}
