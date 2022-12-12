package main

import (
	"github.com/BurntSushi/toml"
	"github.com/dualm/mq/tibco"
	"log"
	"os"
	"strconv"
	"time"
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
		Requests []Message
		Reports  []Message
	})
	_, err := toml.DecodeFile(os.Args[1], &option)
	if err != nil {
		log.Fatal(err)
	}

	var (
		errChan  = make(chan error, 1)
		infoChan = make(chan string, 1)
	)

	tib := tibco.NewTibCmSender(&option.Option, infoChan, errChan)
	if err := tib.Run(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tib.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if option.Listener != "" {
		if err := tib.AddListener(option.Listener, ""); err != nil {
			log.Fatal("ADD LISTENER ERROR: ", err)
		}
	}

	for i := range option.Reports {
		_reports := option.Reports[i]

		msg, err := tibco.NewMessage()
		if err != nil {
			log.Fatal(err)
		}

		err = msg.SetSendSubject(option.Option.TargetSubjectName)
		if err != nil {
			log.Fatal(err)
		}

		for j := range _reports.Fields {
			_f := _reports.Fields[j]
			msg, err = AddField(msg, &_f)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = tib.SendReportMessage(msg)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(i)
	}

	//
	//for i := range option.Reports {
	//	_m := option.Reports[i]
	//
	//	log.Println("SEND Report: ", _m)
	//
	//	if err != nil {
	//		log.Println("SEND REPORT ERROR: ", err)
	//	}
	//}

	//for i := range option.Requests {
	//	_r := option.Requests[i]
	//
	//	log.Println("SEND Request ", _r)
	//	_reply, err := tib.SendRequest(_r.Message, _r.Field, "", option.Option.LimitTime)
	//	if err != nil {
	//		log.Println("SEND REQUEST ERROR: ", err)
	//
	//		continue
	//	}
	//
	//	var _rplStr string
	//	_rplStr, err = _reply.ConvertToString()
	//	if err != nil {
	//		log.Println("ERROR: ", err)
	//
	//		continue
	//	}
	//
	//	log.Println("RECV Reply: ", _rplStr)
	//}

	time.Sleep(3 * time.Second)
}

func AddField(msg *tibco.Message, field *Field) (*tibco.Message, error) {
	switch field.Type {
	case "i16":
		n, err := strconv.ParseInt(field.Message, 10, 64)
		if err != nil {
			return nil, err
		}

		err = msg.AddI16(field.FieldName, int16(n), 0)
		if err != nil {
			return nil, err
		}
	case "i32":
		n, err := strconv.ParseInt(field.Message, 10, 64)
		if err != nil {
			return nil, err
		}

		err = msg.AddI32(field.FieldName, int32(n), 0)
		if err != nil {
			return nil, err
		}
	case "string":
		err := msg.AddString(field.FieldName, field.Message, 0)
		if err != nil {
			return nil, err
		}
	case "opaque":
		err := msg.AddOpaque(field.FieldName, []byte(field.Message), 0)
		if err != nil {
			return nil, err
		}
	}

	return msg, nil
}
