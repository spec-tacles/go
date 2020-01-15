package broker

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRWSubscribe(t *testing.T) {
	var rw = &testReadWriter{
		C: make(chan []byte),
	}
	var b = NewRW(rw, rw, cb)

	err := b.Subscribe("foo")
	assert.NoError(t, err)

	d := json.RawMessage(`{"event":"foo","data":["bar"]}`)
	assert.NoError(t, err)
	go rw.Write(d)

	res := <-rcv
	assert.Equal(t, "foo", res.Event)
	assert.EqualValues(t, []byte("[\"bar\"]"), res.Data)
}

func TestRWPublish(t *testing.T) {
	var rw = &testReadWriter{
		C: make(chan []byte),
	}
	var b = NewRW(rw, rw, cb)

	go func() {
		err := b.Publish("foo", []byte("[\"bar\"]"))
		assert.NoError(t, err)
	}()

	expected, err := json.Marshal(&IOPacket{
		Event: "foo",
		Data:  []byte("[\"bar\"]"),
	})
	assert.NoError(t, err)

	res := <-rw.C
	assert.EqualValues(t, expected, res)
}
