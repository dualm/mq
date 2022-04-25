package main

import (
	"context"
	"log"
	"strings"

	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq/rpc"
	"github.com/dualm/zispc"
	"github.com/spf13/viper"
)

func main() {
	machineName := "TEST-M"
	recipe := ""
	spec := "TEST-P"
	flow := "TEST-F"
	lotname := "W1"
	operation := "TEST-O"
	unit := ""

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	infoChan := make(chan string, 10)
	errChan := make(chan error, 10)

	r := rpc.New(infoChan, errChan)
	if paras, err := r.Run(ctx, InitConfig, "config", "SPC"); err != nil {
		log.Fatal(err)
	} else {
		log.Println(paras)
	}

	zispc.UnsetWithS()

	items := make([]zispc.XMLItem, 0, 2)
	items = zispc.AddItemToXML(items, "para", "", lotname, map[string]string{"001": "8", "002": "8", "003": "8"})
	items = zispc.AddItemToXML(items, "HYJ", "", lotname, map[string]string{"001": "12", "002": "12", "003": "12"})

	dvMsg, err := zispc.NewXMLProcessData(
		10,
		machineName,
		recipe,
		unit,
		spec,
		flow,
		lotname,
		lotname,
		"ZLFMM",
		operation,
		nil,
		items,
	).Encode()
	if err != nil {
		log.Fatal(err)
	}

	rsp := make(chan mq.MqResponse, 1)

	r.Send(ctx, rsp, []mq.MqMessage{
		{
			Msg:     dvMsg,
			IsEvent: true,
		},
	})

	re := <-rsp
	log.Println(re)

	if re.Err != nil {
		log.Println(re.Err)
	}

	r.Send(ctx, rsp, []mq.MqMessage{
		{
			Msg:     nil,
			IsEvent: true,
		},
	})

	re = <-rsp

	log.Println(re)

	if re.Err != nil {
		log.Println(re.Err)
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

/*
<Message>
    <Header>
        <MESSAGENAME>DataCollectRequest</MESSAGENAME>
        <EVENTCOMMENT>DataCollectRequest</EVENTCOMMENT>
        <EVENTUSER>DataCollectRequest</EVENTUSER>
        <ORIGINALSOURCESUBJECTNAME>amq.rabbitmq.reply-to.g1h2AA5yZXBseUAyNzIyNTAxMAAAaQcAAACjYbwn+w==.Ah/RUZx2xeNcHnvcyAQG9A==, AMQPCorrelationID[30]</ORIGINALSOURCESUBJECTNAME>
    </Header>
    <Body>
        <FACTORYNAME>ZLFMM</FACTORYNAME>
        <PRODUCTSPECNAME>TEST-P</PRODUCTSPECNAME>
        <PROCESSFLOWNAME>TEST-F</PROCESSFLOWNAME>
        <PROCESSOPERATIONNAME>TEST-O</PROCESSOPERATIONNAME>
        <MACHINENAME>TEST-M</MACHINENAME>
        <MACHINERECIPENAME>-</MACHINERECIPENAME>
        <UNITNAME>-</UNITNAME>
        <LOTNAME>W1</LOTNAME>
        <PRODUCTNAME>W1</PRODUCTNAME>
        <ITEMLIST>
            <ITEM>
                <ITEMNAME>para</ITEMNAME>
                <SITELIST>
                    <SITE>
                        <SITENAME>001</SITENAME>
                        <SITEVALUE>8</SITEVALUE>
                        <SAMPLEMATERIALNAME>W1</SAMPLEMATERIALNAME>
                    </SITE>
                    <SITE>
                        <SITENAME>002</SITENAME>
                        <SITEVALUE>8</SITEVALUE>
                        <SAMPLEMATERIALNAME>W1</SAMPLEMATERIALNAME>
                    </SITE>
                    <SITE>
                        <SITENAME>003</SITENAME>
                        <SITEVALUE>8</SITEVALUE>
                        <SAMPLEMATERIALNAME>W1</SAMPLEMATERIALNAME>
                    </SITE>
                </SITELIST>
            </ITEM>
            <ITEM>
                <ITEMNAME>HYJ</ITEMNAME>
                <SITELIST>
                    <SITE>
                        <SITENAME>001</SITENAME>
                        <SITEVALUE>12</SITEVALUE>
                        <SAMPLEMATERIALNAME>W1</SAMPLEMATERIALNAME>
                    </SITE>
                    <SITE>
                        <SITENAME>002</SITENAME>
                        <SITEVALUE>12</SITEVALUE>
                        <SAMPLEMATERIALNAME>W1</SAMPLEMATERIALNAME>
                    </SITE>
                    <SITE>
                        <SITENAME>003</SITENAME>
                        <SITEVALUE>12</SITEVALUE>
                        <SAMPLEMATERIALNAME>W1</SAMPLEMATERIALNAME>
                    </SITE>
                </SITELIST>
            </ITEM>
        </ITEMLIST>
    </Body>
</Message>
*/
