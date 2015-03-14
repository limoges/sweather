package weather

import "fmt"

type Temperature interface {
	Kelvin() Kelvin
	Celsius() Celsius
	Fahrenheit() Fahrenheit
	Value() float64
}

type Kelvin float64

func (k Kelvin) Kelvin() Kelvin {
	return k
}

func (k Kelvin) Celsius() Celsius {
	return Celsius(float64(k) - 273.15)
}

func (k Kelvin) Fahrenheit() Fahrenheit {
	return Fahrenheit(float64(k)*9/5 - 459.67)
}

func (k Kelvin) Value() float64 {
	return float64(k)
}

func (k Kelvin) String() string {
	return fmt.Sprintf("%.2f K", k)
}

type Celsius float64

func (c Celsius) Kelvin() Kelvin {
	return Kelvin(float64(c) + 273.15)
}

func (c Celsius) Celsius() Celsius {
	return c
}

func (c Celsius) Fahrenheit() Fahrenheit {
	return Fahrenheit(float64(c)*9/5 + 32)
}

func (c Celsius) String() string {
	return fmt.Sprintf("%.0f °C", c)
}

func (c Celsius) Value() float64 {
	return float64(c)
}

type Fahrenheit float64

func (f Fahrenheit) Kelvin() Kelvin {
	return Kelvin((float64(f) + 459.67) * 5 / 9)
}

func (f Fahrenheit) Celsius() Celsius {
	return Celsius((float64(f) - 32) * 5 / 9)
}

func (f Fahrenheit) Fahrenheit() Fahrenheit {
	return f
}

func (f Fahrenheit) Value() float64 {
	return float64(f)
}

func (f Fahrenheit) String() string {
	return fmt.Sprintf("%.0f °F", f)
}
