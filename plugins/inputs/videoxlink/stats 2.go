package videoxlink

import (
	"time"

	"github.com/influxdata/telegraf/metric"
	client "github.com/lukirs95/goxlinkclient"
)

func (xlink *VideoXLink) handleStats(stats client.Stats) {
	now := time.Now()
	systemId := stats.Id()
	uCache, ok := xlink.updateCache[systemId]
	if !ok || uCache.name == "" {
		// we don't add the stats if we can't find a name for tagging
		return
	}
	name := uCache.name

	// SystemStats
	systemStats := stats.SystemStats()
	systemTags := map[string]string{
		"id":   systemId,
		"name": name,
	}
	systemFields := make(map[string]interface{})
	systemFields["sysUpTime"] = systemStats.OSUpTime()
	systemFields["cpuTemp"] = systemStats.CPUTemp()
	systemFields["sysTemp"] = systemStats.SysTemp()

	systemMetric := metric.New(
		"system",
		systemTags,
		systemFields,
		now)
	xlink.Buf.PushBack(systemMetric)

	// PTP Stats
	ptpFields := make(map[string]interface{})
	ptpFields["enabled"] = systemStats.Ptp()
	ptpFields["sync"] = systemStats.PtpSync()
	ptpFields["syncLocal"] = systemStats.PtpSyncLocal()

	ptpMetric := metric.New(
		"ptp",
		systemTags,
		ptpFields,
		now,
	)
	xlink.Buf.PushBack(ptpMetric)

	// Eth Stats
	for _, eth := range stats.EthStats() {
		ethTags := map[string]string{
			"id":        systemId,
			"name":      name,
			"interface": eth.Ident(),
		}
		ethFields := make(map[string]interface{})
		ethFields["rx"] = eth.RX()
		ethFields["tx"] = eth.TX()

		ethMetric := metric.New("interface", ethTags, ethFields, now)
		xlink.Buf.PushBack(ethMetric)
	}

	// Decoder Stats
	for _, decoder := range stats.DecoderStats() {
		decTags := map[string]string{
			"id":      systemId,
			"name":    name,
			"decoder": decoder.Ident(),
		}

		decFields := make(map[string]interface{})
		decFields["statsTime"] = decoder.StatsTime()
		decFields["fromCloud"] = decoder.FromCloud()
		decFields["fromP2P"] = decoder.FromP2P()
		decFields["dropped"] = decoder.Dropped()
		decFields["rtt"] = decoder.RTT()
		decFields["upTime"] = decoder.UpTime()
		decFields["totalDecoded"] = decoder.VideoDTotal()
		decFields["fps"] = decoder.VideoOutFps()
		decFields["rxMbps"] = decoder.RXmbps()
		decFields["txMbps"] = decoder.TXmbps()
		decFields["dMissing"] = decoder.VideoDMissing()
		decFields["dDropped"] = decoder.VideoDDrop()
		decFields["dCorrected"] = decoder.VideoDCorr()
		decFields["rMissing"] = decoder.VideoRMissing()

		decMetric := metric.New("decoder", decTags, decFields, now)
		xlink.Buf.PushBack(decMetric)
	}

	for _, encoder := range stats.EncoderStats() {
		encTags := map[string]string{
			"id":      systemId,
			"name":    name,
			"encoder": encoder.Ident(),
		}

		encFields := make(map[string]interface{})
		encFields["statsTime"] = encoder.StatsTime()
		encFields["upTime"] = encoder.UpTime()
		encFields["fps"] = encoder.VideoInFps()

		encMetric := metric.New("encoder", encTags, encFields, now)
		xlink.Buf.PushBack(encMetric)
	}
}
