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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	infoChan := make(chan string)
	errChan := make(chan error)
	config, err := InitConfig("config")
	if err != nil {
		log.Fatal(err)
	}

	opt := tibco.TibOption{
		FieldName: config.GetString("APS.TibcoFieldName"),
		Service:   config.GetString("APS.TibcoService"),
		Network:   config.GetString("APS.TibcoNetwork"),
		Daemon:    config.GetStringSlice("APS.TibcoDaemon"),
	}

	tib := tibco.New(&opt, infoChan, errChan)

	if err := tib.Run(ctx); err != nil {
		log.Fatal(err)
	}

	rspChan := make(chan mq.MqResponse)

	// req, err := NewRecipeValidationRequest("ACAB661B15010S001", "TEST01", []Parameter{
	// 	{
	// 		ItemName:  "P01",
	// 		ItemValue: "1",
	// 	}, {
	// 		ItemName:  "P02",
	// 		ItemValue: "2",
	// 	},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(string(req[0].Msg))

	tib.Send(
		ctx,
		rspChan,
		[]mq.MqMessage{
			{
				Msg:     []byte(spcMessage),
				IsEvent: true,
			},
		},
		config.GetString("APS.TibcoTargetSubjectName"),
	)

	result := <-rspChan
	if err := result.Err; err != nil {
		log.Fatal(err)
	}

	log.Println(result)
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

var spcMessage string = `
<Header>
          <MESSAGENAME>DataCollectRequest</MESSAGENAME>
          <EVENTCOMMENT>DataCollectRequest</EVENTCOMMENT>
          <EVENTUSER>DataCollectRequest</EVENTUSER>
          <ORIGINALSOURCESUBJECTNAME>_INBOX.TXZRLJTPVB.20220607134044199372</ORIGINALSOURCESUBJECTNAME>
      </Header>
      <Body>
          <FACTORYNAME>TOKEN-3B</FACTORYNAME>
          <PRODUCTSPECNAME>spec-1</PRODUCTSPECNAME>
          <PROCESSFLOWNAME>-</PROCESSFLOWNAME>
          <PROCESSOPERATIONNAME>operation-1</PROCESSOPERATIONNAME>
          <MACHINENAME>MachineName</MACHINENAME>
          <MACHINERECIPENAME>-</MACHINERECIPENAME>
          <UNITNAME>UnitName</UNITNAME>
          <LOTNAME>1RB-E1S210363-29-1</LOTNAME>
          <PRODUCTNAME>1RB-E1S210363-29-1</PRODUCTNAME>
          <ITEMLIST>
              <ITEM>
                  <ITEMNAME>X1</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>409</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
              <ITEM>
                  <ITEMNAME>X2</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>413</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
              <ITEM>
                  <ITEMNAME>X3</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>406</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
              <ITEM>
                  <ITEMNAME>X4</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>408</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
              <ITEM>
                  <ITEMNAME>X5</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>409</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
              <ITEM>
                  <ITEMNAME>X6</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>407</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
              <ITEM>
                  <ITEMNAME>X7</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>403</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
              <ITEM>
                  <ITEMNAME>X8</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>406</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
              <ITEM>
                  <ITEMNAME>X9</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>403</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
              <ITEM>
                  <ITEMNAME>X</ITEMNAME>
                  <SITELIST>
                      <SITE>
                          <SITENAME>036</SITENAME>
                          <SITEVALUE>407</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>038</SITENAME>
                          <SITEVALUE>406</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>007</SITENAME>
                          <SITEVALUE>409</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>025</SITENAME>
                          <SITEVALUE>412</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>016</SITENAME>
                          <SITEVALUE>408</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>019</SITENAME>
                          <SITEVALUE>409</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>028</SITENAME>
                          <SITEVALUE>408</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>033</SITENAME>
                          <SITEVALUE>400</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>035</SITENAME>
                          <SITEVALUE>403</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>037</SITENAME>
                          <SITEVALUE>408</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>001</SITENAME>
                          <SITEVALUE>409</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>012</SITENAME>
                          <SITEVALUE>407</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>015</SITENAME>
                          <SITEVALUE>406</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>017</SITENAME>
                          <SITEVALUE>410</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>039</SITENAME>
                          <SITEVALUE>405</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>044</SITENAME>
                          <SITEVALUE>408</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>010</SITENAME>
                          <SITEVALUE>410</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>013</SITENAME>
                          <SITEVALUE>407</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>014</SITENAME>
                          <SITEVALUE>408</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>021</SITENAME>
                          <SITEVALUE>407</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>023</SITENAME>
                          <SITEVALUE>409</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>030</SITENAME>
                          <SITEVALUE>407</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>042</SITENAME>
                          <SITEVALUE>410</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>004</SITENAME>
                          <SITEVALUE>405</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>011</SITENAME>
                          <SITEVALUE>410</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>031</SITENAME>
                          <SITEVALUE>403</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>041</SITENAME>
                          <SITEVALUE>409</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>022</SITENAME>
                          <SITEVALUE>406</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>029</SITENAME>
                          <SITEVALUE>408</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>018</SITENAME>
                          <SITEVALUE>411</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>020</SITENAME>
                          <SITEVALUE>407</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>027</SITENAME>
                          <SITEVALUE>409</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>032</SITENAME>
                          <SITEVALUE>397</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>034</SITENAME>
                          <SITEVALUE>400</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>043</SITENAME>
                          <SITEVALUE>410</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>002</SITENAME>
                          <SITEVALUE>409</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>003</SITENAME>
                          <SITEVALUE>406</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>008</SITENAME>
                          <SITEVALUE>413</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>009</SITENAME>
                          <SITEVALUE>413</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>024</SITENAME>
                          <SITEVALUE>412</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>026</SITENAME>
                          <SITEVALUE>404</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>040</SITENAME>
                          <SITEVALUE>406</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>045</SITENAME>
                          <SITEVALUE>403</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>005</SITENAME>
                          <SITEVALUE>410</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                      <SITE>
                          <SITENAME>006</SITENAME>
                          <SITEVALUE>408</SITEVALUE>
                          <SAMPLEMATERIALNAME></SAMPLEMATERIALNAME>
                      </SITE>
                  </SITELIST>
              </ITEM>
          </ITEMLIST>
      </Body>
`
