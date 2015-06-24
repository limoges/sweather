package providers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/limoges/weather"
)

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

func (r OpenWeatherMapResponse) Temperature() weather.Temperature {
	return weather.Kelvin(r.Main.Temp)
}

func (o OpenWeatherMap) Temperature(lat, lon string) (weather.Temperature, error) {
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
