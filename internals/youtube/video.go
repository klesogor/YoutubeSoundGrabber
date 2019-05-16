package youtube

import "errors"

const urlSignatureKey = "signature"

type YoutubeVideo struct {
	Id         string
	Title      string
	streamData streamData
	body       string
}

type streamData struct {
	muxedStreams    string
	adaptiveFormats adaptiveFormats
}

type adaptiveFormats struct {
	videoStreams []videoStream
	audioStreams []audioStream
}

type streamBase struct {
	bitrate         int
	itag            int
	clen            int
	init            string
	ctype           string
	codecs          string
	lmt             int
	url             string
	index           string
	projection_type int
	xtags           string
	signature       string
	isDecrypted     bool
}

type videoStream struct {
	streamBase
	quality_label string
	eotf          string
	size          string
	fps           int
}

type audioStream struct {
	streamBase
	audio_channels    int
	audio_sample_rate int
}

func (sb *streamBase) mustDecryptSignature() bool {
	return sb.signature != "" && sb.isDecrypted
}

func (sb streamBase) getContentLength() int {
	return sb.clen
}

func (sb streamBase) getContentType() string {
	return sb.ctype
}

func (sb streamBase) getDownloadUrl() (string, error) {
	if sb.mustDecryptSignature() {
		return "", errors.New("Video signature must be decrypted first!")
	}
	return sb.url, nil
}
