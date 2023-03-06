package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	botPkg "github.com/trickybestia/flopper/internal/bot"
	"github.com/trickybestia/flopper/internal/musicplayer"
)

func main() {
	args := parseArgs()
	config := loadConfig(args.Config)

	log.SetFlags(config.GetLoggerFlags())

	bot, err := botPkg.New(config.Token)

	if err != nil {
		log.Fatalln(err)
	}

	bot.Session.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		switch config.StatusType {
		case "none":
			break
		case "listening":
			bot.Session.UpdateListeningStatus(config.Status)
		case "watching":
			bot.Session.UpdateWatchStatus(0, config.Status)
		case "playing":
			bot.Session.UpdateGameStatus(0, config.Status)
		case "streaming":
			bot.Session.UpdateStreamingStatus(0, config.Status, config.StreamingUrl)
		default:
			log.Fatalf("`%s` is invalid value for StatusType", config.StatusType)
		}
	})

	bot.CommandPrefix = config.CommandPrefix
	bot.CommandSuccessReaction = config.CommandSuccessReaction
	bot.CommandFailReaction = config.CommandFailReaction

	musicPlayer := musicplayer.New(bot)

	registerCommand := func(command botPkg.CommandEntry) {
		bot.Commands = append(bot.Commands, command)
	}

	registerCommand(botPkg.CommandEntry{Command: bot.HelpCommand, Aliases: config.HelpCommandAliases,
		Description: "выводит список доступных команд"})
	registerCommand(botPkg.CommandEntry{Command: musicPlayer.PlayCommand, Aliases: config.PlayCommandAliases,
		Description:     "ставит трек в очередь",
		ArgsDescription: "<название или ссылка>"})
	registerCommand(botPkg.CommandEntry{Command: musicPlayer.ShowCommand, Aliases: config.ShowCommandAliases,
		Description: "показывает очередь треков"})
	registerCommand(botPkg.CommandEntry{Command: musicPlayer.ClearCommand, Aliases: config.ClearCommandAliases,
		Description: "очищает очередь треков"})
	registerCommand(botPkg.CommandEntry{Command: musicPlayer.RemoveCommand, Aliases: config.RemoveCommandAliases,
		Description:     "удаляет трек из очереди",
		ArgsDescription: "<номер трека>"})
	registerCommand(botPkg.CommandEntry{Command: musicPlayer.SkipCommand, Aliases: config.SkipCommandAliases,
		Description: "пропускает текущий трек"})
	registerCommand(botPkg.CommandEntry{Command: musicPlayer.PauseCommand, Aliases: config.PauseCommandAliases,
		Description: "ставит трек на паузу"})
	registerCommand(botPkg.CommandEntry{Command: musicPlayer.ResumeCommand, Aliases: config.ResumeCommandAliases,
		Description: "снимает трек с паузы"})

	if err = bot.Run(); err != nil {
		log.Fatalln(err)
	}
}
