package ytsubtitles

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Track struct {
	BaseURL string `json:"baseUrl"`
	Name    struct {
		SimpleText string `json:"simpleText"`
	} `json:"name"`
	VssID          string `json:"vssId"`
	LanguageCode   string `json:"languageCode"`
	Kind           string `json:"kind"`
	IsTranslatable bool   `json:"isTranslatable"`
}

type YTS struct {
	VideoID      string
	Language     string
	Tracks       map[string]Track
	subtitlesXML struct {
		Text []struct {
			Text  string `xml:",chardata"`
			Start string `xml:"start,attr"`
			Dur   string `xml:"dur,attr"`
		} `xml:"text"`
	}
}

// GetLanguage choice language os available
func (yts *YTS) GetLanguage(lang string) error {
	if lang == "" {
		for _, v := range yts.Tracks {
			lang = v.Name.SimpleText //random
			break
		}
	} else if _, ok := yts.Tracks[lang]; !ok {
		return errors.New(
			fmt.Sprintf("Video %s does't have lang \"%s\"", yts.VideoID, lang))
	}
	yts.Language = yts.Tracks[lang].Name.SimpleText
	resp, err := http.Get(yts.Tracks[lang].BaseURL)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	if err != nil {
		return err
	}
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(bodyByte, &yts.subtitlesXML)
	return err
}

// GetPlainText get subtitles into only plain text, without timestamps
func (yts *YTS) GetPlainText() (text string, err error) {
	if len(yts.subtitlesXML.Text) == 0 {
		return text, errors.New(fmt.Sprintf("subtitles of %s %s is emptty", yts.VideoID, yts.Language))
	}
	for _, v := range yts.subtitlesXML.Text {
		text += fmt.Sprintf("%s\n", v.Text)
	}
	return strings.TrimRight(text, "\n"), err
}

// GetJson subtitles in string JSON
func (yts *YTS) GetJson() (res []byte, err error) {
	if res, err := json.Marshal(yts.subtitlesXML); err != nil {
		return res, err
	}
	return res, err
}

// GetJsonPretty return subtitles in string JSON with 4 spaces for better humans reading
// pretty style JSON print
func (yts *YTS) GetJsonPretty() (res []byte, err error) {
	res, err = json.MarshalIndent(yts.subtitlesXML.Text, "", "    ")
	if err != nil {
		return res, err
	}
	return res, err
}

// Init check on YouTube all available subtitles and languages
func Init(ID string) (yts *YTS, err error) {
	yts = new(YTS)
	regular, err := regexp.Compile(`([a-zA-Z0-9-_]{11})`)
	if err != nil {
		return yts, err
	}
	if regular.Match([]byte(ID)) != true {
		return yts, errors.New(fmt.Sprintf("error parse args requestToYouTybe: \"%s\"", ID))
	}
	yts.VideoID = string(regular.Find([]byte(ID)))

	data, err := requestYT(yts.VideoID)

	if err != nil {
		return yts, err
	}

	regular, err = regexp.Compile(`("captionTracks":.*isTranslatable":(true|false)}])`)
	if err != nil {
		return yts, err
	}
	if regular.Match(data) != true {
		return yts, errors.New(fmt.Sprintf("captions not found on video: \"%s\"", yts.VideoID))
	}

	jsonByteArray := append([]byte{123}, regular.Find(data)...) //add {
	jsonByteArray = append(jsonByteArray, 125)                  //add }

	aux := struct {
		CaptionTracks []struct {
			BaseURL string `json:"baseUrl"`
			Name    struct {
				SimpleText string `json:"simpleText"`
			} `json:"name"`
			VssID          string `json:"vssId"`
			LanguageCode   string `json:"languageCode"`
			Kind           string `json:"kind"`
			IsTranslatable bool   `json:"isTranslatable"`
		} `json:"captionTracks"`
	}{}

	if err = json.Unmarshal(jsonByteArray, &aux); err != nil {
		return yts, err
	}

	yts.Tracks = make(map[string]Track)
	for _, v := range aux.CaptionTracks {
		yts.Tracks[v.Name.SimpleText] = v
	}
	return yts, err
}
func requestYT(ID string) ([]byte, error) {
	res, err := http.Get(fmt.Sprintf("https://youtube.com/watch?v=%s", ID))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf("http StatusCode is %d", res.StatusCode)
		return nil, errors.New(errStr)
	}
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return data, err
	}
	return data, err
}
