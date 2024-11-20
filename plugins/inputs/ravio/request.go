package ravio

import (
	"fmt"
	"net/http"

	"github.com/influxdata/telegraf"
)

func (r *RavIO) gather(acc telegraf.Accumulator, system string) error {
	newStats, err := r.gatherStats(system)
	if err != nil {
		return err
	}

	ptpTags := make(map[string]string)
	ptpTags["source"] = system

	ptpFields := make(map[string]interface{})
	ptpFields["clockStatus"] = newStats.PTPClockStatus()
	ptpFields["clockDrift"] = newStats.PTPClockDrift()
	ptpFields["clockOffset"] = newStats.PTPClockOffset()
	ptpFields["noDelayRespReceived"] = newStats.PTPNoDelayRespRecv()
	acc.AddFields("ptp", ptpFields, ptpTags)

	streamStatusRXPhy1, err := newStats.StreamStatusRXPhy1()
	if err != nil {
		return err
	}
	streamStatusRXPhy2, err := newStats.StreamStatusRXPhy2()
	if err != nil {
		return err
	}
	streamConnLostPhy1, err := newStats.StreamConnectionLostRXPhy1()
	if err != nil {
		return err
	}
	streamConnLostPhy2, err := newStats.StreamConnectionLostRXPhy2()
	if err != nil {
		return err
	}
	streamPacketLostPhy1, err := newStats.StreamPacketLostRXPhy1()
	if err != nil {
		return err
	}
	streamPacketLostPhy2, err := newStats.StreamPacketLostRXPhy2()
	if err != nil {
		return err
	}
	streamWrongTimestampPhy1, err := newStats.StreamWrongTimestampRXPhy1()
	if err != nil {
		return err
	}
	streamWrongTimestampPhy2, err := newStats.StreamWrongTimestampRXPhy2()
	if err != nil {
		return err
	}
	streamTimestampMin, err := newStats.StreamTimestampMinRX()
	if err != nil {
		return err
	}
	streamTimestampMax, err := newStats.StreamTimestampMaxRX()
	if err != nil {
		return err
	}

	for i := range streamStatusRXPhy1 {
		phy1Tags := make(map[string]string)
		phy1Tags["source"] = system
		phy1Tags["stream"] = fmt.Sprintf("%d", i+1)
		phy1Tags["direction"] = "Rx"
		phy1Tags["nic"] = "1"

		phy1Stats := make(map[string]interface{})
		phy1Stats["status"] = streamStatusRXPhy1[i]
		phy1Stats["connectionLoss"] = streamConnLostPhy1[i]
		phy1Stats["packetLoss"] = streamPacketLostPhy1[i]
		phy1Stats["wrongTimestamp"] = streamWrongTimestampPhy1[i]
		phy1Stats["timestampMin"] = streamTimestampMin[i]
		phy1Stats["timestampMax"] = streamTimestampMax[i]
		acc.AddFields("stream", phy1Stats, phy1Tags)

		phy2Tags := make(map[string]string)
		phy2Tags["source"] = system
		phy2Tags["stream"] = fmt.Sprintf("%d", i+1)
		phy2Tags["direction"] = "Rx"
		phy2Tags["nic"] = "2"

		phy2Stats := make(map[string]interface{})
		phy2Stats["status"] = streamStatusRXPhy2[i]
		phy2Stats["connectionLoss"] = streamConnLostPhy2[i]
		phy2Stats["packetLoss"] = streamPacketLostPhy2[i]
		phy2Stats["wrongTimestamp"] = streamWrongTimestampPhy2[i]
		phy2Stats["timestampMin"] = streamTimestampMin[i]
		phy2Stats["timestampMax"] = streamTimestampMax[i]
		acc.AddFields("stream", phy2Stats, phy2Tags)
	}

	return nil
}

func (r *RavIO) gatherStats(system string) (stats, error) {
	url := fmt.Sprintf("http://%s/statistic.htm?C=1", system)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request with url '%s': %w", url, err)
	}

	request.Header.Add("accept", "*/*")

	resp, err := r.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request faild to url '%s': %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received status code %d (%s), expected status 200", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	stats, err := parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not parse response: %w", err)
	}

	return stats, nil
}
