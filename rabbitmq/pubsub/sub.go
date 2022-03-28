package pubsub

import (
	"context"
	"fmt"

	"github.com/dualm/mq/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribe(ctx context.Context, sessions chan chan rabbitmq.Session,
	_, queue, consumer string, message chan<- amqp.Delivery, subChan <-chan struct{},
	_ chan<- string, errChan chan<- error) {
	for session := range sessions {
		sub := <-session

		if _, err := sub.Channel.QueueDeclare(queue, true, true, false, false, nil); err != nil {
			errChan <- fmt.Errorf("cannot consume from exclusive queue: %q, %v", queue, err)

			return
		}

		deliveries, err := sub.Channel.Consume(queue, consumer, false, false, false, false, nil)
		if err != nil {
			errChan <- fmt.Errorf("cannot consume from: %q, %v", queue, err)

			return
		}

		for range subChan {
		LOOP:
			for msg := range deliveries {
				select {
				case <-ctx.Done():
					return
				case message <- msg:
					sub.Channel.Ack(msg.DeliveryTag, false)

					break LOOP
				}
			}
		}
	}
}
