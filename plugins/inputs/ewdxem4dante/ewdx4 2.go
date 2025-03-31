package ewdxem4dante

import (
	_ "embed"
	"fmt"
	"net"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type EWDX4 struct {
	System string `toml:"system"`

	Log telegraf.Logger `toml:"-"`

	writeMU sync.Mutex
	conn    net.Conn
}

func (*EWDX4) SampleConfig() string {
	return sampleConfig
}

func (ewdx *EWDX4) Init() error {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:45", ewdx.System))
	if err != nil {
		return fmt.Errorf("could not create socket for System: %s: %w", ewdx.System, err)
	}

	ewdx.conn = conn
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
