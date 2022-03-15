package mq

import "time"

const (
	RequestWaitingDuration = 20 * time.Second
)

type MqMessage struct {
	Msg         MqMessageBody
	CorraltedId string
	IsEvent     bool
}

type MqMessageBody interface {
	MarshalByte() []byte
}
