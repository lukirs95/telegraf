package meinbergos

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	errEmptyString = errors.New("empty string")
	errUnitUnkown  = errors.New("unit unknown")
)

func timeConvert(raw string) (time.Duration, error) {
	cvm := make(map[string]time.Duration)
	cvm["ns"] = time.Nanosecond
	cvm["us"] = time.Microsecond
	cvm["ms"] = time.Millisecond
	cvm["sec"] = time.Second

	if raw == "" {
		return time.Nanosecond * 0, errEmptyString
	}

	if raw == "unknown" {
		return time.Nanosecond * -1, nil
	}

	split := strings.Split(raw, " ")
	valueAsFloat, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		return time.Nanosecond * 0, err
	}
	unit := split[1]

	if unit == "us" {
		valueAsFloat = valueAsFloat * 1000
		unit = "ns"
	}

	if unit == "ms" {
		valueAsFloat = valueAsFloat * 1000
		unit = "us"
	}

	if unit == "sec" {
		valueAsFloat = valueAsFloat * 1000
		unit = "ms"
	}

	un, ok := cvm[unit]
	if !ok {
		return time.Nanosecond * 0, errUnitUnkown
	}

	value := int(math.Round(valueAsFloat))
	return un * time.Duration(value), nil
}
