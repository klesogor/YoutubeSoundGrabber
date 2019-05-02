package grabber

import (
	"errors"
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

func grabDownloadUrl(req *RequestMessage) {
	httpContent := getHttpContent(req)
	if req.hasError {
		return
	}
	urls := getUrlEncodedAudioUrl(httpContent)
	decodedUrls := urlDecode(urls, req)
	if req.hasError {
		return
	}

	req.audioDownloadUrl = getDownloadUrlFromParams(splitUrlStrings(decodedUrls))
}

func getDownloadUrlFromParams(params []string) string {
	url := ""
	for _, v := range params {
		if strings.Index(v, "url=") == 0 {
			url = v[4:]
		}
	}

	return url
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

func urlDecode(input string, req *RequestMessage) string {
	res, err := url.QueryUnescape(input)
	if err != nil {
		req.handleError(err)
	}

	return res
}

func getUrlEncodedAudioUrl(body string) string {
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

func getHttpContent(req *RequestMessage) string {
	res, err := http.Get(req.videoUrl)
	if err != nil {
		req.handleError(err)
		return ""
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		req.handleError(err)
		return ""
	}
	return string(response)
}

func extractVideoIdFromUrl(req *RequestMessage) {
	url, err := url.Parse(req.videoUrl)
	if err != nil {
		req.handleError(err)
		return
	}
	params := url.Query()
	if strings.Index(url.Opaque, "youtube.com") == -1 {
		req.handleError(errors.New("Invalid url base!"))
		return
	}

	videoId := params.Get("v")
	if videoId == "" {
		req.handleError(err)
		return
	}
	req.videoId = videoId
}
