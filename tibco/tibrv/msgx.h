/*
 * Copyright (c) 1998-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: msgx.h 23502 2006-10-25 18:43:15Z ltherond $
 */

#ifndef _INCLUDED_tibrvmsgx_h
#define _INCLUDED_tibrvmsgx_h

#include "types.h"
#include "status.h"

#if defined(__cplusplus)
extern "C" {
#endif

extern tibrv_status
tibrvMsg_AddMsgEx(
    tibrvMsg            message,
    const char*         fieldName,
    tibrvMsg            value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddIPAddr32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_ipaddr32      value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddIPPort16Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_ipport16      value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddDateTimeEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrvMsgDateTime * value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddBoolEx(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_bool          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddI8Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i8            value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddI8ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i8*     array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddU8Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u8            value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddU8ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u8*     array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddI16Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i16           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddI16ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i16*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddU16Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u16           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddU16ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u16*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddI32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i32           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddI32ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i32*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddU32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u32           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddU32ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u32*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddI64Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i64           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddI64ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i64*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddU64Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u64           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddU64ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u64*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddF32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_f32           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddF32ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_f32*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddF64Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_f64           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddF64ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_f64*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddStringArrayEx(
    tibrvMsg            message,
    const char*         field_name,
    const char**        value,
    tibrv_u32           num_elements,
    tibrv_u16           opt_identifier);

extern tibrv_status
tibrvMsg_AddMsgArrayEx(
    tibrvMsg            message,
    const char*         field_name,
    const tibrvMsg*     value,
    tibrv_u32           num_elements,
    tibrv_u16           opt_identifier);

extern tibrv_status
tibrvMsg_AddStringEx(
    tibrvMsg            message,
    const char*         fieldName,
    const char*         value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddOpaqueEx(
    tibrvMsg            message,
    const char*         fieldName,
    const void*         value,
    tibrv_u32           size,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_AddXmlEx(
    tibrvMsg            message,
    const char*         fieldName,
    const void*         value,
    tibrv_u32           size,
    tibrv_u16           optIdentifier);

/* Finding fields */

extern tibrv_status
tibrvMsg_GetFieldEx(
    tibrvMsg            message,
    const char*         fieldName,
    tibrvMsgField*      field,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetMsgEx(
    tibrvMsg            message,
    const char*         fieldName,
    tibrvMsg*           subMessage,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetIPAddr32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_ipaddr32*     value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetIPPort16Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_ipport16*     value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetDateTimeEx(
    tibrvMsg            message,
    const char*         fieldName,
    tibrvMsgDateTime*   value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetBoolEx(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_bool*         value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetI8Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i8*           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetI8ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i8**    array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetU8Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u8*           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetU8ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u8**    array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetI16Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i16*          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetI16ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i16**   array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetU16Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u16*          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetU16ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u16**   array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetI32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i32*          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetI32ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i32**   array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetU32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u32*          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetU32ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u32**   array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetI64Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i64*          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetI64ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i64**   array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetU64Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u64*          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetU64ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u64**   array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetF32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_f32*          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetF32ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_f32**   array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetF64Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_f64*          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetF64ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_f64**   array,
    tibrv_u32*          numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetStringArrayEx(
    tibrvMsg            message,
    const char*         field_name,
    const char***       ptr_addr,
    tibrv_u32*          num_elements_addr,
    tibrv_u16           opt_identifier);

extern tibrv_status
tibrvMsg_GetMsgArrayEx(
    tibrvMsg            message,
    const char*         field_name,
    const tibrvMsg**    ptr_addr,
    tibrv_u32*          num_elements_addr,
    tibrv_u16           opt_identifier);

extern tibrv_status
tibrvMsg_GetStringEx(
    tibrvMsg            message,
    const char*         fieldName,
    const char**        value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetOpaqueEx(
    tibrvMsg            message,
    const char*         fieldName,
    const void**        value,
    tibrv_u32*          size,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_GetXmlEx(
    tibrvMsg            message,
    const char*         fieldName,
    const void**        value,
    tibrv_u32*          size,
    tibrv_u16           optIdentifier);

/* remove field */

extern tibrv_status
tibrvMsg_RemoveFieldEx(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u16           optIdentifier);

/* update field */

extern tibrv_status
tibrvMsg_UpdateMsgEx(
    tibrvMsg            message,
    const char*         fieldName,
    tibrvMsg            value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateIPAddr32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_ipaddr32      value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateIPPort16Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_ipport16      value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateDateTimeEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrvMsgDateTime * value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateBoolEx(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_bool          value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateI8Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i8            value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateI8ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i8*     array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateU8Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u8            value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateU8ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u8*     array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateI16Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i16           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateI16ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i16*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateU16Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u16           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateU16ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u16*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateI32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i32           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateI32ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i32*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateU32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u32           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateU32ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u32*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateU64Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_u64           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateU64ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_u64*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateI64Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_i64           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateI64ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_i64*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateF32Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_f32           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateF32ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_f32*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateF64Ex(
    tibrvMsg            message,
    const char*         fieldName,
    tibrv_f64           value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateF64ArrayEx(
    tibrvMsg            message,
    const char*         fieldName,
    const tibrv_f64*    array,
    tibrv_u32           numElements,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateStringArrayEx(
    tibrvMsg            message,
    const char*         field_name,
    const char**        value,
    tibrv_u32           num_elements,
    tibrv_u16           opt_identifier);

extern tibrv_status
tibrvMsg_UpdateMsgArrayEx(
    tibrvMsg            message,
    const char*         field_name,
    const tibrvMsg*     value,
    tibrv_u32           num_elements,
    tibrv_u16           opt_identifier);

extern tibrv_status
tibrvMsg_UpdateStringEx(
    tibrvMsg            message,
    const char*         fieldName,
    const char*         value,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateOpaqueEx(
    tibrvMsg            message,
    const char*         fieldName,
    const void*         value,
    tibrv_u32           size,
    tibrv_u16           optIdentifier);

extern tibrv_status
tibrvMsg_UpdateXmlEx(
    tibrvMsg            message,
    const char*         fieldName,
    const void*         value,
    tibrv_u32           size,
    tibrv_u16           optIdentifier);

#ifdef  __cplusplus
}
#endif

#endif /* _INCLUDED_tibrvmsgx_h */
