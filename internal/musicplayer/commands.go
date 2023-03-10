package musicplayer

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/trickybestia/flopper/internal/ytdlp"
)

func (musicPlayer *MusicPlayer) PlayCommand(s *discordgo.Session, m *discordgo.Message, args string) error {
	connection := musicPlayer.getConnection(m.GuildID)

	justConnected := false

	if connection == nil {
		voiceState, err := s.State.VoiceState(m.GuildID, m.Author.ID)

		if err != nil {
			return err
		}

		connection, err = musicPlayer.Connect(voiceState.GuildID, voiceState.ChannelID)

		if err != nil {
			return err
		}

		justConnected = true

		connection.Lock()

		defer connection.Unlock()
	} else if args == "" {
		connection.Lock()

		defer connection.Unlock()

		return connection.Resume()
	}

	url, err := url.Parse(args)

	if err != nil {
		return err
	}

	if url.Scheme == "" {
		args = "ytsearch1:" + args
	}

	info, err := ytdlp.GetInfo(args)

	if err != nil {
		if justConnected {
			connection.Disconnect()
		}

		return err
	}

	connection.Play(info)

	return nil
}

func (musicPlayer *MusicPlayer) SkipCommand(s *discordgo.Session, m *discordgo.Message, args string) error {
	connection := musicPlayer.getConnection(m.GuildID)

	if connection == nil {
		return errors.New("not connected")
	}

	connection.Lock()

	defer connection.Unlock()

	connection.Skip()

	return nil
}

func (musicPlayer *MusicPlayer) RemoveCommand(s *discordgo.Session, m *discordgo.Message, args string) error {
	connection := musicPlayer.getConnection(m.GuildID)

	if connection == nil {
		return errors.New("not connected")
	}

	connection.Lock()

	defer connection.Unlock()

	trackIndex, err := strconv.ParseInt(args, 10, 64)

	if err != nil {
		return err
	}

	trackIndex--

	if trackIndex < 0 || trackIndex >= int64(len(connection.tracks)) {
		return errors.New("invalid track index")
	}

	if trackIndex == 0 {
		connection.Skip()
	} else {
		connection.tracks = append(connection.tracks[:trackIndex], connection.tracks[trackIndex+1:]...)
	}

	return nil
}

func (musicPlayer *MusicPlayer) ShowCommand(s *discordgo.Session, m *discordgo.Message, args string) error {
	connection := musicPlayer.getConnection(m.GuildID)

	message := ""

	if connection == nil {
		message = "???????????? ??????????????"
	} else {
		connection.Lock()

		elapsedTime := connection.playbackController.ElapsedTime()
		message = fmt.Sprintf("1. %s [???????????? ????????????]", InfoToString(connection.tracks[0], 35, &elapsedTime))

		if connection.playbackController.Paused() {
			message += " [??????????]"
		}

		message += "\n"

		for i, track := range connection.tracks[1:] {
			message += fmt.Sprintf("%d. %s\n", i+2, InfoToString(track, 45, nil))
		}

		connection.Unlock()
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{Description: message}})

	return err
}

func (musicPlayer *MusicPlayer) ClearCommand(s *discordgo.Session, m *discordgo.Message, args string) error {
	connection := musicPlayer.getConnection(m.GuildID)

	if connection == nil {
		return errors.New("not connected")
	}

	connection.Lock()

	defer connection.Unlock()

	connection.Disconnect()

	return nil
}

func (musicPlayer *MusicPlayer) PauseCommand(s *discordgo.Session, m *discordgo.Message, args string) error {
	connection := musicPlayer.getConnection(m.GuildID)

	if connection == nil {
		return errors.New("not connected")
	}

	connection.Lock()

	defer connection.Unlock()

	return connection.Pause()
}

func (musicPlayer *MusicPlayer) ResumeCommand(s *discordgo.Session, m *discordgo.Message, args string) error {
	connection := musicPlayer.getConnection(m.GuildID)

	if connection == nil {
		return errors.New("not connected")
	}

	connection.Lock()

	defer connection.Unlock()

	return connection.Resume()
}
