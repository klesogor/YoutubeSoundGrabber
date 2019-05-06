package grabber

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const (
	filepath      = "/home/klesogor/parsed-audio/"
	fileExt       = ".webm"
	fileExtTarget = ".mp3"
)

type AudioDownloader interface {
	DownloadAudioByStream(AudioStream, string) string
}

type YoutubeSegmentedAudioDownloader struct {
	DownloadLimit      uint
	DownloadRangeLimit uint
}

type dataSegment struct {
	offset uint
	data   []byte
}

func (downloader *YoutubeSegmentedAudioDownloader) DownloadAudioByStream(stream *AudioStream, name string, retries int) (string, error) {
	var dec YoutubeDechiper
	url := stream.Base.Url
	if stream.Base.Signature != "" {
		url += "&signature=" + dec.Decrypt(stream.Base.Signature)
	}

	fmt.Printf("Downloading from base url:%v\n", url)
	success := false
	var acc [][]byte
	var err error
	for i := 0; i < retries; i++ {
		acc, err = tryDownload(stream, url, downloader.DownloadRangeLimit)
		if err == nil {
			success = true
			break
		}
	}
	if !success {
		return "", errors.New("Unable to download file!")
	}

	return writeToFile(acc, name)
}

func tryDownload(stream *AudioStream, url string, downloadRange uint) ([][]byte, error) {
	var wg sync.WaitGroup
	totalDownloadsRequired := stream.Base.Clen/downloadRange + 1
	acc, i, start, totalData, size := make([][]byte, totalDownloadsRequired), 0, 0, int(stream.Base.Clen), int(downloadRange)
	var downloadErr error
	for totalData > 0 {
		wg.Add(1)
		go func(offset, start, size int) {
			urlWithRange := url + "&range=" + strconv.Itoa(start) + "-" + strconv.Itoa(start+size)
			response, err := http.Get(urlWithRange)
			if err != nil {
				downloadErr = err
				wg.Done()
				return
			}
			if response.StatusCode != 200 {
			}
			res, _ := ioutil.ReadAll(response.Body)
			acc[offset] = res
			wg.Done()
		}(i, start, size)
		start += size + 1
		totalData -= size
		i += 1
	}
	wg.Wait()
	if downloadErr != nil {
		return nil, downloadErr
	}

	return acc, nil
}

func writeToFile(arr [][]byte, fileName string) (string, error) {
	path := filepath + fileName + fileExt
	handle, err := os.Create(path)
	if err != nil {
		return "", err
	}
	for _, v := range arr {
		handle.Write(v)
	}

	return path, nil
}
