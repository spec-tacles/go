package amqp

import (
	"context"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spec-tacles/go/broker"
	"github.com/stretchr/testify/assert"
)

var connected = false
var a = AMQP{Group: "test"}

func connect() {
	if connected {
		return
	}

	conn, err := amqp091.Dial("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}

	err = a.Init(conn)
	if err != nil {
		panic(err)
	}

	connected = true
}

func TestSubscribe(t *testing.T) {
	connect()

	ctx := context.Background()

	go func() {
		err := a.Subscribe(ctx, []string{"foo"}, broker.Rcv)
		assert.NoError(t, err)
	}()

	event := "foo"
	data := []byte("bar")
	err := a.Publish(ctx, event, data)
	assert.NoError(t, err)

	res := <-broker.Rcv
	assert.Equal(t, event, res.Event)
	assert.EqualValues(t, data, res.Body())

	err = a.Publish(ctx, event, data)
	assert.NoError(t, err)

	res = <-broker.Rcv
	assert.Equal(t, event, res.Event)
	assert.EqualValues(t, data, res.Body())

	err = a.Publish(ctx, event, data)
	assert.NoError(t, err)

	select {
	case d := <-broker.Rcv:
		assert.FailNow(t, "unexpected response from AMQP", d)
	case <-time.After(5 * time.Second):
	}

	go func() {
		err := a.Subscribe(ctx, []string{"foo"}, broker.Rcv)
		assert.NoError(t, err)
	}()

	res = <-broker.Rcv
	assert.Equal(t, event, res.Event)
	assert.EqualValues(t, data, res.Body())
}
