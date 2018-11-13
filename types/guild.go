package types

import "github.com/bwmarrin/snowflake"

// Guild represents a guild on Discord
type Guild struct {
	ID                          snowflake.ID  `json:"id"`
	Name                        string        `json:"name"`
	Icon                        string        `json:"icon"`
	Splash                      string        `json:"splash"`
	Owner                       bool          `json:"owner,omitempty"`
	OwnerID                     snowflake.ID  `json:"owner_id"`
	Permissions                 int           `json:"permissions,omitempty"`
	Region                      string        `json:"region"`
	AFKChannelID                snowflake.ID  `json:"afk_channel_id"`
	AFKTimeout                  int           `json:"afk_timeout"`
	EmbedEnabled                bool          `json:"embed_enabled,omitempty"`
	EmbedChannelID              snowflake.ID  `json:"embed_channel_id,omitempty"`
	VerificationLevel           int           `json:"verification_level"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	ExplicitContentFilter       int           `json:"explicit_content_filter"`
	Roles                       []interface{} `json:"roles"`  // TODO: type
	Emojis                      []interface{} `json:"emojis"` // TODO: type
	Features                    []string      `json:"features"`
	MFALevel                    int           `json:"mfa_level"`
	ApplicationID               snowflake.ID  `json:"application_id"`
	WidgetEnabled               bool          `json:"widget_enabled,omitempty"`
	WidgetChannelID             snowflake.ID  `json:"widget_channel_id,omitempty"`
	SystemChannelID             snowflake.ID  `json:"system_channel_id"`
	JoinedAt                    string        `json:"joined_at,omitempty"`
	Large                       bool          `json:"large,omitempty"`
	Unavailable                 bool          `json:"unavailable,omitempty"`
	MemberCount                 int           `json:"member_count,omitempty"`
	VoiceStates                 []interface{} `json:"voice_states,omitempty"` // TODO: type
	Members                     []interface{} `json:"members,omitempty"`      // TODO: type
	Channels                    []Channel     `json:"channels,omitempty"`
	Presences                   []interface{} `json:"presences,omitempty"` // TODO: type
}

// UnavailableGuild represents an unavailable guild
type UnavailableGuild struct {
	ID          snowflake.ID `json:"id"`
	Unavailable bool         `json:"unavailable"`
}

// message notification levels
const (
	MessageNotificationsAllMessages = iota
	MessageNotificationsOnlyMentions
)

// explicit content filter levels
const (
	ExplicitContentFilterDisabled = iota
	ExplicitContentFilterMembersWithoutRoles
	ExplicitContentFilterAllMembers
)

// MFA levels
const (
	MFALevelNone = iota
	MFALevelElevated
)

// verification levels
const (
	VerificationLevelNone = iota
	VerificationLevelLow
	VerificationLevelMedium
	VerificationLevelHigh
	VerificationLevelVeryHigh
)

// GuildCreate represents a guild create packet
type GuildCreate Guild

// GuildUpdate represents a guild update packet
type GuildUpdate Guild

// GuildDelete represents a guild delete packet
type GuildDelete UnavailableGuild

// GuildBanAdd represents a guild ban add packet
type GuildBanAdd struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    User         `json:"user"`
}

// GuildBanRemove represents a guild ban remove packet
type GuildBanRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    User         `json:"user"`
}

// GuildEmojisUpdate represents a guild emojis update packet
type GuildEmojisUpdate struct {
	GuildID snowflake.ID  `json:"guild_id"`
	Emojis  []interface{} `json:"emojis"` // TODO: type
}

// GuildIntegrationsUpdate represents a guild integrations update packet
type GuildIntegrationsUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
}

// GuildMemberAdd represents a guild member add packet
type GuildMemberAdd struct {
	GuildMember
	GuildID snowflake.ID `json:"guild_id"`
}

// GuildMemberRemove represents a guild member remove packet
type GuildMemberRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    User         `json:"user"`
}

// GuildMemberUpdate represents a guild member update packet
type GuildMemberUpdate struct {
	GuildID snowflake.ID   `json:"guild_id"`
	Roles   []snowflake.ID `json:"roles"`
	User    User           `json:"user"`
	Nick    string         `json:"nick"`
}

// GuildMembersChunk represents a guild members chunk packet
type GuildMembersChunk struct {
	GuildID snowflake.ID  `json:"guild_id"`
	Members []GuildMember `json:"members"`
}

// GuildRoleCreate represents a guild role create packet
type GuildRoleCreate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    interface{}  `json:"role"` // TODO: type
}

// GuildRoleUpdate represents a guild role update packet
type GuildRoleUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    interface{}  `json:"role"` // TODO: type
}

// GuildRoleDelete represents a guild role delete packet
type GuildRoleDelete struct {
	GuildID snowflake.ID `json:"guild_id"`
	RoleID  snowflake.ID `json:"role_id"`
}
