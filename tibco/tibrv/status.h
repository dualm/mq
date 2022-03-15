/*
 * Copyright (c) 1998-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: status.h 45610 2010-03-30 15:55:46Z jpenning $
 */

#ifndef _INCLUDED_tibrvstatus_h
#define _INCLUDED_tibrvstatus_h

#if defined(__cplusplus)
extern "C" {
#endif

typedef enum
{
    TIBRV_OK                        = 0,

    TIBRV_INIT_FAILURE              = 1,
    TIBRV_INVALID_TRANSPORT         = 2,
    TIBRV_INVALID_ARG               = 3,
    TIBRV_NOT_INITIALIZED           = 4,
    TIBRV_ARG_CONFLICT              = 5,

    TIBRV_SERVICE_NOT_FOUND         = 16,
    TIBRV_NETWORK_NOT_FOUND         = 17,
    TIBRV_DAEMON_NOT_FOUND          = 18,
    TIBRV_NO_MEMORY                 = 19,
    TIBRV_INVALID_SUBJECT           = 20,
    TIBRV_DAEMON_NOT_CONNECTED      = 21,
    TIBRV_VERSION_MISMATCH          = 22,
    TIBRV_SUBJECT_COLLISION         = 23,
    TIBRV_VC_NOT_CONNECTED          = 24,

    TIBRV_NOT_PERMITTED             = 27,

    TIBRV_INVALID_NAME              = 30,
    TIBRV_INVALID_TYPE              = 31,
    TIBRV_INVALID_SIZE              = 32,
    TIBRV_INVALID_COUNT             = 33,

    TIBRV_NOT_FOUND                 = 35,
    TIBRV_ID_IN_USE                 = 36,
    TIBRV_ID_CONFLICT               = 37,
    TIBRV_CONVERSION_FAILED         = 38,
    TIBRV_RESERVED_HANDLER          = 39,
    TIBRV_ENCODER_FAILED            = 40,
    TIBRV_DECODER_FAILED            = 41,
    TIBRV_INVALID_MSG               = 42,
    TIBRV_INVALID_FIELD             = 43,
    TIBRV_INVALID_INSTANCE          = 44,
    TIBRV_CORRUPT_MSG               = 45,
    TIBRV_ENCODING_MISMATCH         = 46,

    TIBRV_TIMEOUT                   = 50,
    TIBRV_INTR                      = 51,

    TIBRV_INVALID_DISPATCHABLE      = 52,
    TIBRV_INVALID_DISPATCHER        = 53,

    TIBRV_INVALID_EVENT             = 60,
    TIBRV_INVALID_CALLBACK          = 61,
    TIBRV_INVALID_QUEUE             = 62,
    TIBRV_INVALID_QUEUE_GROUP       = 63,

    TIBRV_INVALID_TIME_INTERVAL     = 64,

    TIBRV_INVALID_IO_SOURCE         = 65,
    TIBRV_INVALID_IO_CONDITION      = 66,
    TIBRV_SOCKET_LIMIT              = 67,

    TIBRV_OS_ERROR                  = 68,

    TIBRV_INSUFFICIENT_BUFFER       = 70,
    TIBRV_EOF                       = 71,
    TIBRV_INVALID_FILE              = 72,
    TIBRV_FILE_NOT_FOUND            = 73,
    TIBRV_IO_FAILED                 = 74,

    TIBRV_NOT_FILE_OWNER            = 80,
    TIBRV_USERPASS_MISMATCH         = 81,

    TIBRV_TOO_MANY_NEIGHBORS        = 90,
    TIBRV_ALREADY_EXISTS            = 91,

    TIBRV_PORT_BUSY                 = 100,
    TIBRV_DELIVERY_FAILED           = 101,
    TIBRV_QUEUE_LIMIT               = 102,

    TIBRV_INVALID_CONTENT_DESC      = 110,
    TIBRV_INVALID_SERIALIZED_BUFFER = 111,
    TIBRV_DESCRIPTOR_NOT_FOUND      = 115,
    TIBRV_CORRUPT_SERIALIZED_BUFFER = 116,

    TIBRV_IPM_ONLY                  = 117

} tibrv_status;

extern const char*
tibrvStatus_GetText(
    tibrv_status        status);


#ifdef  __cplusplus
}
#endif

#endif /* _INCLUDED_tibrvstatus_h */
