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
	Run(ctx context.Context, configID string, initConfig ConfigFunc) error
	Send(ctx context.Context, c chan<- MqResponse, msg []MqMessage)
	Close(ctx context.Context)
}
