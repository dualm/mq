/*
 * Copyright (c) 1998-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: qgroup.h 12313 2005-11-23 18:08:58Z dleshc $
 */

#ifndef _INCLUDED_tibrvqgroup_h
#define _INCLUDED_tibrvqgroup_h

#include "types.h"
#include "status.h"

#if defined(__cplusplus)
extern "C" {
#endif

extern tibrv_status
tibrvQueueGroup_Create(
    tibrvQueueGroup*            eventQueueGroup);

#define \
tibrvQueueGroup_Dispatch(qg)    tibrvQueueGroup_TimedDispatch((qg), \
                                                          TIBRV_WAIT_FOREVER)
#define \
tibrvQueueGroup_Poll(qg)        tibrvQueueGroup_TimedDispatch((qg), \
                                                          TIBRV_NO_WAIT)
extern tibrv_status
tibrvQueueGroup_TimedDispatch(
    tibrvQueueGroup             eventQueueGroup,
    tibrv_f64                   timeout);

extern tibrv_status
tibrvQueueGroup_Destroy(
    tibrvQueueGroup             eventQueueGroup);

extern tibrv_status
tibrvQueueGroup_Add(
    tibrvQueueGroup             eventQueueGroup,
    tibrvQueue                  eventQueue);

extern tibrv_status
tibrvQueueGroup_Remove(
    tibrvQueueGroup             eventQueueGroup,
    tibrvQueue                  eventQueue);

#ifdef  __cplusplus
}
#endif

#endif /* _INCLUDED_tibrvqgroup_h */
