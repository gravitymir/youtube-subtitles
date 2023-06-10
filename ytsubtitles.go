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

// YTS Subtitles is slice of available subtitles
type YTS struct {
	VideoID   string
	Language  string
	Tracks    map[string]Track
	Subtitles struct {
		Text []struct {
			Text  string `xml:",chardata"`
			Start string `xml:"start,attr"`
			Dur   string `xml:"dur,attr"`
		} `xml:"text"`
	}
	data []byte
}

func (yts *YTS) getLanguage(lang string) error {
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
	err = xml.Unmarshal(bodyByte, &yts.Subtitles)
	return err
}

// PlainText get subtitles into only plain text, without timestamps
func (yts *YTS) PlainText(lang string) (text string, err error) {
	if err = yts.getLanguage(lang); err != nil {
		return text, err
	}
	for _, v := range yts.Subtitles.Text {
		text += fmt.Sprintf("%s\n", v.Text)
	}
	return strings.TrimRight(text, "\n"), err
}

// Json subtitles in string JSON
func (yts *YTS) Json(lang string) (sbt []byte, err error) {
	if err = yts.getLanguage(lang); err != nil {
		return sbt, err
	}
	if res, err := json.Marshal(yts.Subtitles); err != nil {
		return res, err
	}
	return sbt, err
}

// JsonPretty return subtitles in string JSON with 4 spaces for better humans reading
// pretty style JSON print
func (yts *YTS) JsonPretty(lang string) (sbt []byte, err error) {
	if err = yts.getLanguage(lang); err != nil {
		return sbt, err
	}
	sbt, err = json.MarshalIndent(yts.Subtitles.Text, "", "    ")
	if err != nil {
		return sbt, err
	}
	return sbt, err
}

// Get check on YouTube all available subtitles and languages
// you can use types of videoID
// 1 https://www.youtube.com/watch?v=videoID
// 2 www.youtube.com/watch?v=videoID
// 3 youtube.com/watch?v=videoID
// 4 https://www.youtube.com/watch?v=videoID&t=215s
// 5 https://youtu.be/videoID?t=215
// 6 videoID
// any string with videoID
func Get(ID string) (yts *YTS, err error) {
	yts = new(YTS)
	regular, err := regexp.Compile(`([a-zA-Z0-9-_]{11})`)
	if err != nil {
		return yts, err
	}
	if regular.Match([]byte(ID)) != true {
		return yts, errors.New(fmt.Sprintf("error parse args requestToYouTybe: \"%s\"", ID))
	}
	yts.VideoID = string(regular.Find([]byte(ID)))

	//yts.data add
	if err = yts.requestYT(yts.VideoID); err != nil {
		return yts, err
	}

	regular, err = regexp.Compile(`("captionTracks":.*isTranslatable":(true|false)}])`)
	if err != nil {
		return yts, err
	}
	if regular.Match(yts.data) != true {
		return yts, errors.New(fmt.Sprintf("captions not found on video: \"%s\"", yts.VideoID))
	}

	yts.data = append([]byte{123}, regular.Find(yts.data)...) //add {
	yts.data = append(yts.data, 125)                          //add }

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

	if err = json.Unmarshal(yts.data, &aux); err != nil {
		return yts, err
	}
	yts.data = nil
	yts.Tracks = make(map[string]Track)
	for _, v := range aux.CaptionTracks {
		yts.Tracks[v.Name.SimpleText] = v
	}
	return yts, err
}
func (yts *YTS) requestYT(ID string) error {
	res, err := http.Get(fmt.Sprintf("https://youtube.com/watch?v=%s", ID))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf("http StatusCode is %d", res.StatusCode)
		return errors.New(errStr)
	}
	if yts.data, err = io.ReadAll(res.Body); err != nil {
		return err
	}
	return err
}
