package meinbergos

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
)

func (m *MeinbergOS) gatherAntennaModule(acc telegraf.Accumulator, system string) error {
	resp, err := m.gatherSystem(system, "api/state/modules/clk1")
	if err != nil {
		return err
	}

	var clkModule responseModulesClk1
	if err := json.Unmarshal(resp, &clkModule); err != nil {
		return fmt.Errorf("could not parse data from %s: %w", system, err)
	}

	// measurement antenna
	tags := make(map[string]string)
	tags["source"] = system

	fields := make(map[string]interface{})
	fields["isConnected"] = clkModule.Data.Receiver.Antenna.IsConnected
	fields["satellitesInView"] = clkModule.Data.Receiver.Satellites.TotalInView
	fields["satellitesUsed"] = clkModule.Data.Receiver.Satellites.TotalUsed

	for _, system := range clkModule.Data.Receiver.Satellites.Systems {
		inView := fmt.Sprintf("%sSatellitesInView", system.Type)
		fields[inView] = system.InView
		used := fmt.Sprintf("%sSatellitesUsed", system.Type)
		fields[used] = system.Used
	}

	antennaMetric := metric.New("antenna", tags, fields, time.Now())

	// measurment time
	// TBD

	acc.AddFields(antennaMetric.Name(), antennaMetric.Fields(), antennaMetric.Tags())
	return nil
}

func (m *MeinbergOS) gatherClock(acc telegraf.Accumulator, system string) error {
	resp, err := m.gatherSystem(system, "api/state/references/global")
	if err != nil {
		return err
	}

	var clock responseReferencesGlobal
	if err := json.Unmarshal(resp, &clock); err != nil {
		return fmt.Errorf("could not parse data from %s: %w", system, err)
	}

	tags := make(map[string]string)
	tags["source"] = system

	fields := make(map[string]interface{})
	if value, err := clock.EstimatedTimeQualityNano(); err != nil {
		acc.AddError(fmt.Errorf("error reading estimated Quality from %s: %s", system, err.Error()))
		fields["estimatedQualityNano"] = value
	} else {
		fields["estimatedQualityNano"] = value
	}
	if value, err := clock.HoldoverOffsetNano(); err != nil {
		acc.AddError(fmt.Errorf("error reading holdover offset from %s: %s", system, err.Error()))
		fields["holdoverOffsetNano"] = value
	} else {
		fields["holdoverOffsetNano"] = value
	}
	if value, err := clock.HoldoverTimeSec(); err != nil {
		acc.AddError(fmt.Errorf("error reading holdovertime from %s: %s", system, err.Error()))
		fields["holdoverTime"] = value
	} else {
		fields["holdoverTime"] = value
	}
	fields["clockState"] = clock.Data.ClockState
	fields["holdoverState"] = clock.Data.HoldoverState

	clockMetric := metric.New("clock", tags, fields, time.Now())
	acc.AddFields(clockMetric.Name(), clockMetric.Fields(), clockMetric.Tags())

	return nil
}

func (m *MeinbergOS) gatherSystem(system string, apiEndpoint string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/%s", system, apiEndpoint)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request with url '%s': %w", url, err)
	}

	if err := m.setRequestAuth(request); err != nil {
		return nil, fmt.Errorf("could not set auth header to url '%s': %w", url, err)
	}

	resp, err := m.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request faild to url '%s': %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received status code %d (%s), expected status 200", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body from url '%s' failed: %w", url, err)
	}

	return b, nil
}

func (m *MeinbergOS) setRequestAuth(request *http.Request) error {
	if m.Username == "" || m.Password == "" {
		return fmt.Errorf("no username or password provided")
	}

	request.SetBasicAuth(m.Username, m.Password)

	return nil
}
