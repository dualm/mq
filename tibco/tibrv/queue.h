/*
 * Copyright (c) 1998-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: queue.h 39400 2009-02-11 20:29:48Z pbrotea $
 */

#ifndef _INCLUDED_tibrvqueue_h
#define _INCLUDED_tibrvqueue_h

#include "types.h"
#include "status.h"

#if defined(__cplusplus)
extern "C" {
#endif

extern tibrv_status
tibrvQueue_Create(
    tibrvQueue*                 eventQueue);

#define \
tibrvQueue_Dispatch(q)          tibrvQueue_TimedDispatch((q),TIBRV_WAIT_FOREVER)

#define \
tibrvQueue_Poll(q)              tibrvQueue_TimedDispatch((q), TIBRV_NO_WAIT)

extern tibrv_status
tibrvQueue_TimedDispatch(
    tibrvQueue                  eventQueue,
    tibrv_f64                   timeout);

tibrv_status
tibrvQueue_TimedDispatchOneEvent(
    tibrvQueue                  queue,
    tibrv_f64                   waitTime);

#define \
tibrvQueue_Destroy(queue)       tibrvQueue_DestroyEx(queue, NULL, NULL)

extern tibrv_status
tibrvQueue_DestroyEx(
    tibrvQueue                  eventQueue,
    tibrvQueueOnComplete        completeCallback,
    const void*                 closure);

/* Properties */

extern tibrv_status
tibrvQueue_GetCount(
    tibrvQueue                  eventQueue,
    tibrv_u32*                  numEvents);

extern tibrv_status
tibrvQueue_GetPriority(
    tibrvQueue                  eventQueue,
    tibrv_u32*                  priority);

extern tibrv_status
tibrvQueue_SetPriority(
    tibrvQueue                  eventQueue,
    tibrv_u32                   priority);

extern tibrv_status
tibrvQueue_GetLimitPolicy(
    tibrvQueue                  eventQueue,
    tibrvQueueLimitPolicy*      policy,
    tibrv_u32*                  maxEvents,
    tibrv_u32*                  discardAmount);

extern tibrv_status
tibrvQueue_SetLimitPolicy(
    tibrvQueue                  eventQueue,
    tibrvQueueLimitPolicy       policy,
    tibrv_u32                   maxEvents,
    tibrv_u32                   discardAmount);

extern tibrv_status
tibrvQueue_SetName(
    tibrvQueue                  eventQueue,
    const char*                 queueName);

extern tibrv_status
tibrvQueue_GetName(
    tibrvQueue                  eventQueue,
    const char**                queueName);

typedef void
(*tibrvQueueHook) (
    tibrvQueue          eventQueue,
    void*               closure);

extern tibrv_status
tibrvQueue_SetHook(
    tibrvQueue                  eventQueue,
    tibrvQueueHook              eventQueueHook,
    void*                       closure);

extern tibrv_status
tibrvQueue_GetHook(
    tibrvQueue                  eventQueue,
    tibrvQueueHook*             eventQueueHook);

#define \
tibrvQueue_RemoveHook(q)        tibrvQueue_SetHook((q), NULL, NULL);

#ifdef  __cplusplus
}
#endif

#endif /* _INCLUDED_tibrvqueue_h */
