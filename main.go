package main

import (
	"fmt"

	"github.com/klesogor/youtube-grabber/grabber"
)

const videoUrl = "https://www.youtube.com/watch?v=GAhiW1Z3GJY"

func main() {
	grabber := grabber.YoutubeDownloadGrabber{}
	stream, err := grabber.GetStreamData(videoUrl)
	if err != nil {
		fmt.Println("Error is not nil!!! Not nil!!!")
	}
	fmt.Printf("%q\n", stream)
}
