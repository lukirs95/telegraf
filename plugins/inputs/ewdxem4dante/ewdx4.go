package ewdxem4dante

import (
	_ "embed"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type EWDX4 struct {
	System   string `toml:"system"`
	Password string `toml:"password"`

	State [4]Channel

	Log telegraf.Logger `toml:"-"`
}

func (*EWDX4) SampleConfig() string {
	return sampleConfig
}

func (ewdx *EWDX4) Init() error {
	createRestyClient(ewdx.Password)

	ewdx.State = [4]Channel{
		{Id: 0},
		{Id: 1},
		{Id: 2},
		{Id: 3},
	}

	return nil
}

func (ewdx *EWDX4) Gather(acc telegraf.Accumulator) error {
	if err := ewdx.gather(acc); err != nil {
		ewdx.Log.Error(err)
	}
	return nil
}

func init() {
	inputs.Add("ewdxem4dante", func() telegraf.Input {
		return &EWDX4{}
	})
}
