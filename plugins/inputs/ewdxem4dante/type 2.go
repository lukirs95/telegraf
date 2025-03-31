package ewdxem4dante

type Name *string
type Gain *int
type Mute *bool
type Trim *int
type Gauge *int
type Type *string
type Divi *int
type RSQI *int
type RSSI *float32
type Frequency *int
type Warnings *[]string

type EWDXEM4DANTE struct {
	Rx1 RX `json:"rx1"`
	Rx2 RX `json:"rx2"`
	Rx3 RX `json:"rx3"`
	Rx4 RX `json:"rx4"`
	M   struct {
		RX1 MRX `json:"rx1"`
		RX2 MRX `json:"rx2"`
		RX3 MRX `json:"rx3"`
		RX4 MRX `json:"rx4"`
	} `json:"m"`
}

type MateTX1 struct {
	Mates struct {
		Mate Mate `json:"tx1"`
	} `json:"mates"`
}

type MateTX2 struct {
	Mates struct {
		Mate Mate `json:"tx2"`
	} `json:"mates"`
}

type MateTX3 struct {
	Mates struct {
		Mate Mate `json:"tx3"`
	} `json:"mates"`
}

type MateTX4 struct {
	Mates struct {
		Mate Mate `json:"tx4"`
	} `json:"mates"`
}

type RX struct {
	Frequency Frequency `json:"frequency"`
	Gain      Gain      `json:"gain"`
	Mute      Mute      `json:"mute"`
	Name      Name      `json:"name"`
	Warnings  Warnings  `json:"warnings"`
}

type MRX struct {
	Divi Divi `json:"divi"`
	RSQI RSQI `json:"rsqi"`
	RSSI RSSI `json:"rssi"`
}

type Mate struct {
	Mute     Mute     `json:"mute"`
	Trim     Trim     `json:"trim"`
	Type     Type     `json:"type"`
	Warnings Warnings `json:"Type"`
	Battery  Battery  `json:"battery"`
}

type Battery struct {
	Gauge Gauge `json:"gauge"`
	Type  Type  `json:"type"`
}
