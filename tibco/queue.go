package tibco

/*
#include <stdlib.h>
#include "tibrv/tibrv.h"
extern void goQueueHook(tibrvQueue queue, void* closure);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type TibrvQueueHook interface {
	CallHook(Queue)
}

//export goQueueHook
func goQueueHook(queue C.tibrvQueue, p unsafe.Pointer) {
	callback := *(*TibrvQueueHook)(p)

	callback.CallHook(Queue{
		tibrvQueue: queue,
	})
}

type Queue struct {
	tibrvQueue C.tibrvQueue
}

func NewQueue() (*Queue, error) {
	var q C.tibrvQueue

	if status := C.tibrvQueue_Create(&q); status != C.TIBRV_OK {
		return nil, fmt.Errorf("Create queue error, %d", status)
	}
	return &Queue{
		tibrvQueue: q,
	}, nil
}

func (q *Queue) Destroy() error {
	if status := C.tibrvQueue_DestroyEx(q.tibrvQueue, nil, nil); status != C.TIBRV_OK {
		return fmt.Errorf("destroy queue error, %d", status)
	}

	return nil
}

// Dispatch, 如果Queue不为空时分发(dispatch)队列头部的event后返回。如果Queue为空，则一直阻塞到Queue收到新的event为止.
func (q *Queue) Dispatch() error {
	if status := C.tibrvQueue_TimedDispatch(q.tibrvQueue, C.TIBRV_WAIT_FOREVER); status != C.TIBRV_OK {
		return fmt.Errorf("dispatch event error, %d", status)
	}

	return nil
}

func (q *Queue) Len() (int, error) {
	var n C.uint

	if status := C.tibrvQueue_GetCount(q.tibrvQueue, &n); status != C.TIBRV_OK {
		return 0, fmt.Errorf("get queue length error, %d", status)
	}

	return int(n), nil

}

func (q *Queue) GetHook() (C.tibrvQueueHook, error) {
	var f C.tibrvQueueHook

	if status := C.tibrvQueue_GetHook(q.tibrvQueue, &f); status != C.TIBRV_OK {
		return nil, fmt.Errorf("get queue hook error, %d", status)
	}

	return f, nil
}

func (q *Queue) GetName() (string, error) {
	var _cName *C.char

	if status := C.tibrvQueue_GetName(q.tibrvQueue, &_cName); status != C.TIBRV_OK {
		return "", fmt.Errorf("get queue name error, %d", status)
	}

	return C.GoString(_cName), nil
}

func (q *Queue) DispatchPoll() error {
	if status := C.tibrvQueue_TimedDispatch(q.tibrvQueue, 0); status != C.TIBRV_OK {
		return fmt.Errorf("dispatch event error, %d", status)
	}

	return nil
}

func (q *Queue) RemoveHook() error {
	if status := C.tibrvQueue_SetHook(q.tibrvQueue, nil, nil); status != C.TIBRV_OK {
		return fmt.Errorf("remove hook error, %d", status)
	}

	return nil
}

func (q *Queue) SetHook(hook TibrvQueueHook) error {
	if status := C.tibrvQueue_SetHook(q.tibrvQueue, C.tibrvQueueHook(C.goQueueHook), unsafe.Pointer(&hook)); status != C.TIBRV_OK {
		return fmt.Errorf("set hook error, %d", status)
	}

	return nil
}

func (q *Queue) SetName(name string) error {
	_cName := C.CString(name)
	defer C.free(unsafe.Pointer(_cName))

	if status := C.tibrvQueue_SetName(q.tibrvQueue, _cName); status != C.TIBRV_OK {
		return fmt.Errorf("set queue name error, %d", status)
	}

	return nil
}
