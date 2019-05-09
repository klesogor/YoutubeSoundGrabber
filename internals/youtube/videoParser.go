package youtube

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const baseVideoUrl = "https://youtube.com/watch?v="

type videoRawData struct {
	Title        string `json:"title"`
	AdaptiveFmts string `json:"adaptive_fmts"`
	MuxedStreams string `json:"url_encoded_fmt_stream_map"`
	body         []byte
}

type videoDataParser struct {
}

func (vp *videoDataParser) parseVideoData(vid string) (*videoRawData, error) {
	res, err := http.Get(baseVideoUrl + vid)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	bodyString := string(body)
	data, err := parseArgsJson(bodyString)
	if err != nil {
		return nil, err
	}
	data.body = body
	return data, nil
}

func parseArgsJson(body string) (*videoRawData, error) {
	indexOfArgs := strings.Index(body, "args:")
	if indexOfArgs == -1 {
		return nil, errors.New("No args found in body")
	}

	s, err := parseJsonString(body[:indexOfArgs])
	if err != nil {
		return nil, err
	}
	var data videoRawData
	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func parseJsonString(s string) (string, error) {
	cur, arr := 0, 0
	var args strings.Builder
	if s[0] != '{' && s[0] != '[' {
		return "", errors.New("String should start with [ or {")
	}
	for _, v := range s {
		args.WriteRune(v)
		switch v {
		case '{':
			cur += 1
		case '[':
			arr += 1
		case '}':
			cur -= 1
		case ']':
			arr -= 1
		}

		if cur == 0 && arr == 0 {
			return args.String(), nil
		}
	}

	return args.String(), errors.New("Unable to parse json due to enclosing delimiter")
}
