package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq"
	"github.com/patrickmn/go-cache"
	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribe(ctx context.Context, sessions chan chan rabbitmq.Session,
	_, queue, consumer string, message chan<- amqp.Delivery, subChan <-chan rabbitmq.Subscription,
	infoChan chan<- string, errChan chan<- error) {
	_cache := cache.New(time.Minute, time.Hour)

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

	LOOP:
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-deliveries:
				if !ok {
					infoChan <- "session closed"

					break LOOP
				}

				if err := sub.Channel.Ack(msg.DeliveryTag, false); err != nil {
					errChan <- fmt.Errorf("ack message error, quit session, Error: %w", err)

					break LOOP
				}

				_rspChan, found := _cache.Get(msg.CorrelationId)
				if !found {
					errChan <- fmt.Errorf("receive a unexpected message, Message: %s", string(msg.Body))

					continue
				}

				_cache.Delete(msg.CorrelationId)

				go rabbitmq.SendResponse(mq.MqResponse{
					Msg: msg.Body,
					Err: nil,
				}, _rspChan.(chan<- mq.MqResponse), errChan)
			case newSub := <-subChan:
				_cache.Set(newSub.CorrelationID, newSub.RspChan, cache.DefaultExpiration)
			}
		}
	}
}
