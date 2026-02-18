package prtgnm

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/influxdata/telegraf"
)

type Datapoint struct {
	Name         string      `json:"name"`
	NameRaw      string      `json:"name_raw"`
	LastValue    string      `json:"lastvalue"`
	LastValueRaw interface{} `json:"lastvalue_raw"`
}

func (d *Datapoint) LastValueConverted() float64 {
	if i, ok := d.LastValueRaw.(float64); ok {
		return math.Pow(10, i/10) / 1000
	}
	return 0
}

type Response struct {
	Channels []Datapoint `json:"channels"`
}

var fiberRegex = regexp.MustCompile(`^subslot \d/\d+ transceiver \d+ ([RT]x) ([WL]AN) ?(\d?)(.*)$`)
var rateRegex = regexp.MustCompile(`^(Discards|Errors|Traffic) (In|Out)$`)

func (r *prtgnm) gather(acc telegraf.Accumulator, s sensor) error {
	urlBuilder := url.URL{}
	urlBuilder.Scheme = r.ApiScheme
	urlBuilder.Host = r.ApiHost
	urlBuilder.Path = "/api/table.json"
	query := urlBuilder.Query()
	query.Set("content", "channels")
	query.Set("sortby", "name")
	query.Set("columns", "name,lastvalue")
	query.Set("username", r.Username)
	query.Set("password", r.Password)
	query.Set("id", strconv.Itoa(s.Id))

	urlBuilder.RawQuery = query.Encode()

	request, err := http.NewRequest("GET", urlBuilder.String(), nil)
	if err != nil {
		return fmt.Errorf("could not create request with url '%s': %w", urlBuilder.String(), err)
	}

	request.Header.Add("Accept", "application/json")

	resp, err := r.client.Do(request)
	if err != nil {
		return fmt.Errorf("could not make request with url '%s': %w", urlBuilder.String(), err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("received status code %d (%s), expected status 200", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var datapoints Response
	if err := json.NewDecoder(resp.Body).Decode(&datapoints); err != nil {
		return fmt.Errorf("could not decode response body: %w", err)
	}

	for _, dp := range datapoints.Channels {

		tags := map[string]string{
			"location": s.Location,
		}

		fields := map[string]interface{}{}

		matchPower := fiberRegex.FindStringSubmatch(dp.Name)
		if matchPower != nil {
			if len(matchPower) == 4 {
				tags["port"] = fmt.Sprintf("%s%s", matchPower[2], matchPower[3])
			} else if len(matchPower) == 5 {
				tags["port"] = fmt.Sprintf("%s%s%s", matchPower[2], matchPower[3], matchPower[4])
			}

			tags["direction"] = matchPower[1]

			fields["power"] = dp.LastValueConverted()
		}

		matchRate := rateRegex.FindStringSubmatch(dp.Name)
		if matchRate != nil {
			if matchRate[2] == "In" {
				tags["direction"] = "Rx"
			} else {
				tags["direction"] = "Tx"
			}

			toLower := strings.ToLower(matchRate[1])

			if toLower == "traffic" {
				i, ok := dp.LastValueRaw.(float64)
				if ok {
					fields[toLower] = i * 8
				}
			} else {
				fields[toLower] = dp.LastValueRaw
			}
		}

		if len(fields) > 0 {
			acc.AddFields("prtg_nm", fields, tags)
		}
	}

	return nil
}
