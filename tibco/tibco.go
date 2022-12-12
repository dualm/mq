package tibco

type MqMessage struct {
	Msg           []byte
	CorrelationID string
	IsEvent       bool
}

type MqResponse struct {
	Msg []byte
	Err error
}
