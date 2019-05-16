package internals

type VideoRequest struct {
	VideoUrl string
	VideoId  string
	Context  interface{}
}

type AudioResponse struct {
	VideoId   string
	AudioName string
	Data      []byte
	Context   interface{}
}

type ErrorResponse struct {
	VideoId string
	Error   error
	Context interface{}
}

type AudioCache interface {
	GetAudioFromCache(VideoId string) (string, error)
	PersistAudioInCache(content []byte) error
}

type FileSystem interface {
	Create(Path string) (FileSystemWriter, error)
	Open(Path string) (FileSystemReader, error)
	Remove(Path string) error
}

type FileSystemWriter interface {
	Write(Bytes []byte)
}

type FileSystemReader interface {
	ReadAll() []byte
}
