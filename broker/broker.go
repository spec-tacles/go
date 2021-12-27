package broker

import (
	"context"
	"errors"
)

// ErrDisconnected occurs when trying to do something that requires a connection but one was
// unavailable
var ErrDisconnected = errors.New("disconnected from the broker")

type Message interface {
	Event() string
	Body() interface{}
	Reply(ctx context.Context, data interface{}) error
	Ack(ctx context.Context) error
}

// Broker is an interface describing message brokers
type Broker interface {
	Publish(ctx context.Context, event string, data interface{}) error
	Subscribe(ctx context.Context, events []string, messages chan<- Message) error
}
