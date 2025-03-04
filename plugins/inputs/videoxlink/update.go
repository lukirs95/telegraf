package videoxlink

import (
	client "github.com/lukirs95/goxlinkclient"
)

func (xlink *VideoXLink) handleUpdate(update client.XLink) {
	systemId := update.Ident()

	var cache *updateCache
	if c, ok := xlink.updateCache[systemId]; ok {
		cache = c
	} else {
		cache = &updateCache{
			eth:     make(map[string]*ethUpdateCache),
			decoder: make(map[string]*decoderUpdateCache),
			encoder: make(map[string]*encoderUpdateCache),
		}
		xlink.updateCache[systemId] = cache
	}

	if name, OK := update.GetName(); OK {
		cache.name = name
	}

	for _, eth := range update.GetInterfaces() {
		var ethCache *ethUpdateCache
		if c, ok := cache.eth[eth.Ident()]; ok {
			ethCache = c
		} else {
			ethCache = &ethUpdateCache{}
			cache.eth[eth.Ident()] = ethCache
		}

		if stat, OK := eth.IsEnabled(); OK {
			ethCache.adminEnabled = stat
		}
		if stat, OK := eth.IsLinkUp(); OK {
			ethCache.up = stat
		}
		if stat, OK := eth.IsActive(); OK {
			ethCache.activeUplink = stat
		}
		if stat, OK := eth.IsDefaultUplink(); OK {
			ethCache.primaryUplink = stat
		}
		if stat, OK := eth.IsBackupUplink(); OK {
			ethCache.secondaryUplink = stat
		}

	}

	for _, decoder := range update.GetDecoders() {
		var decoderCache *decoderUpdateCache
		if c, ok := cache.decoder[decoder.Ident()]; ok {
			decoderCache = c
		} else {
			decoderCache = &decoderUpdateCache{}
			cache.decoder[decoder.Ident()] = decoderCache
		}

		if stat, OK := decoder.IsConnected(); OK {
			decoderCache.conected = stat
		}
		if stat, OK := decoder.IsRunning(); OK {
			decoderCache.running = stat
		}
		if stat, OK := decoder.IsEnabled(); OK {
			decoderCache.enabled = stat
		}
		if stat, OK := decoder.IsVideoEnabled(); OK {
			decoderCache.videoEnabled = stat
		}
		if stat, OK := decoder.IsAudioEnabled(); OK {
			decoderCache.audioEnabled = stat
		}

		if stat, OK := decoder.HasVideoSignal(); OK {
			decoderCache.hasVideoSignal = stat
		}

		if stat, OK := decoder.HasAudioSignal(); OK {
			decoderCache.hasAudioSignal = stat
		}
	}

	for _, encoder := range update.GetEncoders() {
		var encoderCache *encoderUpdateCache
		if c, ok := cache.encoder[encoder.Ident()]; ok {
			encoderCache = c
		} else {
			encoderCache = &encoderUpdateCache{}
			cache.encoder[encoder.Ident()] = encoderCache
		}

		if stat, OK := encoder.IsConnected(); OK {
			encoderCache.connected = stat
		}
		if stat, OK := encoder.IsRunning(); OK {
			encoderCache.running = stat
		}
		if stat, OK := encoder.IsEnabled(); OK {
			encoderCache.enabled = stat
		}
		if stat, OK := encoder.IsVideoEnabled(); OK {
			encoderCache.videoEnabled = stat
		}
		if stat, OK := encoder.IsAudioEnabled(); OK {
			encoderCache.audioEnabled = stat
		}

		if stat, OK := encoder.HasVideoSignal(); OK {
			encoderCache.hasVideoSignal = stat
		}

		if stat, OK := encoder.HasAudioSignal(); OK {
			encoderCache.hasAudioSignal = stat
		}
	}

	for _, metric := range cache.Metric(systemId) {
		xlink.Buf.PushBack(metric)
	}
}
