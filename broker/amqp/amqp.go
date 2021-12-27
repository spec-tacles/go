package amqp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/spec-tacles/go/broker"
)

// ErrorNoRes occurs when no response is returned from the server on an RPC call
var ErrNoRes = errors.New("no response from server")

// ErrRpcQueueAssertionFailure occurrs when the anon RPC queue fails to create
var ErrRpcQueueAssertionFailure = errors.New("failed to create anonymous rpc queue")

type AMQPMessage struct {
	amqp    *AMQP
	event   string
	rcvChan *amqp091.Channel
	d       amqp091.Delivery
}

func (m *AMQPMessage) Event() string {
	return m.event
}

func (m *AMQPMessage) Body() (data interface{}) {
	_ = broker.Decode(m.d.Body, &data)
	return
}

func (m *AMQPMessage) Reply(ctx context.Context, data interface{}) error {
	return m.amqp.Publish(ctx, m.d.ReplyTo, data)
}

func (m *AMQPMessage) Ack(ctx context.Context) error {
	return m.rcvChan.Ack(m.d.DeliveryTag, false)
}

// AMQP is a broker for AMQP clients. Probably most useful for RabbitMQ.
type AMQP struct {
	conn        *amqp091.Connection
	publishChan *amqp091.Channel
	rpcQueue    amqp091.Queue
	rpcConsumer <-chan amqp091.Delivery

	Group    string
	Subgroup string
	Timeout  time.Duration
}

// Connect will connect this client to the AQMP Server
func (a *AMQP) Connect(ctx context.Context, url string) error {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return err
	}
	a.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	a.publishChan = ch

	err = a.setupRPC()
	return err
}

func (a *AMQP) setupRPC() error {
	ch, err := a.conn.Channel()
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(
		a.Group,
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// setup RPC callback queue
	rpc, err := ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	a.rpcQueue = rpc

	if err != nil {
		return ErrRpcQueueAssertionFailure
	}
	msgs, err := ch.Consume(rpc.Name, "", true, true, false, false, nil)
	if err != nil {
		return ErrRpcQueueAssertionFailure
	}
	a.rpcConsumer = msgs

	return nil
}

// Close implements io.Closer and Closes the Channel & Connection of this Client
func (a *AMQP) Close() (err error) {
	err = a.conn.Close()
	return
}

// Publish sends data to AMQP
func (a *AMQP) Publish(ctx context.Context, event string, data interface{}) error {
	b, err := broker.Encode(data)
	if err != nil {
		return err
	}

	return a.publish(event, amqp091.Publishing{
		Body:       b,
		Expiration: strconv.FormatInt(a.Timeout.Milliseconds(), 10),
	})
}

func (a *AMQP) publish(event string, opts amqp091.Publishing) error {
	if a.publishChan == nil {
		return broker.ErrDisconnected
	}

	return a.publishChan.Publish(
		a.Group,
		event,
		false,
		false,
		opts,
	)
}

// Subscribe will make this client consume for the specific event
func (a *AMQP) Subscribe(ctx context.Context, events []string, messages chan<- broker.Message) (err error) {
	ch, err := a.conn.Channel()
	if err != nil {
		return err
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error)
	defer close(errs)

	for _, event := range events {
		go func(e string) {
			errs <- a.subscribeSingle(ctx, ch, e, messages)
		}(event)
	}

	return <-errs
}

func (a *AMQP) subscribeSingle(ctx context.Context, ch *amqp091.Channel, event string, messages chan<- broker.Message) (err error) {
	subgroup := a.Subgroup
	if subgroup != "" {
		subgroup += ":"
	}
	queueName := fmt.Sprintf("%s:%s%s", a.Group, subgroup, event)

	_, err = ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	err = ch.QueueBind(queueName, event, a.Group, false, nil)
	if err != nil {
		return
	}

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return
	}

	errs := make(chan *amqp091.Error)
	ch.NotifyClose(errs)

	var consumerTag string
	for {
		select {
		case <-ctx.Done():
			err = ch.Cancel(consumerTag, false)
			if err != nil {
				return
			}

		case err = <-errs:
			return

		case d, ok := <-msgs:
			if !ok {
				return
			}

			consumerTag = d.ConsumerTag
			messages <- &AMQPMessage{
				amqp:  a,
				event: event,
				d:     d,
			}
		}
	}
}

func (a *AMQP) Call(event string, opts amqp091.Publishing) ([]byte, error) {
	correlation := uuid.New().String()
	opts.CorrelationId = correlation
	opts.ReplyTo = a.rpcQueue.Name

	err := a.publish(event, opts)
	if err != nil {
		return nil, err
	}

	for d := range a.rpcConsumer {
		if correlation == d.CorrelationId {
			return d.Body, nil
		}
	}

	return nil, ErrNoRes
}
