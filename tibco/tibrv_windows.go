//go:build windows
// +build windows

package tibco

// #cgo LDFLAGS: -L${SRCDIR}/lib_windows -ltibrvcm
// #cgo LDFLAGS: -L${SRCDIR}/lib_windows -ltibrv
/*
#include <stdlib.h>
#include "tibrv/cm.h"
extern void goCmCallback(tibrvcmEvent event, tibrvMsg message, void* closure);
*/
import "C"