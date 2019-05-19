package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	video, tempPath := `https://www.youtube.com/watch?v=a7rhmLC_Zv4&list=RDMMa7rhmLC_Zv4&index=1`, "/home/klesogor/"

	/*audio.Data = res
	var converter internals.FFMPEGConverter
	converted, err := converter.Convert(audio, internals.ConvertingSettings{TargetFormat: internals.MP3, PreserveVideo: false}) */
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
