package ravio

import (
	"bufio"
	"bytes"
	"io"
)

func parse(raw io.Reader) (stats, error) {
	newStats := make(stats)
	statsScanner := bufio.NewScanner(raw)
	statsScanner.Split(bufio.ScanLines)

	sep := []byte(":")
	for statsScanner.Scan() {
		rawString := statsScanner.Bytes()
		key, value, found := bytes.Cut(rawString, sep)
		if found {
			newStats[string(key)] = string(value[1:])
		}
	}

	return newStats, nil
}
