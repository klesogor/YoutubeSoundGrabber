package grabber

import (
	"fmt"
	"os"
	"os/exec"
)

func ConvertToMp3(name, path string) string {
	filePath := filepath + name + fileExtTarget
	convertAndSave(filePath, path)
	removeFile(path)
	return filePath
}

func convertAndSave(toPath, fromPath string) {
	cmd := exec.Command("ffmpeg", "-vn", "-i", fromPath, "-y", toPath)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
	}
}

func removeFile(filepath string) {
	os.Remove(filepath)
}
