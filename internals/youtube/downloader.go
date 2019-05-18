package youtube

import (
	"net/http"
	"strconv"
)

func Downloader(url string, start, end int) ([]byte, error) {
	res, err := http.Get(url + "&range=" + strconv.Itoa(start) + "-" + strconv.Itoa(end))
	if err != nil {
		return nil, err
	}
	return parseBody(res)
}
