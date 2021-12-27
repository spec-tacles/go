package redis

import (
	"context"
	"testing"
	"time"

	"github.com/mediocregopher/radix/v4"
	"github.com/spec-tacles/go/broker"
	"github.com/stretchr/testify/assert"
)

var connected = false
var r = NewRedis("test", "")

func connect() {
	if connected {
		return
	}

	ctx := context.Background()

	conn, err := radix.PoolConfig{}.New(ctx, "tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}

	if err = conn.Do(ctx, radix.Cmd(nil, "FLUSHDB")); err != nil {
		panic(err)
	}

	err = r.Connect(ctx, "localhost:6379")
	if err != nil {
		panic(err)
	}
	connected = true
}

func TestSubscribe(t *testing.T) {
	connect()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := r.Subscribe(ctx, []string{"foo"}, broker.Rcv)
		assert.ErrorIs(t, err, context.Canceled)
	}()

	event := "foo"
	data := "bar"
	err := r.Publish(ctx, event, data)
	assert.NoError(t, err)

	res := <-broker.Rcv
	assert.NoError(t, res.Ack(ctx))
	assert.Equal(t, event, res.Event())
	assert.EqualValues(t, data, res.Body())

	err = r.Publish(ctx, event, data)
	assert.NoError(t, err)

	res = <-broker.Rcv
	assert.NoError(t, res.Ack(ctx))
	assert.Equal(t, event, res.Event())
	assert.EqualValues(t, data, res.Body())

	// unsubscribe by canceling context
	cancel()
	ctx = context.Background()

	err = r.Publish(ctx, event, data)
	assert.NoError(t, err)

	select {
	case d := <-broker.Rcv:
		assert.FailNow(t, "unexpected response from Redis", d)
	case <-time.After(5 * time.Second):
	}

	go func() {
		err := r.Subscribe(ctx, []string{"foo"}, broker.Rcv)
		assert.NoError(t, err)
	}()

	res = <-broker.Rcv
	assert.NoError(t, res.Ack(ctx))
	assert.Equal(t, event, res.Event())
	assert.EqualValues(t, data, res.Body())
}

func TestAutoclaim(t *testing.T) {
	connect()

	ctx := context.Background()
	otherCtx, otherCancel := context.WithCancel(ctx)

	otherRedis := NewRedis("test", "")
	err := otherRedis.Connect(otherCtx, "localhost:6379")
	assert.NoError(t, err)

	go func() {
		err := otherRedis.Subscribe(otherCtx, []string{"foo"}, broker.Rcv)
		assert.ErrorIs(t, err, context.Canceled)
	}()

	err = r.Publish(ctx, "foo", "bar")
	assert.NoError(t, err)

	msg := <-broker.Rcv
	assert.Equal(t, "bar", msg.Body())

	otherCancel()

	go func() {
		err := r.Subscribe(ctx, []string{"foo"}, broker.Rcv)
		assert.NoError(t, err)
	}()

	msg = <-broker.Rcv
	assert.NoError(t, msg.Ack(ctx))
	assert.Equal(t, "bar", msg.Body())
}
