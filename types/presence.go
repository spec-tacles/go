package types

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

// StatusUpdate represents the packet sent by the client to update its status
type StatusUpdate struct {
	Since  *int64    `json:"since"` // Unix timestamp
	Game   *Activity `json:"game"`
	Status string    `json:"status"`
	AFK    bool      `json:"afk"`
}

// PresenceStatus represents a presence's status
type PresenceStatus string

// Presence statuses
const (
	PresenceStatusIdle      PresenceStatus = "idle"
	PresenceStatusDND       PresenceStatus = "dnd"
	PresenceStatusOnline    PresenceStatus = "online"
	PresenceStatusOffline   PresenceStatus = "offline"
	PresenceStatusInvisible PresenceStatus = "invisible"
)

// PresenceUpdate represents a presence update packet
type PresenceUpdate struct {
	User         User           `json:"user"`
	Roles        []snowflake.ID `json:"roles"`
	Game         Activity       `json:"activity"`
	GuildID      snowflake.ID   `json:"guild_id"`
	Status       PresenceStatus `json:"status"`
	Activities   []Activity     `json:"activities"`
	ClientStatus ClientStatus   `json:"client_status"`
	PremiumSince time.Time      `json:"premium_since,omitempty"`
	Nick         string         `json:"nick,omitempty"`
}

// ActivityType represents an activity's type
type ActivityType int

// Activity types
const (
	ActivityTypeGame ActivityType = iota
	ActivityTypeStreaming
	ActivityTypeListening
	ActivityTypeCustom
)

// ActivityFlag represents an activity's flags
type ActivityFlag int

// Activity flags
const (
	ActivityFlagInstance ActivityFlag = 1 << iota
	ActivityFlagJoin
	ActivityFlagSpectate
	ActivityFlagJoinRequest
	ActivityFlagSync
	ActivityFlagPlay
)

// Activity represents an activity as sent as part of other packets
type Activity struct {
	Name          string        `json:"name"`
	Type          ActivityType  `json:"type"`
	URL           string        `json:"url,omitempty"`
	Timestamps    Timestamps    `json:"timestamps,omitempty"`
	ApplicationID snowflake.ID  `json:"application_id"`
	Details       string        `json:"details,omitempty"`
	State         string        `json:"state,omitempty"`
	Emoji         ActivityEmoji `json:"emoji,omitempty"`
	Party         Party         `json:"party,omitempty"`
	Assets        Assets        `json:"assets,omitempty"`
	Secrets       Secrets       `json:"secrets,omitempty"`
	Instance      bool          `json:"instance,omitempty"`
	Flags         ActivityFlag  `json:"flags,omitempty"`
}

// Timestamps represents the starting and ending timestamp of an activity
type Timestamps struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

// Party represents an activity's current party information
type Party struct {
	ID   string `json:"id,omitempty"`
	Size []int  `json:"size,omitempty"`
}

// Assets represents an activity's images and their hover texts
type Assets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

// Secrets represents an activity's secrets for Rich Presence joining and spectating
type Secrets struct {
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
	Match    string `json:"match,omitempty"`
}

// ActivityEmoji represents the emoji shown in a custom status
type ActivityEmoji struct {
	Name     string       `json:"name"`
	ID       snowflake.ID `json:"id,omitempty"`
	Animated bool         `json:"animated,omitempty"`
}

// ClientStatus represents the client's status on each platform
type ClientStatus struct {
	Desktop PresenceStatus `json:"desktop,omitempty"`
	Mobile  PresenceStatus `json:"mobile,omitempty"`
	Web     PresenceStatus `json:"web,omitempty"`
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
