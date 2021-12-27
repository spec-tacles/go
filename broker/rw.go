package broker

import (
	"context"
	"errors"
	"io"

	"github.com/ugorji/go/codec"
)

// RWBroker is a broker that uses a Go Reader and Writer
type RWBroker struct {
	R io.Reader
	W io.Writer
}

var ErrCannotReply = errors.New("cannot reply")

// IOPacket represents a JSON packet transmitted through an RW broker
type IOPacket struct {
	E string      `codec:"event"`
	D interface{} `codec:"data"`
}

func (p *IOPacket) Event() string {
	return p.E
}

func (p *IOPacket) Body() interface{} {
	return p.D
}

func (p *IOPacket) Reply(ctx context.Context, data interface{}) error {
	return ErrCannotReply
}

func (p *IOPacket) Ack(context.Context) error {
	return nil
}

// Close implements io.Closer
func (b *RWBroker) Close() error {
	return nil
}

// Connect implements Broker interface
func (b *RWBroker) Connect(ctx context.Context, url string) error {
	return nil
}

// Publish writes data to the writer
func (b *RWBroker) Publish(ctx context.Context, event string, data interface{}) error {
	return codec.NewEncoder(b.W, &codecHandle).Encode(IOPacket{event, data})
}

// Subscribe implements Broker interface
func (b *RWBroker) Subscribe(ctx context.Context, events []string, messages chan<- Message) (err error) {
	eMap := make(map[string]struct{}, len(events))
	for _, event := range events {
		eMap[event] = struct{}{}
	}

	decoder := codec.NewDecoder(b.R, &codecHandle)
	for {
		pk := &IOPacket{}
		if err = decoder.Decode(pk); err != nil {
			return
		}

		if _, ok := eMap[pk.E]; ok {
			messages <- pk
		}
	}
}
