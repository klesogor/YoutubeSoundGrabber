package grabber

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	jsonVariableForUrls = "\"adaptive_fmts\":"
	delimiter           = "\\u0026"
	otherDelims         = ";,"
)

type DownloadGrabber interface {
	GetDownloadUrl(url string) string
}

type YoutubeDownloadGrabber struct {
	StreamData *StreamData
}

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

func extractStreamDataFromAdaptiveFormats(adaptiveFormats string) (*StreamData, error) {

	return &StreamData{}, nil
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

func splitUrlStrings(s string) []string {
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
