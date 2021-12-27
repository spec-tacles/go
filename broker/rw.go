package broker

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
)

// RWBroker is a broker that uses a Go Reader and Writer
type RWBroker struct {
	R io.Reader
	W io.Writer
}

var ErrCannotReply = errors.New("cannot reply")

// IOPacket represents a JSON packet transmitted through an RW broker
type IOPacket struct {
	E string `json:"event"`
	D []byte `json:"data"`
}

func (p *IOPacket) Event() string {
	return p.E
}

func (p *IOPacket) Body() []byte {
	return p.D
}

func (p *IOPacket) Reply(ctx context.Context, data []byte) error {
	return ErrCannotReply
}

func (p *IOPacket) Ack(context.Context) error {
	return nil
}

// NewRW creates a new Read/Write broker
func NewRW(r io.Reader, w io.Writer) *RWBroker {
	b := &RWBroker{
		R: r,
		W: w,
	}

	return b
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
func (b *RWBroker) Publish(ctx context.Context, event string, data []byte) (err error) {
	pk, err := json.Marshal(&IOPacket{event, data})
	if err != nil {
		return
	}

	_, err = b.W.Write(pk)
	return
}

// Subscribe implements Broker interface
func (b *RWBroker) Subscribe(ctx context.Context, events []string, messages chan<- Message) error {
	eMap := make(map[string]struct{}, len(events))
	for _, event := range events {
		eMap[event] = struct{}{}
	}

	decoder := json.NewDecoder(b.R)
	pk := &IOPacket{}
	for {
		if err := decoder.Decode(pk); err != nil {
			if err != io.EOF {
				log.Printf("error decoding JSON: %s", err)
			}
			break
		}

		if _, ok := eMap[pk.E]; ok {
			messages <- pk
		}
	}

	return nil
}
