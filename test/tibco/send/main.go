package main

import (
	"github.com/dualm/mq/tibco"
	"log"
	"time"
)

func main() {
	opt := tibco.TibOption{
		FieldName:         "DATA",
		Service:           "",
		Network:           "",
		Daemon:            []string{},
		TargetSubjectName: "a",
		SourceSubjectName: "",
		PooledMessage:     true,
	}

	infoC := make(chan string)
	errC := make(chan error)

	tib := tibco.NewTibSender(&opt, infoC, errC)

	if err := tib.Run(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := tib.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	for j := 0; j < 10; j++ {
		_n := time.Now()
		for i := 0; i < 10; i++ {
			err := tib.SendReports([]string{"1", "2", "3", "4", "5"})
			if err != nil {
				panic(err)
			}
		}
		log.Println(time.Since(_n).Nanoseconds())
	}
}
