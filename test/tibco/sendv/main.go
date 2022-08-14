package main

import (
	"context"
	"log"
	"strconv"

	"github.com/dualm/mq/tibco"
)

func main() {
	opt := tibco.TibOption{
		FieldName: "Message",
		Service:   "",
		Network:   "",
		Daemon:    []string{},
	}

	infoC := make(chan string)
	errC := make(chan error)

	tib := tibco.New(&opt, infoC, errC)

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

	var msgs = make([]*tibco.Message, 0)

	for i := range [10]struct{}{} {
		msg, err := tibco.NewMessage()
		if err != nil {
			log.Fatal(err)
		}

		msg.SetSendSubject("a")

		if err = msg.AddString("Message", strconv.Itoa(i), 0);err != nil {
			log.Fatal(err)
		}

		msgs = append(msgs, msg)
	}

	transport.Sendv(msgs)
}
