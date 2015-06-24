package weather

import (
	"math"
	"testing"
)

type valuable interface {
	Value() float64
}

func assert(t *testing.T, i int, got, expected valuable) {
	if !nearlyEqual(got.Value(), expected.Value(), 1e-4) {
		t.Errorf("[%d] got %v, expected %v", i, got.Value(), expected.Value())
	}
}

func nearlyEqual(a, b, epsilon float64) bool {
	// Copied and translated straight from
	// http://floating-point-gui.de/errors/comparison/
	diff := math.Abs(a - b)
	const (
		MinNormal float64 = 2.2250738585072014E-308 // IEEE-754 Double Precision
		MaxValue  float64 = math.MaxFloat64
	)
	if a == b {
		// shortcut, handles infinities
		return true

	} else if a == 0 || b == 0 || diff < MinNormal {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return diff < (epsilon * MinNormal)

	} else {
		// Relative error
		absA := math.Abs(a)
		absB := math.Abs(b)
		return diff/math.Min(absA+absB, MaxValue) < epsilon
	}
}

func TestTemperature(t *testing.T) {

	tests := []struct {
		k Kelvin
		c Celsius
		f Fahrenheit
	}{
		{k: Kelvin(245.123), c: Celsius(-28.027), f: Fahrenheit(-18.4486)},
	}

	for i, test := range tests {
		assert(t, i, test.k.Kelvin(), test.k)
		assert(t, i, test.k.Celsius(), test.c)
		assert(t, i, test.k.Fahrenheit(), test.f)

		assert(t, i, test.c.Kelvin(), test.k)
		assert(t, i, test.c.Celsius(), test.c)
		assert(t, i, test.c.Fahrenheit(), test.f)

		assert(t, i, test.f.Kelvin(), test.k)
		assert(t, i, test.f.Celsius(), test.c)
		assert(t, i, test.f.Fahrenheit(), test.f)
	}
}
