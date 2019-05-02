package grabber

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	filepath   = "/home/klesogor/parsed-audio/"
	fileExt    = ".webm"
	fileExtMp3 = ".mp3"
)

func extractToFile(req *RequestMessage) {
	fileContents := downloadFile(req)
	if req.hasError {
		return
	}
	req.tempAudioPath = writeToFile(fileContents, req.videoId+fileExt, req)
}

func writeToFile(content, name string, req *RequestMessage) string {
	path := filepath + name
	file, err := os.Create(path)
	if err != nil {
		req.handleError(err)
		return ""
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	w.WriteString(content)
	w.Flush()

	return path
}

func downloadFile(req *RequestMessage) string {
	res, err := http.Get(req.audioDownloadUrl)
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
