package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Command func(*discordgo.Session, *discordgo.Message, string) error

type CommandEntry struct {
	Command         Command
	Aliases         []string // first name is considered main
	Description     string
	ArgsDescription string
}

type Bot struct {
	Session                *discordgo.Session
	Commands               []CommandEntry
	CommandPrefix          string
	CommandSuccessReaction string
	CommandFailReaction    string
}

func New(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	bot := Bot{
		Session:  session,
		Commands: make([]CommandEntry, 0),
	}

	session.AddHandler(bot.onMessage)

	return &bot, nil
}

func (bot *Bot) Run() error {
	if err := bot.Session.Open(); err != nil {
		return err
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	bot.Session.Close()

	return nil
}
