/*
 * Copyright (c) 1998-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: types.h 41785 2009-08-19 20:34:37Z mwoehlke $
 */

#ifndef _INCLUDED_tibrvtypes_h
#define _INCLUDED_tibrvtypes_h


/* Subject constraints */

#define TIBRV_SUBJECT_MAX               (255)
#define TIBRV_SUBJECT_TOKEN_MAX         (127)

/*
 * 64 bit integer support.
 *
 * If the type your compiler uses for a 64 bit integer
 * is not covered below, define TIBRV_I64_TYPE to be
 * the type your compiler uses for 64 bit integers.
 */

#if !defined(TIBRV_I64_TYPE)
# if defined(_WIN32) || defined(__vms)
#  define TIBRV_I64_TYPE __int64
# elif defined(__alpha)
#  define TIBRV_I64_TYPE long
# else
#  define TIBRV_I64_TYPE long long
# endif
#endif



typedef signed char             tibrv_i8;
typedef unsigned char           tibrv_u8;

typedef short                   tibrv_i16;
typedef unsigned short          tibrv_u16;

typedef int                     tibrv_i32;
typedef unsigned int            tibrv_u32;

typedef TIBRV_I64_TYPE          tibrv_i64;
typedef unsigned TIBRV_I64_TYPE tibrv_u64;

typedef float                   tibrv_f32;
typedef double                  tibrv_f64;

typedef unsigned short          tibrv_ipport16;
typedef unsigned int            tibrv_ipaddr32;

typedef enum {
    TIBRV_FALSE                 = 0,
    TIBRV_TRUE                  = 1
} tibrv_bool;

typedef struct __tibrvMsg*      tibrvMsg;

typedef struct __tibrvSerializedBuffer* tibrvSerializedBuffer;

typedef tibrv_u32               tibrvId;

#define TIBRV_INVALID_ID        (0)

typedef tibrvId                 tibrvEvent;
typedef tibrvEvent              tibrvPollEvent;
typedef tibrvId                 tibrvQueue;
typedef tibrvQueue              tibrvPollQueue;
typedef tibrvId                 tibrvTransport;
typedef tibrvId                 tibrvQueueGroup;
typedef tibrvId                 tibrvDispatchable;
typedef tibrvId                 tibrvDispatcher;

typedef tibrvId                 tibrvContentDesc;

typedef tibrv_u32               tibrvEventType;

#define TIBRV_TIMER_EVENT       (1)
#define TIBRV_IO_EVENT          (2)
#define TIBRV_LISTEN_EVENT      (3)

#ifdef __cplusplus
#define TIBRV_DEFAULT_QUEUE     static_cast<tibrvId>(1)
#define TIBRV_PROCESS_TRANSPORT static_cast<tibrvId>(10)
#else
#define TIBRV_DEFAULT_QUEUE     ((tibrvId)1)
#define TIBRV_PROCESS_TRANSPORT ((tibrvId)10)
#endif


typedef enum {
    TIBRVQUEUE_DISCARD_NONE     = 0,
    TIBRVQUEUE_DISCARD_NEW      = 1,
    TIBRVQUEUE_DISCARD_FIRST    = 2,
    TIBRVQUEUE_DISCARD_LAST     = 3
} tibrvQueueLimitPolicy;

#define TIBRVQUEUE_DEFAULT_POLICY   (TIBRVQUEUE_DISCARD_NONE)
#define TIBRVQUEUE_DEFAULT_PRIORITY (1)

#ifdef __cplusplus
#define TIBRV_WAIT_FOREVER      static_cast<tibrv_f64>(-1.0)
#define TIBRV_NO_WAIT           static_cast<tibrv_f64>(+0.0)
#else
#define TIBRV_WAIT_FOREVER      ((tibrv_f64)-1.0)
#define TIBRV_NO_WAIT           ((tibrv_f64)+0.0)
#endif

