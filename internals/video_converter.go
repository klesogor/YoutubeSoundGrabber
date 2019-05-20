package internals

import (
	"bytes"
	"os/exec"
	"strings"
)

const (
	_ = iota
	MP3
	MP4
	WEBM
)

type Convertable interface {
	GetContentType() string
	GetData() []byte
}

type Converter interface {
	Convert(c Convertable, settings ConvertingSettings) ([]byte, error)
}

type FFMPEGConverter struct{}

func (f *FFMPEGConverter) Convert(c Convertable, settings ConvertingSettings) ([]byte, error) {
	options := []string{"-stdin", "-f", formatToExt(ctypeToFormat(c.GetContentType())), "-i", "-"}
	if !settings.PreserveVideo {
		options = append(options, "-vn")
	}
	options = append(options, "-f", formatToExt(settings.TargetFormat), "pipe:1")
	cmd := exec.Command("ffmpeg", options...)
	cmd.Stdin = bytes.NewReader(c.GetData())
	return cmd.Output()
}

func formatToExt(f int) string {
	switch f {
	case MP3:
		return "mp3"
	case MP4:
		return "mp4"
	case WEBM:
		return "webm"
	default:
		return ""
	}
}

func ctypeToFormat(ctype string) int {
	if strings.Index(ctype, "mp3") != -1 {
		return MP3
	} else if strings.Index(ctype, "mp4") != -1 {
		return MP4
	} else if strings.Index(ctype, "webm") != -1 {
		return WEBM
	}

	return -1
}

type ConvertingSettings struct {
	TargetFormat  int
	PreserveVideo bool
}
