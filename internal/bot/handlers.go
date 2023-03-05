package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if command := tryRemovePrefix(m.Content, bot.CommandPrefix); command != nil {
		if match := tryFindCommand(*command, bot.Commands); match != nil {
			command, args := match.command, match.args

			if err := command(bot.Session, m.Message, args); err != nil {
				log.Printf("Command `%s` failed with error `%s`", match.prefix, err)

				if bot.CommandFailReaction != "" {
					s.MessageReactionAdd(m.ChannelID, m.ID, bot.CommandFailReaction)
				}
			} else if bot.CommandSuccessReaction != "" {
				s.MessageReactionAdd(m.ChannelID, m.ID, bot.CommandSuccessReaction)
			}
		}
	}
}
