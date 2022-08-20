package main

import (
	"context"
	"log"
	"time"

	"github.com/dualm/mq"
	"github.com/dualm/mq/tibco"
)

func main() {
	opt := tibco.TibOption{
		FieldName:         "Message",
		Service:           "",
		Network:           "",
		Daemon:            []string{},
		TargetSubjectName: "a",
	}

	infoC := make(chan string)
	errC := make(chan error)

	tib := tibco.NewTibSender(&opt, infoC, errC)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := tib.Run(ctx); err != nil {
		log.Fatal(err)
	}

	// async
	go func() {
		for {
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
			)

			<-resp

			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
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
			)

			<-resp

			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			tib.SendReport("aaaa")

			time.Sleep(time.Second)
		}
	}()

	<-ctx.Done()
}
