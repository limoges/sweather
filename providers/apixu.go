package providers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/limoges/weather"
)

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

func (r ApixuResponse) Temperature() weather.Temperature {
	return weather.Celsius(r.Current.TempC)
}

func (a Apixu) Temperature(lat, lon string) (weather.Temperature, error) {
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
