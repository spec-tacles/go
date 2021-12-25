package broker

import "context"

type Message interface {
	Body() []byte
	Reply(ctx context.Context, data []byte) error
	Ack(ctx context.Context) error
}

// EventHandler represents a function that handles an event
type EventHandler = func(string, Message)
