package pubsub

import (
	"context"
	"fmt"

	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func publish(ctx context.Context, sessions chan chan rabbitmq.Session, exchange, routingKey, replyQueue string,
	messages <-chan mq.MqMessage, infoChan chan<- string, errChan chan<- error) {
	for session := range sessions {
		var (
			running bool
			reading = messages
			pending = make(chan mq.MqMessage, 1)
			confirm = make(chan amqp.Confirmation, 1)
		)

		pub := <-session

		if err := pub.Channel.Confirm(false); err != nil {
			infoChan <- "rabbitmq/pubsub publisher confirms not supported"

			close(confirm)
		} else {
			pub.Channel.NotifyPublish(confirm)
		}

	Publish:
		for {
			var body mq.MqMessage
			select {
			case <-ctx.Done():
				return
			case confirmed, ok := <-confirm:
				if !ok {
					break Publish
				}

				if !confirmed.Ack {
					infoChan <- fmt.Sprintf("rabbitmq/pubsub nack message %d, body: %q", confirmed.DeliveryTag, body.Msg)
				}

				reading = messages
			case body = <-pending:
				err := pub.Channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
					Body:          body.Msg,
					CorrelationId: body.CorraltedId,
					ReplyTo: func() string {
						if body.IsEvent {
							return ""
						}

						return replyQueue
					}(),
				})

				if err != nil {
					errChan <- fmt.Errorf("rabbitmq/pubsub pub message error, Error: %w", err)

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
