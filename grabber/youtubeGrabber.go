package grabber

type MessageHandler func(ResponseMessage)
type FileMessageHandler func(ResponseFileMessage)

type YoutubeGrabber interface {
	Handle(url string, handlers Handlers)
}

type ResponseMessage struct {
	Message string
	Err     error
}

type Handlers struct {
	MessageHandler MessageHandler
	FileHandler    FileMessageHandler
}

type SimpleYoutubeGrabber struct {
	fanIn chan<- RequestMessage
}

type ResponseFileMessage struct {
	FilePath string
	VideoId  string
}

type RequestMessage struct {
	videoUrl         string
	videoId          string
	audioDownloadUrl string
	tempAudioPath    string
	cachedAudioPath  string
	hasError         bool
	handler          Handlers
}

func (mes *RequestMessage) handleError(err error) {
	mes.hasError = true
	go mes.handler.MessageHandler(ResponseMessage{Message: "error ocured", Err: err})
}

func NewHandler(workers int) YoutubeGrabber {
	com := make(chan RequestMessage, workers*2)
	for i := 0; i < workers; i++ {
		go runWorker(com)
	}

	return SimpleYoutubeGrabber{fanIn: com}
}

func runWorker(in <-chan RequestMessage) {
	for {
		message := <-in
		extractVideoIdFromUrl(&message)
		if message.hasError {
			continue
		}
		grabDownloadUrl(&message)
		if message.hasError {
			continue
		}
		extractToFile(&message)
		if message.hasError {
			continue
		}
		convertToMp3(&message)
		message.handler.FileHandler(ResponseFileMessage{FilePath: message.cachedAudioPath, VideoId: message.videoId})
	}
}

func (grabber SimpleYoutubeGrabber) Handle(url string, handler Handlers) {
	mes := RequestMessage{videoUrl: url}
	grabber.fanIn <- mes
}
