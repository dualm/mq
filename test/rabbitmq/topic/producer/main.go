package main

import (
	"context"
	"github.com/dualm/mq/rabbitmq"
	"log"
	"sync"
	"time"
)

func main() {
	producer, err := rabbitmq.NewTopicProducer(
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
			AutoDelete: true,
			Durable:    false,
			NoWait:     false,
			Args:       nil,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//for i := 0; i < 1000; i++ {
	//err = producer.SendReport(
	//	ctx,
	//	"key.eap.test",
	//	false, false,
	//	[]byte("Hello World"),
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//}

	log.Println("send")
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			_, err := producer.SendRequest(
				ctx,
				"key.eap.test",
				false, false,
				[]byte("Hi"),
			)
			if err != nil {
				log.Println(err)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
