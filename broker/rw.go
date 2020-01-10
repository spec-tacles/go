package broker

import (
	"encoding/json"
	"log"
	"io"
	"sync"
)

// RWBroker is a broker that uses a Go Reader and Writer
type RWBroker struct {
	R io.Reader
	W io.Writer

	callback  EventHandler
	events    map[string]struct{}
	eventsMux sync.RWMutex
}

// IOPacket represents a JSON packet transmitted through an RW broker
type IOPacket struct {
	Event string `json:"event"`
	Data  []byte `json:"data"`
}

// NewRW creates a new Read/Write broker
func NewRW(r io.Reader, w io.Writer, callback EventHandler) *RWBroker {
	b := &RWBroker{
		R: r,
		W: w,

		callback:  callback,
		events:    make(map[string]struct{}),
		eventsMux: sync.RWMutex{},
	}

	go func() {
		decoder := json.NewDecoder(r)
		pk := &IOPacket{}
		for {
			if err := decoder.Decode(pk); err != nil {
				log.Printf("error decoding JSON: %s", err)
				break
			}

			if b.callback == nil {
				continue
			}

			b.eventsMux.RLock()
			defer b.eventsMux.RUnlock()

			if _, ok := b.events[pk.Event]; ok {
				go b.callback(pk.Event, pk.Data)
			}
		}
	}()

	return b
}

// Close implements io.Closer
func (b *RWBroker) Close() error {
	return nil
}

// Connect implements Broker interface
func (b *RWBroker) Connect(url string) error {
	return nil
}

// Publish writes data to the writer
func (b *RWBroker) Publish(event string, data []byte) (err error) {
	pk, err := json.Marshal(&IOPacket{event, data})
	if err != nil {
		return
	}

	_, err = b.W.Write(pk)
	return
}

// Subscribe implements Broker interface
func (b *RWBroker) Subscribe(event string) error {
	b.eventsMux.Lock()
	defer b.eventsMux.Unlock()

	b.events[event] = struct{}{}
	return nil
}

// Unsubscribe implements Broker interface
func (b *RWBroker) Unsubscribe(event string) error {
	b.eventsMux.Lock()
	defer b.eventsMux.Unlock()

	delete(b.events, event)
	return nil
}

// SetCallback implements Broker interface
func (b *RWBroker) SetCallback(handler EventHandler) {
	b.callback = handler
}

// NotifyClose implements Broker interface
func (b *RWBroker) NotifyClose(chan error) error {
	return nil
}
