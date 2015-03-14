package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/limoges/weather"
)

func main() {
	err := mainWithError()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mainWithError() error {

	var (
		lat       string
		lon       string
		token     string
		err       error
		providers []weather.Provider
	)

	if lat, err = strictEnv("LATITUDE"); err != nil {
		return err
	}

	if lon, err = strictEnv("LONGITUDE"); err != nil {
		return err
	}

	if token, err = strictEnv("DARKSKY_TOKEN"); err != nil {
		return err
	} else {
		providers = append(providers, weather.Darksky{Token: token})
	}

	if token, err = strictEnv("WUNDERGROUND_TOKEN"); err != nil {
		return err
	} else {
		providers = append(providers, weather.Wunderground{Token: token})
	}

	if token, err = strictEnv("OPENWEATHERMAP_TOKEN"); err != nil {
		return err
	} else {
		providers = append(providers, weather.OpenWeatherMap{Token: token})
	}

	if token, err = strictEnv("APIXU_TOKEN"); err != nil {
		return err
	} else {
		providers = append(providers, weather.Apixu{Token: token})
	}

	if len(providers) < 1 {
		return errors.New("no provider available")
	}

	redundant := weather.MultiProvider(providers)
	t, err := redundant.Temperature(lat, lon)
	if err != nil {
		return err
	}

	fmt.Println(t.Celsius())
	return nil
}

func strictEnv(identifier string) (string, error) {
	s := os.Getenv(identifier)
	if s == "" {
		return "", errors.New(fmt.Sprintf("%s not set, check `env`", identifier))
	}
	return s, nil
}
