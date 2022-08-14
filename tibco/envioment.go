package tibco

/*
#include "tibrv/tibrv.h"
*/
import "C"
import (
	"fmt"
)

func tibrvOpen() error {
	if status := C.tibrv_Open(); status != C.TIBRV_OK {
		return fmt.Errorf("Tibco Open Error, code: %d", status)
	}

	return nil
}

func tibrvClose() error {
	if status := C.tibrv_Close(); status != C.TIBRV_OK {
		return fmt.Errorf("Tibco Close Error, code: %d", status)
	}

	return nil
}
