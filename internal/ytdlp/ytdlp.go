package ytdlp

import (
	"encoding/json"
	"os/exec"
)

type Info struct {
	WebpageUrl string `json:"webpage_url"`
	Title      string
	Duration   Duration
	AudioUrl   string `json:"-"`
}

func GetInfo(ytdlpArg string) (*Info, error) {
	cmd := exec.Command("yt-dlp", "-j", "--no-playlist", ytdlpArg)

	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	var info Info

	if err = json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	audioUrl, err := getAudioUrl(info.WebpageUrl)

	if err != nil {
		return nil, err
	}

	info.AudioUrl = *audioUrl

	return &info, nil
}

func getAudioUrl(url string) (*string, error) {
	cmd := exec.Command("yt-dlp", "-x", "--get-url", "--no-playlist", url)

	output, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	result := string(output)

	return &result, nil
}
