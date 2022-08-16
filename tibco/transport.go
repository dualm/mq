//go:build windows || linux
// +build windows linux

package tibco

/*
#include <stdlib.h>
#include "tibrv/tibrv.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Transport struct {
	tibrvTransport C.tibrvTransport
}

func NewTransport(service, network string, daemon []string) (*Transport, error) {
	var transport C.tibrvTransport
	var err error

	_cService := C.CString(service)
	_cNetwork := C.CString(network)

	defer C.free(unsafe.Pointer(_cService))
	defer C.free(unsafe.Pointer(_cNetwork))

	for i := range daemon {
		_cDaemon := C.CString(daemon[i])

		transport, err = newTransport(
			_cService,
			_cNetwork,
			_cDaemon,
		)
		if err != nil {
			C.free(unsafe.Pointer(_cDaemon))

			continue
		}

		C.free(unsafe.Pointer(_cDaemon))

		return &Transport{
			tibrvTransport: transport,
		}, nil
	}

	if len(daemon) != 0 {
		return nil, fmt.Errorf("daemons cannot connect")
	}

	transport, err = newTransport(
		_cService,
		_cNetwork,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Create Transport error, %w", err)
	}

	return &Transport{
		tibrvTransport: transport,
	}, nil
}

func (transport *Transport) Destroy() error {
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

func (transport *Transport) Send(msg *Message) error {
	if status := C.tibrvTransport_Send(transport.tibrvTransport, msg.tibrvMsg); status != C.TIBRV_OK {
		return fmt.Errorf("Send error, code: %d", status)
	}

	return nil
}

func (transport *Transport) Sendv(msg []*Message) error {
	msgs := make([]C.tibrvMsg, len(msg))
	for i := range msg {
		msgs[i] = msg[i].tibrvMsg
	}

	if status := C.tibrvTransport_Sendv(transport.tibrvTransport, &msgs[0], C.uint(len(msgs))); status != C.TIBRV_OK {
		return fmt.Errorf("Sendv error, code: %d", status)
	}

	return nil
}

func (transport *Transport) SendReply(reply *Message, request *Message) error {
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
