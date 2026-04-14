package ewdxem4dante

import (
	"errors"
	"fmt"
	"sync"

	"github.com/influxdata/telegraf"
)

func UrlSigStrength(base string, channel int) string {
	return fmt.Sprintf("%s/channel/%d/signalStrengthIndicator", base, channel)
}

func UrlSigQuality(base string, channel int) string {
	return fmt.Sprintf("%s/channel/%d/signalQualityIndicator", base, channel)
}

func UrlDiversity(base string, channel int) string {
	return fmt.Sprintf("%s/channel/%d/diversityIndicator", base, channel)
}

func UrlReceiver(base string, channel int) string {
	return fmt.Sprintf("%s/channel/%d", base, channel)
}

func UrlTransmitter(base string, channel int) string {
	return fmt.Sprintf("%s/transmitters/%d", base, channel)
}

func UrlBattery(base string, channel int) string {
	return fmt.Sprintf("%s/transmitters/%d/battery", base, channel)
}

func GetSigStrength(base string, channel *Channel) error {
	res, err := rest.R().
		SetResult(SignalStrength{}).
		Get(UrlSigStrength(base, channel.Id))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if res.IsError() {
		return fmt.Errorf("request returned error")
	}

	sigStrength := res.Result().(*SignalStrength)

	channel.SignalStrength = sigStrength.Value

	return nil
}

func GetSigQuality(base string, channel *Channel) error {
	res, err := rest.R().
		SetResult(SignalQuality{}).
		Get(UrlSigQuality(base, channel.Id))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if res.IsError() {
		return fmt.Errorf("request returned error")
	}

	sigQuality := res.Result().(*SignalQuality)

	channel.SignalQuality = sigQuality.Value

	return nil
}

func GetDiversity(base string, channel *Channel) error {
	res, err := rest.R().
		SetResult(Diversity{}).
		Get(UrlDiversity(base, channel.Id))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if res.IsError() {
		return fmt.Errorf("request returned error")
	}

	diversity := res.Result().(*Diversity)
	channel.Diversity = diversity.Value

	return nil
}

func GetReceiver(base string, channel *Channel) error {
	res, err := rest.R().
		SetResult(Receiver{}).
		Get(UrlReceiver(base, channel.Id))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if res.IsError() {
		return fmt.Errorf("request returned error")
	}

	receiver := res.Result().(*Receiver)
	channel.Receiver = *receiver

	return nil
}

func GetTransmitter(base string, channel *Channel) error {
	res, err := rest.R().
		SetResult(Transmitter{}).
		Get(UrlTransmitter(base, channel.Id))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if res.IsError() {
		return fmt.Errorf("request returned error")
	}

	transmitter := res.Result().(*Transmitter)
	channel.Transmitter = *transmitter

	return nil
}

func GetBattery(base string, channel *Channel) error {
	res, err := rest.R().
		SetResult(Battery{}).
		Get(UrlBattery(base, channel.Id))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if res.IsError() {
		return fmt.Errorf("request returned error")
	}

	battery := res.Result().(*Battery)
	channel.Battery = *battery

	return nil
}

func (c *Channel) Accumulate(acc telegraf.Accumulator, system string) error {
	tags := make(map[string]string)
	tags["source"] = system
	tags["receiver"] = fmt.Sprintf("Rx%d", c.Id)

	radioFields := make(map[string]any)
	radioFields["divi"] = c.Diversity
	radioFields["rsqi"] = c.SignalQuality
	radioFields["rssi"] = c.SignalStrength
	acc.AddFields("radio", radioFields, tags)

	receiverFields := make(map[string]any)
	receiverFields["gain"] = c.Receiver.Gain
	receiverFields["mute"] = c.Receiver.Mute
	receiverFields["name"] = c.Receiver.Name
	acc.AddFields("receiver", receiverFields, tags)

	transmitterFields := make(map[string]any)
	transmitterFields["mute"] = c.Transmitter.Mute
	transmitterFields["trim"] = c.Transmitter.Trim
	transmitterFields["type"] = c.Transmitter.Type
	acc.AddFields("transmitter", transmitterFields, tags)

	batteryFields := make(map[string]any)
	batteryFields["gauge"] = c.Battery.Gauge
	batteryFields["type"] = c.Battery.Type
	acc.AddFields("battery", batteryFields, tags)

	return nil
}

func (ewdx *EWDX4) gather(acc telegraf.Accumulator) error {
	url := fmt.Sprintf("https://%s/api", ewdx.System)

	var errAcc error
	for i := range ewdx.State {
		c := &ewdx.State[i]
		var (
			wg  sync.WaitGroup
			mu  sync.Mutex
			err error
		)

		collect := func(e error) {
			if e != nil {
				mu.Lock()
				err = errors.Join(err, e)
				mu.Unlock()
			}
		}

		wg.Add(6)
		go func() { defer wg.Done(); collect(GetSigStrength(url, c)) }()
		go func() { defer wg.Done(); collect(GetSigQuality(url, c)) }()
		go func() { defer wg.Done(); collect(GetDiversity(url, c)) }()
		go func() { defer wg.Done(); collect(GetReceiver(url, c)) }()
		go func() { defer wg.Done(); collect(GetTransmitter(url, c)) }()
		go func() { defer wg.Done(); collect(GetBattery(url, c)) }()

		wg.Wait()
		errAcc = errors.Join(errAcc, err)
		errAcc = errors.Join(errAcc, c.Accumulate(acc, ewdx.System))
	}

	return errAcc
}
