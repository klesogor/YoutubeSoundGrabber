package grabber

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	jsonVariableForUrls = "\"adaptive_fmts\":"
	delimiter           = "\\u0026"
	otherDelims         = ";,"
	formatsDelim        = "bitrate"
)

type DownloadGrabber interface {
	GetAudioDownloadUrl(url string) (AudioStream, error)
}

type YoutubeDownloadGrabber struct {
	StreamData *StreamData
}

type StreamData struct {
	muxedStreams    string
	adaptiveFormats AdaptiveFormats
}

type AdaptiveFormats struct {
	VideoStreams []VideoStream
	AudioStreams []AudioStream
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
	base          StreamBase
	quality_label string
	eotf          string
	size          string
	fps           uint
}

type AudioStream struct {
	base              StreamBase
	audio_channels    uint
	audio_sample_rate uint
}

func (downloadGrabber *YoutubeDownloadGrabber) GetDownloadUrl(url string) (AudioStream, error) {
	adaptiveFormats, err := extractAdaptiveFormats(url)
	if err != nil {
		return AudioStream{}, err
	}
	streamData := extractStreamDataFromAdaptiveFormats(adaptiveFormats)
	downloadGrabber.StreamData = streamData
	return streamData.adaptiveFormats.AudioStreams[0], nil
}

func extractAdaptiveFormats(videoUrl string) (string, error) {
	content, err := getVideoPage(videoUrl)
	if err != nil {
		return "", err
	}
	formats := getUrlEncodedAdaptiveFormats(content)
	decodedFormats, err := url.QueryUnescape(formats)
	if err != nil {
		return "", err
	}
	return decodedFormats, nil
}

func getUrlEncodedAdaptiveFormats(body string) string {
	indexOfJson := strings.Index(body, jsonVariableForUrls)
	substr := body[indexOfJson+len(jsonVariableForUrls)+1:]
	var res strings.Builder
	for _, v := range substr {
		if v == '"' {
			return res.String()
		}
		res.WriteRune(v)
	}

	return res.String()
}

func trimStrings(s ...string) []string {
	res := make([]string, 0, len(s))
	for _, v := range s {
		res = append(res, strings.TrimSpace(v))
	}

	return res
}

func getVideoPage(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

func extractStreamDataFromAdaptiveFormats(adaptiveFormats string) *StreamData {
	splittedFormats := splitAdaptiveFormats(adaptiveFormats)
	groupedFormats := groupAdaptiveFormats(splittedFormats)
	return createStreamDataFromGroupedFormats(groupedFormats)
}

func splitAdaptiveFormats(s string) []string {
	res := make([]string, 0, 20)
	splitted := strings.Split(s, delimiter)
	for _, v := range splitted {
		splittedParam := strings.FieldsFunc(v, func(r rune) bool {
			return strings.Index(otherDelims, string(r)) != -1
		})
		res = append(res, splittedParam...)
	}

	return res
}

func groupAdaptiveFormats(splittedFormats []string) []map[string]string {
	res := make([]map[string]string, 0)
	var cur map[string]string
	for _, v := range splittedFormats {
		pair := strings.Split(v, "=")
		if pair[0] == formatsDelim {
			if len(cur) > 0 {
				res = append(res, cur)
				cur = make(map[string]string)
			}
		}
		if len(pair) == 2 {
			cur[pair[0]] = pair[1]
		} else {
			cur[pair[0]] = ""
		}
	}

	return res
}

func createStreamDataFromGroupedFormats(data []map[string]string) *StreamData {
	videoStreams, audioStreams := make([]VideoStream, 4), make([]AudioStream, 4)
	for _, v := range data {
		res := parseStreamFromDataMap(v)
		switch res.(type) {
		case VideoStream:
			videoStreams = append(videoStreams, res.(VideoStream))
			break
		case AudioStream:
			audioStreams = append(audioStreams, res.(AudioStream))
		}
	}

	return &StreamData{adaptiveFormats: AdaptiveFormats{VideoStreams: videoStreams, AudioStreams: audioStreams}}
}

func parseStreamFromDataMap(data map[string]string) interface{} {
	bitrate, _ := strconv.ParseUint(data["bitrate"], 10, 32)
	itag, _ := strconv.ParseUint(data["itag"], 10, 32)
	clen, _ := strconv.ParseUint(data["clen"], 10, 32)
	init := data["init"]
	ctype := data["type"]
	codecs := data["codecs"]
	lmt, _ := strconv.ParseUint(data["lmt"], 10, 32)
	url := data["url"]
	index := data["index"]
	projection_type, _ := strconv.ParseUint(data["projection_type"], 10, 32)
	xtags := data["xtags"]
	signature := data["s"]
	base := StreamBase{
		bitrate:         uint(bitrate),
		itag:            uint(itag),
		clen:            uint(clen),
		init:            init,
		ctype:           ctype,
		codecs:          codecs,
		lmt:             uint(lmt),
		url:             url,
		index:           index,
		projection_type: uint(projection_type),
		xtags:           xtags,
		signature:       signature}
	if size, err := data["size"]; !err {
		fps, _ := strconv.ParseUint(data["fps"], 10, 32)
		return VideoStream{base: base, fps: uint(fps), quality_label: data["quality_label"], size: size, eotf: data["eotf"]}
	}
	chans, _ := strconv.ParseUint(data["audio_channels"], 10, 32)
	rate, _ := strconv.ParseUint(data["audio_sample_rate"], 10, 32)
	return AudioStream{base: base, audio_channels: uint(chans), audio_sample_rate: uint(rate)}
}
