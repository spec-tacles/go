package broker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Packet struct {
	Event string
	Data  []byte
}

var connected = false
var rcv = make(chan *Packet)
var a = NewAMQP("test", "", func(event string, data []byte) {
	rcv <- &Packet{event, data}
})

func connect() {
	if connected {
		return
	}

	err := a.Connect("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}
	connected = true
}

func TestSubscribe(t *testing.T) {
	connect()

	go func() {
		err := a.Subscribe("foo")
		assert.NoError(t, err)
	}()

	event := "foo"
	data := []byte("bar")
	err := a.Publish(event, data)
	assert.NoError(t, err)

	res := <-rcv
	assert.Equal(t, event, res.Event)
	assert.EqualValues(t, data, res.Data)

	err = a.Publish(event, data)
	assert.NoError(t, err)

	res = <-rcv
	assert.Equal(t, event, res.Event)
	assert.EqualValues(t, data, res.Data)

	err = a.Unsubscribe(event)
	assert.NoError(t, err)

	err = a.Publish(event, data)
	assert.NoError(t, err)

	select {
	case d := <-rcv:
		assert.FailNow(t, "unexpected response from AMQP", d)
	case <-time.After(5 * time.Second):
	}

	go func() {
		err := a.Subscribe("foo")
		assert.NoError(t, err)
	}()

	res = <-rcv
	assert.Equal(t, event, res.Event)
	assert.EqualValues(t, data, res.Data)
}

func TestClose(t *testing.T) {
	connect()

	closes := make(chan error)
	err := a.NotifyClose(closes)
	assert.NoError(t, err)

	err = a.Close()
	assert.NoError(t, err)
	assert.NoError(t, <-closes)

	err = a.Publish("foo", []byte("bar"))
	assert.Error(t, err)
}
