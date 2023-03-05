package bot

import (
	"strings"
	"unicode/utf8"
)

func tryRemovePrefix(s string, prefix string) *string {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, prefix) {
		s = strings.TrimSpace(strings.TrimPrefix(s, prefix))

		return &s
	} else {
		return nil
	}
}

type commandFindResult struct {
	command Command
	prefix  string
	args    string
}

func tryFindCommand(command string, commands map[string]Command) *commandFindResult {
	type BestMatch struct {
		runes  int
		prefix string
		Command
	}

	var bestMatch *BestMatch

	for k, v := range commands {
		runes := utf8.RuneCountInString(k)

		if (bestMatch == nil || runes > bestMatch.runes) && strings.HasPrefix(command, k) {
			bestMatch = &BestMatch{
				runes:   runes,
				Command: v,
				prefix:  k}
		}
	}

	if bestMatch != nil {
		args := strings.TrimSpace(command[len(bestMatch.prefix):])
		return &commandFindResult{
			command: bestMatch.Command,
			prefix:  bestMatch.prefix,
			args:    args}
	} else {
		return nil
	}
}
