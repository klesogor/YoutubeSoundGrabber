package youtube

type youtubeDownloader interface {
	tryDownload(url string) ([]byte, error)
}

type downloadable interface {
	getContentLength() int
	getContentType() string
	getDownloadUrl() (string, error)
}

type signatureDecryptor interface {
	decryptSignature(sign string) string
}
