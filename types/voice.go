package types

import "github.com/bwmarrin/snowflake"

// VoiceState represents a voice state on Discord
type VoiceState struct {
	GuildID   snowflake.ID `json:"guild_id,omitempty"`
	ChannelID snowflake.ID `json:"channel_id"`
	UserID    snowflake.ID `json:"user_id"`
	Member    GuildMember  `json:"member,omitempty"`
	SessionID string       `json:"session_id"`
	Deaf      bool         `json:"deaf"`
	Mute      bool         `json:"mute"`
	SelfDeaf  bool         `json:"self_deaf"`
	SelfMute  bool         `json:"self_mute"`
	Suppress  bool         `json:"suppress"`
}

// VoiceStateUpdate represents a voice state update packet
type VoiceStateUpdate VoiceState

// VoiceServerUpdate represents a voice server update packet
type VoiceServerUpdate struct {
	Token    string       `json:"token"`
	GuildID  snowflake.ID `json:"guild_id"`
	Endpoint string       `json:"endpoint"`
}
