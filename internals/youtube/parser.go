package youtube

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/klesogor/youtube-grabber/internals"
)

func (player *PlayerConfig) DownloadAudio() (StreamData, error) {
	audio, err := player.getBestAudioStream()
	if err != nil {
		return StreamData{}, err
	}
	res, err := internals.DownloadSegmented(audio, internals.SegmentedDownloadConfig{ChunkSize: 989898, ConcurrencyLimit: 10}, Downloader)
	if err != nil {
		return StreamData{}, err
	}

	audio.data = res

	return audio, nil
}

type StreamData struct {
	Clen      int
	Url       string
	Ctype     string
	Signature string
	Bitrate   int
	data      []byte
}

func (s StreamData) GetDownloadUrl() string {
	if s.Signature != "" {
		return s.Url + "&signature="
	}

	return s.Url
}

func (s StreamData) GetContentLengh() int {
	return s.Clen
}

func (s StreamData) GetContentType() string {
	return s.Ctype
}

func (s StreamData) GetData() []byte {
	return s.data
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
	Sts        int `json:"sts"`
	streamData []StreamData
}

func (c *PlayerConfig) getBestAudioStream() (StreamData, error) {
	if c.streamData == nil {
		data, err := c.getStreamData()
		if err != nil {
			return StreamData{}, err
		}
		c.streamData = data
	}
	sort.Slice(c.streamData, func(i, j int) bool {
		return getAudioScore(c.streamData[i]) > getAudioScore(c.streamData[j])
	})

	return c.streamData[0], nil
}

func (c *PlayerConfig) getStreamData() ([]StreamData, error) {
	if c.streamData == nil {
		data, err := parseStreamData(c.Args.AdaptiveFmts)
		if err != nil {
			return nil, err
		}
		c.streamData = data
	}
	return c.streamData, nil
}

func parseStreamDataFromUrlString(s string) StreamData {
	res := make(map[string]string)
	splitted := strings.Split(s, "\u0026")
	for _, v := range splitted {
		pair := strings.SplitN(v, "=", 2)
		if len(pair) >= 2 {
			res[pair[0]] = pair[1]
		}
	}
	decodedUrl, _ := url.QueryUnescape(res["url"])

	return StreamData{
		Clen:      getIntOrDefault(res["clen"]),
		Bitrate:   getIntOrDefault(res["bitrate"]),
		Signature: res["signature"],
		Url:       decodedUrl,
		Ctype:     res["type"]}
}

func GetPlayerConfig(url string) (*PlayerConfig, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := parseBody(res)
	if err != nil {
		return nil, err
	}
	return parseVideoForomHttpResponse(body)
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
