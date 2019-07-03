package broker

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

// ErrDisconnected occurs when trying to do something that requires a connection but one was
// unavailable
var ErrDisconnected = errors.New("disconnected from the broker")

// AMQP is a broker for AMQP clients. Probably most useful for RabbitMQ.
type AMQP struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	receiveCallback EventHandler
	consumerTags    map[string]string

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

// Publish sends data to aqmp
func (a *AMQP) Publish(event string, data []byte) error {
	if a.channel == nil {
		return ErrDisconnected
	}

	return a.channel.Publish(
		a.Group,
		event,
		false,
		false,
		amqp.Publishing{
			Body:        data,
			ContentType: "application/json",
		},
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
