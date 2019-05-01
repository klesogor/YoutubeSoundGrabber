package main

import (
	"fmt"

	"github.com/klesogor/youtube-grabber/grabber"
)

func main() {
	url := grabber.GrabDownloadUrl(`https://www.youtube.com/watch?v=dV4nRwg4PRw&list=RDRsPgV5ijaNg&index=4`)
	downloaded := grabber.ExtractToFile(url, "TestId")
	grabber.ConvertToMp3(downloaded)
	fmt.Println("Downloaded and converted")
}
