package tibco

/*
#include <stdlib.h>
#include "tibrv/cm.h"
extern void goCmCallback(tibrvcmEvent event, tibrvMsg message, void* closure);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type TibrvCmEventCallback interface {
	CmCallBack(CmEvent, Message)
}

//export goCmCallback
func goCmCallback(cmListener C.tibrvcmEvent, message C.tibrvMsg, p unsafe.Pointer) {
	cmCallBack := *(*TibrvCmEventCallback)(p)

	cmCallBack.CmCallBack(
		CmEvent{
			tibrvCmEvent: cmListener,
		},
		Message{
			tibrvMsg: message,
		},
	)
}

type CmEvent struct {
	tibrvCmEvent C.tibrvcmEvent
}

func NewCmListener(queue *Queue, transport *CmTransport, subject string, p TibrvCmEventCallback) (*CmEvent, error) {
	_cSubject := C.CString(subject)
	defer C.free(unsafe.Pointer(_cSubject))

	var (
		q     C.uint
		event C.tibrvcmEvent
	)

	if queue == nil {
		q = C.TIBRV_DEFAULT_QUEUE
	} else {
		q = queue.tibrvQueue
	}

	if status := C.tibrvcmEvent_CreateListener(&event, q, C.tibrvcmEventCallback(C.goCmCallback), transport.tibrvCmTransport, _cSubject, unsafe.Pointer(&p)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("create new cm listener error, %d", status)
	}

	return &CmEvent{
		tibrvCmEvent: event,
	}, nil
}

func (e *CmEvent) Close(cancelAgreements bool) error {
	var _cCancelAgreements C.tibrv_bool
	if cancelAgreements {
		_cCancelAgreements = C.TIBRV_TRUE
	} else {
		_cCancelAgreements = C.TIBRV_FALSE
	}

	if status := C.tibrvcmEvent_DestroyEx(e.tibrvCmEvent, _cCancelAgreements, nil); status != C.TIBRV_OK {
		return fmt.Errorf("close cmevent error, %d", status)
	}

	return nil
}

// ConfirmMsg 显式地对一个消息的接收进行确认。仅用于重写了默认的自动确认函数的场景.
func (e *CmEvent) ConfirmMsg(msg C.tibrvMsg) error {
	if status := C.tibrvcmEvent_ConfirmMsg(e.tibrvCmEvent, msg); status != C.TIBRV_OK {
		return fmt.Errorf("confirm event error, %d", status)
	}

	return nil
}

func (e *CmEvent) GetListenerSubject() (string, error) {
	var subject *C.char

	if status := C.tibrvcmEvent_GetListenerSubject(e.tibrvCmEvent, &subject); status != C.TIBRV_OK {
		return "", fmt.Errorf("get listener subject error, %d", status)
	}

	return C.GoString(subject), nil
}

func (e *CmEvent) GetListenerTransport() (*CmTransport, error) {
	var transport C.tibrvcmTransport

	if status := C.tibrvcmEvent_GetListenerTransport(e.tibrvCmEvent, &transport); status != C.TIBRV_OK {
		return nil, fmt.Errorf("get listener cmTransport error, %d", status)
	}

	return &CmTransport{
		tibrvCmTransport: transport,
	}, nil
}

func (e *CmEvent) GetQueue() (*Queue, error) {
	var queue C.tibrvQueue

	if status := C.tibrvcmEvent_GetQueue(e.tibrvCmEvent, &queue); status != C.TIBRV_OK {
		return nil, fmt.Errorf("get queue error, %d", status)
	}

	return &Queue{tibrvQueue: queue}, nil
}

func (e *CmEvent) SetExplicitConfirm() error {
	if status := C.tibrvcmEvent_SetExplicitConfirm(e.tibrvCmEvent); status != C.TIBRV_OK {
		return fmt.Errorf("confirm message error, %d", status)
	}

	return nil
}
