package rpc

import (
	"context"
	"fmt"

	"github.com/dualm/mq"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/dualm/mq/rabbitmq"
)

func publish(ctx context.Context,
	sessions chan chan rabbitmq.Session,
	queue, cQueue string,
	messageChan <-chan mq.MqMessage,
	responseChan chan<- mq.MqResponse, infoChan chan<- string, errChan chan<- error) {
	for session := range sessions {
		var (
			running bool
			reading = messageChan
			pending = make(chan mq.MqMessage, 1)
			confirm = make(chan amqp.Confirmation, 1)
			getting = make(<-chan amqp.Delivery, 1)
			body    mq.MqMessage
		)

		pub := <-session

		if err := pub.Channel.Confirm(false); err != nil {
			infoChan <- "rabbitmq/rpc publisher confirms not supported"

			close(confirm)
		} else {
			pub.Channel.NotifyPublish(confirm)
		}

	Publish:
		for {
			select {
			case <-ctx.Done():
				return
			case confirmed, ok := <-confirm:
				if !ok {
					break Publish
				}

				if !confirmed.Ack {
					infoChan <- fmt.Sprintf("rabbitmq/rpc nack message %d, body: %q", confirmed.DeliveryTag, body.Msg)
				}

				reading = messageChan
			case body = <-pending:
				if !body.IsEvent {
					msgs, err := pub.Channel.Consume(cQueue, "", true, false, false, false, nil)
					if err != nil {
						errChan <- err

						continue
					}

					getting = msgs
				}

				err := pub.Channel.Publish("", queue, false, false, amqp.Publishing{
					Body:    body.Msg,
					ReplyTo: cQueue,
				})

				if err != nil {
					responseChan <- mq.MqResponse{
						Msg: nil,
						Err: fmt.Errorf("rabbitmq/rpc pub message error, Error: %w", err),
					}

					pub.Close()

					break Publish
				}

				if !body.IsEvent {
					i := <-getting
					responseChan <- mq.MqResponse{
						Msg: i.Body,
						Err: nil,
					}
				} else {
					responseChan <- mq.MqResponse{
						Msg: nil,
						Err: nil,
					}
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
