package tibco

// #cgo LDFLAGS: -L${SRCDIR}/tibrv -ltibrvcm
/*
#include <stdlib.h>
#include "tibrv/cm.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type CmTransport struct {
	tibrvCmTransport C.tibrvcmTransport
}

// NewCmTransport，创建一个tibrvcmTransport，其中，cmName为""时系统会生成唯一字符串作为CM Transport标识。
func NewCmTransport(transport *Transport, cmName string, requestOld bool, ledgerName string, syncLedger bool, relayAgent string) (*CmTransport, error) {
	var cmTransport C.tibrvcmTransport
	_cCmName := C.CString(cmName)
	if cmName == "" {
		_cCmName = nil
	}

	_cLedgerName := C.CString(ledgerName)
	if ledgerName == "" {
		_cLedgerName = nil
	}

	_cRelayAgent := C.CString(relayAgent)
	if relayAgent == "" {
		_cRelayAgent = nil
	}

	if status := C.tibrvcmTransport_Create(&cmTransport, transport.tibrvTransport, _cCmName, TibrvBool(requestOld), _cLedgerName, TibrvBool(syncLedger), _cRelayAgent); status != C.TIBRV_OK {
		return nil, fmt.Errorf("create cmTransport error, %d", status)
	}

	return &CmTransport{
		tibrvCmTransport: cmTransport,
	}, nil
}

func (t *CmTransport) Destroy() error {
	if status := C.tibrvcmTransport_Destroy(t.tibrvCmTransport); status != C.TIBRV_OK {
		return fmt.Errorf("destroy transport error, %d", status)
	}

	return nil
}

func (t *CmTransport) Send(message *Message) error {
	if status := C.tibrvcmTransport_Send(t.tibrvCmTransport, message.tibrvMsg); status != C.TIBRV_OK {
		return fmt.Errorf("send cm message error, %d", status)
	}

	return nil
}

func (t *CmTransport) SendReply(reply, request *Message) error {
	if status := C.tibrvcmTransport_SendReply(t.tibrvCmTransport, reply.tibrvMsg, request.tibrvMsg); status != C.TIBRV_OK {
		return fmt.Errorf("send reply error, %d", status)
	}

	return nil
}

func (t *CmTransport) SendRequest(request *Message, timeout float64) (*Message, error) {
	var reply C.tibrvMsg

	if status := C.tibrvcmTransport_SendRequest(t.tibrvCmTransport, request.tibrvMsg, &reply, C.double(timeout)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("send request error, %d", status)
	}

	return &Message{
		tibrvMsg: reply,
	}, nil
}

// AddListener, 预注册一个预期中的listener。
func (t *CmTransport) AddListener(cmName, subject string) error {
	_cCmName := C.CString(cmName)
	_cSubject := C.CString(subject)

	defer C.free(unsafe.Pointer(_cCmName))
	defer C.free(unsafe.Pointer(_cSubject))

	if status := C.tibrvcmTransport_AddListener(t.tibrvCmTransport, _cCmName, _cSubject); status != C.TIBRV_OK {
		return fmt.Errorf("add listener error, %d", status)
	}

	return nil
}

func (t *CmTransport) AllowListener(cmName string) error {
	_cCmName := C.CString(cmName)
	defer C.free(unsafe.Pointer(_cCmName))

	if status := C.tibrvcmTransport_AllowListener(t.tibrvCmTransport, _cCmName); status != C.TIBRV_OK {
		return fmt.Errorf("allow listener error, %d", status)
	}

	return nil
}

func (t *CmTransport) DisallowListener(cmName string) error {
	_cCmName := C.CString(cmName)
	defer C.free(unsafe.Pointer(_cCmName))

	if status := C.tibrvcmTransport_DisallowListener(t.tibrvCmTransport, _cCmName); status != C.TIBRV_OK {
		return fmt.Errorf("disallow listener error, %d", status)
	}

	return nil
}

func (t *CmTransport) ExpireMessage(subject string, sequence uint64) error {
	_cSubject := C.CString(subject)
	defer C.free(unsafe.Pointer(_cSubject))

	if status := C.tibrvcmTransport_ExpireMessages(t.tibrvCmTransport, _cSubject, C.tibrv_u64(sequence)); status != C.TIBRV_OK {
		return fmt.Errorf("expire message error, %d", status)
	}

	return nil
}

func (t *CmTransport) GetDefaultCMTimeLimit() (float64, error) {
	var timeLimit C.tibrv_f64

	if status := C.tibrvcmTransport_GetDefaultCMTimeLimit(t.tibrvCmTransport, &timeLimit); status != C.TIBRV_OK {
		return 0, fmt.Errorf("get default cm time limit error, %d", status)
	}

	return float64(timeLimit), nil
}

func (t *CmTransport) GetLedgerName() (string, error) {
	var ledgerName *C.char

	if status := C.tibrvcmTransport_GetLedgerName(t.tibrvCmTransport, &ledgerName); status != C.TIBRV_OK {
		return "", fmt.Errorf("get ledger name error, %d", status)
	}

	return C.GoString(ledgerName), nil
}

func (t *CmTransport) GetName() (string, error) {
	var name *C.char

	if status := C.tibrvcmTransport_GetName(t.tibrvCmTransport, &name); status != C.TIBRV_OK {
		return "", fmt.Errorf("get name error, %d", status)
	}

	return C.GoString(name), nil
}

func (t *CmTransport) GetRelayAgent() (string, error) {
	var relayAgent *C.char

	if status := C.tibrvcmTransport_GetRelayAgent(t.tibrvCmTransport, &relayAgent); status != C.TIBRV_OK {
		return "", fmt.Errorf("get relay agent name error, %d", status)
	}

	return C.GoString(relayAgent), nil
}

func (t *CmTransport) GetRequestOld() (bool, error) {
	var requestOld C.tibrv_bool

	if status := C.tibrvcmTransport_GetRequestOld(t.tibrvCmTransport, &requestOld); status != C.TIBRV_OK {
		return false, fmt.Errorf("get request old error, %d", status)
	}

	if requestOld == C.TIBRV_TRUE {
		return true, nil
	}

	return false, nil
}

func (t *CmTransport) GetSyncLedger() (bool, error) {
	var syncLedger C.tibrv_bool

	if status := C.tibrvcmTransport_GetSyncLedger(t.tibrvCmTransport, &syncLedger); status != C.TIBRV_OK {
		return false, fmt.Errorf("get sync ledger error, %d", status)
	}

	if syncLedger == C.TIBRV_TRUE {
		return true, nil
	}

	return false, nil
}

func (t *CmTransport) GetTransport() (*Transport, error) {
	var _t C.tibrvTransport

	if status := C.tibrvcmTransport_GetTransport(t.tibrvCmTransport, &_t); status != C.TIBRV_OK {
		return nil, fmt.Errorf("get transport error, %d", status)
	}

	return &Transport{
		tibrvTransport: _t,
	}, nil
}

func (t *CmTransport) RemoveListener(cmName, subject string) error {
	_cCmName := C.CString(cmName)
	_cSubject := C.CString(subject)

	defer C.free(unsafe.Pointer(_cCmName))
	defer C.free(unsafe.Pointer(_cSubject))

	if status := C.tibrvcmTransport_RemoveListener(t.tibrvCmTransport, _cCmName, _cSubject); status != C.TIBRV_OK {
		return fmt.Errorf("remove listener error, %d", status)
	}

	return nil
}

func (t *CmTransport) RemoveSendState(subject string) error {
	_cSubject := C.CString(subject)
	defer C.free(unsafe.Pointer(_cSubject))

	if status := C.tibrvcmTransport_RemoveSendState(t.tibrvCmTransport, _cSubject); status != C.TIBRV_OK {
		return fmt.Errorf("remove send state error, %d", status)
	}

	return nil
}

// todo
func (t *CmTransport) ReviewLedger() error {
	panic("not implemented")
}

func (t *CmTransport) SetDefaultCMTimeLimit(limit float64) error {
	if status := C.tibrvcmTransport_SetDefaultCMTimeLimit(t.tibrvCmTransport, C.double(limit)); status != C.TIBRV_OK {
		return fmt.Errorf("set default cm time limit error, %d", status)
	}

	return nil
}

func (t *CmTransport) SetPublisherInactivityDiscardInterval(interval int32) error {
	if status := C.tibrvcmTransport_SetPublisherInactivityDiscardInterval(t.tibrvCmTransport, C.tibrv_i32(interval)); status != C.TIBRV_OK {
		return fmt.Errorf("set publiser inactivity discard interval error, %d", status)
	}

	return nil
}

func (t *CmTransport) SyncLedger() error {
	if status := C.tibrvcmTransport_SyncLedger(t.tibrvCmTransport); status != C.TIBRV_OK {
		return fmt.Errorf("sync ledger error, %d", status)
	}

	return nil
}
