package musicplayer

import (
	"github.com/bwmarrin/discordgo"
)

func (musicPlayer *MusicPlayer) onVoiceStateUpdate(s *discordgo.Session, state *discordgo.VoiceStateUpdate) {
	if state.ChannelID == "" { // Disconnected
		connection := musicPlayer.getConnection(state.GuildID)

		if connection == nil {
			return
		}

		connection.Lock()
		defer connection.Unlock()

		if len(connection.tracks) == 0 {
			return
		}

		connection.Disconnect()
	}
}
