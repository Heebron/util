package util

import "testing"

func TestAllDigitsEmpty(t *testing.T) {
	if IsASCIIDigits("") {
		t.Fail()
	}
}

func TestAllDigits(t *testing.T) {
	if !IsASCIIDigits("0123456789") {
		t.Fail()
	}

	if IsASCIIDigits("a") {
		t.Fail()
	}
}
