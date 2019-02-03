package types

import (
	"encoding/json"

	"github.com/bwmarrin/snowflake"
)

// Gateway close event codes
const (
	CloseUnknownError = 4000 + iota
	CloseUnknownOpCode
	CloseDecodeError
	CloseNotAuthenticated
	CloseAuthenticationFailed
	CloseAlreadyAuthenticated
	_
	CloseInvalidSeq
	CloseRateLimited
	CloseSessionTimeout
	CloseInvalidShard
	CloseShardingRequired
)

// SendPacket represents a JSON packet sent over the Discord gateway
type SendPacket struct {
	Op   GatewayOp   `json:"op"`
	Data interface{} `json:"d"`
}

// Seq represents the seq of a gateway packet
type Seq uint64

// ReceivePacket represents a JSON packet received over the Discord gateway
type ReceivePacket struct {
	Op    GatewayOp       `json:"op"`
	Data  json.RawMessage `json:"d"`
	Seq   Seq             `json:"s,omitempty"`
	Event GatewayEvent    `json:"t,omitempty"`
}

// Identify represents an identify packet
type Identify struct {
	Token          string              `json:"token"`
	Properties     *IdentifyProperties `json:"properties"`
	Compress       bool                `json:"compress,omitempty"`
	LargeThreshold int                 `json:"large_threshold,omitempty"`
	Shard          []int               `json:"shard,omitempty"`
	Presence       *Activity           `json:"presence,omitempty"`
}

// IdentifyProperties represents properties sent in an identify packet
type IdentifyProperties struct {
	OS      string `json:"$os"`
	Browser string `json:"$browser"`
	Device  string `json:"$device"`
}

// Resume represents a resume packet
type Resume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       Seq    `json:"seq"`
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
	HeartbeatInterval int64    `json:"heartbeat_interval"`
	Trace             []string `json:"_trace"`
}

// Ready represents a ready packet
type Ready struct {
	Version   int           `json:"v"`
	User      interface{}   `json:"user"`   // TODO: type with user
	Guilds    []interface{} `json:"guilds"` // TODO: type with guild
	SessionID string        `json:"session_id"`
	Trace     []string      `json:"_trace"`
}

// Resumed represents a resumed packet
type Resumed struct {
	Trace []string `json:"_trace"`
}
