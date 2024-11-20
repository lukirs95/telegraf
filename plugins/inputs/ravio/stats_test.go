package ravio

import (
	"os"
	"testing"
)

func TestStats(t *testing.T) {
	readFile, err := os.Open("test.htm")
	if err != nil {
		t.Error(err)
	}
	defer readFile.Close()

	newStats, err := parse(readFile)
	if err != nil {
		t.Error(err)
	}

	streamStatus, err := newStats.StreamStatusRXPhy1()
	if err != nil {
		t.Error(err)
	}

	if streamStatus[0] != 1 {
		t.Errorf("wanted Stream 1 to be 1")
	}
	if streamStatus[1] != 0 {
		t.Errorf("wanted Stream 2 to be 0")
	}
	if streamStatus[31] != 2 {
		t.Errorf("wanted Stream 32 to be 2")
	}

	streamErrors, err := newStats.StreamErrorStatusRXPhy1()
	if err != nil {
		t.Error(err)
	}

	if streamErrors[0] != 2097152 {
		t.Errorf("wanted Stream 1 to have 2097152")
	}
	if streamErrors[1] != 0 {
		t.Errorf("wanted Stream 2 to hvae 0")
	}
	if streamErrors[31] != 93 {
		t.Errorf("wanted Stream 32 to have 93")
	}
}
