package ravio

import (
	"fmt"
	"strconv"
	"strings"
)

type stats map[string]string

func (s stats) PTPClockStatus() int {
	value, err := strconv.ParseInt(
		s["STATISTIC_PARAM_PTP_CLOCK_STATUS"],
		10,
		64,
	)
	if err != nil {
		return -1
	}
	return int(value)
}

func (s stats) PTPClockDrift() int64 {
	value, err := strconv.ParseInt(
		s["STATISTIC_PARAM_PTP_CLOCK_DRIFT"],
		10,
		64,
	)
	if err != nil {
		return -1
	}
	return value
}

func (s stats) PTPClockOffset() int64 {
	value, err := strconv.ParseInt(
		s["STATISTIC_PARAM_PTP_CLOCK_OFFSET"],
		10,
		64,
	)
	if err != nil {
		return -1
	}
	return value
}

func (s stats) PTPClockOffsetMSSec() int64 {
	value, err := strconv.ParseInt(
		s["STATISTIC_PARAM_PTP_CLOCK_OFFSET_MS_SEC"],
		10,
		64,
	)
	if err != nil {
		return -1
	}
	return value
}

func (s stats) PTPClockOffsetMSNSec() int64 {
	value, err := strconv.ParseInt(
		s["STATISTIC_PARAM_PTP_CLOCK_OFFSET_MS_NSEC"],
		10,
		64,
	)
	if err != nil {
		return -1
	}
	return value
}

func (s stats) PTPClockOffsetSMSec() int64 {
	value, err := strconv.ParseInt(
		s["STATISTIC_PARAM_PTP_CLOCK_OFFSET_SM_SEC"],
		10,
		64,
	)
	if err != nil {
		return -1
	}
	return value
}

func (s stats) PTPClockOffsetSMNSec() int64 {
	value, err := strconv.ParseInt(
		s["STATISTIC_PARAM_PTP_CLOCK_OFFSET_SM_NSEC"],
		10,
		64,
	)
	if err != nil {
		return -1
	}
	return value
}

func (s stats) PTPNoDelayRespRecv() int64 {
	value, err := strconv.ParseInt(
		s["STATISTIC_PARAM_PTP_NO_PTP_DELAY_RESP_RECV"],
		10,
		64,
	)
	if err != nil {
		return -1
	}
	return value
}

func (s stats) StreamStatusRXPhy1() ([]int, error) {
	values := make([]int, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_AUDIO_CLIENT_STREAM_STATUS_PHY1"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamPresent, err := strconv.ParseInt(before, 10, 32)
		if err != nil {
			return values, err
		}
		values = append(values, int(streamPresent))
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamStatusRXPhy2() ([]int, error) {
	values := make([]int, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_AUDIO_CLIENT_STREAM_STATUS_PHY2"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamPresent, err := strconv.ParseInt(before, 10, 32)
		if err != nil {
			return values, err
		}
		values = append(values, int(streamPresent))
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamStatusActive() ([]bool, error) {
	magicNumberString := s["STATISTIC_PARAM_STATUS3"]
	magicNumber, err := strconv.ParseInt(magicNumberString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse magicNumber string to int")
	}
	values := make([]bool, 0)
	for i := 0; i < 32; i++ {
		status := (((magicNumber) >> (31 - i)) & 0x1) == (1 & 0xFFFFFFFF)
		values = append(values, status)
	}

	return values, nil
}

func (s stats) StreamErrorStatusRXPhy1() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_AUDIO_CLIENT_STREAM_ERROR_STATUS_PHY1"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamErrorStatusRXPhy2() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_AUDIO_CLIENT_STREAM_ERROR_STATUS_PHY2"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamConnectionLostRXPhy1() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_CONNECTION_LOST_PHY1"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamConnectionLostRXPhy2() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_CONNECTION_LOST_PHY2"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamPacketLostRXPhy1() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_PACKET_LOST_PHY1"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamPacketLostRXPhy2() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_PACKET_LOST_PHY2"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamWrongTimestampRXPhy1() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_WRONG_TIMESTAMP_PHY1"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamWrongTimestampRXPhy2() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_WRONG_TIMESTAMP_PHY2"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamTimestampMinRX() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_AUDIO_CLIENT_STREAM_TIMESTAMP_MIN"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}

func (s stats) StreamTimestampMaxRX() ([]int64, error) {
	values := make([]int64, 0)
	var found bool = true
	var before string
	var after string = s["STATISTIC_PARAM_AUDIO_CLIENT_STREAM_TIMESTAMP_MAX"]
	for found {
		before, after, found = strings.Cut(after, "#")
		streamErrors, err := strconv.ParseInt(before, 10, 64)
		if err != nil {
			return values, err
		}
		values = append(values, streamErrors)
	}
	if len(values) != 32 {
		return values, fmt.Errorf("could not parse all values, only got '%d'", len(values))
	}
	return values, nil
}
