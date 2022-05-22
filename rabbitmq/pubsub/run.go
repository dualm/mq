package pubsub

import (
	"context"
	"fmt"

	"github.com/dualm/common"
	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func New(infoChan chan<- string, errChan chan<- error) mq.Mq {
	return &pubSub{
		reqChan:  make(chan mq.MqMessage, rabbitmq.ChanBufferSize),
		subChan:  make(chan rabbitmq.Subscription, rabbitmq.ChanBufferSize),
		rspChan:  make(chan amqp.Delivery, rabbitmq.ChanBufferSize),
		errChan:  errChan,
		infoChan: infoChan,
	}
}

type pubSub struct {
	reqChan  chan mq.MqMessage
	subChan  chan rabbitmq.Subscription
	rspChan  chan amqp.Delivery
	errChan  chan<- error
	infoChan chan<- string
}

func (ps *pubSub) Run(ctx context.Context, initConfig mq.ConfigFunc, configID string, keys string) (map[string]string, error) {
	conf, err := initConfig(configID)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq/pubsub init config error, Error: %w", err)
	}

	if conf == nil {
		return nil, fmt.Errorf("rabbitmq/pubsub nil config")
	}

	url := fmt.Sprintf(
		rabbitmq.URLFORMAT,
		common.GetString(conf, keys, rabbitmq.RbtUsername),
		common.GetString(conf, keys, rabbitmq.RbtPassword),
		common.GetString(conf, keys, rabbitmq.RbtHost),
		common.GetString(conf, keys, rabbitmq.RbtPort),
	)

	vhost := common.GetString(conf, keys, rabbitmq.RbtVHost)
	targetExchange := common.GetString(conf, keys, rabbitmq.RbtTargetExchange)
	routingKey := common.GetString(conf, keys, rabbitmq.RbtTargetRoutingKey)
	rspQueue := common.GetString(conf, keys, rabbitmq.RbtClientQueue)

	go func() {
		publish(
			ctx, redial(ctx, url, targetExchange, vhost, ps.infoChan, ps.errChan),
			targetExchange, routingKey, rspQueue, ps.reqChan, ps.infoChan, ps.errChan)
	}()

	go func() {
		subscribe(
			ctx, redial(ctx, url, targetExchange, vhost, ps.infoChan, ps.errChan),
			targetExchange, rspQueue, rspQueue, ps.rspChan, ps.subChan, ps.infoChan, ps.errChan)
	}()

	return map[string]string{
		"url":            url,
		"VHost":          vhost,
		"TargetExchange": targetExchange,
		"RoutingKey":     routingKey,
		"RspQueue":       rspQueue,
	}, nil
}

func (ps *pubSub) Send(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) <-chan struct{} {
	c := make(chan struct{})

	go func() {
		ps.send(ctx, responseChan, msg)
		c <- struct{}{}
	}()

	return c
}

func (ps *pubSub) send(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) {
	for i := range msg {
		if len(msg[i].Msg) == 0 {
			rabbitmq.SendResponse(mq.MqResponse{}, responseChan, ps.errChan)

			continue
		}

		if msg[i].IsEvent {
			ps.sendEvent(msg[i])

			rabbitmq.SendResponse(mq.MqResponse{}, responseChan, ps.errChan)
		} else {
			ps.sendRequest(ctx, msg[i], responseChan)
		}
	}
}

func (ps *pubSub) sendEvent(msg mq.MqMessage) {
	ps.reqChan <- msg
}

func (ps *pubSub) sendRequest(ctx context.Context, msg mq.MqMessage, rsp chan<- mq.MqResponse) {
	select {
	case <-ctx.Done():
		ps.errChan <- fmt.Errorf("send subscription error, Error: %w", ctx.Err())
	case ps.subChan <- rabbitmq.Subscription{CorrelationID: msg.CorrelationID, RspChan: rsp}:
		select {
		case <-ctx.Done():
			ps.errChan <- fmt.Errorf("send request error, Error: %w", ctx.Err())
		case ps.reqChan <- msg:
			return
		}
	}
}

func (ps *pubSub) Close(_ context.Context) {
}
