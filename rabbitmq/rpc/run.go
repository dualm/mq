package rpc

import (
	"context"
	"fmt"

	"github.com/dualm/mq"

	"github.com/dualm/mq/rabbitmq"

	"github.com/spf13/viper"
)

func Run(ctx context.Context, configId string, initconfig func(id string) *viper.Viper) {
	conf := initconfig(configId)
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
		publish(redial(ctx, url, queue, vhost), queue, requestChan)
	}()
}

func Send(_ chan<- []byte, msg mq.MqMessage) {
	requestChan <- msg
}
