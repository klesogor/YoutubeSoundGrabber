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
)

type DownloadGrabber interface {
	GetStreamData(url string) (*StreamData, error)
}

type YoutubeDownloadGrabber struct {
}

type StreamData struct {
	MuxedStreams    string
	AdaptiveFormats AdaptiveFormats
}

type AdaptiveFormats struct {
	VideoStreams []VideoStream
	AudioStreams []AudioStream
}

type StreamBase struct {
	Bitrate         uint
	Itag            uint
	Clen            uint
	Init            string
	Ctype           string
	Codecs          string
	Lmt             uint
	Url             string
	Index           string
	Projection_type uint
	Xtags           string
	Signature       string
}

type VideoStream struct {
	Base          StreamBase
	Quality_label string
	Eotf          string
	Size          string
	Fps           uint
}

type AudioStream struct {
	Base              StreamBase
	Audio_channels    uint
	Audio_sample_rate uint
}

func (downloadGrabber *YoutubeDownloadGrabber) GetStreamData(url string) (*StreamData, error) {
	adaptiveFormats, err := extractAdaptiveFormats(url)
	if err != nil {
		return nil, err
	}
	streamData := extractStreamDataFromAdaptiveFormats(adaptiveFormats)

	return streamData, nil
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
	trimedFormats := trimStrings(splittedFormats...)
	groupedFormats := groupAdaptiveFormats(trimedFormats)
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
	//elements order is randomized every time, so we peek first element to get delimeter
	formatsDelim := strings.Split(splittedFormats[0], "=")[0]
	res := make([]map[string]string, 0)
	cur := make(map[string]string)
	for _, v := range splittedFormats {
		pair := strings.SplitN(v, "=", 2)
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
	videoStreams, audioStreams := make([]VideoStream, 0, 4), make([]AudioStream, 0, 4)
	for _, v := range data {
		res := parseStreamFromDataMap(v)
		switch res.(type) {
		case VideoStream:
			videoStreams = append(videoStreams, res.(VideoStream))
			break
		case AudioStream:
			audioStreams = append(audioStreams, res.(AudioStream))
			break
		}
	}

	return &StreamData{AdaptiveFormats: AdaptiveFormats{VideoStreams: videoStreams, AudioStreams: audioStreams}}
}

func parseStreamFromDataMap(data map[string]string) interface{} {
	bitrate, _ := strconv.ParseUint(data["bitrate"], 10, 32)
	itag, _ := strconv.ParseUint(data["itag"], 10, 32)
	clen, _ := strconv.ParseUint(data["clen"], 10, 32)
	init := data["init"]
	ctype := data["type"]
	codecs := data[" codecs"]
	lmt, _ := strconv.ParseUint(data["lmt"], 10, 32)
	url := data["url"]
	index := data["index"]
	projection_type, _ := strconv.ParseUint(data["projection_type"], 10, 32)
	xtags := data["xtags"]
	signature := data["s"]
	base := StreamBase{
		Bitrate:         uint(bitrate),
		Itag:            uint(itag),
		Clen:            uint(clen),
		Init:            init,
		Ctype:           ctype,
		Codecs:          codecs,
		Lmt:             uint(lmt),
		Url:             url,
		Index:           index,
		Projection_type: uint(projection_type),
		Xtags:           xtags,
		Signature:       signature}
	if strings.Index(ctype, "video") != -1 {
		fps, _ := strconv.ParseUint(data["fps"], 10, 32)
		return VideoStream{Base: base, Fps: uint(fps), Quality_label: data["quality_label"], Size: data["size"], Eotf: data["eotf"]}
	}
	chans, _ := strconv.ParseUint(data["audio_channels"], 10, 32)
	rate, _ := strconv.ParseUint(data["audio_sample_rate"], 10, 32)
	return AudioStream{Base: base, Audio_channels: uint(chans), Audio_sample_rate: uint(rate)}
}
