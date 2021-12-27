package broker

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRWSubscribe(t *testing.T) {
	ctx := context.Background()
	r, w := io.Pipe()
	b := RWBroker{r, w}

	go func() {
		err := b.Subscribe(ctx, []string{"foo"}, Rcv)
		assert.NoError(t, err)
	}()

	go func() {
		err := b.Publish(ctx, "foo", "bar")
		assert.NoError(t, err)
	}()

	res := <-Rcv
	assert.Equal(t, "foo", res.Event())
	assert.EqualValues(t, "bar", res.Body())
}
