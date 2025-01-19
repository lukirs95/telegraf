package videoxlink

import (
	"time"

	"github.com/influxdata/telegraf/metric"
	client "github.com/lukirs95/goxlinkclient"
)

func (xlink *VideoXLink) handleUpdate(update client.XLink) {
	systemId := update.Ident()
	if name, OK := update.GetName(); OK {
		xlink.sysNameMap[systemId] = name
	}
	name := xlink.sysNameMap[systemId]

	for _, eth := range update.GetInterfaces() {
		ethTags := map[string]string{
			"name":      name,
			"id":        systemId,
			"interface": eth.Ident(),
		}

		ethFields := make(map[string]interface{})

		if stat, OK := eth.IsEnabled(); OK {
			ethFields["adminEnabled"] = stat
		}
		if stat, OK := eth.IsLinkUp(); OK {
			ethFields["up"] = stat
		}
		if stat, OK := eth.IsActive(); OK {
			ethFields["activeUplink"] = stat
		}
		if stat, OK := eth.IsDefaultUplink(); OK {
			ethFields["primaryUplink"] = stat
		}
		if stat, OK := eth.IsBackupUplink(); OK {
			ethFields["secondaryUplink"] = stat
		}

		xlink.Buf.PushBack(metric.New("interface", ethTags, ethFields, time.Now()))
	}

	for _, decoder := range update.GetDecoders() {
		decoderTags := map[string]string{
			"name":    name,
			"id":      systemId,
			"decoder": decoder.Ident(),
		}

		decoderFields := make(map[string]interface{})

		if stat, OK := decoder.IsConnected(); OK {
			decoderFields["connected"] = stat
		}
		if stat, OK := decoder.IsRunning(); OK {
			decoderFields["running"] = stat
		}
		if stat, OK := decoder.IsEnabled(); OK {
			decoderFields["enabled"] = stat
		}
		if stat, OK := decoder.IsVideoEnabled(); OK {
			decoderFields["videoEnabled"] = stat
		}
		if stat, OK := decoder.IsAudioEnabled(); OK {
			decoderFields["audioEnabled"] = stat
		}

		if stat, OK := decoder.HasVideoSignal(); OK {
			decoderFields["hasVideoSignal"] = stat
		}

		if stat, OK := decoder.HasAudioSignal(); OK {
			decoderFields["hasAudioSignal"] = stat
		}

		decoderStat := metric.New(
			"decoder",
			decoderTags, decoderFields, time.Now())
		xlink.Buf.PushBack(decoderStat)
	}

	for _, encoder := range update.GetEncoders() {
		encoderTags := map[string]string{
			"name":    name,
			"id":      systemId,
			"encoder": encoder.Ident(),
		}

		encoderFields := make(map[string]interface{})

		if stat, OK := encoder.IsConnected(); OK {
			encoderFields["connected"] = stat
		}
		if stat, OK := encoder.IsRunning(); OK {
			encoderFields["running"] = stat
		}
		if stat, OK := encoder.IsEnabled(); OK {
			encoderFields["enabled"] = stat
		}
		if stat, OK := encoder.IsVideoEnabled(); OK {
			encoderFields["videoEnabled"] = stat
		}
		if stat, OK := encoder.IsAudioEnabled(); OK {
			encoderFields["audioEnabled"] = stat
		}

		if stat, OK := encoder.HasVideoSignal(); OK {
			encoderFields["hasVideoSignal"] = stat
		}

		if stat, OK := encoder.HasAudioSignal(); OK {
			encoderFields["hasAudioSignal"] = stat
		}

		encoderStat := metric.New(
			"encoder",
			encoderTags, encoderFields, time.Now())
		xlink.Buf.PushBack(encoderStat)
	}
}
