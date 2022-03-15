package tibco

/*
#include "tibrv/tibrv.h"
*/
import "C"

type event struct {
	tibrvEvent C.tibrvEvent
}

