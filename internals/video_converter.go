package internals

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
)

const (
	_ = iota
	MP3
	MP4
	WEBP
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
		options = append(options, "-nv")
	}
	options = append(options, "-f", formatToExt(settings.TargetFormat), "pipe:1")

	cmd := exec.Command("ffmpeg", options...)
	cmd.Stdin = bytes.NewReader(c.GetData())
	buffer, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(buffer)
}

func formatToExt(f int) string {
	switch f {
	case 1:
		return "mp3"
	case 2:
		return "mp4"
	case 3:
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
	} else if strings.Index(ctype, "webp") != -1 {
		return WEBP
	}

	return -1
}

type ConvertingSettings struct {
	TargetFormat  int
	PreserveVideo bool
}
