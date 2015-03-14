package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
)

type Provider interface {
	Temperature(lat, lon string) (Temperature, error)
}

func median(v []float64) float64 {
	sort.Float64s(v)
	// [0 1 (2) 3 4]
	if len(v)%2 == 1 {
		return v[len(v)/2]
	}
	// interpolate the 2 middle
	return (v[len(v)/2] + v[len(v)/2+1]) / 2
}

type AggregateError []error

func (e AggregateError) Error() string {
	return fmt.Sprint([]error(e))
}

type MultiProvider []Provider

func (r MultiProvider) Temperatures(lat, lon string) ([]Temperature, error) {
	// We'll use a channel to gather all the results and process them after.
	ts := make(chan Temperature, len(r))
	es := make(chan error, len(r))

	// Collect the temperatures
	for _, provider := range r {
		go func(p Provider) {
			t, err := p.Temperature(lat, lon)
			if err != nil {
				es <- err
				return
			}
			ts <- t
		}(provider)
	}

	// Merge the results
	var values []Temperature
	var errors []error

	for i := 0; i < len(r); i++ {
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

func (r MultiProvider) Temperature(lat, lon string) (Temperature, error) {

	temps, err := r.Temperatures(lat, lon)
	if err != nil {
		return nilTemp, err
	}

	values := make([]float64, len(temps))
	for i, t := range temps {
		values[i] = t.Celsius().Value()
	}

	t := Celsius(median(values))
	return t, nil
}

var nilTemp = Celsius(0)

type RequestError error

func NewRequestError(err error) error {
	return errors.New(fmt.Sprintf("request error: %s", err))
}

type DecodeError error

func NewDecodeError(err error) error {
	return errors.New(fmt.Sprintf("decode error: %s", err))
}

type Apixu struct {
	Token string
}

func (a Apixu) String() string {
	return "Apixu"
}

type ApixuResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (a Apixu) Temperature(lat, lon string) (Temperature, error) {
	var r ApixuResponse
	url := fmt.Sprintf("https://api.apixu.com/v1/current.json?key=%s&q=%s,%s",
		a.Token, lat, lon)
	//log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nilTemp, NewRequestError(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nilTemp, NewDecodeError(err)
	}

	t := r.Temperature()
	//log.Printf("OpenWeatherMap: %s,%s: %v\n", lat, lon, t.Celsius())
	return t, nil
}

type OpenWeatherMap struct {
	Token string
}

func (o OpenWeatherMap) String() string {
	return "OpenWeatherMap"
}

type OpenWeatherMapResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func (r OpenWeatherMapResponse) Temperature() Temperature {
	return Kelvin(r.Main.Temp)
}

func (o OpenWeatherMap) Temperature(lat, lon string) (Temperature, error) {
	var r OpenWeatherMapResponse
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&APPID=%s",
		lat, lon, o.Token)
	//log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nilTemp, NewRequestError(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nilTemp, NewDecodeError(err)
	}

	t := r.Temperature()
	//log.Printf("OpenWeatherMap: %s,%s: %v\n", lat, lon, t.Celsius())
	return t, nil
}

type Wunderground struct {
	Token string
}

func (w Wunderground) String() string {
	return "Wunderground"
}

type WundergroundResponse struct {
	CurrentObservation struct {
		TempF float64 `json:"temp_f"`
	} `json:"current_observation"`
}

func (r WundergroundResponse) Temperature() Temperature {
	return Fahrenheit(r.CurrentObservation.TempF)
}

func (w Wunderground) Temperature(lat, lon string) (Temperature, error) {
	var r WundergroundResponse
	url := fmt.Sprintf("https://api.wunderground.com/api/%s/conditions/q/%s,%s.json", w.Token, lat, lon)
	//log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nilTemp, NewRequestError(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nilTemp, NewDecodeError(err)
	}

	t := r.Temperature()
	//log.Printf("Wunderground: %s,%s: %v\n", lat, lon, t.Celsius())
	return t, nil
}

type Darksky struct {
	Token string
}

func (d Darksky) String() string {
	return "Darksky"
}

type DarkskyResponse struct {
	Currently struct {
		Temperature float64 `json:"temperature"`
	} `json:"currently"`
}

func (r DarkskyResponse) Temperature() Temperature {
	return Fahrenheit(r.Currently.Temperature)
}

func (d Darksky) Temperature(lat, lon string) (Temperature, error) {
	var r DarkskyResponse
	url := fmt.Sprintf("https://api.darksky.net/forecast/%s/%s,%s", d.Token, lat, lon)
	//log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		return nilTemp, NewRequestError(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nilTemp, NewDecodeError(err)
	}

	t := r.Temperature()
	// log.Printf("Darksky: %s,%s: %v\n", lat, lon, t.Celsius())
	return t, nil
}
