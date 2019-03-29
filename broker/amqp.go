package broker

import (
	"fmt"

	"github.com/streadway/amqp"
)

// AMQP is a broker for AMQP clients. Probably most useful for RabbitMQ.
type AMQP struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	receiveCallback func(string, []byte)
	consumerTags    map[string]string

	Group    string
	Subgroup string
}

// NewAMQP creates a new AMQP broker.
func NewAMQP(group string, subgroup string, receiveCallback func(string, []byte)) *AMQP {
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
func (a *AMQP) Close() error {
	err := a.channel.Close()
	if err != nil {
		return err
	}
	err = a.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// Publish sends data to aqmp
func (a *AMQP) Publish(event string, data []byte) error {
	err := a.channel.Publish(
		a.Group,
		event,
		false,
		false,
		amqp.Publishing{
			Body:        data,
			ContentType: "application/json",
		},
	)
	return err
}

// Subscribe will make this client consume for the specific event
func (a *AMQP) Subscribe(events ...string) error {
	for i := range events {
		event := events[i]
		subgroup := ""
		if a.Subgroup == "" {
			subgroup = a.Subgroup + ":"
		}
		queueName := fmt.Sprintf("%s:%s%s", a.Group, subgroup, event)

		_, err := a.channel.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			return err
		}

		err = a.channel.QueueBind(queueName, a.Group, event, false, nil)

		if err != nil {
			return err
		}

		msgs, err := a.channel.Consume(queueName, "", false, false, false, false, nil)

		if err != nil {
			return err
		}

		firstMessage := <-msgs
		a.consumerTags[firstMessage.ConsumerTag] = event
		firstMessage.Ack(false)
		a.receiveCallback(event, firstMessage.Body)

		go func(receiveCallback func(string, []byte)) {
			for d := range msgs {
				d.Ack(false)
				receiveCallback(event, d.Body)
			}
		}(a.receiveCallback)
	}

	return nil
}

// Unsubscribe will make this client cancel the subscription for specific events
func (a *AMQP) Unsubscribe(event string) error {
	err := a.channel.Cancel(a.consumerTags[event], false)

	if err != nil {
		return err
	}

	a.consumerTags[event] = ""

	return nil
}
