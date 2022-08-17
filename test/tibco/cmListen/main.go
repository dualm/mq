package main

import (
	"context"
	"log"

	"github.com/dualm/mq/tibco"
)

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

	cmTransport, err := tibco.NewCmTransport(transport, "a", true, "ledger", true, "relay")
	if err != nil {
		log.Fatal(err)
	}
	defer cmTransport.Close()

	cmListener, err := tibco.NewCmListener(nil, cmTransport, "a", new(callback))
	if err != nil {
		log.Fatal(err)
	}

	q, err := cmListener.GetQueue()
	if err != nil {
		log.Fatal(err)
	}

	for {
		q.Dispatch()
	}
}

type callback struct {
}

func (c *callback) CmCallBack(event tibco.CmEvent, msg tibco.Message) {
	log.Println("1")
}
