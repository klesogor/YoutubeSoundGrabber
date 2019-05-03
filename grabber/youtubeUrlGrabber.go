package grabber

type StreamData struct {
	muxedStreams    string
	adaptiveFormats AdaptiveFormats
}

type AdaptiveFormats struct {
	videoStreams []VideoStream
	audioStreams []AudioStream
}

type StreamBase struct {
	bitrate         uint
	itag            uint
	clen            uint
	init            string
	ctype           string
	codecs          string
	lmt             uint
	url             string
	index           string
	projection_type uint
	xtags           string
	signature       string
}

type VideoStream struct {
	quality_label string
	eotf          string
	size          string
	fps           uint
}

type AudioStream struct {
	audio_channels    uint
	audio_sample_rate uint
}

func extractAdaptiveFormats(url string) (string, error) {
	return "", nil
}

func extractStreamDataFromAdaotiveFormats(adaptiveFormats string) (*StreamData, error) {
	return &StreamData{}, nil
}
