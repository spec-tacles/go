package redis

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/mediocregopher/radix/v4"
	"github.com/mediocregopher/radix/v4/resp/resp3"
	"github.com/spec-tacles/go/broker"
)

const streamDataKey = "data"

// RedisMessage represents a message received from the Redis broker
type RedisMessage struct {
	r     *Redis
	id    radix.StreamEntryID
	event string
	body  string
}

type RedisActor interface {
	Do(context.Context, radix.Action) error
}

func (m *RedisMessage) Event() string {
	return m.event
}

// Body returns the body of the message
func (m *RedisMessage) Body() (data interface{}) {
	_ = broker.Decode([]byte(m.body), &data)
	return
}

// Reply sends a RPC response back to the original client
func (m *RedisMessage) Reply(ctx context.Context, data interface{}) error {
	b, err := broker.Encode(data)
	if err != nil {
		return err
	}

	key := m.event + m.id.String()
	return m.r.actor.Do(ctx, radix.Cmd(nil, "PUBLISH", key, string(b)))
}

// Ack acknowledges receipt of the message
func (m *RedisMessage) Ack(ctx context.Context) error {
	return m.r.actor.Do(ctx, radix.Cmd(nil, "XACK", m.event, m.r.Group, m.id.String()))
}

// Redis is a broker that uses Redis streams
type Redis struct {
	actor RedisActor

	Config        radix.PoolConfig
	Group         string
	Name          string
	MaxChunk      uint64
	BlockInterval time.Duration

	// UnackTimeout is the amount of time a client is allowed to wait before acknowledging a stream
	// item. In Redis terms, this is the amount of time an item is allowed to spend in the PEL
	// before being claimed by another client.
	UnackTimeout time.Duration

	// PendingTimeout is the amount of time an item can live in the queue. Items older than this
	// duration (relative to current server time) are evicted from the queue and not processed. This
	// relies on your server time being somewhat synced with your Redis server.
	PendingTimeout time.Duration
}

// NewRedis creates a new Redis broker
func NewRedis(actor RedisActor, group string) *Redis {
	return &Redis{
		actor: actor,

		Group:          group,
		Name:           strconv.FormatInt(rand.Int63(), 16),
		MaxChunk:       10,
		BlockInterval:  3 * time.Second,
		UnackTimeout:   15 * time.Second,
		PendingTimeout: 1 * time.Hour,
	}
}

// Publish publishes a message to the broker
func (r *Redis) Publish(ctx context.Context, event string, data interface{}) error {
	if r.actor == nil {
		return broker.ErrDisconnected
	}

	b, err := broker.Encode(data)
	if err != nil {
		return err
	}

	var action radix.Action
	if r.UnackTimeout != 0 {
		minTime := strconv.FormatInt(time.Now().Add(-r.PendingTimeout).UnixMilli(), 10)
		action = radix.Cmd(
			nil,
			"XADD", event,
			"MINID", "~", minTime,
			"*",
			streamDataKey, string(b),
		)
	} else {
		action = radix.Cmd(
			nil,
			"XADD", event,
			"*",
			streamDataKey, string(b),
		)
	}

	return r.actor.Do(ctx, action)
}

// Subscribe subscribes this broker to an event
func (r *Redis) Subscribe(ctx context.Context, events []string, messages chan<- broker.Message) error {
	if r.actor == nil {
		return broker.ErrDisconnected
	}

	for _, event := range events {
		err := r.actor.Do(ctx, radix.Cmd(nil, "XGROUP", "CREATE", event, r.Group, "0", "MKSTREAM"))

		var redisError resp3.SimpleError
		if errors.As(err, &redisError) && strings.HasPrefix(redisError.S, "BUSYGROUP") {
			continue
		}

		if err != nil {
			return err
		}
	}

	return r.listen(ctx, events, messages)[0]
}

func (r *Redis) listen(ctx context.Context, events []string, messages chan<- broker.Message) [2]error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error)
	defer close(errs)

	go func() {
		errs <- r.listenXread(ctx, events, messages)
	}()

	go func() {
		errs <- r.listenXautoclaim(ctx, events, messages)
	}()

	return [2]error{<-errs, <-errs}
}

func (r *Redis) listenXread(ctx context.Context, events []string, messages chan<- broker.Message) (err error) {
	var (
		data      []radix.StreamEntries
		streamIds []string
	)

	for {
		eventCount := len(events)
		if len(streamIds) != eventCount {
			streamIds = make([]string, eventCount)
			for i := 0; i < eventCount; i++ {
				streamIds[i] = ">"
			}
		}

		action := radix.FlatCmd(&data, "XREADGROUP",
			"GROUP", r.Group, r.Name,
			"COUNT", strconv.FormatUint(r.MaxChunk, 10),
			"BLOCK", strconv.FormatInt(r.BlockInterval.Milliseconds(), 10),
			"STREAMS", events, streamIds,
		)
		err = r.actor.Do(ctx, action)

		if err != nil {
			return
		}

		for _, entry := range data {
			r.handleData(&entry.Entries, entry.Stream, messages)
		}
	}
}

func (r *Redis) listenXautoclaim(ctx context.Context, events []string, messages chan<- broker.Message) (err error) {
	var data radix.StreamEntries

	for {
		timeout := strconv.FormatInt(r.UnackTimeout.Milliseconds(), 10)
		start := "0-0"

		for _, event := range events {
			action := radix.Cmd(&data, "XAUTOCLAIM", event, r.Group, r.Name, timeout, start)
			err = r.actor.Do(ctx, action)

			if err != nil {
				return
			}

			start = data.Stream
			r.handleData(&data.Entries, event, messages)
		}

		time.Sleep(r.BlockInterval)
	}
}

func (r *Redis) handleData(data *[]radix.StreamEntry, event string, msgs chan<- broker.Message) {
	for _, entry := range *data {
		for _, v := range entry.Fields {
			k, v := v[0], v[1]
			if k != streamDataKey {
				continue
			}

			msgs <- &RedisMessage{
				r:     r,
				event: event,
				body:  v,
				id:    entry.ID,
			}
		}
	}
}
