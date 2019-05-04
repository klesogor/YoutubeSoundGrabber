package main

import (
	"fmt"

	"github.com/klesogor/youtube-grabber/grabber"
)

const videoUrl = "https://www.youtube.com/watch?v=GAhiW1Z3GJY"

func main() {
	grabber := grabber.YoutubeDownloadGrabber{}
	stream := grabber.GetDownloadUrl(videoUrl)
	fmt.Printf("%v\n", stream)
}
