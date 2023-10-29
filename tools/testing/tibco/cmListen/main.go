package main

import (
	"github.com/BurntSushi/toml"
	"github.com/dualm/mq/tibco"
	"log"
	"os"
	"os/signal"
)

type Message struct {
	Fields []Field
}

type Field struct {
	FieldName string
	Type      string
	Message   string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("参数不足, 请提供包含通信参数的文件路径")
	}

	option := new(struct {
		Listener string
		Option   tibco.TibCmOption
	})
	_, err := toml.DecodeFile(os.Args[1], &option)
	if err != nil {
		log.Fatal(err)
	}

	var (
		errChan  = make(chan error, 1)
		infoChan = make(chan string, 1)
	)

	tib, err := tibco.NewTibCmListener(&option.Option, infoChan, errChan)
	if err != nil {
		log.Fatal(err)
	}

	if err := tib.Listen(option.Option.TargetSubjectName, new(callback)); err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		select {
		case <-sigChan:
			if err := tib.Close(); err != nil {
				log.Fatal(err)
			}

			return
		case info := <-infoChan:
			log.Println("info", info)
		case err := <-errChan:
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

	b, err := msg.GetOpaque("Data/", 0)
	if err != nil {
		log.Println(err)

		return
	}

	log.Println(string(b))

	defer func() {
		_ = msg.Close()
	}()

}
