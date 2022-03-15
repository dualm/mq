package tibco

/*
#include "tibrv/tibrv.h"
*/
import "C"
import "fmt"

type queue struct {
	tibrvQueue C.tibrvQueue
}

func tibCreateQueue() (*queue, error) {
	var q C.tibrvQueue
	if status := C.tibrvQueue_Create(&q); status != C.TIBRV_OK {
		return nil, fmt.Errorf("Create queue error: %d", status)
	}
	return &queue{
		tibrvQueue: q,
	}, nil
}

func DestroyQueue(queue queue) {
	C.tibrvQueue_DestroyEx(queue.tibrvQueue, nil, nil)
}
