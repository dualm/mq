package tibco

/*
#include <stdlib.h>
#include "tibrv/cm.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type MsgField struct {
	field C.tibrvMsgField
}

type Message struct {
	tibrvMsg C.tibrvMsg
}

func NewMessage() (*Message, error) {
	var _msg C.tibrvMsg

	if status := C.tibrvMsg_Create(&_msg); status != C.TIBRV_OK {
		return nil, fmt.Errorf("Create message error, code: %d", status)
	}

	return &Message{
		tibrvMsg: _msg,
	}, nil
}

func (msg *Message) CopyCreate(m Message) error {
	var _newMsg C.tibrvMsg
	if status := C.tibrvMsg_CreateCopy(m.tibrvMsg, &_newMsg); status != C.TIBRV_OK {
		return fmt.Errorf("CreateCopy message error, code: %d", status)
	}

	msg.tibrvMsg = _newMsg

	return nil
}

func (msg *Message) Close() error {
	if status := C.tibrvMsg_Destroy(msg.tibrvMsg); status != C.TIBRV_OK {
		return fmt.Errorf("Destroy message error, code: %d", status)
	}

	msg.tibrvMsg = nil

	return nil
}

func (msg *Message) Detach() error {
	if status := C.tibrvMsg_Detach(msg.tibrvMsg); status != C.TIBRV_OK {
		return fmt.Errorf("Detach message error, code: %d", status)
	}

	return nil
}

func (msg *Message) Expand(additionalStorage uint32) error {
	if status := C.tibrvMsg_Expand(msg.tibrvMsg, C.int(additionalStorage)); status != C.TIBRV_OK {
		return fmt.Errorf("Expand message error, code: %d", status)
	}

	return nil
}

func (msg *Message) GetAsByte() ([]byte, error) {
	b := C.CBytes(make([]byte, 0))
	if status := C.tibrvMsg_GetAsBytes(msg.tibrvMsg, &b); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetAsByte error, code: %d", status)
	}

	return *(*[]byte)(b), nil
}

func (msg *Message) GetAsBytesCopy(n uint32) ([]byte, error) {
	b := C.CBytes(make([]byte, 0))
	if status := C.tibrvMsg_GetAsBytesCopy(msg.tibrvMsg, b, C.uint(n)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetAsBytesCopy error, code: %d", status)
	}

	return *(*[]byte)(b), nil
}

func (msg Message) GetByteSize() (uint32, error) {
	var n C.uint
	if status := C.tibrvMsg_GetByteSize(msg.tibrvMsg, &n); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetByteSize error, code: %d", status)
	}

	return uint32(n), nil
}

func (msg Message) GetClosure() ([]byte, error) {
	b := C.CBytes(make([]byte, 0))
	if status := C.tibrvMsg_GetClosure(msg.tibrvMsg, &b); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetClosure error, code: %d", status)
	}

	return *(*[]byte)(b), nil
}

func (msg Message) GetEvent() (*Event, error) {
	var _event C.tibrvEvent
	if status := C.tibrvMsg_GetEvent(msg.tibrvMsg, &_event); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetEvent error, code: %d", status)
	}

	return &Event{
		tibrvEvent: _event,
	}, nil
}

func (msg Message) GetField(fieldName string, fieldId uint16) (*MsgField, error) {
	var _field C.tibrvMsgField

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetFieldEx(msg.tibrvMsg, _cFieldName, &_field, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetField error, code: %d", status)
	}

	return &MsgField{
			field: _field,
		},
		nil
}

func (msg Message) GetBool(fieldName string, fieldId uint16) (bool, error) {
	var b C.tibrv_bool

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetBoolEx(msg.tibrvMsg, _cFieldName, &b, C.ushort(fieldId)); status != C.TIBRV_OK {
		return false, fmt.Errorf("GetBool error, code: %d", status)
	}

	if b == C.TIBRV_FALSE {
		return false, nil
	}

	return true, nil
}

func (msg Message) GetF32(fieldName string, fieldId uint16) (float32, error) {
	var f C.tibrv_f32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetF32Ex(msg.tibrvMsg, _cFieldName, &f, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetF32 error, code: %d", status)
	}

	return float32(f), nil
}

func (msg Message) GetF64(fieldName string, fieldId uint16) (float64, error) {
	var f C.tibrv_f64

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetF64Ex(msg.tibrvMsg, _cFieldName, &f, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetF64 error, code: %d", status)
	}

	return float64(f), nil
}

func (msg Message) GetI8(fieldName string, fieldId uint16) (int8, error) {
	var i C.tibrv_i8

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetI8Ex(msg.tibrvMsg, _cFieldName, &i, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetI8 error, code: %d", status)
	}

	return int8(i), nil
}

func (msg Message) GetI16(fieldName string, fieldId uint16) (int16, error) {
	var i C.tibrv_i16

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetI16Ex(msg.tibrvMsg, _cFieldName, &i, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetI16 error, code: %d", status)
	}

	return int16(i), nil
}

func (msg Message) GetI32(fieldName string, fieldId uint16) (int32, error) {
	var i C.tibrv_i32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetI32Ex(msg.tibrvMsg, _cFieldName, &i, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetI32 error, code: %d", status)
	}

	return int32(i), nil
}

func (msg Message) GetI64(fieldName string, fieldId uint16) (int64, error) {
	var i C.tibrv_i64

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetI64Ex(msg.tibrvMsg, _cFieldName, &i, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetI64 error, code: %d", status)
	}

	return int64(i), nil
}

func (msg Message) GetU8(fieldName string, fieldId uint16) (uint8, error) {
	var i C.tibrv_u8

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetU8Ex(msg.tibrvMsg, _cFieldName, &i, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetU8 error, code: %d", status)
	}

	return uint8(i), nil
}

func (msg Message) GetU16(fieldName string, fieldId uint16) (uint16, error) {
	var i C.tibrv_u16

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetU16Ex(msg.tibrvMsg, _cFieldName, &i, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetU16 error, code: %d", status)
	}

	return uint16(i), nil
}

func (msg Message) GetU32(fieldName string, fieldId uint16) (uint32, error) {
	var i C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetU32Ex(msg.tibrvMsg, _cFieldName, &i, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetU32 error, code: %d", status)
	}

	return uint32(i), nil
}

func (msg Message) GetU64(fieldName string, fieldId uint16) (uint64, error) {
	var i C.tibrv_u64

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetU64Ex(msg.tibrvMsg, _cFieldName, &i, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetU64 error, code: %d", status)
	}

	return uint64(i), nil
}

func (msg Message) GetIPAddr32(fieldName string, fieldId uint16) (uint32, error) {
	var ip C.tibrv_ipaddr32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetIPAddr32Ex(msg.tibrvMsg, _cFieldName, &ip, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetIPAddr32 error, code: %d", status)
	}

	return uint32(ip), nil
}

func (msg Message) GetIPPort16(fieldName string, fieldId uint16) (uint16, error) {
	var port C.tibrv_ipport16

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetIPPort16Ex(msg.tibrvMsg, _cFieldName, &port, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetIPPort16 error, code: %d", status)
	}

	return uint16(port), nil
}

func (msg Message) GetF32Array(fieldName string, fieldId uint16) ([]float32, error) {
	f := make([]float32, 0)
	_p := (*C.float)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetF32ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetF64Array(fieldName string, fieldId uint16) ([]float64, error) {
	f := make([]float64, 0)
	_p := (*C.double)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetF64ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetI8Array(fieldName string, fieldId uint16) ([]int8, error) {
	f := make([]int8, 0)
	_p := (*C.schar)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetI8ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetI16Array(fieldName string, fieldId uint16) ([]int16, error) {
	f := make([]int16, 0)
	_p := (*C.short)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetI16ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetI32Array(fieldName string, fieldId uint16) ([]int32, error) {
	f := make([]int32, 0)
	_p := (*C.int)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetI32ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetI64Array(fieldName string, fieldId uint16) ([]int64, error) {
	f := make([]int64, 0)
	_p := (*C.longlong)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetI64ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetU8Array(fieldName string, fieldId uint16) ([]uint8, error) {
	f := make([]uint8, 0)
	_p := (*C.uchar)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetU8ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetU16Array(fieldName string, fieldId uint16) ([]uint16, error) {
	f := make([]uint16, 0)
	_p := (*C.ushort)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetU16ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetU32Array(fieldName string, fieldId uint16) ([]uint32, error) {
	f := make([]uint32, 0)
	_p := (*C.uint)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetU32ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetU64Array(fieldName string, fieldId uint16) ([]uint64, error) {
	f := make([]uint64, 0)
	_p := (*C.ulonglong)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetU64ArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetMsgArray(fieldName string, fieldId uint16) ([]C.tibrvMsg, error) {
	f := make([]C.tibrvMsg, 0)
	_p := (*C.tibrvMsg)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetMsgArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

func (msg Message) GetStringArray(fieldName string, fieldId uint16) ([]string, error) {
	f := make([]string, 0)
	_p := (**C.char)(unsafe.Pointer(&f[0]))
	var n C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetStringArrayEx(msg.tibrvMsg, _cFieldName, &_p, &n, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf(" error, code: %d", status)
	}

	return f, nil
}

// GetString get the value of a field as character string.
func (msg Message) GetString(fieldName string, fieldId uint16) (string, error) {
	var _s *C.char

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetStringEx(msg.tibrvMsg, _cFieldName, &_s, C.ushort(fieldId)); status != C.TIBRV_OK {
		return "", fmt.Errorf("GetString error, code: %d", status)
	}

	return C.GoString(_s), nil
}

func (msg Message) GetMsg(fieldName string, fieldId uint16) (*Message, error) {
	var _msg C.tibrvMsg

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetMsgEx(msg.tibrvMsg, _cFieldName, &_msg, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetMsg error, code: %d", status)
	}

	return &Message{
		tibrvMsg: _msg,
	}, nil
}

// GetOpaque get the value of a field as a opaque byte sequence
func (msg Message) GetOpaque(fieldName string, fieldId uint16) ([]byte, error) {
	b := C.CBytes(make([]byte, 0))
	var l C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetOpaqueEx(msg.tibrvMsg, _cFieldName, &b, &l, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetClosure error, code: %d", status)
	}

	return *(*[]byte)(b), nil
}

// GetXml get the value of a field as an XML byte sequence
func (msg Message) GetXml(fieldName string, fieldId uint16) ([]byte, error) {
	b := C.CBytes(make([]byte, 0))
	var l C.tibrv_u32

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetXmlEx(msg.tibrvMsg, _cFieldName, &b, &l, C.ushort(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetXml error, code: %d", status)
	}

	return *(*[]byte)(b), nil
}

// GetDataTime gets the value of a field as a Rendezvous date time value.
func (msg Message) GetDateTime(fieldName string, fieldId uint16) (int64, error) {
	var t C.tibrvMsgDateTime

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetDateTimeEx(msg.tibrvMsg, _cFieldName, &t, C.ushort(fieldId)); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetDateTime error, code: %d", status)
	}

	return int64(t.sec), nil
}

// GetFieldByIndex get a field from a message by an index.
func (msg Message) GetFieldByIndex(fieldId uint16) (*MsgField, error) {
	var _field C.tibrvMsgField

	if status := C.tibrvMsg_GetFieldByIndex(msg.tibrvMsg, &_field, C.uint(fieldId)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetFieldByIndex error, code: %d", status)
	}

	return &MsgField{
		field: _field,
	}, nil
}

// GetFieldInstance get a specified instance of a field from a message.
// instances with same field names.
func (msg Message) GetFieldInstance(fieldName string, instance uint32) (*MsgField, error) {
	var _field C.tibrvMsgField

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_GetFieldInstance(msg.tibrvMsg, _cFieldName, &_field, C.uint(instance)); status != C.TIBRV_OK {
		return nil, fmt.Errorf("GetFieldInstance error, code: %d", status)
	}

	return &MsgField{
		field: _field,
	}, nil
}

// GetNumFields extract the number of fields in a message
func (msg Message) GetNumFields() (uint32, error) {
	var n C.tibrv_u32

	if status := C.tibrvMsg_GetNumFields(msg.tibrvMsg, &n); status != C.TIBRV_OK {
		return 0, fmt.Errorf("GetNumFields error, code: %d", status)
	}

	return uint32(n), nil
}

// GetReplySubject extract the subject from a message
func (msg Message) GetReplySubject() (string, error) {
	var s *C.char
	if status := C.tibrvMsg_GetReplySubject(msg.tibrvMsg, &s); status != C.TIBRV_OK {
		return "", fmt.Errorf("GetReplySubject error, code: %d", status)
	}

	return C.GoString(s), nil
}

// GetSendSubject extract the subject from a message
func (msg Message) GetSendSubject() (string, error) {
	var s *C.char

	if status := C.tibrvMsg_GetSendSubject(msg.tibrvMsg, &s); status != C.TIBRV_OK {
		return "", fmt.Errorf("GetSendSubject error, code: %d", status)
	}

	return C.GoString(s), nil
}

// MarkReferences mark and clear references
func (msg *Message) MarkReferences() error {
	if status := C.tibrvMsg_MarkReferences(msg.tibrvMsg); status != C.TIBRV_OK {
		return fmt.Errorf("MarkReferences error, code: %d", status)
	}

	return nil
}

// RemoveField remove a field from a message.
func (msg *Message) RemoveField(fieldName string, fieldId uint16) error {
	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_RemoveFieldEx(msg.tibrvMsg, _cFieldName, C.ushort(fieldId)); status != C.TIBRV_OK {
		return fmt.Errorf("RemoveField error, code: %d", status)
	}

	return nil
}

// RemoveFieldInstance remove a field instance of a field from a message.
func (msg *Message) RemoveFieldInstance(fieldName string, instance uint32) error {
	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_RemoveFieldInstance(msg.tibrvMsg, _cFieldName, C.uint(instance)); status != C.TIBRV_OK {
		return fmt.Errorf("RemoveFieldInstance error, code: %d", status)
	}

	return nil
}

// Reset clear a message, preparing it for re-use
func (msg *Message) Reset() error {
	if status := C.tibrvMsg_Reset(msg.tibrvMsg); status != C.TIBRV_OK {
		return fmt.Errorf("Reset error, code: %d", status)
	}

	return nil
}

// SetReplySubject set the reply subject for a message.
func (msg *Message) SetReplySubject(replySubject string) error {
	_cReplySubject := C.CString(replySubject)
	defer C.free(unsafe.Pointer(_cReplySubject))

	if status := C.tibrvMsg_SetReplySubject(msg.tibrvMsg, _cReplySubject); status != C.TIBRV_OK {
		return fmt.Errorf("SetReplySubject error, code: %d", status)
	}

	return nil
}

// SetReplySubject set the reply subject for a message.
func (msg *Message) SetSendSubject(sendSubject string) error {
	_cSendSubject := C.CString(sendSubject)
	defer C.free(unsafe.Pointer(_cSendSubject))

	if status := C.tibrvMsg_SetSendSubject(msg.tibrvMsg, _cSendSubject); status != C.TIBRV_OK {
		return fmt.Errorf("SetReplySubject error, code: %d", status)
	}

	return nil
}

// UpdateField update a field within a message
func (msg *Message) UpdateField(field *MsgField) error {
	if status := C.tibrvMsg_UpdateField(msg.tibrvMsg, &field.field); status != C.TIBRV_OK {
		return fmt.Errorf("UpdateField error, code: %d", status)
	}

	return nil
}

// UpdateBool
func (msg *Message) UpdateBool(fieldName string, b bool, fieldId uint16) error {
	var _b C.tibrv_bool

	if b {
		_b = C.TIBRV_TRUE
	} else {
		_b = C.TIBRV_FALSE
	}

	_cFieldName := C.CString(fieldName)
	defer C.free(unsafe.Pointer(_cFieldName))

	if status := C.tibrvMsg_UpdateBoolEx(msg.tibrvMsg, _cFieldName, _b, C.ushort(fieldId)); status != C.TIBRV_OK {
		return fmt.Errorf("UpdateBool error, code: %d", status)
	}

	return nil
}

// AddString add a field containing a string
func (msg *Message) AddString(fieldName, value string, fieldId uint16) error {
	_cFieldName := C.CString(fieldName)
	_cValue := C.CString(value)

	defer C.free(unsafe.Pointer(_cFieldName))
	defer C.free(unsafe.Pointer(_cValue))

	if status := C.tibrvMsg_AddStringEx(msg.tibrvMsg, _cFieldName, _cValue, C.ushort(fieldId)); status != C.TIBRV_OK {
		return fmt.Errorf("AddString error, code : %d", status)
	}

	return nil
}

func (msg *Message) GetCMSender() (string, error) {
	var name *C.char

	if status := C.tibrvMsg_GetCMSender(msg.tibrvMsg, &name); status != C.TIBRV_OK {
		return "", fmt.Errorf("get cm sender error, %d", status)
	}

	return C.GoString(name), nil
}

func (msg *Message) GetCMSequence() (uint64, error) {
	var n C.tibrv_u64

	if status := C.tibrvMsg_GetCMSequence(msg.tibrvMsg, &n); status != C.TIBRV_OK {
		return 0, fmt.Errorf("get cm sequence error, %d", status)
	}

	return uint64(n), nil
}

func (msg *Message) GetCMTimeLimit() (float64, error) {
	var limit C.double

	if status := C.tibrvMsg_GetCMTimeLimit(msg.tibrvMsg, &limit); status != C.TIBRV_OK {
		return 0, fmt.Errorf("get cm time limit error, %d", status)
	}

	return float64(limit), nil
}

func (msg *Message) SetCMTimeLimit(limit float64) error {
	if status := C.tibrvMsg_SetCMTimeLimit(msg.tibrvMsg, C.double(limit)); status != C.TIBRV_OK {
		return fmt.Errorf("set cm time limit error, %d", status)
	}

	return nil
}

func TibrvBool(b bool) C.tibrv_bool {
	if b {
		return C.TIBRV_TRUE
	}

	return C.TIBRV_FALSE
}
