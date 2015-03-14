package weather

import (
	"reflect"
	"testing"
)

func assert(t *testing.T, got, expected interface{}) {
	if reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, expected %v")
	}
}

func TestTemperature(t *testing.T) {

	tests := []struct {
		k Kelvin
		c Celsius
		f Fahrenheit
	}{
		{k: Kelvin(245.123), c: Celsius(0.2), f: Fahrenheit(0)},
	}

	for _, test := range tests {
		assert(t, test.k.Kelvin(), test.k)
		assert(t, test.k.Celsius(), test.c)
		assert(t, test.k.Fahrenheit(), test.f)

		assert(t, test.c.Kelvin(), test.k)
		assert(t, test.c.Celsius(), test.c)
		assert(t, test.c.Fahrenheit(), test.f)

		assert(t, test.f.Kelvin(), test.k)
		assert(t, test.f.Celsius(), test.c)
		assert(t, test.f.Fahrenheit(), test.f)
	}
}
