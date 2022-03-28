package rpc

import (
	"context"
	"fmt"

	"github.com/dualm/mq"

	"github.com/dualm/mq/rabbitmq"

	"github.com/spf13/viper"
)

func NewRpc() *Rpc {
	return &Rpc{
		requestChan: make(chan mq.MqMessage, rabbitmq.ChanBufferSize),
		infoChan:    make(chan<- string),
		errChan:     make(chan<- error),
	}
}

type Rpc struct {
	requestChan chan mq.MqMessage
	infoChan    chan<- string
	errChan     chan<- error
}

func (r *Rpc) Run(ctx context.Context, configId string, initconfig func(id string) *viper.Viper,
	infoChan chan<- string, errChan chan<- error) {
	r.infoChan = infoChan
	r.errChan = errChan

	conf := initconfig(configId)
	if conf == nil {
		errChan <- fmt.Errorf("nil config")

		return
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

	go func() {
		publish(ctx, redial(ctx, url, queue, vhost, infoChan, errChan),
			queue, r.requestChan, infoChan, errChan)
	}()
}

func (r *Rpc) Send(_ chan<- []byte, msg mq.MqMessage) {
	r.requestChan <- msg
}
