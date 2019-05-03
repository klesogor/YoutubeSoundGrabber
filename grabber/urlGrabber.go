package grabber

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func grabDownloadUrl(req *RequestMessage) {
	httpContent := getHttpContent(req)
	if req.hasError {
		return
	}
	urls := ""
	decodedUrls := urlDecode(urls, req)
	if req.hasError {
		return
	}

	req.audioDownloadUrl = getDownloadUrlFromParams(splitUrlStrings(decodedUrls))
}

func getDownloadUrlFromParams(params []string) string {
	for _, v := range params {
		if isRightUrl(v) {
			return v[4:]
		}
	}

	return ""
}
func isRightUrl(url string) bool {
	return strings.Index(url, "url=") == 0 && strings.Index(url, "mime=audio%2Fwebm") != -1
}

func urlDecode(input string, req *RequestMessage) string {
	res, err := url.QueryUnescape(input)
	if err != nil {
		req.handleError(err)
	}

	return res
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
	if strings.Index(req.videoUrl, "youtube.com") == -1 {
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
