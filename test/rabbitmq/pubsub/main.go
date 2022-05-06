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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	infoChan := make(chan string, 10)
	errChan := make(chan error, 10)

	pb := pubsub.New(infoChan, errChan)
	if paras, err := pb.Run(ctx, InitConfig, "config", "MES-DEV"); err != nil {
		log.Fatal(err)
	} else {
		log.Println(paras)
	}

	rsp := make(chan mq.MqResponse)

	msg := NewAreYouThereReq("TestEQP", "ZLFMM")

	log.Println(string(msg.Msg))

	ctxQuery, cancalQuery := context.WithTimeout(ctx, 20*time.Second)
	defer cancalQuery()

	pb.Send(ctxQuery, rsp, []mq.MqMessage{
		msg,
		NewMachineStateChangeEve("Z1TVIS02", "ZLFMM", 1),
	})

	pb.Send(ctx, rsp, []mq.MqMessage{
		NewMachineStateChangeEve("Z1TVIS02", "ZLFMM", 1),
	})

	for {
		select {
		case <-ctx.Done():
			return
		case <-ctxQuery.Done():
			return
		case info := <-infoChan:
			log.Println("info", info)
		case err := <-errChan:
			log.Fatal("Error", err)
		case re := <-rsp:
			log.Println("response", string(re.Msg))
		}
	}
}

func InitConfig(configId string) (*viper.Viper, error) {
	conf := viper.New()

	conf.SetConfigType("toml")
	conf.SetConfigName(strings.ToLower(configId))
	conf.AddConfigPath(".")

	err := conf.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return conf, nil
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

type machineStateChangedEve struct {
	XMLName          xml.Name `xml:"Message"`
	Header           Header   `xml:"Header"`
	MachineName      string   `xml:"Body>MACHINENAME"`
	MachineStateName string   `xml:"Body>MACHINESTATENAME"`
	ReasonCode       string   `xml:"Body>REASONCODE"`
}

func (eve *machineStateChangedEve) MarshalByte() []byte {
	return MarshalByte(eve)
}

func NewMachineStateChangeEve(machineName, factoryName string, state uint16) mq.MqMessage {
	messageName := "MachineStateChanged"
	header, _ := NewHeader(messageName, machineName, factoryName)
	eve := machineStateChangedEve{
		Header:      header,
		MachineName: machineName,
		MachineStateName: func(state uint16) string {
			switch state {
			case 0:
				return "ETC"
			case 4:
				return "IDLE"
			case 2:
				return "RUN"
			case 8, 10, 12:
				return "DOWN"
			default:
				return ""
			}
		}(state),
		ReasonCode: "",
	}

	return mq.MqMessage{
		Msg:         eve.MarshalByte(),
		CorraltedId: "",
		IsEvent:     true,
	}
}
