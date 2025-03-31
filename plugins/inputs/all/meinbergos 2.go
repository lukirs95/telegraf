//go:build !custom || inputs || inputs.meinbergos

package all

import _ "github.com/influxdata/telegraf/plugins/inputs/meinbergos" // register plugin
