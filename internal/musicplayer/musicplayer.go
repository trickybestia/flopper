package musicplayer

import (
	"errors"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/trickybestia/flopper/internal/bot"
	"github.com/trickybestia/flopper/internal/dgvoice"
	"github.com/trickybestia/flopper/internal/ytdlp"
)

type MusicPlayerVoiceConnection struct {
	sync.Mutex
	musicPlayer        *MusicPlayer
	voiceConnection    *discordgo.VoiceConnection
	playbackController *dgvoice.PlaybackController
	tracks             []*ytdlp.Info
}

type MusicPlayer struct {
	sync.Mutex
	bot         *bot.Bot
	connections map[string]*MusicPlayerVoiceConnection
}

func New(bot *bot.Bot) *MusicPlayer {
	musicPlayer := &MusicPlayer{
		bot:         bot,
		connections: make(map[string]*MusicPlayerVoiceConnection),
	}

	bot.Session.AddHandler(musicPlayer.onVoiceStateUpdate)

	return musicPlayer
}

func (musicPlayer *MusicPlayer) Connect(guildID string, channelID string) (*MusicPlayerVoiceConnection, error) {
	if musicPlayer.getConnection(guildID) != nil {
		return nil, errors.New("already connected")
	}

	voiceConnection, err := musicPlayer.bot.Session.ChannelVoiceJoin(guildID, channelID, false, true)

	if err != nil {
		return nil, err
	}

	connection := &MusicPlayerVoiceConnection{
		musicPlayer:        musicPlayer,
		voiceConnection:    voiceConnection,
		playbackController: dgvoice.NewPlaybackController(),
	}

	musicPlayer.connections[guildID] = connection

	return connection, nil
}

func (musicPlayer *MusicPlayer) getConnection(guildID string) *MusicPlayerVoiceConnection {
	musicPlayer.Lock()

	defer musicPlayer.Unlock()

	if connection, ok := musicPlayer.connections[guildID]; ok {
		return connection
	}

	return nil
}

func (connection *MusicPlayerVoiceConnection) Play(info *ytdlp.Info) {
	connection.tracks = append(connection.tracks, info)

	if len(connection.tracks) == 1 {
		go connection.playNextTrack()
	}
}

func (connection *MusicPlayerVoiceConnection) Pause() error {
	return connection.playbackController.Pause()
}

func (connection *MusicPlayerVoiceConnection) Resume() error {
	return connection.playbackController.Resume()
}

func (connection *MusicPlayerVoiceConnection) Skip() {
	connection.playbackController.Skip()
}

func (connection *MusicPlayerVoiceConnection) Disconnect() {
	if len(connection.tracks) == 0 {
		connection.disconnectInternal()

		return
	}

	connection.tracks = connection.tracks[:1]

	connection.playbackController.Skip()
}

func (connection *MusicPlayerVoiceConnection) disconnectInternal() {
	connection.voiceConnection.Disconnect()

	connection.musicPlayer.Lock()

	defer connection.musicPlayer.Unlock()

	delete(connection.musicPlayer.connections, connection.voiceConnection.GuildID)
}

func (connection *MusicPlayerVoiceConnection) playNextTrack() {
	dgvoice.PlayAudio(connection.voiceConnection, connection.tracks[0].AudioUrl, connection.playbackController)

	connection.onTrackEnd()
}

func (connection *MusicPlayerVoiceConnection) onTrackEnd() {
	connection.Lock()

	defer connection.Unlock()

	if len(connection.tracks) == 1 {
		connection.disconnectInternal()

		return
	}

	connection.tracks = connection.tracks[1:]

	connection.playbackController.Reset()

	go connection.playNextTrack()
}
