package gateway

import "errors"

// Errors
var (
	ErrGatewayAbsent           = errors.New("gateway information hasn't been fetched")
	ErrHeartbeatUnacknowledged = errors.New("heartbeat was never acknowledged")
	ErrMaxRetriesExceeded      = errors.New("max retries exceeded")
)
