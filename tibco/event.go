package tibco

/*
#include <stdlib.h>
#include "tibrv/tibrv.h"
#include "tibrv/events.h"
extern void goCallback(tibrvEvent event, tibrvMsg message, void* closure);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type EventType uint32

const (
	ListenEvent EventType = iota
	TimerEvent
	IOEvent
)

type TibrvEventCallback interface {
	CallBack(Event, Message)
}

//export goCallback
func goCallback(event C.tibrvEvent, message C.tibrvMsg, p unsafe.Pointer) {
	callback := *(*TibrvEventCallback)(p)

	callback.CallBack(Event{
		tibrvEvent: event,
	}, Message{
		tibrvMsg: message,
	})
}

type Event struct {
	tibrvEvent C.tibrvEvent
}

// NewListener，创建一个message event。如果queue为nil，则使用默认Queue.
func NewListener(queue *Queue, transport *Transport, subject string, callback TibrvEventCallback) (*Event, error) {
	var event C.tibrvEvent
	var q C.uint

	if queue == nil {
		q = C.TIBRV_DEFAULT_QUEUE
	} else {
		q = queue.tibrvQueue
	}

	_cSubject := C.CString(subject)
	defer C.free(unsafe.Pointer(_cSubject))

	if status := C.tibrvEvent_CreateListener(&event, q, C.tibrvEventCallback(C.goCallback), transport.tibrvTransport, _cSubject, unsafe.Pointer(&callback)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("create listener error, %d", status)
	}

	return &Event{
		tibrvEvent: event,
	}, nil
}

func (e *Event) Close() error {
	if status := C.tibrvEvent_DestroyEx(e.tibrvEvent, nil); status != C.TIBRV_OK {
		return fmt.Errorf("event destroy error, %d", status)
	}

	return nil
}

func (e *Event) GetListenerSubject() (string, error) {
	var subject *C.char

	if status := C.tibrvEvent_GetListenerSubject(e.tibrvEvent, &subject); status != C.TIBRV_OK {
		return "", fmt.Errorf("get listener subject error, %d", status)
	}

	return C.GoString(subject), nil
}

func (e *Event) GetListenerTransport() (*Transport, error) {
	var t C.tibrvTransport

	if status := C.tibrvEvent_GetListenerTransport(e.tibrvEvent, &t); status != C.TIBRV_OK {
		return nil, fmt.Errorf("get listener transport error, %d", status)
	}

	return &Transport{tibrvTransport: t}, nil
}

func (e *Event) GetType() (EventType, error) {
	var t C.uint

	if status := C.tibrvEvent_GetType(e.tibrvEvent, &t); status != C.TIBRV_OK {
		return 255, fmt.Errorf("get event type error, %d", status)
	}

	return EventType(t), nil
}

func (e *Event) GetQueue() (*Queue, error) {
	var q C.tibrvQueue

	if status := C.tibrvEvent_GetQueue(e.tibrvEvent, &q); status != C.TIBRV_OK {
		return nil, fmt.Errorf("get event queue error, %d", status)
	}

	return &Queue{tibrvQueue: q}, nil
}

func (e *Event) GetTimeInterval() (float64, error) {
	var i C.double

	if status := C.tibrvEvent_GetTimerInterval(e.tibrvEvent, &i); status != C.TIBRV_OK {
		return -1, fmt.Errorf("get interval error, %d", status)
	}

	return float64(i), nil
}
