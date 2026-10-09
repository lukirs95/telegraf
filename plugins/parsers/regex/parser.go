package regex

import (
	"bytes"
	"fmt"
	"maps"
	"regexp"
	"strconv"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/parsers"
)

type CaptureGroup struct {
	// index is 1 based
	Index     int    `toml:"index"`
	IsTag     bool   `toml:"is_tag"`
	TagKey    string `toml:"tag_key"`
	FieldName string `toml:"field_name"`
	DataType  string `toml:"data_type"`
}

type Parser struct {
	RegexString   string            `toml:"regex"`
	CaptureGroups []CaptureGroup    `toml:"capture_groups"`
	MetricName    string            `toml:"-"`
	DefaultTags   map[string]string `toml:"-"`
	Regex         *regexp.Regexp
}

func (v *Parser) Init() error {
	for _, g := range v.CaptureGroups {
		switch g.DataType {
		case "", "int", "integer":
			g.DataType = "int"
		case "float", "long":
			g.DataType = "float"
		case "str", "string":
			g.DataType = "string"
		case "base64":
			g.DataType = "base64"
		case "bool", "boolean":
			g.DataType = "bool"
		case "auto_integer", "auto_float":
			// Do nothing both are valid
		default:
			return fmt.Errorf("unknown datatype %q, index %d", g.DataType, g.Index)
		}
		if g.FieldName == "" {
			g.FieldName = "value"
		}
	}

	reg, err := regexp.Compile(v.RegexString)
	if err != nil {
		return fmt.Errorf("invalid regex string: %v", v.RegexString)
	}
	v.Regex = reg

	return nil
}

func (v *Parser) Parse(buf []byte) ([]telegraf.Metric, error) {
	vStr := string(bytes.TrimSpace(bytes.Trim(buf, "\x00")))

	matches := v.Regex.FindStringSubmatch(vStr)

	tags := map[string]string{}
	maps.Copy(tags, v.DefaultTags)
	fields := map[string]interface{}{}

	for i, capture := range matches {
		for _, cg := range v.CaptureGroups {
			if cg.Index != i {
				continue
			}

			if cg.IsTag {
				tags[cg.TagKey] = capture
				continue
			}

			var value interface{}
			var err error
			switch cg.DataType {
			case "int":
				value, err = strconv.Atoi(capture)
			case "float":
				value, err = strconv.ParseFloat(capture, 64)
			case "string":
				value = capture
			case "bool":
				value, err = strconv.ParseBool(capture)
			case "auto_integer":
				value, err = strconv.Atoi(capture)
				if err != nil {
					value = capture
					err = nil
				}
			case "auto_float":
				value, err = strconv.ParseFloat(capture, 64)
				if err != nil {
					value = capture
					err = nil
				}
			}
			if err != nil {
				return nil, err
			}

			fields[cg.FieldName] = value

		}
	}

	m := metric.New(v.MetricName, tags,
		fields, time.Now().UTC())

	return []telegraf.Metric{m}, nil
}

func (v *Parser) ParseLine(line string) (telegraf.Metric, error) {
	metrics, err := v.Parse([]byte(line))

	if err != nil {
		return nil, err
	}

	if len(metrics) < 1 {
		return nil, fmt.Errorf("can not parse the line: %s, for data format: value", line)
	}

	return metrics[0], nil
}

func (v *Parser) SetDefaultTags(tags map[string]string) {
	v.DefaultTags = tags
}

func init() {
	parsers.Add("regex",
		func(defaultMetricName string) telegraf.Parser {
			return &Parser{
				MetricName: defaultMetricName,
			}
		},
	)
}
