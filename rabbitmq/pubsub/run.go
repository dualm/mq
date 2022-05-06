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
		eveChan:  make(chan mq.MqMessage, rabbitmq.ChanBufferSize),
		reqChan:  make(chan mq.MqMessage, rabbitmq.ChanBufferSize),
		subChan:  make(chan struct{}, rabbitmq.ChanBufferSize),
		rspChan:  make(chan amqp.Delivery, rabbitmq.ChanBufferSize),
		errChan:  errChan,
		infoChan: infoChan,
	}
}

type pubSub struct {
	eveChan  chan mq.MqMessage
	reqChan  chan mq.MqMessage
	subChan  chan struct{}
	rspChan  chan amqp.Delivery
	errChan  chan<- error
	infoChan chan<- string
}

func (ps *pubSub) Run(ctx context.Context, initConfig mq.ConfigFunc, configID string, keys ...string) (map[string]string, error) {
	conf, err := initConfig(configID)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq/pubsub init config error, Error: %w", err)
	}

	if conf == nil {
		return nil, fmt.Errorf("rabbitmq/pubsub nil config")
	}

	url := fmt.Sprintf(
		rabbitmq.URLFORMAT,
		common.GetString(conf, rabbitmq.RbtUsername, keys...),
		common.GetString(conf, rabbitmq.RbtPassword, keys...),
		common.GetString(conf, rabbitmq.RbtHost, keys...),
		common.GetString(conf, rabbitmq.RbtPort, keys...),
	)

	vhost := common.GetString(conf, rabbitmq.RbtVHost, keys...)
	targetExchange := common.GetString(conf, rabbitmq.RbtTargetExchange, keys...)
	routingKey := common.GetString(conf, rabbitmq.RbtTargetRoutingKey, keys...)
	rspQueue := common.GetString(conf, rabbitmq.RbtClientQueue, keys...)

	// event
	go func() {
		publish(
			ctx, redial(ctx, url, targetExchange, vhost, ps.infoChan, ps.errChan),
			targetExchange, routingKey, rspQueue, ps.eveChan, ps.infoChan, ps.errChan)
	}()

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

func (ps *pubSub) Send(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) {
	go ps.send(ctx, responseChan, msg)
}

func (ps *pubSub) send(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) {
	for i := range msg {
		if len(msg[i].Msg) == 0 {
			if responseChan != nil {
				responseChan <- mq.MqResponse{}
			}

			continue
		}

		if msg[i].IsEvent {
			ps.sendEvent(msg[i])

			if responseChan != nil {
				responseChan <- mq.MqResponse{}
			}
		} else {
			go ps.sendRequest(ctx, msg[i], responseChan)
		}
	}
}

func (ps *pubSub) sendEvent(msg mq.MqMessage) {
	ps.eveChan <- msg
}

func (ps *pubSub) sendRequest(ctx context.Context, msg mq.MqMessage, rsp chan<- mq.MqResponse) {
	// 发送数据
	ps.subChan <- struct{}{}
	ps.reqChan <- msg

	select {
	case <-ctx.Done():
		rsp <- mq.MqResponse{
			Msg: nil,
			Err: fmt.Errorf("rabbitmq/pubsub SendRequest timeout: %s", msg.CorraltedId),
		}

		return
	case delivery := <-ps.rspChan:
		if delivery.CorrelationId == msg.CorraltedId {
			rsp <- mq.MqResponse{
				Msg: delivery.Body,
				Err: nil,
			}

			return
		}

		ps.rspChan <- delivery
	}
}

func (ps *pubSub) Close(_ context.Context) {
}
