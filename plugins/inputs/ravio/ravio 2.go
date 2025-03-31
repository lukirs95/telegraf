package ravio

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"sync"

	"github.com/influxdata/telegraf"
	httpconfig "github.com/influxdata/telegraf/plugins/common/http"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type RavIO struct {
	Systems []string `toml:"systems"`

	Log telegraf.Logger `toml:"-"`

	httpconfig.HTTPClientConfig
	client *http.Client
}

func (*RavIO) SampleConfig() string {
	return sampleConfig
}

func (r *RavIO) Init() error {

	// Create the client
	ctx := context.Background()
	r.InsecureSkipVerify = true

	client, err := r.HTTPClientConfig.CreateClient(ctx, r.Log)
	if err != nil {
		return err
	}

	r.client = client
	return nil
}

func (r *RavIO) Gather(acc telegraf.Accumulator) error {
	var wg sync.WaitGroup

	for _, system := range r.Systems {
		wg.Add(1)
		go func(system string) {
			defer wg.Done()
			if err := r.gather(acc, system); err != nil {
				acc.AddError(fmt.Errorf("could not get stats from %s: %w", system, err))
			}
		}(system)
	}

	wg.Wait()
	return nil
}

func init() {
	inputs.Add("ravio", func() telegraf.Input {
		return &RavIO{}
	})
}
