package main

import (
	"encoding/xml"
	"fmt"
	"github.com/dualm/mq"
	"github.com/dualm/mq/tibco"
	"log"
	"time"

	"github.com/spf13/viper"
)

var v *viper.Viper

func main() {
	var err error
	v, err = InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	tibcoService := v.GetString("TibcoService")
	tibcoNetwork := v.GetString("TibcoNetwork")
	tibcoDaemons := v.GetStringSlice("TibcoDaemon")

	if tibcoService == "" || tibcoNetwork == "" {
		log.Fatal("Tibco参数未配置，请联系管理员")
	}

	if err := tibco.TibInit(tibcoService, tibcoNetwork, tibcoDaemons); err != nil {
		panic(err)
	}
	defer tibco.TibDestroy()

	rspChan := make(chan []byte)
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

	log.Println(string(req[0].Msg.MarshalByte()))

	go tibco.Send(v.GetString("TargetSubjectName"), v.GetString("TibcoFieldName"), rspChan, req)

	result := <-rspChan

	_, err = RecipeValidationDecoder(result)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Recipe Verified")
}

type RecipeValidationReuqest struct {
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

func (req *RecipeValidationReuqest) MarshalByte() []byte {
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
	_rcpReq := RecipeValidationReuqest{
		Header:        header,
		MachineName:   header.MachineName,
		LotName:       lotName,
		CarrierName:   "",
		RecipeName:    recipeName,
		ParameterList: parameters,
	}

	return []mq.MqMessage{
		{
			Msg:     &_rcpReq,
			IsEvent: false,
		},
	}, nil
}

func RecipeValidationDecoder(b []byte) (*RecipeValidationReuqest, error) {
	if b == nil {
		return nil, fmt.Errorf("invalid recipe validation result")
	}

	rsp := new(RecipeValidationReuqest)

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

func InitConfig() (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func CheckConfig(conf *viper.Viper, key, description string) string {
	value := conf.GetString(key)
	if value == "" {
		log.Fatalln(description + "未设置，请咨询管理员")
	}

	return value
}

func getTransactionId() string {
	t := time.Now()
	micro := t.Nanosecond() / 1000

	return fmt.Sprintf("%s%06d", t.Format("20060102150405"), micro)
}
