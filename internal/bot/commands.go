package bot

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) HelpCommand(s *discordgo.Session, m *discordgo.Message, args string) error {
	getHelpMessage := func(commandEntry CommandEntry, alias string) string {
		helpMessage := fmt.Sprintf("%s%s ", bot.CommandPrefix, alias)

		if commandEntry.Description != "" {
			if commandEntry.ArgsDescription != "" {
				helpMessage += commandEntry.ArgsDescription + " "
			}

			helpMessage += commandEntry.Description
		} else {
			helpMessage += "описание отсутствует"
		}

		return helpMessage
	}

	response := ""

	if args == "" {
		for _, commandEntry := range bot.Commands {
			response += getHelpMessage(commandEntry, commandEntry.Aliases[0]) + "\n"
		}
	} else {
		for _, commandEntry := range bot.Commands {
			for _, alias := range commandEntry.Aliases {
				if alias == args {
					response += getHelpMessage(commandEntry, alias)

					goto EXIT_FOR
				}
			}
		}

		return errors.New("command does not exist")

	EXIT_FOR:
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{Description: response}})

	return err
}
