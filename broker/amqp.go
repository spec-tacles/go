package broker

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// ErrDisconnected occurs when trying to do something that requires a connection but one was
// unavailable
var ErrDisconnected = errors.New("disconnected from the broker")

// ErrorNoRes occurs when no response is returned from the server on an RPC call
var ErrNoRes = errors.New("no response from server")

// ErrRpcQueueAssertionFailure occurrs when the anon RPC queue fails to create
var ErrRpcQueueAssertionFailure = errors.New("failed to create anonymous rpc queue")

// AMQP is a broker for AMQP clients. Probably most useful for RabbitMQ.
type AMQP struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	receiveCallback EventHandler
	consumerTags    map[string]string
	rpcqueue        amqp.Queue
	rpcconsumer     <-chan amqp.Delivery

	Group    string
	Subgroup string
}

// NewAMQP creates a new AMQP broker.
func NewAMQP(group string, subgroup string, receiveCallback EventHandler) *AMQP {
	return &AMQP{
		receiveCallback: receiveCallback,
		consumerTags:    make(map[string]string),

		Group:    group,
		Subgroup: subgroup,
	}
}

// Connect will connect this client to the AQMP Server
func (a *AMQP) Connect(url string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	a.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	a.channel = ch
	err = ch.ExchangeDeclare(
		a.Group,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	// setup RPC callback queue
	rpc, err := a.channel.QueueDeclare(
		"rpc",
		true,
		false,
		false,
		false,
		nil,
	)
	a.rpcqueue = rpc

	if err != nil {
		return ErrRpcQueueAssertionFailure
	}
	msgs, err := a.channel.Consume("rpc", "", false, false, false, false, nil)
	if err != nil {
		return ErrRpcQueueAssertionFailure
	}
	a.rpcconsumer = msgs

	return err
}

// Close implements io.Closer and Closes the Channel & Connection of this Client
func (a *AMQP) Close() (err error) {
	if a.channel == nil {
		return ErrDisconnected
	}

	err = a.channel.Close()
	if err != nil {
		return
	}

	err = a.conn.Close()
	return
}

// Publish sends data to AMQP
func (a *AMQP) Publish(event string, data []byte) error {
	if a.channel == nil {
		return ErrDisconnected
	}

	return a.publish(event, amqp.Publishing{
		Body:        data,
		ContentType: "application/json",
	})
}

// PublishOptions sends data to AMQP
func (a *AMQP) PublishOptions(opts PublishOptions) error {
	return a.publish(opts.Event, amqp.Publishing{
		Body:        opts.Data,
		ContentType: "application/json",
		Expiration:  strconv.FormatInt(opts.Timeout.Milliseconds(), 10),
	})
}

func (a *AMQP) publish(event string, opts amqp.Publishing) error {
	return a.channel.Publish(
		a.Group,
		event,
		false,
		false,
		opts,
	)
}

// Subscribe will make this client consume for the specific event
func (a *AMQP) Subscribe(event string) (err error) {
	if a.channel == nil {
		return ErrDisconnected
	}

	subgroup := a.Subgroup
	if subgroup != "" {
		subgroup += ":"
	}
	queueName := fmt.Sprintf("%s:%s%s", a.Group, subgroup, event)

	_, err = a.channel.QueueDeclare(
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

	err = a.channel.QueueBind(queueName, event, a.Group, false, nil)
	if err != nil {
		return
	}

	msgs, err := a.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return
	}

	firstMessage := <-msgs
	a.consumerTags[event] = firstMessage.ConsumerTag
	err = firstMessage.Ack(false)
	if err != nil {
		return
	}
	go a.receiveCallback(event, firstMessage.Body)

	for d := range msgs {
		err = d.Ack(false)
		if err != nil {
			return
		}

		go a.receiveCallback(event, d.Body)
	}
	return
}

func (a *AMQP) Call(event string, opts amqp.Publishing) ([]byte, error) {
	correlation := uuid.New().String()
	opts.CorrelationId = correlation
	opts.ReplyTo = "rpc"

	err := a.publish(event, opts)
	if err != nil {
		return nil, err
	}

	for d := range a.rpcconsumer {
		if correlation == d.CorrelationId {
			d.Ack(false)
			return d.Body, nil
		}
	}

	return nil, ErrNoRes
}

// Unsubscribe will make this client cancel the subscription for specific events
func (a *AMQP) Unsubscribe(event string) error {
	if a.channel == nil {
		return ErrDisconnected
	}

	err := a.channel.Cancel(a.consumerTags[event], false)
	if err != nil {
		return err
	}

	delete(a.consumerTags, event)
	return nil
}

// SetCallback sets the function to be called when events are received
func (a *AMQP) SetCallback(handler EventHandler) {
	a.receiveCallback = handler
}

// NotifyClose notifies the given channel of connection closures
func (a *AMQP) NotifyClose(rcv chan error) error {
	if a.conn == nil {
		return ErrDisconnected
	}

	closes := make(chan *amqp.Error) // gets closed automatically
	a.conn.NotifyClose(closes)
	go func() {
		for amqpErr := range closes {
			rcv <- errors.New(amqpErr.Error())
		}

		rcv <- nil
	}()
	return nil
}
