package util

import (
	"testing"
)

func TestNewVector2D0(t *testing.T) {
	v := NewVector2D(1.0, 0.0)

	if v.Magnitude() != 1 {
		t.Fail()
	}
	if v.ThetaDegrees() != 0.0 {
		t.Fail()
	}
}

func TestNewVector2D90(t *testing.T) {
	v := NewVector2D(0.0, 1.0)

	if v.Magnitude() != 1 {
		t.Fail()
	}
	if v.ThetaDegrees() != 90.0 {
		t.Fail()
	}
}

func TestNewVector2D180(t *testing.T) {
	v := NewVector2D(-1.0, 0.0)

	if v.Magnitude() != 1 {
		t.Fail()
	}
	if v.ThetaDegrees() != 180.0 {
		t.Fail()
	}
}

func TestNewVector2D270(t *testing.T) {
	v := NewVector2D(0.0, -1.0)

	if v.Magnitude() != 1 {
		t.Fail()
	}

	if v.ThetaDegrees() != 270.0 {
		t.Fail()
	}
}

func TestNormalizeToUnit(t *testing.T) {
	// Test case 1: Vector2D with non-unit magnitude
	v1 := NewVector2D(3.0, 4.0)
	originalMagnitude := v1.Magnitude()
	originalTheta := v1.ThetaRadians()

	if originalMagnitude != 5.0 {
		t.Errorf("Expected magnitude 5.0, got %f", originalMagnitude)
	}

	v1.NormalizeToUnit()

	// Check magnitude is now 1.0
	if mag := v1.Magnitude(); mag < 0.9999 || mag > 1.0001 {
		t.Errorf("Expected magnitude 1.0 after normalization, got %f", mag)
	}

	// Check direction is preserved
	if theta := v1.ThetaRadians(); !AlmostEqual(theta, originalTheta) {
		t.Errorf("Direction changed after normalization. Expected %f, got %f", originalTheta, theta)
	}

	// Test case 2: Vector2D already with unit magnitude
	v2 := NewVector2D(1.0, 0.0)
	v2.NormalizeToUnit()

	if mag := v2.Magnitude(); mag != 1.0 {
		t.Errorf("Expected magnitude to remain 1.0, got %f", mag)
	}

	// Test case 3: Zero vector
	v4 := NewVector2D(0.0, 0.0)
	v4.NormalizeToUnit()

	if v4.dx != 0.0 || v4.dy != 0.0 {
		t.Errorf("Expected zero vector to remain unchanged, got {%f, %f}", v4.dx, v4.dy)
	}
}

func TestCopy(t *testing.T) {
	// Test case 1: Copy has the same values as the original
	original := NewVector2D(3.0, 4.0)
	theCopy := original.Copy()

	// Check that theCopy has the same values as original
	if theCopy.dx != original.dx || theCopy.dy != original.dy {
		t.Errorf("Copy values don't match original. Expected {%f, %f}, got {%f, %f}",
			original.dx, original.dy, theCopy.dx, theCopy.dy)
	}

	// Test case 2: Modifying the theCopy doesn't affect the original
	theCopy.dx = 5.0
	theCopy.dy = 6.0

	if original.dx != 3.0 || original.dy != 4.0 {
		t.Errorf("Original was modified when theCopy was changed. Expected {3, 4}, got {%f, %f}",
			original.dx, original.dy)
	}

	// Test case 3: Copy is a different object than the original
	if original == theCopy {
		t.Errorf("Copy and original are the same object, expected different objects")
	}

	// Test case 4: Copy of zero vector
	zeroVector := NewVector2D(0.0, 0.0)
	zeroCopy := zeroVector.Copy()

	if zeroCopy.dx != 0.0 || zeroCopy.dy != 0.0 {
		t.Errorf("Zero vector theCopy has non-zero values. Got {%f, %f}", zeroCopy.dx, zeroCopy.dy)
	}
}

