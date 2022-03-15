/*
 * Copyright (c) 1998-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: handler.h 12313 2005-11-23 18:08:58Z dleshc $
 */

#ifndef _INCLUDED_tibrvhandler_h
#define _INCLUDED_tibrvhandler_h

#include "types.h"

#if defined(__cplusplus)
extern "C" {
#endif

/*
 * Used by the tibrvMsg framework for preservation of pointer data
 * and proper maintenance of references.
 */
typedef enum {

  tibrvMsgData_Primitive        = 0,
  tibrvMsgData_MallocBlock      = 1,
  tibrvMsgData_SubMessage       = 2,
  tibrvMsgData_WireReference    = 3

} tibrvMsgDataType;


/* Encoder/Decoder/Converter */

typedef tibrv_status
(*tibrvMsgData_Encoder)(
    char**                      wire_buffer,
    tibrv_u32                   mem_available,
    tibrvMsgField*              field);

typedef tibrv_status
(*tibrvMsgData_Decoder)(
    char**                      wire_buffer,
    tibrvMsgField*              field,
    tibrvMsgDataType*           decoded_type);

typedef tibrv_status
(*tibrvMsgData_Converter)(
    tibrvMsgField*              field,
    tibrv_u8                    destination_type,
    tibrvMsgDataType*           converted_type);

/* Encoder/Decoder/Converter type assignment */

extern tibrv_status
tibrvMsg_SetHandlers(
    tibrv_u8                    type,
    tibrvMsgData_Encoder        encoder,
    tibrvMsgData_Decoder        decoder,
    tibrvMsgData_Converter      converter);

/* Used to aquire memory within Decoders and Converters */

extern void*
tibrvMsgData_Malloc(tibrv_u32       size);

/* Encoder and Decoder helper functions. */

extern tibrv_u32
tibrvMsgData_ByteSize(
    tibrv_u32                   content_size);

extern tibrv_status
tibrvMsgData_CopyBytes(
    char**                      buffer,
    const void*                 src,
    tibrv_u32                   size);

extern tibrv_status
tibrvMsgData_GetBytes(
    char**                      buffer,
    const void**                src,
    tibrv_u32*                  size);

/* 
 * Used to set and get field data sizes
 * when using memcpy and buffer advancement
 * directly. These functions can not be used
 * in conjunction with CopyBytes and GetBytes
 * helper functions.
 */

extern tibrv_status
tibrvMsgData_SetSize(
    char**                      buffer,
    tibrv_u32                   size);

extern tibrv_status
tibrvMsgData_GetSize(
    char**                      buffer,
    tibrv_u32*                  size);

#ifdef  __cplusplus
}
#endif

#endif /* _INCLUDED_tibrvhandler_h */
