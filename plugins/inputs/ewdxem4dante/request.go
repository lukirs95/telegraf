package ewdxem4dante

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/influxdata/telegraf"
)

var errNoAnswer = errors.New("no answer from server in given time window")

func (ewdx *EWDX4) gather(acc telegraf.Accumulator) error {
	data := EWDXEM4DANTE{}

	resp, err := ewdx.request(data)
	if err != nil {
		if errors.Is(err, errNoAnswer) {
			ewdx.Log.Error(err)
			return nil
		}
		return err
	}

	if err := json.Unmarshal(resp, &data); err != nil {
		return fmt.Errorf("could not decode response: %w", err)
	}

	// all about radio receiving
	rx1Tags := make(map[string]string)
	rx1Tags["source"] = ewdx.System
	rx1Tags["receiver"] = "Rx1"
	rx1RadioFields := make(map[string]interface{})
	rx1RadioFields["divi"] = *data.M.RX1.Divi
	rx1RadioFields["rsqi"] = *data.M.RX1.RSQI
	rx1RadioFields["rssi"] = *data.M.RX1.RSSI
	rx1RadioFields["frequency"] = *data.Rx1.Frequency
	acc.AddFields("radio", rx1RadioFields, rx1Tags)

	rx2Tags := make(map[string]string)
	rx2Tags["source"] = ewdx.System
	rx2Tags["receiver"] = "Rx2"
	rx2RadioFields := make(map[string]interface{})
	rx2RadioFields["divi"] = *data.M.RX2.Divi
	rx2RadioFields["rsqi"] = *data.M.RX2.RSQI
	rx2RadioFields["rssi"] = *data.M.RX2.RSSI
	rx2RadioFields["frequency"] = *data.Rx2.Frequency
	acc.AddFields("radio", rx2RadioFields, rx2Tags)

	rx3Tags := make(map[string]string)
	rx3Tags["source"] = ewdx.System
	rx3Tags["receiver"] = "Rx3"
	rx3RadioFields := make(map[string]interface{})
	rx3RadioFields["divi"] = *data.M.RX3.Divi
	rx3RadioFields["rsqi"] = *data.M.RX3.RSQI
	rx3RadioFields["rssi"] = *data.M.RX3.RSSI
	rx3RadioFields["frequency"] = *data.Rx3.Frequency
	acc.AddFields("radio", rx3RadioFields, rx3Tags)

	rx4Tags := make(map[string]string)
	rx4Tags["source"] = ewdx.System
	rx4Tags["receiver"] = "Rx4"
	rx4RadioFields := make(map[string]interface{})
	rx4RadioFields["divi"] = *data.M.RX3.Divi
	rx4RadioFields["rsqi"] = *data.M.RX3.RSQI
	rx4RadioFields["rssi"] = *data.M.RX3.RSSI
	rx4RadioFields["frequency"] = *data.Rx3.Frequency
	acc.AddFields("radio", rx4RadioFields, rx4Tags)

	rx1ReceiverFields := make(map[string]interface{})
	rx1ReceiverFields["gain"] = *data.Rx1.Gain
	rx1ReceiverFields["mute"] = *data.Rx1.Mute
	rx1ReceiverFields["name"] = *data.Rx1.Name
	acc.AddFields("receiver", rx1ReceiverFields, rx1Tags)

	rx2ReceiverFields := make(map[string]interface{})
	rx2ReceiverFields["gain"] = *data.Rx2.Gain
	rx2ReceiverFields["mute"] = *data.Rx2.Mute
	rx2ReceiverFields["name"] = *data.Rx2.Name
	acc.AddFields("receiver", rx2ReceiverFields, rx2Tags)

	rx3ReceiverFields := make(map[string]interface{})
	rx3ReceiverFields["gain"] = *data.Rx3.Gain
	rx3ReceiverFields["mute"] = *data.Rx3.Mute
	rx3ReceiverFields["name"] = *data.Rx3.Name
	acc.AddFields("receiver", rx3ReceiverFields, rx3Tags)

	rx4ReceiverFields := make(map[string]interface{})
	rx4ReceiverFields["gain"] = *data.Rx4.Gain
	rx4ReceiverFields["mute"] = *data.Rx4.Mute
	rx4ReceiverFields["name"] = *data.Rx4.Name
	acc.AddFields("receiver", rx4ReceiverFields, rx4Tags)

	if !slices.Contains(*data.Rx1.Warnings, "NoLink") {
		tx := MateTX1{}
		resp, err := ewdx.request(tx)
		if err != nil {
			if errors.Is(err, errNoAnswer) {
				return nil
			}
			return err
		}

		tx.Mates.Mate.Mute = new(bool)
		tx.Mates.Mate.Type = new(string)
		tx.Mates.Mate.Trim = new(int)
		tx.Mates.Mate.Battery.Gauge = new(int)
		tx.Mates.Mate.Battery.Type = new(string)
		*tx.Mates.Mate.Mute = true
		*tx.Mates.Mate.Type = ""
		*tx.Mates.Mate.Trim = 0
		*tx.Mates.Mate.Battery.Gauge = 0
		*tx.Mates.Mate.Battery.Type = "no data"

		if err := json.Unmarshal(resp, &tx); err != nil {
			return fmt.Errorf("could not decode tx1 response: %w", err)
		}

		transmitterFields := make(map[string]interface{})
		transmitterFields["mute"] = *tx.Mates.Mate.Mute
		transmitterFields["trim"] = *tx.Mates.Mate.Trim
		transmitterFields["type"] = *tx.Mates.Mate.Type
		acc.AddFields("transmitter", transmitterFields, rx1Tags)

		batteryFields := make(map[string]interface{})
		batteryFields["gauge"] = *tx.Mates.Mate.Battery.Gauge
		batteryFields["type"] = *tx.Mates.Mate.Battery.Type
		acc.AddFields("battery", batteryFields, rx1Tags)
	}

	if !slices.Contains(*data.Rx2.Warnings, "NoLink") {
		tx := MateTX2{}
		resp, err := ewdx.request(tx)
		if err != nil {
			if errors.Is(err, errNoAnswer) {
				return nil
			}
			return err
		}

		if err := json.Unmarshal(resp, &tx); err != nil {
			return fmt.Errorf("could not decode tx2 response: %w", err)
		}

		tx.Mates.Mate.Mute = new(bool)
		tx.Mates.Mate.Type = new(string)
		tx.Mates.Mate.Trim = new(int)
		tx.Mates.Mate.Battery.Gauge = new(int)
		tx.Mates.Mate.Battery.Type = new(string)
		*tx.Mates.Mate.Mute = true
		*tx.Mates.Mate.Type = ""
		*tx.Mates.Mate.Trim = 0
		*tx.Mates.Mate.Battery.Gauge = 0
		*tx.Mates.Mate.Battery.Type = "no data"

		transmitterFields := make(map[string]interface{})
		transmitterFields["mute"] = *tx.Mates.Mate.Mute
		transmitterFields["trim"] = *tx.Mates.Mate.Trim
		transmitterFields["type"] = *tx.Mates.Mate.Type
		acc.AddFields("transmitter", transmitterFields, rx2Tags)

		batteryFields := make(map[string]interface{})
		batteryFields["gauge"] = *tx.Mates.Mate.Battery.Gauge
		batteryFields["type"] = *tx.Mates.Mate.Battery.Type
		acc.AddFields("battery", batteryFields, rx2Tags)
	}

	if !slices.Contains(*data.Rx3.Warnings, "NoLink") {
		tx := MateTX3{}
		resp, err := ewdx.request(tx)
		if err != nil {
			if errors.Is(err, errNoAnswer) {
				return nil
			}
			return err
		}

		tx.Mates.Mate.Mute = new(bool)
		tx.Mates.Mate.Type = new(string)
		tx.Mates.Mate.Trim = new(int)
		tx.Mates.Mate.Battery.Gauge = new(int)
		tx.Mates.Mate.Battery.Type = new(string)
		*tx.Mates.Mate.Mute = true
		*tx.Mates.Mate.Type = ""
		*tx.Mates.Mate.Trim = 0
		*tx.Mates.Mate.Battery.Gauge = 0
		*tx.Mates.Mate.Battery.Type = "no data"

		if err := json.Unmarshal(resp, &tx); err != nil {
			return fmt.Errorf("could not decode tx3 response: %w", err)
		}

		transmitterFields := make(map[string]interface{})
		transmitterFields["mute"] = *tx.Mates.Mate.Mute
		transmitterFields["trim"] = *tx.Mates.Mate.Trim
		transmitterFields["type"] = *tx.Mates.Mate.Type
		acc.AddFields("transmitter", transmitterFields, rx3Tags)

		batteryFields := make(map[string]interface{})
		batteryFields["gauge"] = *tx.Mates.Mate.Battery.Gauge
		batteryFields["type"] = *tx.Mates.Mate.Battery.Type
		acc.AddFields("battery", batteryFields, rx3Tags)
	}

	if !slices.Contains(*data.Rx4.Warnings, "NoLink") {
		tx := MateTX4{}
		resp, err := ewdx.request(tx)
		if err != nil {
			if errors.Is(err, errNoAnswer) {
				return nil
			}
			return err
		}

		tx.Mates.Mate.Mute = new(bool)
		tx.Mates.Mate.Type = new(string)
		tx.Mates.Mate.Trim = new(int)
		tx.Mates.Mate.Battery.Gauge = new(int)
		tx.Mates.Mate.Battery.Type = new(string)
		*tx.Mates.Mate.Mute = true
		*tx.Mates.Mate.Type = ""
		*tx.Mates.Mate.Trim = 0
		*tx.Mates.Mate.Battery.Gauge = 0
		*tx.Mates.Mate.Battery.Type = "no data"

		if err := json.Unmarshal(resp, &tx); err != nil {
			return fmt.Errorf("could not decode tx4 response: %w", err)
		}

		transmitterFields := make(map[string]interface{})
		transmitterFields["mute"] = *tx.Mates.Mate.Mute
		transmitterFields["trim"] = *tx.Mates.Mate.Trim
		transmitterFields["type"] = *tx.Mates.Mate.Type
		acc.AddFields("transmitter", transmitterFields, rx4Tags)

		batteryFields := make(map[string]interface{})
		batteryFields["gauge"] = *tx.Mates.Mate.Battery.Gauge
		batteryFields["type"] = *tx.Mates.Mate.Battery.Type
		acc.AddFields("battery", batteryFields, rx4Tags)
	}

	return nil
}

func (ewdx *EWDX4) request(data any) ([]byte, error) {
	ewdx.writeMU.Lock()
	defer ewdx.writeMU.Unlock()
	encoder := json.NewEncoder(ewdx.conn)
	encoder.Encode(data)

	buf := make([]byte, 2048)

	ewdx.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := ewdx.conn.Read(buf)
	if err != nil {
		if errors.Is(err, os.ErrDeadlineExceeded) {
			return nil, errNoAnswer
		}
		return nil, err
	}
	return buf[:n], nil
}
