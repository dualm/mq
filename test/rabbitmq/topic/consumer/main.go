package main

import (
	"context"
	"github.com/dualm/mq/rabbitmq"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	errChan := make(chan error)

	consumer, err := rabbitmq.NewTopicConsumer(
		&rabbitmq.DialOption{
			Username: "spcadm",
			Password: "spcadm",
			Host:     "10.1.14.188",
			Port:     "5671",
			VHost:    "/spc",
		},
		&rabbitmq.ExchangeOption{
			Name:       "exchange.eap.test",
			Kind:       "topic",
			AutoDelete: true,
			Durable:    false,
			NoWait:     false,
			Args:       nil,
		},
		&rabbitmq.ExchangeOption{
			Name:       "exchange.eap.test",
			Kind:       "topic",
			AutoDelete: false,
			Durable:    true,
			NoWait:     false,
			Args:       nil,
		},
		&rabbitmq.ConsumerOption{
			Name:      "",
			AutoAck:   false,
			Exclusive: false,
			NoWait:    false,
			Args:      nil,
		},
		errChan,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := consumer.Recv(&rabbitmq.QueueOption{
		Name:       "",
		AutoDelete: true,
		Durable:    false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	},
		"key.eap.test",
	)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	for {
		select {
		case <-c:
			os.Exit(0)
		case e := <-errChan:
			log.Println(e)
		case m := <-msgs:
			go func() {
				if m.ReplyTo != "" {
					ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
					err = consumer.SendReply(ctx, m, false, false, []byte("Hello"+string(m.Body)))
					if err != nil {
						log.Fatal(err)
					}
					cancel()
					log.Println(*m)
				} else {
					log.Println(string(m.Body))
				}
			}()
		}
	}
}