#if defined(__cplusplus)
extern "C" {
#endif

typedef void
(*tibrvEventCallback) (
    tibrvEvent          event,
    tibrvMsg            message,
    void*               closure);

typedef void
(*tibrvEventVectorCallback) (
    tibrvMsg            messages[],
    tibrv_u32           numMessages);

/* Optional destruction complete callbacks */
typedef void
(*tibrvEventOnComplete) (
    tibrvEvent          event,
    void*               closure);

typedef void
(*tibrvQueueOnComplete) (
    tibrvQueue          queue,
    void*               closure);

#ifdef  __cplusplus
}
#endif

typedef enum
{
    TIBRV_IO_READ               = 1,
    TIBRV_IO_WRITE              = 2,
    TIBRV_IO_EXCEPTION          = 4

} tibrvIOType;

typedef enum
{
    TIBRV_TRANSPORT_DEFAULT_BATCH = 0,
    TIBRV_TRANSPORT_TIMER_BATCH   = 1

} tibrvTransportBatchMode;

#define TIBRVMSG_DATETIME_SEC_MAX ( (((tibrv_i64)0x80) << 32) - 1)
#define TIBRVMSG_DATETIME_SEC_MIN (-(((tibrv_i64)0x80) << 32))

#define TIBRVMSG_DATETIME_STRING_SIZE	(32)

typedef struct tibrvMsgDateTime
{
    tibrv_i64                   sec;
    tibrv_u32                   nsec;
} tibrvMsgDateTime;

typedef union
{
    tibrvMsg                    msg;

    const char*                 str;
    const void*                 buf;

    const void*                 array;

    tibrv_bool                  boolean;

    tibrv_i8                    i8;
    tibrv_u8                    u8;
    tibrv_i16                   i16;
    tibrv_u16                   u16;
    tibrv_i32                   i32;
    tibrv_u32                   u32;
    tibrv_i64                   i64;
    tibrv_u64                   u64;
    tibrv_f32                   f32;
    tibrv_f64                   f64;

    tibrv_ipport16              ipport16;
    tibrv_ipaddr32              ipaddr32;

    tibrvMsgDateTime            date;

} tibrvLocalData;

#define TIBRVMSG_FIELDNAME_MAX (127)

typedef struct tibrvMsgField
{
    const char*                 name;
    tibrv_u32                   size;
    tibrv_u32                   count;
    tibrvLocalData              data;
    tibrv_u16                   id;
    tibrv_u8                    type;

} tibrvMsgField;

/* Defined types */

#define TIBRVMSG_MSG            (1)

#define TIBRVMSG_DATETIME       (3)

#define TIBRVMSG_OPAQUE         (7)

#define TIBRVMSG_STRING         (8)

#define TIBRVMSG_BOOL           (9)

#define TIBRVMSG_I8             (14)
#define TIBRVMSG_U8             (15)
#define TIBRVMSG_I16            (16)
#define TIBRVMSG_U16            (17)
#define TIBRVMSG_I32            (18)
#define TIBRVMSG_U32            (19)
#define TIBRVMSG_I64            (20)
#define TIBRVMSG_U64            (21)

#define TIBRVMSG_F32            (24)
#define TIBRVMSG_F64            (25)

#define TIBRVMSG_IPPORT16       (26)
#define TIBRVMSG_IPADDR32       (27)

#define TIBRVMSG_ENCRYPTED      (32)

#define TIBRVMSG_NONE           (22)

/* Array support */
#define TIBRVMSG_I8ARRAY        (34)
#define TIBRVMSG_U8ARRAY        (35)
#define TIBRVMSG_I16ARRAY       (36)
#define TIBRVMSG_U16ARRAY       (37)
#define TIBRVMSG_I32ARRAY       (38)
#define TIBRVMSG_U32ARRAY       (39)
#define TIBRVMSG_I64ARRAY       (40)
#define TIBRVMSG_U64ARRAY       (41)

#define TIBRVMSG_F32ARRAY       (44)
#define TIBRVMSG_F64ARRAY       (45)

#define TIBRVMSG_XML            (47)

#define TIBRVMSG_STRINGARRAY    (48)
#define TIBRVMSG_MSGARRAY       (49)

#define TIBRVMSG_USER_FIRST     (128)
#define TIBRVMSG_USER_LAST      (255)

/* 128 - 255 are available for user defined types */

#define TIBRVMSG_NO_TAG (0)

#endif /* _INCLUDED_tibrvtypes_h */
