package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/limoges/weather"
	"github.com/limoges/weather/providers"
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
		lat     string
		lon     string
		token   string
		err     error
		ps      []weather.Provider
		verbose bool
		logger  *log.Logger
	)

	flag.BoolVar(&verbose, "v", false, "print step by step information")
	flag.StringVar(&lat, "lat", "", "latitude")
	flag.StringVar(&lon, "lon", "", "latitude")
	flag.Parse()

	logger = log.New(ioutil.Discard, "weather", log.Lshortfile)
	if verbose {
		logger.SetOutput(os.Stdout)
	}

	if lat == "" {
		if lat, err = strictEnv("LATITUDE"); err != nil {
			return err
		}
	}

	if lon == "" {
		if lon, err = strictEnv("LONGITUDE"); err != nil {
			return err
		}
	}

	if token, err = strictEnv("DARKSKY_TOKEN"); err != nil {
		return err
	} else {
		ps = append(ps, providers.Darksky{Token: token})
	}

	if token, err = strictEnv("WUNDERGROUND_TOKEN"); err != nil {
		return err
	} else {
		ps = append(ps, providers.Wunderground{Token: token})
	}

	if token, err = strictEnv("OPENWEATHERMAP_TOKEN"); err != nil {
		return err
	} else {
		ps = append(ps, providers.OpenWeatherMap{Token: token})
	}

	if token, err = strictEnv("APIXU_TOKEN"); err != nil {
		return err
	} else {
		ps = append(ps, providers.Apixu{Token: token})
	}

	if len(ps) < 1 {
		return errors.New("no providers available")
	}

	redundant := providers.MultiProvider{
		Providers: ps,
		Logger:    logger,
	}
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
