package broker

import (
	"context"
	"errors"
	"io"
)

// ErrDisconnected occurs when trying to do something that requires a connection but one was
// unavailable
var ErrDisconnected = errors.New("disconnected from the broker")

type Message interface {
	Event() string
	Body() []byte
	Reply(ctx context.Context, data []byte) error
	Ack(ctx context.Context) error
}

// Broker is an interface describing message brokers
type Broker interface {
	io.Closer
	Connect(ctx context.Context, url string) error
	Publish(ctx context.Context, event string, data []byte) error
	Subscribe(ctx context.Context, events []string, messages chan<- Message) error
	NotifyClose(chan error) error
}
