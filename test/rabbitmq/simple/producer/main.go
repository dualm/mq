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
	queueOpt := &rabbitmq.QueueOption{
		Name:       "eap_test",
		AutoDelete: true,
		Durable:    true,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}

	dialOption := &rabbitmq.DialOption{
		Username: "spcadm",
		Password: "spcadm",
		Host:     "10.1.14.188",
		Port:     "5671",
		VHost:    "/spc",
	}

	sender, err := rabbitmq.NewSimpleProducer(dialOption, queueOpt)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := sender.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for i := 0; i < 1000; i++ {
		go func() {
			if err := sender.Send(ctx, false, false, []byte("Hello world")); err != nil {
				log.Fatal(err)
			}
		}()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case <-c:
		os.Exit(0)
	}
}
