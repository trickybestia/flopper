package bot

import (
	"errors"
	"strings"
	"unicode/utf8"
)

func removePrefix(s string, prefix string) (string, error) {
	s = strings.TrimSpace(s)

	if !strings.HasPrefix(s, prefix) {
		return "", errors.New("string doesn't have a prefix")
	}

	s = strings.TrimSpace(strings.TrimPrefix(s, prefix))

	return s, nil
}

type commandFindResult struct {
	command Command
	alias   string
	args    string
}

func (bot *Bot) findCommand(commandLine string) (*commandFindResult, error) {
	type BestMatch struct {
		runesInAlias int
		alias        string
		command      Command
	}

	var bestMatch *BestMatch

	for _, v := range bot.Commands {
		for _, alias := range v.Aliases {
			runes := utf8.RuneCountInString(alias)

			if (bestMatch == nil || runes > bestMatch.runesInAlias) && strings.HasPrefix(commandLine, alias) {
				bestMatch = &BestMatch{
					runesInAlias: runes,
					command:      v.Command,
					alias:        alias}
			}
		}
	}

	if bestMatch == nil {
		return nil, errors.New("command not found")
	}

	args := strings.TrimSpace(commandLine[len(bestMatch.alias):])
	result := commandFindResult{
		command: bestMatch.command,
		alias:   bestMatch.alias,
		args:    args}

	return &result, nil
}
