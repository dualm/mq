package main

import (
	"context"
	"log"

	"github.com/dualm/mq"
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

	resp := make(chan mq.MqResponse)
	tib.Send(
		ctx,
		resp,
		[]mq.MqMessage{
			{
				Msg:           []byte("normal send"),
				CorrelationID: "",
				IsEvent:       true,
			},
		},
		"a",
	)

	<-resp
}
