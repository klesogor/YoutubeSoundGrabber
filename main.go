package main

import (
	"fmt"

	"github.com/klesogor/youtube-grabber/grabber"
)

const videoUrl = "https://www.youtube.com/watch?v=BUWuDdfe7v4"

func main() {
	Grabber := grabber.YoutubeDownloadGrabber{}
	stream, err := Grabber.GetStreamData(videoUrl)
	if err != nil {
		fmt.Println("Error is not nil!!! Not nil!!!")
	}
	downloader := grabber.YoutubeSegmentedAudioDownloader{DownloadLimit: 20, DownloadRangeLimit: 98989}
	var astream grabber.AudioStream
	for _, v := range stream.AdaptiveFormats.AudioStreams {
		if v.Base.Ctype == "audio/webm" {
			astream = v
			break
		}
	}
	path := downloader.DownloadAudioByStream(&astream, "GAhiW1Z3GJY")
	finalPath := grabber.ConvertToMp3("GAhiW1Z3GJY", path)
	fmt.Printf("%v\n", finalPath)
}
