package youtube

const urlSignatureKey = "signature"

type bitrateable interface {
	getBitrate() int
	getStreamBase() *streamBase
}

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
}

func (sb *streamBase) getBitrate() int {
	return sb.bitrate
}

func (v *youtubeVideo) getBestQualityAudioStream() (*audioStream, error) {
	audioStreams := v.streamData.adaptiveFormats.audioStreams
	if audioStreams != nil && len(audioStreams) > 0 {

	}
	return nil, nil
}
