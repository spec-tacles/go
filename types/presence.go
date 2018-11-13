package types

import "github.com/bwmarrin/snowflake"

// PresenceUpdate represents a presence update packet
type PresenceUpdate struct {
	User       User           `json:"user"`
	Roles      []snowflake.ID `json:"roles"`
	Game       Activity       `json:"activity"`
	GuildID    snowflake.ID   `json:"guild_id"`
	Status     string         `json:"status"`
	Activities []Activity     `json:"activities"`
}

// presence statuses
const (
	PresenceStatusIdle    = "idle"
	PresenceStatusDND     = "dnd"
	PresenceStatusOnline  = "online"
	PresenceStatusOffline = "offline"
)

// activity types
const (
	ActivityTypeGame = iota
	ActivityTypeStreaming
	ActivityTypeListening
)

// activity flags
const (
	ActivityFlagInstance = 1 << iota
	ActivityFlagJoin
	ActivityFlagSpectate
	ActivityFlagJoinRequest
	ActivityFlagSync
	ActivityFlagPlay
)

// Activity represents an activity as sent as part of other packets
type Activity struct {
	Name       string `json:"name"`
	Type       int    `json:"type"`
	URL        string `json:"url,omitempty"`
	Timestamps struct {
		Start int `json:"start,omitempty"`
		End   int `json:"end,omitempty"`
	} `json:"timestamps,omitempty"`
	ApplicationID snowflake.ID `json:"application_id"`
	Details       string       `json:"details,omitempty"`
	State         string       `json:"state,omitempty"`
	Party         struct {
		ID   string `json:"id,omitempty"`
		Size []int  `json:"size,omitempty"`
	} `json:"party,omitempty"`
	Assets struct {
		LargeImage string `json:"large_image,omitempty"`
		LargeText  string `json:"large_text,omitempty"`
		SmallImage string `json:"small_image,omitempty"`
		SmallText  string `json:"small_text,omitempty"`
	} `json:"assets,omitempty"`
	Secrets struct {
		Join     string `json:"join,omitempty"`
		Spectate string `json:"spectate,omitempty"`
		Match    string `json:"match,omitempty"`
	} `json:"secrets,omitempty"`
	Instance bool `json:"instance,omitempty"`
	Flags    int  `json:"flags,omitempty"`
}

// TypingStart represents a typing start packet
type TypingStart struct {
	ChannelID snowflake.ID `json:"channel_id"`
	GuildID   snowflake.ID `json:"guild_id,omitempty"`
	UserID    snowflake.ID `json:"user_id"`
	Timestamp int          `json:"timestamp"`
}

// UserUpdate represents a user update packet
type UserUpdate User
