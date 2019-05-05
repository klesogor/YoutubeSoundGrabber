package grabber

import (
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

func (downloader *YoutubeSegmentedAudioDownloader) DownloadAudioByStream(stream *AudioStream, name string) string {
	var wg sync.WaitGroup
	totalDownloadsRequired := stream.Base.Clen/downloader.DownloadRangeLimit + 1
	acc, i, start, totalData, size := make([][]byte, totalDownloadsRequired), 0, 0, int(stream.Base.Clen), int(downloader.DownloadRangeLimit)
	for totalData > 0 {
		wg.Add(1)
		go func(offset, start, size int) {
			urlWithRange := stream.Base.Url + "&range=" + strconv.Itoa(start) + "-" + strconv.Itoa(start+size)
			fmt.Printf("range: %d-%d\n", start, start+size)
			response, _ := http.Get(urlWithRange)
			res, _ := ioutil.ReadAll(response.Body)
			acc[offset] = res
			wg.Done()
		}(i, start, size)
		start += size + 1
		totalData -= size
		i += 1
	}
	wg.Wait()
	return writeToFile(acc, name)
}

func writeToFile(arr [][]byte, fileName string) string {
	path := filepath + fileName + fileExt
	handle, err := os.Create(path)
	if err != nil {
		return ""
	}
	for _, v := range arr {
		handle.Write(v)
	}

	return path
}
