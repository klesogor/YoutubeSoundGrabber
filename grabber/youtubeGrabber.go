package grabber;

import (
	"net/url"
)

type YoutubeGrabber interface {
	Handle(videoUrl string, context interface{})
}

type SimpleYoutubeGrabber struct {
	FanIn chan RequestMessage
	FanOut chan ResponseFileMessage
	Err chan ResponseErrorMessage
}

type ResponseErrorMessage struct {
	Context interface{}
	Error error 
}

type ResponseFileMessage struct {
	Context interface{}
	AudioPath string
}

type RequestMessage struct {
	VideoUrl string
	VideoId  string
	Context  interface{}
}

func NewHandler(workers int) SimpleYoutubeGrabber {
	com := make(chan RequestMessage, workers*2)
	resp := make(chan ResponseFileMessage, workers*2)
	err := make(chan ResponseErrorMessage, workers*2)
	for i := 0; i < workers; i++ {
		go runWorker(com,resp,err)
	}

	return SimpleYoutubeGrabber{FanIn: com, FanOut: resp, Err: err}
}

func runWorker(in <-chan RequestMessage, out chan<- ResponseFileMessage, Err chan<- ResponseErrorMessage) {
	Grabber := YoutubeDownloadGrabber{}
	downloader := YoutubeSegmentedAudioDownloader{DownloadLimit: 20, DownloadRangeLimit: 98989}
	for {
		mes := <-in
		stream, err := Grabber.GetStreamData(mes.VideoUrl)
		if err != nil {
			Err <- ResponseErrorMessage{Context: mes.Context,Error: err}
		}
		var astream AudioStream
		for _, v := range stream.AdaptiveFormats.AudioStreams {
			if v.Base.Ctype == "audio/webm" {
				astream = v
				break
			}
		}
		path, err := downloader.DownloadAudioByStream(&astream, mes.VideoId)
		if err != nil {
			Err <- ResponseErrorMessage{Context: mes.Context,Error: err}
		}
		finalPath := ConvertToMp3(mes.VideoId, path)
		out<- ResponseFileMessage{Context: mes.Context,AudioPath: finalPath}
	}
}

func (grabber SimpleYoutubeGrabber) Handle(videoUrl string, context interface{}) {
	params, err := url.ParseQuery(videoUrl)
	if err != nil{
		grabber.Err<- ResponseErrorMessage{Context: context, Error: err}
	}
	mes := RequestMessage{VideoUrl: videoUrl,VideoId: params.Get("v"), Context: context}
	grabber.FanIn <- mes
}
