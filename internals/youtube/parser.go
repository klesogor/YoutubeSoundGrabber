package youtube

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type StreamData struct {
	clen      int
	ctype     string
	signature string
	bitrate   int
}

type PlayerConfig struct {
	Assets struct {
		CSS string `json:"css"`
		Js  string `json:"js"`
	} `json:"assets"`
	Args struct {
		ApiaryHost                    string `json:"apiary_host"`
		Authuser                      int    `json:"authuser"`
		Cr                            string `json:"cr"`
		UserDisplayImage              string `json:"user_display_image"`
		CsiPageType                   string `json:"csi_page_type"`
		WatchAjaxToken                string `json:"watch_ajax_token"`
		GapiHintParams                string `json:"gapi_hint_params"`
		XhrApiaryHost                 string `json:"xhr_apiary_host"`
		Ps                            string `json:"ps"`
		Hl                            string `json:"hl"`
		EnabledEngageTypes            string `json:"enabled_engage_types"`
		FmtList                       string `json:"fmt_list"`
		URLEncodedFmtStreamMap        string `json:"url_encoded_fmt_stream_map"`
		HostLanguage                  string `json:"host_language"`
		ShowMiniplayerButton          string `json:"show_miniplayer_button"`
		C                             string `json:"c"`
		Fexp                          string `json:"fexp"`
		Cbrver                        string `json:"cbrver"`
		UseMiniplayerUI               string `json:"use_miniplayer_ui"`
		InnertubeAPIKey               string `json:"innertube_api_key"`
		Fflags                        string `json:"fflags"`
		InnertubeAPIVersion           string `json:"innertube_api_version"`
		Timestamp                     string `json:"timestamp"`
		LoaderURL                     string `json:"loaderUrl"`
		Itct                          string `json:"itct"`
		Watermark                     string `json:"watermark"`
		ExternalFullscreen            bool   `json:"external_fullscreen"`
		VideoID                       string `json:"video_id"`
		Author                        string `json:"author"`
		AdaptiveFmts                  string `json:"adaptive_fmts"`
		InnertubeContextClientVersion string `json:"innertube_context_client_version"`
		ShowContentThumbnail          bool   `json:"show_content_thumbnail"`
		Ucid                          string `json:"ucid"`
		VssHost                       string `json:"vss_host"`
		UseFastSizingOnWatchDefault   bool   `json:"use_fast_sizing_on_watch_default"`
		Cos                           string `json:"cos"`
		UserDisplayName               string `json:"user_display_name"`
		ThumbnailURL                  string `json:"thumbnail_url"`
		Enablecsi                     string `json:"enablecsi"`
		AccountPlaybackToken          string `json:"account_playback_token"`
		ExternalPlayVideo             string `json:"external_play_video"`
		LengthSeconds                 string `json:"length_seconds"`
		Ldpj                          string `json:"ldpj"`
		PlayerResponse                string `json:"player_response"`
		Cbr                           string `json:"cbr"`
		Cver                          string `json:"cver"`
		Title                         string `json:"title"`
		ApiaryHostFirstparty          string `json:"apiary_host_firstparty"`
		TransparentBackground         string `json:"transparent_background"`
		Idpj                          string `json:"idpj"`
		Enablejsapi                   string `json:"enablejsapi"`
		Ssl                           string `json:"ssl"`
	} `json:"args"`
	Attrs struct {
		ID string `json:"id"`
	} `json:"attrs"`
	Sts int `json:"sts"`
}

func parseVideoForomHttpResponse(body []byte) (*PlayerConfig, error) {
	bodyString := string(body)
	configJson, err := getConfig(bodyString)
	if err != nil {
		return nil, err
	}
	var config PlayerConfig
	err = json.Unmarshal([]byte(configJson), &config)
	return &config, err
}

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
	return getBetween(body, "ytplayer.config=", ";ytplayer.load=")
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
