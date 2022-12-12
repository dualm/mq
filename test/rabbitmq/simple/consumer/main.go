package main

import (
	"github.com/dualm/mq/rabbitmq"
	"log"
	"os"
	"os/signal"
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

	consumerOpt := &rabbitmq.ConsumerOption{
		Name:      "",
		AutoAck:   false,
		Exclusive: false,
		NoWait:    false,
		Args:      nil,
	}

	errChan := make(chan error)

	recv, err := rabbitmq.NewSimpleConsumer(dialOption, queueOpt, consumerOpt, errChan)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := recv.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	c, err := recv.Recv()
	if err != nil {
		log.Fatal(err)
	}

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, os.Interrupt)

	for {
		select {
		case e := <-errChan:
			log.Fatal("Error", e)
		case b := <-c:
			log.Println(string(b))
		case <-quitChan:
			os.Exit(0)
		}
	}
}
