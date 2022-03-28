package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq/pubsub"
	"github.com/spf13/viper"
)

func main() {
	pb := pubsub.NewPubsub()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	infoChan := make(chan string, 10)
	errChan := make(chan error, 10)

	pb.Run(ctx, "config", InitConfig, infoChan, errChan)

	rsp := make(chan []byte)
	msg := NewAreYouThereReq("TestEQP", "ZLFMM")

	log.Println(string(msg.Msg))

	ctxQuery, cancalQuery := context.WithTimeout(ctx, 20*time.Second)
	defer cancalQuery()

	go pb.Send(ctxQuery, rsp, []mq.MqMessage{
		msg,
	})

	for {
		select {
		case <-ctx.Done():
			return
		case <- ctxQuery.Done():
			return
		case info := <-infoChan:
			log.Println(info)
		case err := <-errChan:
			log.Fatal(err)
		case re := <-rsp:
			log.Println(string(re))
			_, err := AreYouThereDecoder(re)
			if err != nil {
				log.Fatal(err)
			}

			return
		}
	}
}

func InitConfig(configId string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigType("toml")
	conf.SetConfigName(strings.ToLower(configId))
	conf.AddConfigPath(".")
	err := conf.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	return conf
}

type AreYouThereReq struct {
	XMLName       xml.Name `xml:"Message"`
	Header        Header   `xml:"Header"`
	MachineName   string   `xml:"Body>MACHINENAME"`
	ReturnCode    string   `xml:"Return>RETURNCODE"`
	ReturnMessage string   `xml:"Return>RETURNMESSAGE"`
}

func (req *AreYouThereReq) MarshalByte() []byte {
	return MarshalByte(req)
}

func NewAreYouThereReq(machineName, factoryName string) mq.MqMessage {
	messageName := "AreYouThereRequest"
	header, id := NewHeader(messageName, machineName, factoryName)
	areYouThereReq := AreYouThereReq{
		Header:      header,
		MachineName: machineName,
	}

	return mq.MqMessage{
		Msg:         areYouThereReq.MarshalByte(),
		CorraltedId: id,
		IsEvent:     false,
	}
}

func AreYouThereDecoder(b []byte) (*AreYouThereReq, error) {
	if b == nil {
		return nil, fmt.Errorf("nil message")
	}
	rsp := new(AreYouThereReq)

	err := xml.Unmarshal(b, rsp)
	if err != nil {
		return nil, err
	}

	if rsp.ReturnCode != "0" {
		return nil,
			fmt.Errorf("EQP: %s, ReturnCode: %s, ReturnMessage: %s",
				rsp.MachineName,
				rsp.ReturnCode,
				rsp.ReturnMessage,
			)
	}

	return rsp, nil
}

func GetTransactionId() string {
	t := time.Now()
	micro := t.Nanosecond() / 1000

	return fmt.Sprintf("%s%06d", t.Format("20060102150405"), micro)
}

type Header struct {
	XMLName                   xml.Name `xml:"Header"`
	MessageName               string   `xml:"MESSAGENAME"`
	ShopName                  string   `xml:"SHOPNAME"`
	MachineName               string   `xml:"MACHINENAME"`
	TransactionId             string   `xml:"TRANSACTIONID"`
	OriginalSourceSubjectName string   `xml:"ORIGINALSOURCESUBJECTNAME"`
	SourceSubjectName         string   `xml:"SOURCESUBJECTNAME"`
	TargetSubjectName         string   `xml:"TARGETSUBJECTNAME"`
	EventUser                 string   `xml:"EVENTUSER"`
	EventComment              string   `xml:"EVENTCOMMENT"`
}

func NewHeader(messageName, machineName, factoryName string) (Header, string) {
	id := GetRandStr(10)

	return Header{
		MessageName:               messageName,
		ShopName:                  factoryName,
		MachineName:               machineName,
		TransactionId:             GetTransactionId(),
		OriginalSourceSubjectName: id + "-" + machineName + "-" + messageName,
		SourceSubjectName:         "ZL.FMM.MES.TEST.8MMPI91",
		TargetSubjectName:         "ZL.FMM.MES.TEST.PEXsvr",
		EventUser:                 machineName,
		EventComment:              messageName,
	}, machineName + "-" + id
}

func GetRandStr(l int) string {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))

	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}

	return string(bytes)
}

func MarshalByte(v interface{}) []byte {
	b, err := xml.MarshalIndent(v, "", "        ")
	if err != nil {
		return nil
	}

	return b

}
