package main

import (
	"fmt"
	"net/url"

	"github.com/klesogor/youtube-grabber/grabber"
)

const videoUrl = "https://www.youtube.com/watch?v=GAhiW1Z3GJY"

func main() {
	Grabber := grabber.YoutubeDownloadGrabber{}
	stream, err := Grabber.GetStreamData(videoUrl)
	if err != nil {
		fmt.Println("Error is not nil!!! Not nil!!!")
	}
	url, _ := url.ParseQuery(videoUrl)
	downloader := grabber.YoutubeSegmentedAudioDownloader{DownloadLimit: 20, DownloadRangeLimit: 98989}
	var astream grabber.AudioStream
	for _, v := range stream.AdaptiveFormats.AudioStreams {
		if v.Base.Ctype == "audio/webm" {
			astream = v
			break
		}
	}
	path := downloader.DownloadAudioByStream(&astream, url.Get("v"))
	finalPath := grabber.ConvertToMp3(url.Get("v"), path)
	fmt.Printf("%v\n", finalPath)
}
