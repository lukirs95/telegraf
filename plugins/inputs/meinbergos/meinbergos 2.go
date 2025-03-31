package meinbergos

import (
	"context"
	_ "embed"
	"net/http"
	"sync"

	"github.com/influxdata/telegraf"
	httpconfig "github.com/influxdata/telegraf/plugins/common/http"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type MeinbergOS struct {
	Systems []string `toml:"systems"`

	Username string `toml:"username"`
	Password string `toml:"password"`

	Log telegraf.Logger `toml:"-"`

	httpconfig.HTTPClientConfig
	client *http.Client
}

func (*MeinbergOS) SampleConfig() string {
	return sampleConfig
}

func (m *MeinbergOS) Init() error {

	// Create the client
	ctx := context.Background()
	m.InsecureSkipVerify = true

	client, err := m.HTTPClientConfig.CreateClient(ctx, m.Log)
	if err != nil {
		return err
	}

	m.client = client
	return nil
}

func (m *MeinbergOS) Gather(acc telegraf.Accumulator) error {
	var wg sync.WaitGroup

	for _, system := range m.Systems {
		wg.Add(3)
		go func(system string) {
			defer wg.Done()
			if err := m.gatherAntennaModule(acc, system); err != nil {
				acc.AddError(err)
			}
		}(system)
		go func(system string) {
			defer wg.Done()
			if err := m.gatherClock(acc, system); err != nil {
				acc.AddError(err)
			}
		}(system)
		go func(system string) {
			defer wg.Done()
		}(system)
	}

	wg.Wait()
	return nil
}

func init() {
	inputs.Add("meinbergos", func() telegraf.Input {
		return &MeinbergOS{}
	})
}
