package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"github.com/dualm/mq"
	"github.com/dualm/mq/tibco"

	"github.com/spf13/viper"
)

func main() {
	infoChan := make(chan string)
	errChan := make(chan error)

	tib := tibco.New(infoChan, errChan)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if m, err := tib.Run(ctx, InitConfig, "", ""); err != nil {
		log.Fatal(err)
	} else {
		log.Println(m)
	}

	rspChan := make(chan mq.MqResponse)

	req, err := NewRecipeValidationRequest("ACAB661B15010S001", "TEST01", []Parameter{
		{
			ItemName:  "P01",
			ItemValue: "1",
		}, {
			ItemName:  "P02",
			ItemValue: "2",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(req[0].Msg))

	go tib.Send(ctx, rspChan, req)

	result := <-rspChan
	if err := result.Err; err != nil {
		log.Fatal(err)
	}

	_, err = RecipeValidationDecoder(result.Msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Recipe Verified")
}

type RecipeValidationRequest struct {
	XMLName       xml.Name    `xml:"Message"`
	Header        Header      `xml:"Header"`
	MachineName   string      `xml:"Body>MACHINENAME"`
	LotName       string      `xml:"Body>LOTNAME"`
	CarrierName   string      `xml:"Body>CARRIERNAME"`
	RecipeName    string      `xml:"Body>RECIPENAME"`
	ParameterList []Parameter `xml:"Body>PARAMETERLIST>PARAMETER"`
	ReturnCode    string      `xml:"Return>RETURNCODE"`
	ReturnMessage string      `xml:"Return>RETURNMESSAGE"`
}

func (req *RecipeValidationRequest) MarshalByte() []byte {
	return MarshalByte(req)
}

func MarshalByte(v interface{}) []byte {
	b, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		return nil
	}

	return b

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
	Listener                  string   `xml:"listener"`
}

func NewHeader(messageName string) Header {
	v, err := InitConfig("config")
	if err != nil {
		log.Fatal(err)
	}

	return Header{
		MessageName:               messageName,
		ShopName:                  v.GetString("ShopName"),
		MachineName:               v.GetString("MachineName"),
		TransactionId:             getTransactionId(),
		OriginalSourceSubjectName: v.GetString("OriginalSourceSubjectName"),
		SourceSubjectName:         v.GetString("SourceSubjectName"),
		TargetSubjectName:         v.GetString("TargetSubjectName"),
		EventUser:                 v.GetString("EventUser"),
		EventComment:              v.GetString("EventComment"),
		Listener:                  v.GetString("Listener"),
	}
}

type Parameter struct {
	XMLName   xml.Name `xml:"PARAMETER"`
	ItemName  string   `xml:"ITEMNAME"`
	ItemValue string   `xml:"ITEMVALUE"`
}

func NewRecipeValidationRequest(lotName, recipeName string, parameters []Parameter) ([]mq.MqMessage, error) {
	header := NewHeader("RecipeParameterValidationRequest")
	_rcpReq := RecipeValidationRequest{
		Header:        header,
		MachineName:   header.MachineName,
		LotName:       lotName,
		CarrierName:   "",
		RecipeName:    recipeName,
		ParameterList: parameters,
	}

	return []mq.MqMessage{
		{
			Msg:     _rcpReq.MarshalByte(),
			IsEvent: false,
		},
	}, nil
}

func RecipeValidationDecoder(b []byte) (*RecipeValidationRequest, error) {
	if b == nil {
		return nil, fmt.Errorf("invalid recipe validation result")
	}

	rsp := new(RecipeValidationRequest)

	if err := xml.Unmarshal(b, rsp); err != nil {
		return nil, err
	}

	if rsp.ReturnCode != "0" {
		return nil,
			NewRecipeReturnCodeError(
				rsp.MachineName,
				rsp.ReturnCode,
				rsp.ReturnMessage,
				func(parameterList []Parameter) []string {
					re := make([]string, len(parameterList))

					for i := range parameterList {
						re[i] = parameterList[i].ItemName
					}

					return re
				}(rsp.ParameterList),
			)
	}

	return rsp, nil
}

func NewRecipeReturnCodeError(machineName, returnCode, returnMesssage string, ngParameters []string) error {
	return fmt.Errorf(
		"EQPId: %s, ReturnCode: %s, ReturnMessage: %s",
		machineName,
		returnCode,
		returnMesssage,
	)
}

func InitConfig(conf string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(conf)
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func getTransactionId() string {
	t := time.Now()
	micro := t.Nanosecond() / 1000

	return fmt.Sprintf("%s%06d", t.Format("20060102150405"), micro)
}
