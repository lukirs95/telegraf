package prtgnm

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

type sensor struct {
	Location string `toml:"location"`
	Id       int    `toml:"id"`
}

type prtgnm struct {
	ApiScheme string   `toml:"apiScheme"`
	ApiHost   string   `toml:"apiHost"`
	Username  string   `toml:"username"`
	Password  string   `toml:"password"`
	Sensors   []sensor `toml:"sensors"`

	Log telegraf.Logger `toml:"-"`

	httpconfig.HTTPClientConfig
	client *http.Client
}

func (*prtgnm) SampleConfig() string {
	return sampleConfig
}

func (r *prtgnm) Init() error {

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

func (r *prtgnm) Gather(acc telegraf.Accumulator) error {
	var wg sync.WaitGroup

	for _, s := range r.Sensors {
		wg.Add(1)
		go func(sensor sensor) {
			defer wg.Done()

			if err := r.gather(acc, sensor); err != nil {
				acc.AddError(fmt.Errorf("could not get stats from %s: %w", sensor.Location, err))
			}
		}(s)
	}

	wg.Wait()
	return nil
}

func init() {
	inputs.Add("prtgnm", func() telegraf.Input {
		return &prtgnm{}
	})
}
