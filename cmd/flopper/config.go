package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type config struct {
	Token string

	Status       string
	StatusType   string
	StreamingUrl string

	CommandPrefix          string
	CommandSuccessReaction string
	CommandFailReaction    string

	LogTimeFormat string

	HelpCommandAliases   []string
	PlayCommandAliases   []string
	ShowCommandAliases   []string
	ClearCommandAliases  []string
	RemoveCommandAliases []string
	SkipCommandAliases   []string
	PauseCommandAliases  []string
	ResumeCommandAliases []string
}

func (config *config) GetLoggerFlags() int {
	switch config.LogTimeFormat {
	case "local":
		return log.Ldate | log.Ltime
	case "utc":
		return log.Ldate | log.Ltime | log.LUTC
	case "none":
		return 0
	default:
		log.Fatalf("`%s` is invalid value for LogTimeFormat", config.LogTimeFormat)
	}

	return 0
}

func loadConfig(path string) config {
	var config config

	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
