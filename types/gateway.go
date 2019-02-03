package types

// Gateway represents a GET /gateway response
type Gateway struct {
	URL string `json:"url"`
}

// GatewayBot represents a GET /gateway/bot response
type GatewayBot struct {
	URL               string            `json:"url"`
	Shards            int               `json:"shards"`
	SessionStartLimit SessionStartLimit `json:"session_start_limit"`
}

// SessionStartLimit represents a GatewayBot's session start limit
type SessionStartLimit struct {
	Total      int `json:"total"`
	Remaining  int `json:"remaining"`
	ResetAfter int `json:"reset_after"`
}

// GatewayOp represents a gateway packet's operation code
type GatewayOp uint8

// Operation codes
const (
	GatewayOpDispatch GatewayOp = iota
	GatewayOpHeartbeat
	GatewayOpIdentify
	GatewayOpStatusUpdate
	GatewayOpVoiceStateUpdate
	_
	GatewayOpResume
	GatewayOpReconnect
	GatewayOpRequestGuildMembers
	GatewayOpInvalidSession
	GatewayOpHello
	GatewayOpHeartbeatACK
)

// GatewayEvent represents a gateway packet's event name
type GatewayEvent string

// Gateway events
const (
	GatewayEventNone    GatewayEvent = ""
	GatewayEventReady   GatewayEvent = "READY"
	GatewayEventResumed GatewayEvent = "RESUMED"
)
