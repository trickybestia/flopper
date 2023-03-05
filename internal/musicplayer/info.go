package musicplayer

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/davidscholberg/go-durationfmt"
	"github.com/trickybestia/flopper/internal/ytdlp"
)

func InfoToString(info *ytdlp.Info, maxTitleLength int, elapsedTime *time.Duration) string {
	title := ""

	if utf8.RuneCountInString(info.Title) >= maxTitleLength {
		i := 0

		for _, char := range info.Title {
			if i == maxTitleLength {
				break
			}

			title += string(char)

			i += 1
		}

		title += "…"
	} else {
		title += info.Title
	}

	formattedDuration := ""

	duration := time.Duration(info.Duration)

	if duration.Hours() >= 1 {
		if elapsedTime == nil {
			formattedDuration, _ = durationfmt.Format(duration, " [%0h:%0m:%0s]")
		} else {
			formattedDuration, _ = durationfmt.Format(*elapsedTime, " [%0h:%0m:%0s / ")
			formattedElapsedTime, _ := durationfmt.Format(duration, "%0h:%0m:%0s]")

			formattedDuration += formattedElapsedTime
		}
	} else {
		if elapsedTime == nil {
			formattedDuration, _ = durationfmt.Format(duration, " [%0m:%0s]")
		} else {
			formattedDuration, _ = durationfmt.Format(*elapsedTime, " [%0m:%0s / ")
			formattedElapsedTime, _ := durationfmt.Format(duration, "%0m:%0s]")

			formattedDuration += formattedElapsedTime
		}
	}

	result := fmt.Sprintf("[%s](%s) %s", title, info.WebpageUrl, formattedDuration)

	return result
}
