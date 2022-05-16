package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/dualm/mq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	URLFORMAT           = "amqp://%s:%s@%s:%s"
	RbtUsername         = "Username"
	RbtPassword         = "Password"
	RbtHost             = "Host"
	RbtPort             = "Port"
	RbtTargetExchange   = "TargetExchange"
	RbtTargetRoutingKey = "TargetRoutingKey"
	RbtQueue            = "Queue"
	RbtClientQueue      = "EapQueue"
	RbtVHost            = "Vhost"

	ChanBufferSize = 0
)

type Session struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

type Subscription struct {
	Name    string
	RspChan chan<- mq.MqResponse
}

func (s *Session) Close() error {
	if s.Connection == nil {
		return nil
	}

	return s.Connection.Close()
}

func SendResponse(rsp mq.MqResponse, rspChan chan<- mq.MqResponse, errChan chan<- error) {
	subCtx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	select {
	case <-subCtx.Done():
		if err := subCtx.Err(); err == context.Canceled {
			errChan <- fmt.Errorf("send response canceled")
		} else if err == context.DeadlineExceeded {
			errChan <- fmt.Errorf("send response time out")
		}
	case rspChan <- rsp:
		return
	}
}
