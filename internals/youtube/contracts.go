package youtube

type youtubeDownloader interface {
	tryDownload(url string) ([]byte, error)
}

type signatureDecryptor interface {
	decryptSignature(sign string) string
}
