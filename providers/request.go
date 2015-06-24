package providers

import (
	"errors"
	"fmt"

	"github.com/limoges/weather"
)

var nilTemp = weather.Kelvin(0) // absolute zero, obviously ;)

type AggregateError []error

func (e AggregateError) Error() string {
	return fmt.Sprint([]error(e))
}

type RequestError error

func NewRequestError(err error) error {
	return errors.New(fmt.Sprintf("request error: %s", err))
}

type DecodeError error

func NewDecodeError(err error) error {
	return errors.New(fmt.Sprintf("decode error: %s", err))
}
