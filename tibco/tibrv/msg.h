/*
 * Copyright (c) 1998-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: msg.h 42697 2009-11-10 17:56:31Z mwoehlke $
 */

#ifndef _INCLUDED_tibrvmsg_h
#define _INCLUDED_tibrvmsg_h

#include "types.h"
#include "status.h"
#include "handler.h"
#include "msgx.h"

#if defined(__cplusplus)
extern "C" {
#endif

/* Creation, Reuse and Referencing */


extern tibrv_status
tibrvMsg_Create(
    tibrvMsg*           message);

extern tibrv_status
tibrvMsg_CreateEx(
    tibrvMsg*           message,
    tibrv_u32           initialStorage);

extern tibrv_status
tibrvMsg_Destroy(
    tibrvMsg            message);

extern tibrv_status
tibrvMsg_Detach(
    tibrvMsg            message);

extern tibrv_status
tibrvMsg_Reset(
    tibrvMsg            message);

extern tibrv_status
tibrvMsg_Expand(
    tibrvMsg            message,
    tibrv_i32           additionalStorage);

/* Publishing messages */
extern tibrv_status
tibrvMsg_SetSendSubject(
    tibrvMsg            message,
    const char*         subject);

extern tibrv_status
tibrvMsg_GetSendSubject(
    tibrvMsg            message,
    const char**        subject);

extern tibrv_status
tibrvMsg_SetReplySubject(
    tibrvMsg            message,
    const char*         replySubject);

extern tibrv_status
tibrvMsg_GetReplySubject(
    tibrvMsg            message,
    const char**        replySubject);

extern tibrv_status
tibrvMsg_GetEvent(
    tibrvMsg            message,
    tibrvEvent*         eventId);

extern tibrv_status
tibrvMsg_GetClosure(
    tibrvMsg            message,
    void**               closure);
/* Attributes */

extern tibrv_status
tibrvMsg_GetNumFields(
    tibrvMsg            message,
    tibrv_u32*          numFields);

extern tibrv_status
tibrvMsg_GetByteSize(
    tibrvMsg            message,
    tibrv_u32*          byteSize);

extern tibrv_status
tibrvMsg_ConvertToString(
    tibrvMsg            message,
    const char**        string);

/* Fields */

/* Adding fields */

#define tibrvMsg_AddMsg(msg, fieldName, value) \
        tibrvMsg_AddMsgEx(msg, fieldName, value, 0)

#define tibrvMsg_AddIPAddr32(msg, fieldName, value) \
        tibrvMsg_AddIPAddr32Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddIPPort16(msg, fieldName, value) \
        tibrvMsg_AddIPPort16Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddDateTime(msg, fieldName, value) \
        tibrvMsg_AddDateTimeEx(msg, fieldName, value, 0)

#define tibrvMsg_AddBool(msg, fieldName, value) \
        tibrvMsg_AddBoolEx(msg, fieldName, value, 0)

#define tibrvMsg_AddI8(msg, fieldName, value) \
        tibrvMsg_AddI8Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddU8(msg, fieldName, value) \
        tibrvMsg_AddU8Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddI16(msg, fieldName, value) \
        tibrvMsg_AddI16Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddU16(msg, fieldName, value) \
        tibrvMsg_AddU16Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddI32(msg, fieldName, value) \
        tibrvMsg_AddI32Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddU32(msg, fieldName, value) \
        tibrvMsg_AddU32Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddI64(msg, fieldName, value) \
        tibrvMsg_AddI64Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddU64(msg, fieldName, value) \
        tibrvMsg_AddU64Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddF32(msg, fieldName, value) \
        tibrvMsg_AddF32Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddF64(msg, fieldName, value) \
        tibrvMsg_AddF64Ex(msg, fieldName, value, 0)

#define tibrvMsg_AddString(msg, fieldName, value) \
        tibrvMsg_AddStringEx(msg, fieldName, value, 0)

#define tibrvMsg_AddOpaque(msg, fieldName, value, length) \
        tibrvMsg_AddOpaqueEx(msg, fieldName, value, length, 0)

#define tibrvMsg_AddXml(msg, fieldName, value, length) \
        tibrvMsg_AddXmlEx(msg, fieldName, value, length, 0)

#define tibrvMsg_AddI8Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddI8ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddU8Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddU8ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddI16Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddI16ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddU16Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddU16ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddI32Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddI32ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddU32Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddU32ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddI64Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddI64ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddU64Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddU64ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddF32Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddF32ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddF64Array(msg, fieldName, array, numElements) \
        tibrvMsg_AddF64ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddStringArray(msg, fieldName, array, numElements) \
        tibrvMsg_AddStringArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_AddMsgArray(msg, fieldName, array, numElements) \
        tibrvMsg_AddMsgArrayEx(msg, fieldName, array, numElements, 0)

extern tibrv_status
tibrvMsg_AddField(
    tibrvMsg            message,
    tibrvMsgField*      field);

/* Finding fields */

#define tibrvMsg_GetField(msg, fieldName, field) \
        tibrvMsg_GetFieldEx(msg, fieldName, field, 0)

#define tibrvMsg_GetMsg(msg, fieldName, value) \
        tibrvMsg_GetMsgEx(msg, fieldName, value, 0)

#define tibrvMsg_GetIPAddr32(msg, fieldName, value) \
        tibrvMsg_GetIPAddr32Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetIPPort16(msg, fieldName, value) \
        tibrvMsg_GetIPPort16Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetDateTime(msg, fieldName, value) \
        tibrvMsg_GetDateTimeEx(msg, fieldName, value, 0)

#define tibrvMsg_GetBool(msg, fieldName, value) \
        tibrvMsg_GetBoolEx(msg, fieldName, value, 0)

#define tibrvMsg_GetI8(msg, fieldName, value) \
        tibrvMsg_GetI8Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetU8(msg, fieldName, value) \
        tibrvMsg_GetU8Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetI16(msg, fieldName, value) \
        tibrvMsg_GetI16Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetU16(msg, fieldName, value) \
        tibrvMsg_GetU16Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetI32(msg, fieldName, value) \
        tibrvMsg_GetI32Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetU32(msg, fieldName, value) \
        tibrvMsg_GetU32Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetI64(msg, fieldName, value) \
        tibrvMsg_GetI64Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetU64(msg, fieldName, value) \
        tibrvMsg_GetU64Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetF32(msg, fieldName, value) \
        tibrvMsg_GetF32Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetF64(msg, fieldName, value) \
        tibrvMsg_GetF64Ex(msg, fieldName, value, 0)

#define tibrvMsg_GetString(msg, fieldName, value) \
        tibrvMsg_GetStringEx(msg, fieldName, value, 0)

#define tibrvMsg_GetOpaque(msg, fieldName, value, length) \
        tibrvMsg_GetOpaqueEx(msg, fieldName, value, length, 0)

#define tibrvMsg_GetXml(msg, fieldName, value, length) \
        tibrvMsg_GetXmlEx(msg, fieldName, value, length, 0)

#define tibrvMsg_GetI8Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetI8ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetU8Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetU8ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetI16Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetI16ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetU16Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetU16ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetI32Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetI32ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetU32Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetU32ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetI64Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetI64ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetU64Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetU64ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetF32Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetF32ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetF64Array(msg, fieldName, array, numElements) \
        tibrvMsg_GetF64ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetStringArray(msg, fieldName, array, numElements) \
        tibrvMsg_GetStringArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_GetMsgArray(msg, fieldName, array, numElements) \
        tibrvMsg_GetMsgArrayEx(msg, fieldName, array, numElements, 0)

extern tibrv_status
tibrvMsg_GetFieldInstance(
    tibrvMsg            message,
    const char*         fieldName,
    tibrvMsgField*      fieldAddr,
    tibrv_u32           instance);

/* Getting fields by field number */

extern tibrv_status
tibrvMsg_GetFieldByIndex(
    tibrvMsg            message,
    tibrvMsgField*      field,
    tibrv_u32           fieldIndex);

/* Removing fields */

#define \
tibrvMsg_RemoveField(msg, fieldName) \
        tibrvMsg_RemoveFieldEx(msg, fieldName, 0)

extern tibrv_status
tibrvMsg_RemoveFieldInstance(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u32           instance);

/* Updating fields */

#define tibrvMsg_UpdateMsg(msg, fieldName, value) \
        tibrvMsg_UpdateMsgEx(msg, fieldName, value, 0)

#define tibrvMsg_UpdateIPAddr32(msg, fieldName, value) \
        tibrvMsg_UpdateIPAddr32Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateIPPort16(msg, fieldName, value) \
        tibrvMsg_UpdateIPPort16Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateDateTime(msg, fieldName, value) \
        tibrvMsg_UpdateDateTimeEx(msg, fieldName, value, 0)

#define tibrvMsg_UpdateBool(msg, fieldName, value) \
        tibrvMsg_UpdateBoolEx(msg, fieldName, value, 0)

#define tibrvMsg_UpdateI8(msg, fieldName, value) \
        tibrvMsg_UpdateI8Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateU8(msg, fieldName, value) \
        tibrvMsg_UpdateU8Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateI16(msg, fieldName, value) \
        tibrvMsg_UpdateI16Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateU16(msg, fieldName, value) \
        tibrvMsg_UpdateU16Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateI32(msg, fieldName, value) \
        tibrvMsg_UpdateI32Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateU32(msg, fieldName, value) \
        tibrvMsg_UpdateU32Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateI64(msg, fieldName, value) \
        tibrvMsg_UpdateI64Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateU64(msg, fieldName, value) \
        tibrvMsg_UpdateU64Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateF32(msg, fieldName, value) \
        tibrvMsg_UpdateF32Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateF64(msg, fieldName, value) \
        tibrvMsg_UpdateF64Ex(msg, fieldName, value, 0)

#define tibrvMsg_UpdateString(msg, fieldName, value) \
        tibrvMsg_UpdateStringEx(msg, fieldName, value, 0)

#define tibrvMsg_UpdateOpaque(msg, fieldName, value, length) \
        tibrvMsg_UpdateOpaqueEx(msg, fieldName, value, length, 0)

#define tibrvMsg_UpdateXml(msg, fieldName, value, length) \
        tibrvMsg_UpdateXmlEx(msg, fieldName, value, length, 0)

#define tibrvMsg_UpdateI8Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateI8ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateU8Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateU8ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateI16Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateI16ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateU16Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateU16ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateI32Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateI32ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateU32Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateU32ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateI64Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateI64ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateU64Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateU64ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateF32Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateF32ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateF64Array(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateF64ArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateStringArray(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateStringArrayEx(msg, fieldName, array, numElements, 0)

#define tibrvMsg_UpdateMsgArray(msg, fieldName, array, numElements) \
        tibrvMsg_UpdateMsgArrayEx(msg, fieldName, array, numElements, 0)

extern tibrv_status
tibrvMsg_UpdateField(
    tibrvMsg                    message,
    tibrvMsgField*              field);

/* Archival */

extern tibrv_status
tibrvMsg_CreateFromBytes(
    tibrvMsg*                   message,
    const void*                 bytes);

extern tibrv_status
tibrvMsg_GetAsBytes(
    tibrvMsg                    message,
    const void**                bytePtr);

extern tibrv_status
tibrvMsg_GetAsBytesCopy(
    tibrvMsg                    message,
    void*                       bytePtr,
    tibrv_u32                   byteSize);

/* Copying */
extern tibrv_status
tibrvMsg_CreateCopy(
    const tibrvMsg              message,
    tibrvMsg*                   copy);

/* Creation with data */
extern tibrv_status
tibrvMsg_Write(
    tibrvMsg                   *msgRef,
    const char*                 sendSubject,
    const char*                 replySubject,
    tibrvMsgField*              fields,
    tibrv_u32                   numFields);

extern tibrv_status
tibrvMsg_WriteOpaque(
    tibrvMsg                   *msgRef,
    const char*                 sendSubject,
    const char*                 replySubject,
    const char*                 data,
    tibrv_u32                   dataLength);

/* Frame references */
extern tibrv_status
tibrvMsg_MarkReferences(
    tibrvMsg                    message);

extern tibrv_status
tibrvMsg_ClearReferences(
    tibrvMsg                    message);

/* Current time */
extern tibrv_status
tibrvMsg_GetCurrentTime(
    tibrvMsgDateTime *          dateTime);

extern tibrv_status
tibrvMsg_GetCurrentTimeString(char	*local,
			      char	*gmt);

extern tibrv_status
_tibrvMsg_SetCommId(
    tibrvMsg                    message,
    tibrv_u16                   daemonId);

extern tibrv_status
_tibrvMsg_GetCommId(
    tibrvMsg                    message,
    tibrv_u16*                  daemonId);

#ifdef  __cplusplus
}
#endif

#endif /* _INCLUDED_tibrvmsg_h */
