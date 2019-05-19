package main

import (
	"fmt"
	"log"
	"os"

	"github.com/klesogor/youtube-grabber/internals"
	"github.com/klesogor/youtube-grabber/internals/youtube"
)

func main() {
	video, tempPath := `https://www.youtube.com/watch?v=a7rhmLC_Zv4&list=RDMMa7rhmLC_Zv4&index=1`, "/home/klesogor/"
	player, err := youtube.GetPlayerConfig(video)
	if err != nil {
		log.Fatal(err.Error())
	}
	audio, err := player.GetBestAudioStream()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(audio)
	res, err := internals.DownloadSegmented(audio, internals.SegmentedDownloadConfig{ChunkSize: 989898, ConcurrencyLimit: 10}, youtube.Downloader)
	if err != nil {
		log.Fatal(err.Error())
	}
	/*audio.Data = res
	var converter internals.FFMPEGConverter
	converted, err := converter.Convert(audio, internals.ConvertingSettings{TargetFormat: internals.MP3, PreserveVideo: false})
	*/
	file, err := os.Create(tempPath + "test.mp3")
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Downloaded youtube audio!")
}
