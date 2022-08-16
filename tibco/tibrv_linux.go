//go:build linux
// +build linux

package tibco

// #cgo LDFLAGS: -L${SRCDIR}/lib_linux -ltibrvcm64
// #cgo LDFLAGS: -L${SRCDIR}/lib_linux -ltibrv64
/*
#include <stdlib.h>
#include "tibrv/cm.h"
extern void goCmCallback(tibrvcmEvent event, tibrvMsg message, void* closure);
*/
import "C"