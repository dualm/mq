package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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

	b, err := SendToSpc(r, mq.MqMessage{Msg: dvMsg, IsEvent: true})
	if err != nil {
		log.Println("Error", err)
	}

	log.Println("response", string(b))
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

func SendToSpc(m mq.Mq, msg mq.MqMessage) ([]byte, error) {
	ctxSend, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	rsp := make(chan mq.MqResponse, 1)

	select {
	case <-ctxSend.Done():
		return nil, fmt.Errorf("send to SPC time out, Error: %w", ctxSend.Err())
	case <-m.Send(ctxSend, rsp, []mq.MqMessage{msg}):
		ctxRecv, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		return WaitRequest(ctxRecv, rsp)
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
			log.Println("Get Response", string(r.Msg))
		}

		return r.Msg, nil
	}
}
