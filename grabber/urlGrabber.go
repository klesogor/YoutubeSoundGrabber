package grabber

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const jsonVariableForUrls = "\"adaptive_fmts\":"
const delimiter = "\\u0026"

func GrabDownloadUrl(videoUrl string) string {
	urls := urlDecode(getUrlEncodedAudioUrl(getHttpContent(videoUrl)))
	splitted := strings.Split(urls, delimiter)
	for _, v := range splitted {
		fmt.Println(v)
	}
	return "not implemented yet!"
}

func urlDecode(input string) string {
	res, err := url.QueryUnescape(input)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func getUrlEncodedAudioUrl(body string) string {
	indexOfJson := strings.Index(body, jsonVariableForUrls)
	substr := body[indexOfJson+len(jsonVariableForUrls)+1:]
	var res strings.Builder
	for _, v := range substr {
		if v == '"' {
			return res.String()
		}
		res.WriteRune(v)
	}
	return res.String()
}

func getHttpContent(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(response)
}
