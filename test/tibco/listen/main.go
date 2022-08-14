package main

import (
	"context"
	"log"

	"github.com/dualm/mq/tibco"
)

var reC = make(chan string, 1)

func main() {
	opt := tibco.TibOption{
		FieldName: "",
		Service:   "",
		Network:   "",
		Daemon:    []string{},
	}

	infoC := make(chan string)
	errC := make(chan error)

	tib := tibco.NewTibSender(&opt, infoC, errC)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := tib.Run(ctx); err != nil {
		log.Fatal(err)
	}
	defer tib.Close(ctx)

	transport, err := tibco.NewTransport("", "", []string{})
	if err != nil {
		log.Fatal(err)
	}

	listener, err := tibco.NewListener(nil, transport, "a", new(callback))
	if err != nil {
		log.Fatal(err)
	}

	_, err = tibco.NewListener(nil, transport, "b", new(callback))
	if err != nil {
		log.Fatal(err)
	}

	q, err := listener.GetQueue()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			q.Dispatch()
		}
	}()

	for {
		select {
		case info := <-infoC:
			log.Println("info", info)
		case err := <-errC:
			log.Println("error", err)
		case re := <-reC:
			log.Println(re)
		}
	}
}

type callback struct {
}

func (c *callback) CallBack(event tibco.Event, msg tibco.Message) {
	re, err := msg.GetString("DATA", 0)
	if err != nil {
		log.Fatal(err)
	}
	defer msg.Destroy()

	reC <- re
}
