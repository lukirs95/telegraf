package meinbergos

import (
	"errors"
	"testing"
	"time"
)

type testCase struct {
	test string
	want time.Duration
	err  error
}

func TestTimeConvert(t *testing.T) {
	cases := []testCase{
		// whole numbers
		{test: "1 ns", want: time.Nanosecond, err: nil},
		{test: "1 us", want: time.Microsecond, err: nil},
		{test: "1 ms", want: time.Millisecond, err: nil},
		{test: "1 sec", want: time.Second, err: nil},
		// fractions
		{test: "1.4 ns", want: time.Nanosecond, err: nil},
		{test: "1.4 us", want: time.Nanosecond * 1400, err: nil},
		{test: "1.4 ms", want: time.Microsecond * 1400, err: nil},
		{test: "1.4 sec", want: time.Millisecond * 1400, err: nil},
		{test: "1.49 ns", want: time.Nanosecond * 1, err: nil},
		{test: "1.49 us", want: time.Nanosecond * 1490, err: nil},
		{test: "1.49 ms", want: time.Microsecond * 1490, err: nil},
		{test: "1.49 sec", want: time.Millisecond * 1490, err: nil},
		// high number
		{test: "1000 sec", want: time.Millisecond * 1000 * 1000, err: nil},
		// empty
		{test: "", want: time.Nanosecond * 0, err: errEmptyString},
		// unknown unit
		{test: "14 s", want: time.Nanosecond * 0, err: errUnitUnkown},
		// "unknown" string
		{test: "unknown", want: time.Nanosecond * -1, err: nil},
	}

	for _, test := range cases {
		got, err := timeConvert(test.test)
		if !errors.Is(err, test.err) {
			t.Errorf("got error testing '%s', wanted: '%s', got: '%s'", test.test, test.err.Error(), err.Error())
		}
		if test.want != got {
			t.Errorf("wanted: %d, got: %d", test.want, got)
		}
	}
}
