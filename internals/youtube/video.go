package youtube

import "errors"

const urlSignatureKey = "signature"

type youtubeVideo struct {
	id         string
	title      string
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

func (sb *streamBase) decryptSignature(decryptor signatureDecryptor) {
	sb.url += urlSignatureKey + "=" + decryptor.decryptSignature(sb.signature)
	sb.isDecrypted = true
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

func (v *youtubeVideo) getBestQualityAudioStream() (downloadable, error) {
	audioStreams := v.streamData.adaptiveFormats.audioStreams
	var res downloadable
	var maxBitRate int
	if audioStreams != nil && len(audioStreams) > 0 {
		for _, v := range audioStreams {
			if v.bitrate > maxBitRate {
				maxBitRate = v.bitrate
				res = v
			}
		}
		return res, nil
	}
	return nil, nil
}
