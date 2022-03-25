package mq

import "time"

const (
	RequestWaitingDuration = 20 * time.Second
)

type MqMessage struct {
	Msg         []byte
	CorraltedId string
	IsEvent     bool
}
