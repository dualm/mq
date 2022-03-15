package pubsub

import (
	"context"
	"fmt"

	"github.com/dualm/mq"
	"github.com/dualm/mq/rabbitmq"
	"github.com/spf13/viper"
)

func Run(ctx context.Context, configID string, initconfig func(id string) *viper.Viper) {
	conf := initconfig(configID)

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
		publish(redial(ctx, url, targetExchange, vhost), targetExchange, routingKey, rspQueue, eveChan)
	}()

	go func() {
		publish(redial(ctx, url, targetExchange, vhost), targetExchange, routingKey, rspQueue, reqChan)
	}()

	go func() {
		subscribe(redial(ctx, url, targetExchange, vhost), targetExchange, rspQueue, rspQueue, rspChan)
	}()
}

// todo 增加ctx.
func Send(c chan<- []byte, msg []mq.MqMessage) {
	if len(msg) == 0 {
		return
	}

	for i := range msg {
		if msg[i].IsEvent {
			send(msg[i])
		} else {
			go sendRequest(mq.RequestWaitingDuration, msg[i], c)
		}
	}
}
