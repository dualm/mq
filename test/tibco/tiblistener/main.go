package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/dualm/mq/tibco"
)

func main() {
	opt := tibco.TibOption{
		FieldName:         "DATA",
		Service:           "",
		Network:           "",
		Daemon:            []string{},
		TargetSubjectName: "",
		SourceSubjectName: "",
		PooledMessage:     false,
	}

	infoC := make(chan string, 1)
	errC := make(chan error, 1)
	tib, err := tibco.NewTibListener(&opt, infoC, errC)
	if err != nil {
		log.Fatal(err)
	}
	defer tib.Destroy()

	tib.Listen("a", nil, new(callback))
	tib.Listen("b", nil, new(callback))

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-sigChan:
			return
		case info := <-infoC:
			log.Println("info", info)
		case err := <-errC:
			log.Println("error", err)
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

	log.Println(re)
}
