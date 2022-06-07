package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
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
	if paras, err := pb.Run(ctx, InitConfig, "config", "MES-PRD"); err != nil {
		log.Fatal(err)
	} else {
		log.Println(paras)
	}

	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_1", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

		// _, err = SendToMes(pb, NewAreYouThereReq("TestEQP_1", "ZLFMM"))
		// if err != nil {
		// 	log.Println("Error", err)
		// }

	// 	msg = NewMachineStateChangeEve("TestEQP_1", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}
	// }()

	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_2", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	_, err = SendToMes(pb, NewAreYouThereReq("TestEQP_2", "ZLFMM"))
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	msg = NewMachineStateChangeEve("TestEQP_2", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}
	// }()

	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_3", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	_, err = SendToMes(pb, NewAreYouThereReq("TestEQP_3", "ZLFMM"))
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	msg = NewMachineStateChangeEve("TestEQP_3", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}
	// }()

	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_4", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	_, err = SendToMes(pb, NewAreYouThereReq("TestEQP_4", "ZLFMM"))
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	msg = NewMachineStateChangeEve("TestEQP_4", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// }()

	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_5", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	_, err = SendToMes(pb, NewAreYouThereReq("TestEQP_5", "ZLFMM"))
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	msg = NewMachineStateChangeEve("TestEQP_5", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// }()
	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_6", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	_, err = SendToMes(pb, NewAreYouThereReq("TestEQP_6", "ZLFMM"))
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	msg = NewMachineStateChangeEve("TestEQP_6", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// }()
	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_7", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	_, err = SendToMes(pb, NewAreYouThereReq("TestEQP_7", "ZLFMM"))
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	msg = NewMachineStateChangeEve("TestEQP_7", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// }()
	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_8", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	_, err = SendToMes(pb, NewAreYouThereReq("TestEQP_8", "ZLFMM"))
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	msg = NewMachineStateChangeEve("TestEQP_8", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// }()
	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_9", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	_, err = SendToMes(pb, NewAreYouThereReq("TestEQP_9", "ZLFMM"))
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	msg = NewMachineStateChangeEve("TestEQP_9", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// }()
	// go func() {
	// 	msg := NewAreYouThereReq("TestEQP_10", "ZLFMM")

	// 	_, err := SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	_, err = SendToMes(pb, NewAreYouThereReq("TestEQP_10", "ZLFMM"))
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// 	msg = NewMachineStateChangeEve("TestEQP_10", "ZLFMM", 1)

	// 	_, err = SendToMes(pb, msg)
	// 	if err != nil {
	// 		log.Println("Error", err)
	// 	}

	// }()

	infoReq := NewPieceInfoDownloadReq("OFFLINEFILE", "ZLFMM", "ABAB012B5V030S012P1", "")

	re, err := SendToMes(pb, infoReq)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(re))

	log.Println(PieceInfoDownloadDecode(re))

	<-ctx.Done()
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
	lock.Lock()
	defer lock.Unlock()

	messageName := "AreYouThereRequest"
	header, id := NewHeader(messageName, machineName, factoryName)
	areYouThereReq := AreYouThereReq{
		Header:      header,
		MachineName: machineName,
	}

	return mq.MqMessage{
		Msg:           areYouThereReq.MarshalByte(),
		CorrelationID: id,
		IsEvent:       false,
	}
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

var lock sync.Mutex

func GetRandStr(l int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

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
		Msg:           eve.MarshalByte(),
		CorrelationID: "",
		IsEvent:       true,
	}
}

func WaitRequest(ctx context.Context, rsp <-chan mq.MqResponse) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("request timeout, Error: %w", ctx.Err())
	case r, ok := <-rsp:
		if !ok {
			return nil, fmt.Errorf("invalid response channel")
		}
		if r.Err != nil {
			return nil, fmt.Errorf("request error, Error: %w", r.Err)
		}

		if len(r.Msg) != 0 {
			log.Println("Get Response")
		}

		return r.Msg, nil
	}
}

func SendToMes(m mq.Mq, msg mq.MqMessage) ([]byte, error) {
	ctxSend, cancal := context.WithTimeout(context.Background(), time.Second*20)
	defer cancal()

	rsp := make(chan mq.MqResponse, 1)

	select {
	case <-ctxSend.Done():
		return nil, fmt.Errorf("send to MES time out, Error: %w", ctxSend.Err())
	case <-m.Send(ctxSend, rsp, []mq.MqMessage{msg}):
		log.Println("Send Mes Message", string(msg.Msg))

		ctxRecv, cancal := context.WithTimeout(context.Background(), time.Second*20)
		defer cancal()

		return WaitRequest(ctxRecv, rsp)
	}
}

type PieceInfoDownloadRequest struct {
	XMLName              xml.Name   `xml:"Message"`
	Header               Header `xml:"Header"`
	MachineName          string     `xml:"Body>MACHINENAME"`
	LotName              string     `xml:"Body>LOTNAME"`
	ProductType          string     `xml:"Body>PRODUCTTYPE"`
	CarrierName          string     `xml:"Body>CARRIERNAME"`
	ProcessOperationName string     `xml:"Body>PROCESSOPERATIONNAME"`
	ProductSpecName      string     `xml:"Body>PRODUCTSPECNAME"`
	ProductRequestName   string     `xml:"Body>PRODUCTREQUESTNAME"`
	ProductionType       string     `xml:"Body>PRODUCTIONTYPE"`
	MachineRecipeName    string     `xml:"Body>MACHINERECIPENAME"`
	LotJudge             string     `xml:"Body>LOTJUDGE"`
	LotGrade             string     `xml:"Body>LOTGRADE"`
	Length               string     `xml:"Body>LENGTH"`
	Location             string     `xml:"Body>LOCATION"`
	Result               string     `xml:"Body>RESULT"`
	ResultDescription    string     `xml:"Body>RESULTDESCRIPTION"`
	ReturnCode           string     `xml:"Return>RETURNCODE"`
	ReturnMessage        string     `xml:"Return>RETURNMESSAGE"`
}

func (req *PieceInfoDownloadRequest) MarshalByte() []byte {
	return MarshalByte(req)
}

func NewPieceInfoDownloadReq(machineName, factoryName, pieceName, carrierName string) mq.MqMessage {
	messageName := "PieceInfoDownloadRequest"

	header, id := NewHeader(messageName, machineName, factoryName)

	req := PieceInfoDownloadRequest{
		Header:               header,
		MachineName:          machineName,
		LotName:              pieceName,
		ProductType:          "PRODUCTION",
		CarrierName:          carrierName,
		ProcessOperationName: "",
		ProductSpecName:      "",
		ProductRequestName:   "",
		ProductionType:       "",
		MachineRecipeName:    "",
		LotJudge:             "",
		LotGrade:             "",
		Length:               "",
		Location:             "",
		Result:               "",
		ResultDescription:    "",
		ReturnCode:           "",
		ReturnMessage:        "",
	}

	return mq.MqMessage{
		Msg:           req.MarshalByte(),
		CorrelationID: id,
		IsEvent:       false,
	}
}

func PieceInfoDownloadDecode(b []byte) (*PieceInfoDownloadRequest, error) {
	if b == nil {
		log.Fatal("nil response")
	}

	rsp := new(PieceInfoDownloadRequest)

	if err := xml.Unmarshal(b, &rsp); err != nil {
		return nil, err
	}

	if rsp.ReturnCode != "0" || rsp.Result != "0" {
		log.Fatal(rsp.ReturnCode)
	}

	return rsp, nil
}
