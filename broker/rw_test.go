package broker

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testReadWriter struct {
	C chan []byte
}

func (r *testReadWriter) Read(d []byte) (int, error) {
	return copy(d, <-r.C), nil
}

func (r *testReadWriter) Write(d []byte) (int, error) {
	r.C <- d
	return len(d), nil
}

func TestRWSubscribe(t *testing.T) {
	ctx := context.Background()
	var rw = &testReadWriter{
		C: make(chan []byte),
	}
	var b = NewRW(rw, rw)

	err := b.Subscribe(ctx, []string{"foo"}, Rcv)
	assert.NoError(t, err)

	d := json.RawMessage(`{"event":"foo","data":["bar"]}`)
	assert.NoError(t, err)
	go rw.Write(d)

	res := <-Rcv
	assert.Equal(t, "foo", res.Event)
	assert.EqualValues(t, []byte("[\"bar\"]"), res.Body())
}

func TestRWPublish(t *testing.T) {
	ctx := context.Background()
	var rw = &testReadWriter{
		C: make(chan []byte),
	}
	var b = NewRW(rw, rw)

	go func() {
		err := b.Publish(ctx, "foo", []byte("[\"bar\"]"))
		assert.NoError(t, err)
	}()

	expected, err := json.Marshal(&IOPacket{
		E: "foo",
		D: []byte("[\"bar\"]"),
	})
	assert.NoError(t, err)

	res := <-rw.C
	assert.EqualValues(t, expected, res)
}
