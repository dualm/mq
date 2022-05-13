package rabbitmq

import (
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

	ChanBufferSize = 256
)

type Session struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func (s *Session) Close() error {
	if s.Connection == nil {
		return nil
	}

	return s.Connection.Close()
}
