package pubsub

import (
	"context"
	"fmt"

	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribe(ctx context.Context, sessions chan chan rabbitmq.Session,
	_, queue, consumer string, message chan<- amqp.Delivery, subChan <-chan rabbitmq.Subscription,
	_ chan<- string, errChan chan<- error) {
	_m := make(map[string]chan<- mq.MqResponse)
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

		for newSub := range subChan {
			_m[newSub.Name] = newSub.RspChan
		LOOP:
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-deliveries:
					// 发送消息接收确认异常不能影响消息处理
					if err := sub.Channel.Ack(msg.DeliveryTag, false); err != nil {
						errChan <- fmt.Errorf("ack message error, Error: %w", err)
					}

					_rspChan, ok := _m[msg.CorrelationId]
					if !ok {
						errChan <- fmt.Errorf("receive a unexpected message, Message: %s", string(msg.Body))

						continue
					}

					delete(_m, msg.CorrelationId)

					go rabbitmq.SendResponse(ctx, mq.MqResponse{
						Msg: msg.Body,
						Err: nil,
					}, _rspChan, errChan)

					break LOOP
				}
			}
		}
	}
}
