package main

import (
	"log"

	"github.com/trickybestia/flopper/internal/bot"
	"github.com/trickybestia/flopper/internal/musicplayer"
)

func main() {
	args := parseArgs()
	config := loadConfig(args.Config)

	log.SetFlags(config.GetLoggerFlags())

	bot, err := bot.New(config.Token)

	if err != nil {
		log.Fatalln(err)
	}

	bot.CommandPrefix = config.CommandPrefix
	bot.CommandSuccessReaction = config.CommandSuccessReaction
	bot.CommandFailReaction = config.CommandFailReaction

	musicPlayer := musicplayer.New(bot)

	bot.RegisterCommand(musicPlayer.PlayCommand, config.PlayCommandNames)
	bot.RegisterCommand(musicPlayer.ShowCommand, config.ShowCommandNames)
	bot.RegisterCommand(musicPlayer.ClearCommand, config.ClearCommandNames)
	bot.RegisterCommand(musicPlayer.RemoveCommand, config.RemoveCommandNames)
	bot.RegisterCommand(musicPlayer.SkipCommand, config.SkipCommandNames)
	bot.RegisterCommand(musicPlayer.PauseCommand, config.PauseCommandNames)
	bot.RegisterCommand(musicPlayer.ResumeCommand, config.RemoveCommandNames)

	if err = bot.Run(); err != nil {
		log.Fatalln(err)
	}
}
