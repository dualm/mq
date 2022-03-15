package pubsub

import (
	"context"
	"log"
	"time"

	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	eveChan = make(chan mq.MqMessage)
	reqChan = make(chan mq.MqMessage)
	subChan = make(chan bool, 10)
	rspChan = make(chan amqp.Delivery)
)

func redial(ctx context.Context, url, exchange, vhost string) chan chan rabbitmq.Session {
	sessions := make(chan chan rabbitmq.Session)

	go func() {
		sess := make(chan rabbitmq.Session)
		defer close(sess)

		for {
			select {
			case sessions <- sess:
			case <-ctx.Done():
				log.Println("shutting down rabbitmq.Session factory")

				return
			}

			conn, err := amqp.DialConfig(url, amqp.Config{
				Vhost: vhost,
			})
			if err != nil {
				log.Printf("cannot (re)dial: %v: %q", err, url)

				time.Sleep(time.Minute)

				continue
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Printf("cannot create channel: %v", err)

				time.Sleep(time.Minute)

				continue
			}

			if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
				log.Printf("cannot declare topic exchange: %v", err)

				time.Sleep(time.Minute)

				continue
			}

			select {
			case sess <- rabbitmq.Session{Connection: conn, Channel: ch}:
			case <-ctx.Done():
				log.Println("shutting down new rabbitmq.Session")

				return
			}
		}
	}()

	return sessions
}

func send(msg mq.MqMessage) {
	eveChan <- msg
}

func sendRequest(duration time.Duration, msg mq.MqMessage, rsp chan<- []byte) {
	// 发送数据
	subChan <- true
	reqChan <- msg
	// 计时开始
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case delivery := <-rspChan:
			if delivery.CorrelationId == msg.CorraltedId {
				log.Printf("Got Response: \n%s", string(delivery.Body))

				rsp <- delivery.Body

				return
			}

			rspChan <- delivery
		case <-ticker.C:
			log.Println("Mes请求消息发送超时: ", msg.CorraltedId)
			rsp <- nil

			return
		}
	}
}
