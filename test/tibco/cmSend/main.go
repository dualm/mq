package main

import (
	"github.com/dualm/mq/tibco"
	"log"
	"strconv"
	"time"
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
		CmName:     "sender",
		LedgerName: "ledger",
		RequestOld: false,
		SyncLedger: false,
		RelayAgent: "",
	}

	infoC := make(chan string)
	errC := make(chan error)

	tib := tibco.NewTibCmSender(&opt, infoC, errC)

	if err := tib.Run(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tib.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := tib.AddListener("listener", "a"); err != nil {
		log.Fatal(err)
	}

	if err := tib.SendReport("register", "Data", "a"); err != nil {
		log.Fatal(err)
	}

	for i := range [10]struct{}{} {
		err := tib.SendReport(strconv.Itoa(i), "Data", "a")
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(3 * time.Second)
}
