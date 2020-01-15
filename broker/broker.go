package broker

import (
	"io"
	"time"
)

// Broker is an interface describing message brokers
type Broker interface {
	io.Closer
	Connect(url string) error
	Publish(event string, data []byte) error
	PublishOptions(PublishOptions) error
	Subscribe(event string) error
	Unsubscribe(event string) error
	SetCallback(handler EventHandler)
	NotifyClose(chan error) error
}

// PublishOptions is the representation for specifying optional properties when publishing
type PublishOptions struct {
	Event   string
	Data    []byte
	Timeout time.Duration
}
