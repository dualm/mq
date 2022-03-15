package rpc

import (
	"github.com/dualm/mq/rabbitmq"

	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribe(sessions chan chan rabbitmq.Session, exchange, queue, consumer string, message chan<- *amqp.Delivery) {
	for session := range sessions {
		sub := <-session

		if _, err := sub.Channel.QueueDeclare(queue, true, true, false, false, nil); err != nil {
			log.Printf("cannot consume from exclusive queue: %q, %v", queue, err)

			return
		}

		deliveries, err := sub.Channel.Consume(queue, consumer, false, false, false, false, nil)
		if err != nil {
			log.Printf("cannot consume from: %q, %v", queue, err)

			return
		}

		for range subChan {
			for msg := range deliveries {
				msg := msg
				message <- &msg

				sub.Channel.Ack(msg.DeliveryTag, false)

				break
			}
		}
	}
}
