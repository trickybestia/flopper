package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type config struct {
	Token string

	CommandPrefix          string
	CommandSuccessReaction string
	CommandFailReaction    string

	PlayCommandNames   []string
	ShowCommandNames   []string
	ClearCommandNames  []string
	RemoveCommandNames []string
	SkipCommandNames   []string
	PauseCommandNames  []string
	ResumeCommandNames []string
}

func loadConfig(path string) config {
	var config config

	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
