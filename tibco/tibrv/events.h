/*
 * Copyright (c) 1998-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: events.h 39740 2009-03-13 17:32:37Z pbrotea $
 */

#ifndef _INCLUDED_tibrvevents_h
#define _INCLUDED_tibrvevents_h

#include "types.h"
#include "status.h"

#if defined(__cplusplus)
extern "C" {
#endif

/* Creation */

extern tibrv_status
tibrvEvent_CreateListener(
    tibrvEvent*                 event,
    tibrvQueue                  queue,
    tibrvEventCallback          callback,
    tibrvTransport              transport,
    const char*                 subject,
    const void*                 closure);

extern tibrv_status
tibrvEvent_CreateVectorListener(
    tibrvEvent*                 event,
    tibrvQueue                  queue,
    tibrvEventVectorCallback    callback,
    tibrvTransport              transport,
    const char*                 subject,
    const void*                 closure);

tibrv_status
tibrvEvent_CreateGroupVectorListener(
    tibrvEvent*                     event,
    tibrvQueue                      queue,
    tibrvEventVectorCallback        cb,
    tibrvTransport                  tport,
    const char*                     subject,
    const void*                     arg,
    const void*                     msgGroupId);

extern tibrv_status
tibrvEvent_CreateTimer(
    tibrvEvent*                 event,
    tibrvQueue                  queue,
    tibrvEventCallback          callback,
    tibrv_f64                   interval,
    const void*                 closure);

extern tibrv_status
tibrvEvent_CreateIO(
    tibrvEvent*                 event,
    tibrvQueue                  queue,
    tibrvEventCallback          callback,
    tibrv_i32                   socketId,
    tibrvIOType                 ioType,
    const void*                 closure);

/* Destruction */

#define \
tibrvEvent_Destroy(event)       tibrvEvent_DestroyEx(event, NULL)

extern tibrv_status
tibrvEvent_DestroyEx(
    tibrvEvent                  event,
    tibrvEventOnComplete         completeCallback);


/* Properties */

extern tibrv_status
tibrvEvent_GetType(
    tibrvEvent                  event,
    tibrvEventType*             type);

extern tibrv_status
tibrvEvent_GetQueue(
    tibrvEvent                  event,
    tibrvQueue*                 queue);

/* Event type specific properties */

extern tibrv_status
tibrvEvent_GetListenerSubject(
    tibrvEvent                  event,
    const char**                subject);

extern tibrv_status
tibrvEvent_GetListenerTransport(
    tibrvEvent                  event,
    tibrvTransport*             transport);

extern tibrv_status
tibrvEvent_GetTimerInterval(
    tibrvEvent                  event,
    tibrv_f64*                  interval);

extern tibrv_status
tibrvEvent_ResetTimerInterval(
    tibrvEvent                  event,
    tibrv_f64                   newInterval);

extern tibrv_status
tibrvEvent_GetIOSource(
    tibrvEvent                  event,
    tibrv_i32*                  source);

extern tibrv_status
tibrvEvent_GetIOType(
    tibrvEvent                  event,
    tibrvIOType*                ioType);

#ifdef  __cplusplus
}
#endif

#endif /* _INCLUDED_tibrvevents_h */
