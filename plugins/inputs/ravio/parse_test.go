package ravio

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	readFile, err := os.Open("test.htm")
	if err != nil {
		t.Error(err)
	}
	defer readFile.Close()

	newStats, err := parse(readFile)
	if err != nil {
		t.Errorf("bla")
	}

	if newStats["STATISTIC_PARAM_PTP_CLOCK_OFFSET"] != "4294967150" {
		t.Errorf("could not read 'STATISTIC_PARAM_PTP_CLOCK_OFFSET', %s", newStats["STATISTIC_PARAM_PTP_CLOCK_OFFSET"])
	}
}
