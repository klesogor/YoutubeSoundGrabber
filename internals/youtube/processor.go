package youtube

import (
	"github.com/klesogor/youtube-grabber/internals"
)

type YoutubeProcessor struct {
	fileSystem internals.FileSystem
	audioCache internals.AudioCache
}

func NewProcessor(workers int) YoutubeProcessor {
	return YoutubeProcessor{}
}

func (ytprocessor *YoutubeProcessor) runWorker(in <-chan internals.VideoRequest, out chan<- internals.AudioResponse, err chan<- internals.ErrorResponse) {

}