func TestTranslate(t *testing.T) {
	// Test case 1: Translate by positive values
	v1 := NewVector2D(3.0, 4.0)
	v1.Translate(2.0, 3.0)

	if v1.dx != 5.0 || v1.dy != 7.0 {
		t.Errorf("Translation by positive values failed. Expected {5, 7}, got {%f, %f}", v1.dx, v1.dy)
	}

	// Test case 2: Translate by negative values
	v2 := NewVector2D(3.0, 4.0)
	v2.Translate(-1.0, -2.0)

	if v2.dx != 2.0 || v2.dy != 2.0 {
		t.Errorf("Translation by negative values failed. Expected {2, 2}, got {%f, %f}", v2.dx, v2.dy)
	}

	// Test case 3: Translate by zero values (no change)
	v3 := NewVector2D(3.0, 4.0)
	v3.Translate(0.0, 0.0)

	if v3.dx != 3.0 || v3.dy != 4.0 {
		t.Errorf("Translation by zero values changed the vector. Expected {3, 4}, got {%f, %f}", v3.dx, v3.dy)
	}

	// Test case 4: Method chaining
	v4 := NewVector2D(1.0, 1.0)
	result := v4.Translate(2.0, 2.0)

	if result != v4 {
		t.Errorf("Method chaining failed. Translate() did not return the vector itself")
	}

	// Test case 5: Multiple translations
	v5 := NewVector2D(0.0, 0.0)
	v5.Translate(1.0, 2.0).Translate(3.0, 4.0)

	if v5.dx != 4.0 || v5.dy != 6.0 {
		t.Errorf("Multiple translations failed. Expected {4, 6}, got {%f, %f}", v5.dx, v5.dy)
	}
}

func TestAccessors(t *testing.T) {
	// Test case 1: Positive values
	v1 := NewVector2D(3.0, 4.0)

	if x := v1.X(); x != 3.0 {
		t.Errorf("X() accessor failed. Expected 3.0, got %f", x)
	}

	if y := v1.Y(); y != 4.0 {
		t.Errorf("Y() accessor failed. Expected 4.0, got %f", y)
	}

	// Test case 2: Negative values
	v2 := NewVector2D(-2.5, -3.5)

	if x := v2.X(); x != -2.5 {
		t.Errorf("X() accessor failed with negative value. Expected -2.5, got %f", x)
	}

	if y := v2.Y(); y != -3.5 {
		t.Errorf("Y() accessor failed with negative value. Expected -3.5, got %f", y)
	}

	// Test case 3: Zero values
	v3 := NewVector2D(0.0, 0.0)

	if x := v3.X(); x != 0.0 {
		t.Errorf("X() accessor failed with zero value. Expected 0.0, got %f", x)
	}

	if y := v3.Y(); y != 0.0 {
		t.Errorf("Y() accessor failed with zero value. Expected 0.0, got %f", y)
	}

	// Test case 4: After modification
	v4 := NewVector2D(1.0, 1.0)
	v4.Translate(2.0, 3.0)

	if x := v4.X(); x != 3.0 {
		t.Errorf("X() accessor failed after modification. Expected 3.0, got %f", x)
	}

	if y := v4.Y(); y != 4.0 {
		t.Errorf("Y() accessor failed after modification. Expected 4.0, got %f", y)
	}
}

func TestScaleXY(t *testing.T) {
	// Test case 1: Scale with positive factors
	v1 := NewVector2D(2.0, 3.0)
	v1.ScaleXY(2.0, 3.0)

	if v1.dx != 4.0 || v1.dy != 9.0 {
		t.Errorf("ScaleXY with positive factors failed. Expected {4, 9}, got {%f, %f}", v1.dx, v1.dy)
	}

	// Test case 2: Scale with negative factors
	v2 := NewVector2D(2.0, 3.0)
	v2.ScaleXY(-1.5, -2.5)

	if v2.dx != -3.0 || v2.dy != -7.5 {
		t.Errorf("ScaleXY with negative factors failed. Expected {-3, -7.5}, got {%f, %f}", v2.dx, v2.dy)
	}

	// Test case 3: Scale with zero factors
	v3 := NewVector2D(2.0, 3.0)
	v3.ScaleXY(0.0, 0.0)

	if v3.dx != 0.0 || v3.dy != 0.0 {
		t.Errorf("ScaleXY with zero factors failed. Expected {0, 0}, got {%f, %f}", v3.dx, v3.dy)
	}

	// Test case 4: Scale with different signs
	v4 := NewVector2D(2.0, 3.0)
	v4.ScaleXY(2.0, -2.0)

	if v4.dx != 4.0 || v4.dy != -6.0 {
		t.Errorf("ScaleXY with different signs failed. Expected {4, -6}, got {%f, %f}", v4.dx, v4.dy)
	}

	// Test case 5: Method chaining
	v5 := NewVector2D(1.0, 1.0)
	result := v5.ScaleXY(2.0, 3.0)

	if result != v5 {
		t.Errorf("Method chaining failed. ScaleXY() did not return the vector itself")
	}

	// Test case 6: Multiple scaling operations
	v6 := NewVector2D(1.0, 1.0)
	v6.ScaleXY(2.0, 3.0).ScaleXY(2.0, 2.0)

	if v6.dx != 4.0 || v6.dy != 6.0 {
		t.Errorf("Multiple scaling operations failed. Expected {4, 6}, got {%f, %f}", v6.dx, v6.dy)
	}
}
