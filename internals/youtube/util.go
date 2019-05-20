package youtube

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func parseBody(resp *http.Response) ([]byte, error) {
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, errors.New("Not 2xx http status!")
	}
	return ioutil.ReadAll(resp.Body)
}

func extractJsonFromString(s string) (string, error) {
	curl, arr := 0, 0
	var res strings.Builder
	index := strings.Index(s, "{")
	if index == -1 {
		return "", errors.New("No opening tag found")
	}
	for _, v := range s[index:] {
		res.WriteRune(v)
		switch v {
		case '{':
			curl++
		case '[':
			arr++
		case ']':
			arr--
		case '}':
			curl--
			if arr == 0 && curl == 0 {
				return res.String(), nil
			}
		}
	}

	return "", errors.New("No json found in string")
}

func getConfig(body string) (string, error) {
	return getBetween(body, "ytplayer.config = ", ";ytplayer.load")
}

func getBetween(s, start, end string) (string, error) {
	offsetStart, lenS := strings.Index(s, start), len(start)
	if offsetStart == -1 {
		return "", errors.New("No substring match in string")
	}
	offsetEnd := strings.Index(s[offsetStart+lenS:], end)
	if offsetEnd == -1 {
		return "", errors.New("No substring match in string")
	}

	return s[offsetStart+lenS : offsetStart+lenS+offsetEnd], nil
}

func getIntOrDefault(val string) int {
	res, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}

	return res
}

func getAudioScore(stream StreamData) int {
	audioBonus := 0
	if strings.Index(stream.Ctype, "audio") != -1 {
		audioBonus = 10000000
	}
	return stream.Bitrate + audioBonus
}

func parseStreamData(adaptiveFmts string) ([]StreamData, error) {
	splitted := strings.Split(adaptiveFmts, ",")
	res := make([]StreamData, 0, 1)
	for _, v := range splitted {
		res = append(res, parseStreamDataFromUrlString(v))
	}
	return res, nil
}
