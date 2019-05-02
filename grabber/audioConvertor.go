package grabber

import (
	"fmt"
	"os"
	"os/exec"
)

func convertToMp3(req *RequestMessage) {
	filePath := filepath + req.videoId + fileExtMp3
	convertAndSave(filePath, req)
	removeFile(req.tempAudioPath)
	req.cachedAudioPath = filePath
}

func convertAndSave(toPath string, req *RequestMessage) {
	cmd := exec.Command("ffmpeg", "-vn", "-i", req.tempAudioPath, "-y", toPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		req.handleError(err)
	}
}

func removeFile(filepath string) {
	os.Remove(filepath)
}
