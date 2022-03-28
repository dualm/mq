package pubsub

import (
	"context"
	"fmt"

	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func NewPubsub() *PubSub {
	return &PubSub{
		eveChan:  make(chan mq.MqMessage, rabbitmq.ChanBufferSize),
		reqChan:  make(chan mq.MqMessage, rabbitmq.ChanBufferSize),
		subChan:  make(chan struct{}, rabbitmq.ChanBufferSize),
		rspChan:  make(chan amqp.Delivery, rabbitmq.ChanBufferSize),
		errChan:  make(chan<- error, rabbitmq.ChanBufferSize),
		infoChan: make(chan<- string, rabbitmq.ChanBufferSize),
	}
}

type PubSub struct {
	eveChan  chan mq.MqMessage
	reqChan  chan mq.MqMessage
	subChan  chan struct{}
	rspChan  chan amqp.Delivery
	errChan  chan<- error
	infoChan chan<- string
}

func (ps *PubSub) Run(ctx context.Context, configID string,
	initconfig func(id string) *viper.Viper, infoChan chan<- string, errChan chan<- error) {
	ps.errChan = errChan
	ps.infoChan = infoChan

	conf := initconfig(configID)
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

	vhost := conf.GetString(rabbitmq.RbtVHost)
	targetExchange := conf.GetString(rabbitmq.RbtTargetExchange)
	routingKey := conf.GetString(rabbitmq.RbtTargetRoutingKey)
	rspQueue := conf.GetString(rabbitmq.RbtEapQueue)

	// event
	go func() {
		publish(
			ctx, redial(ctx, url, targetExchange, vhost, infoChan, errChan),
			targetExchange, routingKey, rspQueue, ps.eveChan, infoChan, errChan)
	}()

	go func() {
		publish(
			ctx, redial(ctx, url, targetExchange, vhost, infoChan, errChan),
			targetExchange, routingKey, rspQueue, ps.reqChan, infoChan, errChan)
	}()

	go func() {
		subscribe(
			ctx, redial(ctx, url, targetExchange, vhost, infoChan, errChan),
			targetExchange, rspQueue, rspQueue, ps.rspChan, ps.subChan, infoChan, errChan)
	}()
}

func (ps *PubSub) Send(ctx context.Context, c chan<- []byte, msg []mq.MqMessage) {
	for i := range msg {
		if msg[i].IsEvent {
			ps.send(msg[i])
		} else {
			go ps.sendRequest(ctx, msg[i], c)
		}
	}
}

func (ps *PubSub) send(msg mq.MqMessage) {
	ps.eveChan <- msg
}

func (ps *PubSub) sendRequest(ctx context.Context, msg mq.MqMessage, rsp chan<- []byte) {
	// 发送数据
	ps.subChan <- struct{}{}
	ps.reqChan <- msg

	select {
	case <-ctx.Done():
		ps.errChan <- fmt.Errorf("Mes请求消息发送超时: %s", msg.CorraltedId)

		rsp <- nil

		return
	case delivery := <-ps.rspChan:
		if delivery.CorrelationId == msg.CorraltedId {
			rsp <- delivery.Body

			return
		}

		ps.rspChan <- delivery
	}
}
