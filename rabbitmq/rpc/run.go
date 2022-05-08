package rpc

import (
	"context"
	"fmt"

	"github.com/dualm/common"
	"github.com/dualm/mq"

	"github.com/dualm/mq/rabbitmq"
)

func New(infoChan chan<- string, errChan chan<- error) mq.Mq {
	return &rpc{
		requestChan: make(chan mq.MqMessage),
		response:    make(chan mq.MqResponse),
		infoChan:    infoChan,
		errChan:     errChan,
	}
}

type rpc struct {
	requestChan chan mq.MqMessage
	response    chan mq.MqResponse
	infoChan    chan<- string
	errChan     chan<- error
}

func (r *rpc) Run(ctx context.Context, initConfig mq.ConfigFunc, configId string, keys string) (map[string]string, error) {
	conf, err := initConfig(configId)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq/rpc init config error, Error: %w", err)
	}

	if conf == nil {
		return nil, fmt.Errorf("rabbitmq/rpc nil config")
	}

	url := fmt.Sprintf(
		rabbitmq.URLFORMAT,
		common.GetString(conf, keys, rabbitmq.RbtUsername),
		common.GetString(conf, keys, rabbitmq.RbtPassword),
		common.GetString(conf, keys, rabbitmq.RbtHost),
		common.GetString(conf, keys, rabbitmq.RbtPort),
	)

	queue := common.GetString(conf, keys, rabbitmq.RbtQueue)
	vhost := common.GetString(conf, keys, rabbitmq.RbtVHost)
	clientQueue := common.GetString(conf, keys, rabbitmq.RbtClientQueue)

	go func() {
		publish(ctx, redial(ctx, url, queue, vhost, r.infoChan, r.errChan),
			queue, clientQueue, r.requestChan, r.response, r.infoChan, r.errChan)
	}()

	return map[string]string{
		"url":         url,
		"VHost":       vhost,
		"Queue":       queue,
		"ClientQueue": clientQueue,
	}, nil
}

func (r *rpc) Send(_ context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) {
	go r.send(context.TODO(), responseChan, msg)
}

func (r *rpc) send(_ context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) {
	for i := range msg {
		if len(msg[i].Msg) == 0 {
			if responseChan != nil {
				responseChan <- mq.MqResponse{}
			}

			continue
		}

		m := msg[i]
		r.requestChan <- m
		responseChan <- <-r.response
	}
}

func (r *rpc) Close(_ context.Context) {
}
