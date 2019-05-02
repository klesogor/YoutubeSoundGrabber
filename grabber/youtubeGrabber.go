package grabber

import uuid "github.com/satori/go.uuid"

type MessageHandler func(ResponseMessage)
type FileMessageHandler func(FileMessageHandler)

type youtubeGrabber interface {
	handle(req RequestMessage)
}

type ResponseMessage struct {
	requestId uuid.UUID
	message   string
	err       error
}

type simpleYoutubeGrabber struct {
	messageHandler MessageHandler
	fileHandler    FileMessageHandler
}

type ResponseFileMessage struct {
	requestId uuid.UUID
	filePath  string
	videoId   string
}

type RequestMessage struct {
	requestId        uuid.UUID
	videoUrl         string
	videoId          string
	audioDownloadUrl string
	tempAudioPath    string
	cachedAudioPath  string
	hasError         bool
	fanOut           chan<- interface{}
}

func (mes *RequestMessage) handleError(err error) {
	mes.hasError = true
	mes.fanOut <- ResponseMessage{message: "error ocured", err: err}
}
