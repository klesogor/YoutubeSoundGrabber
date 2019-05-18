package internals

type ErrorResponse struct {
	VideoId string
	Error   error
	Context interface{}
}

type AudioCache interface {
	GetAudioFromCache(VideoId string) (string, error)
	PersistAudioInCache(content []byte) error
}
