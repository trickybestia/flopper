package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	commandLine, err := removePrefix(m.Content, bot.CommandPrefix)

	if err != nil {
		return
	}

	foundCommand, err := bot.findCommand(commandLine)

	if err != nil {
		return
	}

	if err := foundCommand.command(bot.Session, m.Message, foundCommand.args); err != nil {
		log.Printf("Command `%s` failed with error `%s`", foundCommand.alias, err)

		if bot.CommandFailReaction != "" {
			s.MessageReactionAdd(m.ChannelID, m.ID, bot.CommandFailReaction)
		}
	} else if bot.CommandSuccessReaction != "" {
		s.MessageReactionAdd(m.ChannelID, m.ID, bot.CommandSuccessReaction)
	}
}
