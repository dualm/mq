package tibco

/*
#include "tibrv/tibrv.h"
#include "tibrv/tport.h"
#include "tibrv/types.h"
*/
import "C"
import "fmt"

type Transport struct {
	tibrvTransport C.tibrvTransport
}

func NewTransport() *Transport {
	return &Transport{
		tibrvTransport: 0,
	}
}

func (transport *Transport) create(service, network string, daemon []string) error {
	var _transport C.tibrvTransport
	var err error

	for i := range daemon {
		_transport, err = newTransport(
			C.CString(service),
			C.CString(network),
			C.CString(daemon[i]),
		)
		if err != nil {
			continue
		}

		transport.tibrvTransport = _transport

		return nil
	}

	_transport, err = newTransport(
		C.CString(service),
		C.CString(network),
		nil,
	)
	if err != nil {
		return fmt.Errorf("Create Transport error, %w", err)
	}

	transport.tibrvTransport = _transport

	return nil
}

func (transport *Transport) destroy() error {
	if status := C.tibrvTransport_Destroy(transport.tibrvTransport); status != C.TIBRV_OK {
		return fmt.Errorf("Destroy transport error, code: %d", status)
	}

	return nil
}

func newTransport(service, network, daemon *C.char) (C.tibrvTransport, error) {
	var transport C.tibrvTransport

	if status := C.tibrvTransport_Create(
		&transport,
		service,
		network,
		daemon,
	); status != C.TIBRV_OK {
		return 0, fmt.Errorf("Tibco create transport error, code: %d", status)
	}

	return transport, nil
}

func (transport *Transport) send(msg *Message) error {
	if status := C.tibrvTransport_Send(transport.tibrvTransport, msg.tibrvMsg); status != C.TIBRV_OK {
		return fmt.Errorf("Send error, code: %d", status)
	}

	return nil
}

func (transport *Transport) sendv(msg []Message) error {
	msgs := make([]C.tibrvMsg, len(msg))
	for i := range msg {
		msgs[i] = msg[i].tibrvMsg
	}

	if status := C.tibrvTransport_Sendv(transport.tibrvTransport, &msgs[0], C.uint(len(msgs))); status != C.TIBRV_OK {
		return fmt.Errorf("Sendv error, code: %d", status)
	}

	return nil
}

func (transport *Transport) sendReply(reply *Message, request *Message) error {
	if status := C.tibrvTransport_SendReply(transport.tibrvTransport, reply.tibrvMsg, request.tibrvMsg); status != C.TIBRV_OK {
		return fmt.Errorf("SendReply error, code: %d", status)
	}

	return nil
}

func (transport *Transport) sendRequest(msg *Message, timeOut float64) (*Message, error) {
	var _msg C.tibrvMsg
	if status := C.tibrvTransport_SendRequest(transport.tibrvTransport, msg.tibrvMsg, &_msg, C.double(timeOut)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("SendRequest error, code: %d", status)
	}

	return &Message{
		tibrvMsg: _msg,
	}, nil
}
