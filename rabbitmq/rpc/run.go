package rpc

import (
	"context"
	"fmt"

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

func (r *rpc) Run(ctx context.Context, configId string, initConfig mq.ConfigFunc) error {
	conf, err := initConfig(configId)
	if err != nil {
		return fmt.Errorf("rabbitmq/rpc init config error, Error: %w", err)
	}

	if conf == nil {
		return fmt.Errorf("rabbitmq/rpc nil config")
	}

	url := fmt.Sprintf(
		rabbitmq.URLFORMAT,
		conf.GetString(rabbitmq.RbtUsername),
		conf.GetString(rabbitmq.RbtPassword),
		conf.GetString(rabbitmq.RbtHost),
		conf.GetString(rabbitmq.RbtPort),
	)

	queue := conf.GetString(rabbitmq.RbtQueue)
	vhost := conf.GetString(rabbitmq.RbtVHost)
	clientQueue := conf.GetString(rabbitmq.RbtClientQueue)

	go func() {
		publish(ctx, redial(ctx, url, queue, vhost, r.infoChan, r.errChan),
			queue, clientQueue, r.requestChan, r.response, r.infoChan, r.errChan)
	}()

	return nil
}

func (r *rpc) Send(_ context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) {
	for i := range msg {
		if len(msg[i].Msg) == 0 {
			if responseChan != nil {
				responseChan <- mq.MqResponse{}
			}

			continue
		}

		m := msg[i]
		r.requestChan <- m
		x := <-r.response
		responseChan <- x
	}
}

func (r *rpc) Close(_ context.Context) {
}
