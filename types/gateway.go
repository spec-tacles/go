package types

// Gateway represents a GET /gateway response
type Gateway struct {
	URL string `json:"url"`
}

// GatewayBot represents a GET /gateway/bot response
type GatewayBot struct {
	URL               string `json:"url"`
	Shards            int    `json:"shards"`
	SessionStartLimit struct {
		Total      int `json:"total"`
		Remaining  int `json:"remaining"`
		ResetAfter int `json:"reset_after"`
	} `json:"session_start_limit"`
}
