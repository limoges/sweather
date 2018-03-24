package providers

import (
	"log"

	"github.com/limoges/weather"
)

type MultiProvider struct {
	Providers []weather.Provider
	Logger    *log.Logger
}

func (r MultiProvider) Temperatures(lat, lon string) ([]weather.Temperature, error) {
	// We'll use a channel to gather all the results and process them after.
	ts := make(chan weather.Temperature, len(r.Providers))
	es := make(chan error, len(r.Providers))

	// Collect the temperatures
	for _, provider := range r.Providers {
		go func(p weather.Provider) {
			t, err := p.Temperature(lat, lon)
			if err != nil {
				log.Printf("%v: %v\n", p, err)
				es <- err
				return
			}
			log.Printf("%v: %v\n", p, t)
			ts <- t
		}(provider)
	}

	// Merge the results
	var values []weather.Temperature
	var errors []error

	for i := 0; i < len(r.Providers); i++ {
		select {
		case t := <-ts:
			values = append(values, t)
		case e := <-es:
			errors = append(errors, e)
		}
	}

	if len(values) == 0 {
		return nil, AggregateError(errors)
	}
	return values, nil
}

func (r MultiProvider) Temperature(lat, lon string) (weather.Temperature, error) {

	temps, err := r.Temperatures(lat, lon)
	if err != nil {
		return nilTemp, err
	}

	values := make([]float64, len(temps))
	for i, t := range temps {
		values[i] = t.Celsius().Value()
	}

	t := weather.Celsius(median(values))
	return t, nil
}
