package broker

import "io"

// Broker is an interface describing message brokers
type Broker interface {
	io.Closer
	Connect(url string) error
	Publish(event string, data []byte) error
	Subscribe(event string) error
	Unsubscribe(event string) error
	SetCallback(handler EventHandler)
}
