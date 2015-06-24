package providers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/limoges/weather"
)

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

func (r WundergroundResponse) Temperature() weather.Temperature {
	return weather.Fahrenheit(r.CurrentObservation.TempF)
}

func (w Wunderground) Temperature(lat, lon string) (weather.Temperature, error) {
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
