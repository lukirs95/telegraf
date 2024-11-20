package meinbergos

import "errors"

// http://meinberg/api/state/modules/clk1
type responseModulesClk1 struct {
	Data struct {
		Receiver struct {
			Antenna struct {
				IsConnected bool `json:"isConnected"`
			} `json:"antenna"`
			Satellites struct {
				Mode          int    `json:"mode"`
				OperationMode string `json:"operationMode"`
				Systems       []struct {
					InView int    `json:"inView"`
					Type   string `json:"type"`
					Used   int    `json:"used"`
				} `json:"systems"`
				TotalInView int `json:"totalInView"`
				TotalUsed   int `json:"totalUsed"`
			} `json:"satellites"`
			Time struct {
				GpsWeekNumber        int64  `json:"gpsWeekNumber"`
				GpsWeekSecond        int64  `json:"gpsWeekSecond"`
				IsDaylightSavingTime bool   `json:"isDaylightSavingTime"`
				IsLS59Announced      bool   `json:"isLS59Announced"`
				IsLS61Announced      bool   `json:"isLS61Announced"`
				IsLocalTime          bool   `json:"isLocalTime"`
				OffsetFromUTC        int    `json:"offsetFromUTC"`
				Status               int    `json:"status"`
				Timestamp            string `json:"timestamp"`
			} `json:"time"`
		} `json:"receiver"`
	} `json:"data"`
}

// http://meinberg/api/state/references/global
type responseReferencesGlobal struct {
	Data struct {
		ClockState           string `json:"clockState"`
		EstimatedTimeQuality string `json:"estimatedTimeQuality"`
		HoldoverOffset       string `json:"holdoverOffset"`
		HoldoverState        string `json:"holdoverState"`
		HoldoverTime         string `json:"holdoverTime"`
	} `json:"data"`
}

func (r responseReferencesGlobal) EstimatedTimeQualityNano() (int64, error) {
	d, err := timeConvert(r.Data.EstimatedTimeQuality)
	if err != nil {
		if errors.Is(err, errEmptyString) {
			return int64(d) - 1, nil
		}
		return int64(d) - 1, err
	}
	return d.Nanoseconds(), nil
}

func (r responseReferencesGlobal) HoldoverOffsetNano() (int64, error) {
	d, err := timeConvert(r.Data.HoldoverOffset)
	if err != nil {
		if errors.Is(err, errEmptyString) {
			return int64(d) - 1, nil
		}
		return int64(d) - 1, err
	}
	return d.Nanoseconds(), nil
}

func (r responseReferencesGlobal) HoldoverTimeSec() (int64, error) {
	d, err := timeConvert(r.Data.HoldoverTime)
	if err != nil {
		if errors.Is(err, errEmptyString) {
			return int64(d) - 1, nil
		}
		return int64(d) - 1, err
	}
	return int64(d.Seconds()), nil
}

// http://meinberg/api/state/ptp
type responsePTP struct {
	Data struct {
		Instances []struct {
			Alias          string `json:"alias"`
			DefaultDataset struct {
				ClockAccuracy string `json:"clockAccuracy"`
				ClockClass    int    `json:"clockClass"`
				ClockID       string `json:"clockID"`
				ClockVariance int    `json:"clockVariance"`
				DomainNumber  int    `json:"domainNumber"`
				IsSlaveOnly   bool   `json:"isSlaveOnly"`
				IsTwoStep     bool   `json:"isTwoStep"`
				NumberPorts   int    `json:"numberPorts"`
				Priority1     int    `json:"priority1"`
				Priority2     int    `json:"priority2"`
			} `json:"defaultDataset"`
			IsRunning      bool `json:"isRunning"`
			PacketCounters struct {
				Rx struct {
					Announce                       int64   `json:"announce"`
					AnnouncePerSecond              float32 `json:"announcePerSecond"`
					DelayReq                       int64   `json:"delayReq"`
					DelayReqPerSecond              float32 `json:"delayReqPerSecond"`
					DelayResp                      int64   `json:"delayResp"`
					DelayRespPerSecond             float32 `json:"delayRespPerSecond"`
					FollowUp                       int64   `json:"followUp"`
					FollowUpPerSecond              float32 `json:"followUpPerSecond"`
					Management                     int64   `json:"management"`
					ManagementErrors               int64   `json:"managementErrors"`
					ManagementPerSecond            float32 `json:"managementPerSecond"`
					PeerDelayReq                   int64   `json:"peerDelayReq"`
					PeerDelayReqPerSecond          float32 `json:"peerDelayReqPerSecond"`
					PeerDelayResp                  int64   `json:"peerDelayResp"`
					PeerDelayRespFollowUp          int64   `json:"peerDelayFollowUp"`
					PeerDelayRespFollowUpPerSecond float32 `json:"peerDelayRespFollowUpPerSecond"`
					PeerDelayRespPerSecond         float32 `json:"peerDelayRespPerSecond"`
					Signalling                     int64   `json:"signalling"`
					SignallingPerSecond            float32 `json:"signallingPerSecond"`
					Sync                           int64   `json:"sync"`
					SyncPerSecond                  float32 `json:"syncPerSecond"`
					Total                          int64   `json:"total"`
					TotalPerSecond                 float32 `json:"totalPerSecond"`
				} `json:"rx"`
				Tx struct {
					Announce                       int64   `json:"announce"`
					AnnouncePerSecond              float32 `json:"announcePerSecond"`
					DelayReq                       int64   `json:"delayReq"`
					DelayReqPerSecond              float32 `json:"delayReqPerSecond"`
					DelayResp                      int64   `json:"delayResp"`
					DelayRespPerSecond             float32 `json:"delayRespPerSecond"`
					FollowUp                       int64   `json:"followUp"`
					FollowUpPerSecond              float32 `json:"followUpPerSecond"`
					Management                     int64   `json:"management"`
					ManagementPerSecond            float32 `json:"managementPerSecond"`
					PeerDelayReq                   int64   `json:"peerDelayReq"`
					PeerDelayReqPerSecond          float32 `json:"peerDelayReqPerSecond"`
					PeerDelayResp                  int64   `json:"peerDelayResp"`
					PeerDelayRespFollowUp          int64   `json:"peerDelayFollowUp"`
					PeerDelayRespFollowUpPerSecond float32 `json:"peerDelayRespFollowUpPerSecond"`
					PeerDelayRespPerSecond         float32 `json:"peerDelayRespPerSecond"`
					Signalling                     int64   `json:"signalling"`
					SignallingPerSecond            float32 `json:"signallingPerSecond"`
					Sync                           int64   `json:"sync"`
					SyncPerSecond                  float32 `json:"syncPerSecond"`
					Total                          int64   `json:"total"`
					TotalPerSecond                 float32 `json:"totalPerSecond"`
				} `json:"tx"`
			} `json:"packetCounters"`
			VirtualInterface string `json:"virtualInterface"`
		} `json:"instances"`
		Interfaces []struct {
			CurrentTime         string `json:"currentTime"`
			CurrentTimeNs       int64  `json:"currentTimeNs"`
			CurrentTimeSec      int64  `json:"currentTimeSec"`
			InternalOffset      string `json:"internalOffset"`
			InternalOffsetNs    int64  `json:"internalOffsetNs"`
			InternalOffsetSubNs int64  `json:"internalOffsetSubNs"`
			Name                string `json:"name"`
		} `json:"interfaces"`
	} `json:"data"`
}
