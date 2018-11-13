package types

import "github.com/bwmarrin/snowflake"

// Payload represents a JSON payload sent over the Discord gateway
type Payload struct {
	OP int         `json:"op"`
	D  interface{} `json:"d"`
	S  int         `json:"s"`
	T  string      `json:"t"`
}

// Identify represents an identify packet
type Identify struct {
	Token      string `json:"token"`
	Properties struct {
		OS      string `json:"$os"`
		Browser string `json:"$browser"`
		Device  string `json:"$device"`
	} `json:"properties"`
	Compress       bool        `json:"compress,omitempty"`
	LargeThreshold int         `json:"large_threshold,omitempty"`
	Shard          []int       `json:"shard,omitempty"`
	Presence       interface{} `json:"presence,omitempty"`
}

// Resume represents a resume packet
type Resume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}

// Heartbeat presents a heartbeat packet
type Heartbeat int

// RequestGuildMembers represents a request guild members packet
type RequestGuildMembers struct {
	GuildID snowflake.ID `json:"guild_id"`
	Query   string       `json:"query"`
	Limit   int          `json:"limit"`
}

// UpdateVoiceState represents an update voice state packet
type UpdateVoiceState struct {
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
	SelfMute  bool         `json:"self_mute"`
	SelfDeaf  bool         `json:"self_deaf"`
}

// available statuses
const (
	StatusOnline    = "online"
	StatusDND       = "dnd"
	StatusIdle      = "idle"
	StatusInvisible = "invisible"
	StatusOffline   = "offline"
)

// UpdateStatus represents an update status packet
type UpdateStatus struct {
	Since  int         `json:"since,omitempty"`
	Game   interface{} `json:"game,omitempty"`
	Status string      `json:"status"`
	AFK    bool        `json:"afk"`
}

// Hello represents a hello packet
type Hello struct {
	HeartbeatInterval int      `json:"heartbeat_interval"`
	Trace             []string `json:"_trace"`
}

// Ready represents a ready packet
type Ready struct {
	V         int           `json:"v"`
	User      interface{}   `json:"user"`   // TODO: type with user
	Guilds    []interface{} `json:"guilds"` // TODO: type with guild
	SessionID string        `json:"session_id"`
	Track     []string      `json:"_trace"`
}

// Resumed represents a resumed packet
type Resumed struct {
	Trace []string `json:"_trace"`
}

// InvalidSession represents an invalid session packet
type InvalidSession bool
