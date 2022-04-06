package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/dualm/mq/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func redial(ctx context.Context, url, exchange, vhost string,
	infoChan chan<- string, errChan chan<- error) chan chan rabbitmq.Session {
	sessions := make(chan chan rabbitmq.Session)

	go func() {
		sess := make(chan rabbitmq.Session)
		defer close(sess)

		for {
			select {
			case sessions <- sess:
			case <-ctx.Done():
				infoChan <- "rabbitmq/pubsub shutting down rabbitmq.Session factory"

				return
			}

			conn, err := amqp.DialConfig(url, amqp.Config{
				Vhost: vhost,
			})
			if err != nil {
				infoChan <- fmt.Sprintf("rabbitmq/pubsub cannot (re)dial: %v: %q", err, url)

				time.Sleep(time.Minute)

				continue
			}

			ch, err := conn.Channel()
			if err != nil {
				infoChan <- fmt.Sprintf("rabbitmq/pubsub cannot create channel: %v", err)

				time.Sleep(time.Minute)

				continue
			}

			if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
				infoChan <- fmt.Sprintf("rabbitmq/pubsub cannot declare topic exchange: %v", err)

				time.Sleep(time.Minute)

				continue
			}

			select {
			case sess <- rabbitmq.Session{Connection: conn, Channel: ch}:
			case <-ctx.Done():
				errChan <- fmt.Errorf("rabbitmq/pubsub shutting down new rabbitmq.Session")

				return
			}
		}
	}()

	return sessions
}
