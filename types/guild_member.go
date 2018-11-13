package types

import "github.com/bwmarrin/snowflake"

// GuildMember represents a guild member on Discord
type GuildMember struct {
	User     User           `json:"user"`
	Nick     string         `json:"nick,omitempty"`
	Roles    []snowflake.ID `json:"roles"`
	JoinedAt string         `json:"joined_at"`
	Deaf     bool           `json:"deaf"`
	Mute     bool           `json:"mute"`
}
