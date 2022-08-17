package pubsub

import (
	"context"
	"fmt"

	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type pubSub struct {
	*PubSubOption
	reqChan  chan mq.MqMessage
	subChan  chan rabbitmq.Subscription
	rspChan  chan amqp.Delivery
	errChan  chan<- error
	infoChan chan<- string
}

type PubSubOption struct {
	User           string
	Password       string
	Host           string
	Port           string
	Vhost          string
	TargetExchange string
	RoutingKey     string
	RspQueue       string
}

func New(opt *PubSubOption, infoChan chan<- string, errChan chan<- error) mq.Mq {
	return &pubSub{
		PubSubOption: opt,
		reqChan:      make(chan mq.MqMessage, rabbitmq.ChanBufferSize),
		subChan:      make(chan rabbitmq.Subscription, rabbitmq.ChanBufferSize),
		rspChan:      make(chan amqp.Delivery, rabbitmq.ChanBufferSize),
		errChan:      errChan,
		infoChan:     infoChan,
	}
}

func (ps *pubSub) Run(ctx context.Context) error {
	url := fmt.Sprintf(
		rabbitmq.URLFORMAT,
		ps.User,
		ps.Password,
		ps.Host,
		ps.Port,
	)

	go func() {
		publish(
			ctx, redial(ctx, url, ps.TargetExchange, ps.Vhost, ps.infoChan, ps.errChan),
			ps.TargetExchange, ps.RoutingKey, ps.RspQueue, ps.reqChan, ps.infoChan, ps.errChan)
	}()

	go func() {
		subscribe(
			ctx, redial(ctx, url, ps.TargetExchange, ps.Vhost, ps.infoChan, ps.errChan),
			ps.TargetExchange, ps.RspQueue, ps.RspQueue, ps.rspChan, ps.subChan, ps.infoChan, ps.errChan)
	}()

	return nil
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

func (ps *pubSub) Close(_ context.Context) error {
	return nil
}
