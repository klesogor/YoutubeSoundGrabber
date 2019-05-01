package grabber

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	jsonVariableForUrls = "\"adaptive_fmts\":"
	delimiter           = "\\u0026"
	otherDelims         = ";,"
)

func GrabDownloadUrl(videoUrl string) string {
	urls := getUrlEncodedAudioUrl(getHttpContent(videoUrl))
	return getDownloadUrlFromParams(splitUrlStrings(urlDecode(urls)))
}

func getDownloadUrlFromParams(params []string) string {
	url := ""
	for _, v := range params {
		if strings.Index(v, "url=") == 0 {
			url = v[4:]
		}
	}

	return url
}

func splitUrlStrings(s string) []string {
	res := make([]string, 0, 20)
	splitted := strings.Split(s, delimiter)
	for _, v := range splitted {
		splittedParam := strings.FieldsFunc(v, func(r rune) bool {
			return strings.Index(otherDelims, string(r)) != -1
		})
		res = append(res, splittedParam...)
	}

	return res
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
