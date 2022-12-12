package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/dualm/mq/tibco"
)

func main() {
	opt := tibco.TibCmOption{
		TibOption: &tibco.TibOption{
			FieldName:         "DATA",
			Service:           "",
			Network:           "",
			Daemon:            nil,
			TargetSubjectName: "a",
			SourceSubjectName: "",
			PooledMessage:     true,
		},
		CmName:     "listener",
		LedgerName: "ledger",
		RequestOld: true,
		SyncLedger: true,
		RelayAgent: "",
	}

	infoC := make(chan string)
	errC := make(chan error)

	cmListener, err := tibco.NewTibCmListener(&opt, infoC, errC)
	if err != nil {
		log.Fatal(err)
	}

	if err = cmListener.Listen("a", new(callback)); err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-sigChan:
			if err := cmListener.Close(); err != nil {
				log.Fatal(err)
			}

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

func (c *callback) CmCallBack(event tibco.CmEvent, msg tibco.Message) {
	log.Println("----------Get New Message----------")
	sendSubject, err := msg.GetSendSubject()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("send subject name", sendSubject)

	replySubject, err := msg.GetReplySubject()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("reply subject", replySubject)

	sender, err := msg.GetCMSender()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("sender", sender)

	sequence, err := msg.GetCMSequence()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("sequence", sequence)

	reStr, err := msg.ConvertToString()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("message string", reStr)

	b, err := msg.GetOpaque("Data", 0)
	if err != nil {
		log.Println(err)

		return
	}

	log.Println(string(b))

	defer func() {
		_ = msg.Close()
	}()

}
