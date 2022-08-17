package rpc

import (
	"context"
	"fmt"

	"github.com/dualm/mq"

	"github.com/dualm/mq/rabbitmq"
)

type rpc struct {
	*RpcOption
	requestChan chan mq.MqMessage
	response    chan mq.MqResponse
	infoChan    chan<- string
	errChan     chan<- error
}

type RpcOption struct {
	Username    string
	Password    string
	Host        string
	Port        string
	Queue       string
	VHost       string
	ClientQueue string
}

func New(opt *RpcOption, infoChan chan<- string, errChan chan<- error) mq.Mq {
	return &rpc{
		RpcOption:   opt,
		requestChan: make(chan mq.MqMessage),
		response:    make(chan mq.MqResponse),
		infoChan:    infoChan,
		errChan:     errChan,
	}
}

func (r *rpc) Run(ctx context.Context) error {
	url := fmt.Sprintf(
		rabbitmq.URLFORMAT,
		r.Username,
		r.Password,
		r.Host,
		r.Port,
	)

	go func() {
		publish(ctx, redial(ctx, url, r.Queue, r.VHost, r.infoChan, r.errChan),
			r.Queue, r.ClientQueue, r.requestChan, r.response, r.infoChan, r.errChan)
	}()

	return nil
}

func (r *rpc) Send(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) <-chan struct{} {
	c := make(chan struct{})

	go func() {
		r.send(ctx, responseChan, msg)
		c <- struct{}{}
	}()

	return c
}

func (r *rpc) send(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) {
	for i := range msg {
		if len(msg[i].Msg) == 0 {
			rabbitmq.SendResponse(mq.MqResponse{}, responseChan, r.errChan)

			continue
		}

		m := msg[i]

		select {
		case <-ctx.Done():
			return
		case r.requestChan <- m:
			select {
			case <-ctx.Done():
				return
			case responseChan <- <-r.response:
				continue
			}
		}
	}
}

func (r *rpc) Close(_ context.Context) error {
	return nil
}
