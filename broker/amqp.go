package broker

import (
	"fmt"

	"github.com/streadway/amqp"
)

// AMQP is a broker for AMQP clients. Probably most useful for RabbitMQ.
type AMQP struct {
	publishConn *amqp.Connection
    	publishChannel *amqp.Channel
    	consumeConn *amqp.Connection

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
	pubConn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	a.publishConn = pubConn
	ch, err := pubConn.Channel()
    	a.publishChannel = ch
	if err != nil {
		return err
	}
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

// Close implements io.Closer and closes the publishing Channel & Connection of this Client
func (a *AMQP) Close() error {
	err := a.publishChannel.Close()
	if err != nil {
		return err
	}
	err = a.publishConn.Close()
	if err != nil {
		return err
	}
    	err = a.consumeConn.Close()
	if err != nil {
		return err
	}
	return nil
}

// Publish sends data to amqp
func (a *AMQP) Publish(event string, data []byte) error {
	err := a.publishChannel.Publish(
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
	consumeConn, err := amqp.Dial(url)
    	a.consumeConn = consumeConn
	for i := range events {
		event := events[i]
		subgroup := ""
		if a.Subgroup == "" {
			subgroup = a.Subgroup + ":"
		}
		queueName := fmt.Sprintf("%s:%s%s", a.Group, subgroup, event)
        	ch, err := a.consumeConn.Channel()
        	if err != nil {
			return err
		}

		_, err = ch.QueueDeclare(
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

		err = ch.QueueBind(queueName, a.Group, event, false, nil)

		if err != nil {
			return err
		}

		msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)

		if err != nil {
			return err
		}

		go func(receiveCallback func(string, []byte)) {
			gotConsumerTag := false
			for d := range msgs {
				if gotConsumerTag == false {
					a.consumerTags[d.ConsumerTag] = event
					gotConsumerTag = true
				}
				d.Ack(false)
				receiveCallback(event, d.Body)
			}
		}(a.receiveCallback)
	}

	return nil
}

// Unsubscribe will make this client cancel the subscription for specific events
/*func (a *AMQP) Unsubscribe(event string) error {
	err := a.channel.Cancel(a.consumerTags[event], false)

	if err != nil {
		return err
	}

	a.consumerTags[event] = ""

	return nil
}
*/
