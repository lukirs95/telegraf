package videoxlink

import (
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
)

type updateCache struct {
	name    string
	eth     map[string]*ethUpdateCache
	decoder map[string]*decoderUpdateCache
	encoder map[string]*encoderUpdateCache
}

func (c *updateCache) Metric(id string) []telegraf.Metric {
	metrics := make([]telegraf.Metric, 0)
	for ethId, eth := range c.eth {
		metrics = append(metrics, eth.Metric(c.name, id, ethId))
	}
	for decId, dec := range c.decoder {
		metrics = append(metrics, dec.Metric(c.name, id, decId))
	}
	for encId, enc := range c.encoder {
		metrics = append(metrics, enc.Metric(c.name, id, encId))
	}
	return metrics
}

type ethUpdateCache struct {
	adminEnabled    bool
	up              bool
	activeUplink    bool
	primaryUplink   bool
	secondaryUplink bool
}

func (c *ethUpdateCache) Metric(name string, id string, eth string) telegraf.Metric {
	ethTags := map[string]string{
		"name":      name,
		"id":        id,
		"interface": eth,
	}

	ethFields := make(map[string]interface{})
	ethFields["adminEnabled"] = c.adminEnabled
	ethFields["up"] = c.up
	ethFields["activeUplink"] = c.activeUplink
	ethFields["primaryUplink"] = c.primaryUplink
	ethFields["secondaryUplink"] = c.secondaryUplink

	return metric.New("interface", ethTags, ethFields, time.Now())
}

type decoderUpdateCache struct {
	conected       bool
	running        bool
	enabled        bool
	videoEnabled   bool
	audioEnabled   bool
	hasVideoSignal bool
	hasAudioSignal bool
}

func (c *decoderUpdateCache) Metric(name string, id string, decoder string) telegraf.Metric {
	decoderTags := map[string]string{
		"name":    name,
		"id":      id,
		"decoder": decoder,
	}

	decoderFields := make(map[string]interface{})

	decoderFields["connected"] = c.conected
	decoderFields["running"] = c.running
	decoderFields["enabled"] = c.enabled
	decoderFields["videoEnabled"] = c.videoEnabled
	decoderFields["audioEnabled"] = c.audioEnabled
	decoderFields["hasVideoSignal"] = c.hasVideoSignal
	decoderFields["hasAudioSignal"] = c.hasAudioSignal

	return metric.New(
		"decoder",
		decoderTags, decoderFields, time.Now())
}

type encoderUpdateCache struct {
	connected      bool
	running        bool
	enabled        bool
	videoEnabled   bool
	audioEnabled   bool
	hasVideoSignal bool
	hasAudioSignal bool
}

func (c *encoderUpdateCache) Metric(name string, id string, encoder string) telegraf.Metric {
	encoderTags := map[string]string{
		"name":    name,
		"id":      id,
		"encoder": encoder,
	}

	encoderFields := make(map[string]interface{})

	encoderFields["connected"] = c.connected
	encoderFields["running"] = c.running
	encoderFields["enabled"] = c.enabled
	encoderFields["videoEnabled"] = c.videoEnabled
	encoderFields["audioEnabled"] = c.audioEnabled
	encoderFields["hasVideoSignal"] = c.hasVideoSignal
	encoderFields["hasAudioSignal"] = c.hasAudioSignal

	return metric.New(
		"encoder",
		encoderTags, encoderFields, time.Now())
}
