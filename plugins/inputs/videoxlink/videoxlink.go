package videoxlink

import (
	_ "embed"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/inputs/videoxlink/helper"
	xlinkclient "github.com/lukirs95/goxlinkclient"
)

//go:embed sample.conf
var sampleConfig string

type VideoXLink struct {
	Systems     []string `toml:"systems"`
	updateCache map[string]*updateCache
	statsCache  map[string]*xlinkclient.Stats

	Password string `toml:"password"`

	Log telegraf.Logger `toml:"-"`

	wgMu sync.Mutex
	wg   sync.WaitGroup
	stop func()

	Buf *helper.RingBuffer[telegraf.Metric]
}

func (*VideoXLink) SampleConfig() string {
	return sampleConfig
}

func (x *VideoXLink) Init() error {
	return nil
}

func (m *VideoXLink) Gather(acc telegraf.Accumulator) error {
	// get at least last update state
	for id, system := range m.updateCache {
		for _, metric := range system.Metric(id) {
			m.Buf.PushBack(metric)
		}
	}

	// read out and empty buffer
	tmpBuf := make([]telegraf.Metric, 0)
	for fill := m.Buf.Size(); fill != 0; fill-- {
		if metric, err := m.Buf.PopFront(); err != nil {
			m.Log.Error("read empty buffer")
		} else {
			tmpBuf = append(tmpBuf, metric)
		}
	}
	for _, metric := range tmpBuf {
		acc.AddFields(metric.Name(), metric.Fields(), metric.Tags(), metric.Time())
	}
	return nil
}

func init() {
	inputs.Add("videoxlink", func() telegraf.Input {
		return &VideoXLink{
			Buf:         helper.NewRingBuffer[telegraf.Metric](),
			updateCache: make(map[string]*updateCache),
			statsCache:  make(map[string]*xlinkclient.Stats),
		}
	})
}
