package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/dualm/mq/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func redial(c context.Context, url, queue, vhost string,
	infoChan chan<- string, _ chan<- error) chan chan rabbitmq.Session {
	sessions := make(chan chan rabbitmq.Session)

	go func() {
		sess := make(chan rabbitmq.Session)
		defer close(sess)

		for {
			select {
			case sessions <- sess:
			case <-c.Done():
				infoChan <- ("shutting down rabbitmq.Session factory")

				return
			}

			conn, err := amqp.DialConfig(url, amqp.Config{
				Vhost: vhost,
			})
			if err != nil {
				infoChan <- fmt.Sprintf("cannot (re)dial: %v: %q", err, url)

				time.Sleep(time.Minute)

				continue
			}

			ch, err := conn.Channel()
			if err != nil {
				infoChan <- fmt.Sprintf("cannot create channel: %v", err)

				time.Sleep(time.Minute)

				continue
			}

			_, err = ch.QueueDeclare(queue, true, true, false, false, nil)
			if err != nil {
				infoChan <- fmt.Sprintf("cannot declare queue: %v", err)

				time.Sleep(time.Minute)

				continue
			}

			select {
			case sess <- rabbitmq.Session{Connection: conn, Channel: ch}:
			case <-c.Done():
				infoChan <- "shutting down new rabbitmq.Session"

				return
			}
		}
	}()

	return sessions
}
