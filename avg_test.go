package util

import (
	"math"
	"testing"
)

func TestAverage(t *testing.T) {
	avg := NewAverage(1.0)
	avg.Update(1.0)
	if avg.Get() != 1.0 {
		t.Fatalf("avg(1.0,1.0)=%v", avg.Get())
	}
	avg = NewAverage(1.0)
	avg.Update(-1.0)
	if avg.Get() != 0.0 {
		t.Fatalf("avg(1.0,-1.0)=%v", avg.Get())
	}
	avg = NewAverage(1.0)
	avg.Update(0.0)
	if avg.Get() != 0.5 {
		t.Fatalf("avg(1.0,0.5)=%v", avg.Get())
	}
}
func TestAverage_UpdateGet(t *testing.T) {
	avg := NewAverage(1)
	v := avg.UpdateGet(0)
	if v != 0.5 {
		t.Fatalf("avg(1,0)=%v", v)
	}
}

func TestMovingAverage(t *testing.T) {

	values := []float64{
		1.1666666666666665,
		1.472222222222222,
		1.8935185185185186,
		2.4112654320987654,
		3.0093878600823043,
		3.67448988340192,
		4.395408236168267,
		5.162840196806889,
		5.969033497339074,
	}

	avg := NewMovingAverage(3, 1)

	for i, j := range values {
		avg.Update(float64(i + 2))
		if math.Abs(j-avg.Get()) != 0.0 {
			t.Fatalf("%f!=%f", avg.Get(), j)
		}
	}
}
