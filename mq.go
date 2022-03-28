package mq

type MqMessage struct {
	Msg         []byte
	CorraltedId string
	IsEvent     bool
}
