package providers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/limoges/weather"
)

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

func (r DarkskyResponse) Temperature() weather.Temperature {
	return weather.Fahrenheit(r.Currently.Temperature)
}

func (d Darksky) Temperature(lat, lon string) (weather.Temperature, error) {
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
