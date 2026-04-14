package ewdxem4dante

type Channel struct {
	Id             int
	Receiver       Receiver
	Transmitter    Transmitter
	Battery        Battery
	SignalStrength float32
	SignalQuality  float32
	Diversity      float32
}

type SomeValue struct {
	Value float32 `json:"value"`
}

// {{baseUrl}}/api/channel/:channelId/signalStrengthIndicator
type SignalStrength SomeValue

// {{baseUrl}}/api/channel/:channelId/signalQualityIndicator
type SignalQuality SomeValue

// {{baseUrl}}/api/channel/:channelId/diversityIndicator
type Diversity SomeValue

// {{baseUrl}}/api/channel/:channelId
type Receiver struct {
	Name         string   `json:"name"`
	Mute         bool     `json:"mute"`
	Mates        []string `json:"mates"`
	Gain         int      `json:"gain"`
	Audio        []string `json:"audio"`
	OutputLevel  int      `json:"outputLevel"`
	EventCounter int      `json:"eventCounter"`
}

// {{baseUrl}}/api/transmitters/:transmitterId
type Transmitter struct {
	Type           string `json:"type"`
	Capsule        string `json:"capsule"`
	Version        string `json:"version"`
	Mute           bool   `json:"mute"`
	Identification bool   `json:"identification"`
	LowCut         string `json:"lowcut"`
	Name           string `json:"name"`
	Trim           int    `json:"trim"`
	Led            bool   `json:"led"`
	Lock           bool   `json:"lock"`
}

// {{baseUrl}}/api/transmitters/:transmitterId/battery
type Battery struct {
	Gauge float32 `json:"gauge"`
	Type  string  `json:"type"`
}
