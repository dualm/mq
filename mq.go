package mq

import (
	"context"

	"github.com/spf13/viper"
)

type MqMessage struct {
	Msg         []byte
	CorraltedId string
	IsEvent     bool
}

type MqResponse struct {
	Msg []byte
	Err error
}

type ConfigFunc func(id string) (*viper.Viper, error)

type Mq interface {
	// Run, configID为配置文件的文件名称，keys为各层节点
	Run(ctx context.Context, initConfig ConfigFunc, configID string, nodes string) (map[string]string, error)
	Send(ctx context.Context, responseChan chan<- MqResponse, msg []MqMessage) <-chan struct{}
	Close(ctx context.Context)
}
