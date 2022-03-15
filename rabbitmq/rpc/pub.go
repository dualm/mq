package rpc

import (
	"log"

	"github.com/dualm/mq"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/dualm/mq/rabbitmq"
)

func publish(sessions chan chan rabbitmq.Session, queue string, messages <-chan mq.MqMessage) {
	for session := range sessions {
		var (
			running bool
			reading = messages
			pending = make(chan mq.MqMessage, 1)
			confirm = make(chan amqp.Confirmation, 1)
		)

		pub := <-session

		if err := pub.Channel.Confirm(false); err != nil {
			log.Println("publisher confirms not supported")

			close(confirm)
		} else {
			pub.Channel.NotifyPublish(confirm)
		}

	Publish:
		for {
			var body mq.MqMessage
			select {
			case confirmed, ok := <-confirm:
				if !ok {
					break Publish
				}

				if !confirmed.Ack {
					log.Printf("nack message %d, body: %q", confirmed.DeliveryTag, body.Msg.MarshalByte())
				}

				reading = messages
			case body = <-pending:
				err := pub.Channel.Publish("", queue, false, false, amqp.Publishing{
					Body:    body.Msg.MarshalByte(),
					ReplyTo: queue,
				})

				if err != nil {
					log.Println(err)

					pending <- body
					pub.Close()

					break Publish
				}
			case body, running = <-reading:
				if !running {
					return
				}

				pending <- body
				reading = nil
			}
		}
	}
}
