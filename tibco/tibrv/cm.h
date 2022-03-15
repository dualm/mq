/*
 * Copyright (c) 1996-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * %W%
 */

#ifndef _INCLUDED_tibrvcm_h
#define _INCLUDED_tibrvcm_h

#include "tibrv.h"

#if defined(__cplusplus)
extern "C" {
#endif

extern const char*
tibrvcm_Version(void);

typedef tibrvId  tibrvcmTransport;
typedef tibrvId  tibrvcmEvent;

#define TIBRVCM_DEFAULT_TRANSPORT_TIMELIMIT (0) /* messages do not expire */

extern tibrv_status
tibrvcmTransport_Create(
    tibrvcmTransport*           cmTransport,
    tibrvTransport              transport,
    const char*                 cmName,
    tibrv_bool                  requestOld,
    const char*                 ledgerName,
    tibrv_bool                  syncLedger,
    const char*                 relayAgent);

extern tibrv_status
tibrvcmTransport_Send(
    tibrvcmTransport            cmTransport,
    tibrvMsg                    message);

extern tibrv_status
tibrvcmTransport_SendRequest(
    tibrvcmTransport            cmTransport,
    tibrvMsg                    message,
    tibrvMsg*                   reply,
    tibrv_f64                   idleTimeout);

extern tibrv_status
tibrvcmTransport_SendReply(
    tibrvcmTransport            cmTransport,
    tibrvMsg                    message,
    tibrvMsg                    requestMessage);

extern tibrv_status
tibrvcmTransport_GetTransport(
    tibrvcmTransport            cmTransport,
    tibrvTransport*             transport);

extern tibrv_status
tibrvcmTransport_GetName(
    tibrvcmTransport            cmTransport,
    const char**                cmName);

extern tibrv_status
tibrvcmTransport_GetRelayAgent(
    tibrvcmTransport            cmTransport,
    const char**                relayAgent);

extern tibrv_status
tibrvcmTransport_GetLedgerName(
    tibrvcmTransport            cmTransport,
    const char**                ledgerName);

extern tibrv_status
tibrvcmTransport_GetSyncLedger(
    tibrvcmTransport            cmTransport,
    tibrv_bool*                 syncLedger);

extern tibrv_status
tibrvcmTransport_GetRequestOld(
    tibrvcmTransport            cmTransport,
    tibrv_bool*                 requestOld);

extern tibrv_status
tibrvcmTransport_AllowListener(
    tibrvcmTransport            cmTransport,
    const char*                 cmName);

extern tibrv_status
tibrvcmTransport_DisallowListener(
    tibrvcmTransport            cmTransport,
    const char*                 cmName);

extern tibrv_status
tibrvcmTransport_AddListener(
    tibrvcmTransport            cmTransport,
    const char*                 cmName,
    const char*                 subject);

extern tibrv_status
tibrvcmTransport_RemoveListener(
    tibrvcmTransport            cmTransport,
    const char*                 cmName,
    const char*                 subject);

extern tibrv_status
tibrvcmTransport_RemoveSendState(
    tibrvcmTransport            cmTransport,
    const char*                 subject);

extern tibrv_status
tibrvcmTransport_SyncLedger(
    tibrvcmTransport            cmTransport);

extern tibrv_status
tibrvcmTransport_ConnectToRelayAgent(
    tibrvcmTransport            cmTransport);

extern tibrv_status
tibrvcmTransport_DisconnectFromRelayAgent(
    tibrvcmTransport            cmTransport);

extern tibrv_status
tibrvcmTransport_Destroy(
    tibrvcmTransport            cmTransport);

typedef void
(*tibrvcmTransportOnComplete) (
    tibrvcmTransport		destroyedTransport,
    void*			closure);

extern tibrv_status
tibrvcmTransport_DestroyEx(
    tibrvcmTransport		cmTransport,
    tibrvcmTransportOnComplete	completionFunction,
    void*			closure);

typedef void
(*tibrvcmEventCallback) (
    tibrvcmEvent                event,
    tibrvMsg                    message,
    void*                       closure);

extern tibrv_status
tibrvcmEvent_CreateListener(
    tibrvcmEvent*               cmListener,
    tibrvQueue                  queue,
    tibrvcmEventCallback        callback,
    tibrvcmTransport            cmTransport,
    const char*                 subject,
    const void*                 closure);

extern tibrv_status
tibrvcmEvent_GetQueue(
    tibrvcmEvent                event,
    tibrvQueue*                 queue);

extern tibrv_status
tibrvcmEvent_GetListenerSubject(
    tibrvcmEvent                event,
    const char**                subject);

extern tibrv_status
tibrvcmEvent_GetListenerTransport(
    tibrvcmEvent                event,
    tibrvcmTransport*           transport);

extern tibrv_status
tibrvcmEvent_SetExplicitConfirm(
    tibrvcmEvent                cmListener);

extern tibrv_status
tibrvcmEvent_ConfirmMsg(
    tibrvcmEvent                cmListener,
    tibrvMsg                    message);

#define TIBRVCM_CANCEL          (TIBRV_TRUE)
#define TIBRVCM_PERSIST         (TIBRV_FALSE)

#define \
tibrvcmEvent_Destroy(cmListener, cancelAgreements) \
    tibrvcmEvent_DestroyEx(cmListener, cancelAgreements, NULL)

extern tibrv_status
tibrvcmEvent_DestroyEx(
    tibrvcmEvent                cmListener,
    tibrv_bool                  cancelAgreements,
    tibrvEventOnComplete        completeCallback);

extern tibrv_status
tibrvMsg_GetCMSender(
    tibrvMsg                    message,
    const char**                senderName);

extern tibrv_status
tibrvMsg_GetCMSequence(
    tibrvMsg                    message,
    tibrv_u64*                  sequenceNumber);

extern tibrv_status
tibrvMsg_GetCMTimeLimit(
    tibrvMsg                    message,
    tibrv_f64*                  timeLimit);

extern tibrv_status
tibrvMsg_SetCMTimeLimit(
    tibrvMsg                    message,
    tibrv_f64                   timeLimit);

extern tibrv_status
tibrvcmTransport_GetDefaultCMTimeLimit(
    tibrvcmTransport            cmTransport,
    tibrv_f64*                  timeLimit);

extern tibrv_status
tibrvcmTransport_SetDefaultCMTimeLimit(
    tibrvcmTransport            cmTransport,
    tibrv_f64                   timeLimit);

typedef void*
(*tibrvcmReviewCallback) (
    tibrvcmTransport            cmTransport,
    const char*                 subject,
    tibrvMsg                    message,
    void*                       closure);

extern tibrv_status
tibrvcmTransport_ReviewLedger(
    tibrvcmTransport            cmTransport,
    tibrvcmReviewCallback       callback,
    const char*                 subject,
    const void*                 closure);

extern tibrv_status
tibrvcmTransport_ExpireMessages(
    tibrvcmTransport		cmTransport,
    const char*			subject,
    tibrv_u64			sequenceNumber);

/* Distributed Certified Transports */

#define TIBRVCM_DEFAULT_COMPLETE_TIME           (0)

#define TIBRVCM_DEFAULT_WORKER_WEIGHT           (1)
#define TIBRVCM_DEFAULT_WORKER_TASKS            (1)

#define TIBRVCM_DEFAULT_SCHEDULER_WEIGHT        (1)
#define TIBRVCM_DEFAULT_SCHEDULER_HB            (1.0)
#define TIBRVCM_DEFAULT_SCHEDULER_ACTIVE        (3.5)

#define \
tibrvcmTransport_CreateDistributedQueue(cmTransport, transport, cmName) \
    tibrvcmTransport_CreateDistributedQueueEx( \
        cmTransport,                           \
        transport,                             \
        cmName,                                \
        TIBRVCM_DEFAULT_WORKER_WEIGHT,         \
        TIBRVCM_DEFAULT_WORKER_TASKS,          \
        TIBRVCM_DEFAULT_SCHEDULER_WEIGHT,      \
        TIBRVCM_DEFAULT_SCHEDULER_HB,          \
        TIBRVCM_DEFAULT_SCHEDULER_ACTIVE)

extern tibrv_status
tibrvcmTransport_CreateDistributedQueueEx(
    tibrvcmTransport*           cmTransport,
    tibrvTransport              transport,
    const char*                 cmName,
    tibrv_u32                   workerWeight,
    tibrv_u32                   workerTasks,
    tibrv_u16                   schedulerWeight,
    tibrv_f64                   schedulerHeartbeat,
    tibrv_f64                   schedulerActivation);

extern tibrv_status
tibrvcmTransport_SetCompleteTime(
    tibrvcmTransport            cmTransport,
    tibrv_f64                   completeTime);

extern tibrv_status
tibrvcmTransport_GetCompleteTime(
    tibrvcmTransport            cmTransport,
    tibrv_f64*                  completeTime);

extern tibrv_status
tibrvcmTransport_SetWorkerWeight(
    tibrvcmTransport            cmTransport,
    tibrv_u32                   workerWeight);

extern tibrv_status
tibrvcmTransport_GetWorkerWeight(
    tibrvcmTransport            cmTransport,
    tibrv_u32*                  workerWeight);

extern tibrv_status
tibrvcmTransport_SetWorkerTasks(
    tibrvcmTransport            cmTransport,
    tibrv_u32                   listenerTasks);

extern tibrv_status
tibrvcmTransport_GetWorkerTasks(
    tibrvcmTransport            cmTransport,
    tibrv_u32*                  listenerTasks);

#define TIBRVCMQ_LIMIT_MSGS	0
#define TIBRVCMQ_LIMIT_BYTES	1

extern tibrv_status
tibrvcmTransport_SetTaskBacklogLimit(
    tibrvcmTransport		cmTransport,
    tibrv_u32			limitType,
    tibrv_u32			limitValue);
    
#define tibrvcmTransport_SetTaskBacklogLimitInBytes(cmTransport, limitBySizeInBytes)	\
	tibrvcmTransport_SetTaskBacklogLimit((cmTransport), TIBRVCMQ_LIMIT_BYTES, (limitBySizeInBytes))
#define tibrvcmTransport_SetTaskBacklogLimitInMessages(cmTransport, limitByMessages)	\
	tibrvcmTransport_SetTaskBacklogLimit((cmTransport), TIBRVCMQ_LIMIT_MSGS, (limitByMessages))

extern tibrv_status
tibrvcmTransport_SetPublisherInactivityDiscardInterval(
    tibrvcmTransport            cmTransport, 
    tibrv_i32                   timeout);

extern tibrv_status
tibrvcmTransport_GetUnassignedMessageCount(
    tibrvcmTransport		cmTransport,
    tibrv_u32			*msgCount);

#ifdef  __cplusplus
}
#endif

#endif /* _INCLUDED_tibrvcm_h */
