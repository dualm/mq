package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/dualm/mq/tibco"
)

var (
	tib  *tibco.TibListener
	_str *str
)

type str struct {
	s string
}

func (s *str) Println() {
	log.Println(s.s)
}

func main() {
	opt := tibco.TibOption{
		FieldName:         "Message",
		Service:           "12120",
		Network:           ";225.12.12.12",
		Daemon:            []string{},
		TargetSubjectName: "",
		SourceSubjectName: "",
		PooledMessage:     false,
	}

	_str = &str{
		s: "Hello",
	}

	infoC := make(chan string, 1)
	errC := make(chan error, 1)
	var err error

	tib, err = tibco.NewTibListener(&opt, CallBack, infoC, errC)
	if err != nil {
		log.Fatal(err)
	}

	err = tib.Listen("ZXY.FAB.RMS.DEV.EAPsvr")
	if err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-sigChan:
			log.Println(tib.Close())

			return
		case info := <-infoC:
			log.Println("info", info)
		case err := <-errC:
			log.Println("error", err)
		}
	}
}

func CallBack(_ tibco.Event, msg tibco.Message) {
	_, err := msg.GetString("Message", 0)
	if err != nil {
		log.Fatal(err)
	}

	err = tib.SendReplyMessage("aaa", &msg)
	if err != nil {
		log.Println(err)
	}

	_str.Println()

	defer func() {
		_ = msg.Close()
	}()
}
