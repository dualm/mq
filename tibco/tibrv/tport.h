/*
 * Copyright (c) 1998-2003 TIBCO Software Inc.
 * All rights reserved.
 * TIB/Rendezvous is protected under US Patent No. 5,187,787.
 * For more information, please contact:
 * TIBCO Software Inc., Palo Alto, California, USA
 *
 * $Id: tport.h 47699 2010-07-15 20:15:53Z jpenning $
 */

#ifndef _INCLUDED_tibrvtransport_h
#define _INCLUDED_tibrvtransport_h

#include "types.h"
#include "status.h"

#if defined(__cplusplus)
extern "C" {
#endif

extern tibrv_status
tibrvTransport_Create(
    tibrvTransport*     transport,
    const char*         service,
    const char*         network,
    const char*         daemonStr);

tibrv_status
tibrvTransport_CreateAcceptVc(
    tibrvTransport*     vcTransport,
    const char**        connectSubject,
    tibrvTransport      transport);

tibrv_status
tibrvTransport_CreateConnectVc(
    tibrvTransport*     vcTransport,
    const char*         connectSubject,
    tibrvTransport      transport);

tibrv_status
tibrvTransport_WaitForVcConnection(
    tibrvTransport      vcTransport,
    tibrv_f64           timeout);

extern tibrv_status
tibrvTransport_Send(
    tibrvTransport      transport,
    tibrvMsg            message);

extern tibrv_status
tibrvTransport_Sendv(
    tibrvTransport      transport,
    tibrvMsg*           vec,
    tibrv_u32           len);

extern tibrv_status
tibrvTransport_SendRequest(
    tibrvTransport      transport,
    tibrvMsg            message,
    tibrvMsg*           reply,
    tibrv_f64           idleTimeout);

extern tibrv_status
tibrvTransport_SendReply(
    tibrvTransport      transport,
    tibrvMsg            message,
    tibrvMsg            requestMessage);

extern tibrv_status
tibrvTransport_Destroy(
    tibrvTransport      transport);

extern tibrv_status
tibrvTransport_CreateInbox(
    tibrvTransport      transport,
    char*               subjectString,
    tibrv_u32           subjectLimit);

/* Specific network transport parameters */

extern tibrv_status
tibrvTransport_GetService(
    tibrvTransport      transport,
    const char**        serviceString);

extern tibrv_status
tibrvTransport_GetNetwork(
    tibrvTransport      transport,
    const char**        networkString);

extern tibrv_status
tibrvTransport_GetDaemon(
    tibrvTransport      transport,
    const char**        daemonString);

extern tibrv_status
tibrvTransport_SetDescription(
    tibrvTransport      transport,
    const char*         description);

extern tibrv_status
tibrvTransport_GetDescription(
    tibrvTransport      transport,
    const char**        description);

extern tibrv_status
tibrvTransport_SetSendingWaitLimit(
    tibrvTransport      transport,
    tibrv_u32           numBytes);

extern tibrv_status
tibrvTransport_GetSendingWaitLimit(
    tibrvTransport      transport,
    tibrv_u32*          numBytes);

tibrv_status
tibrvTransport_SetBatchMode(
    tibrvTransport      transport,
    tibrvTransportBatchMode mode);

extern tibrv_status
tibrvTransport_SetBatchSize(
    tibrvTransport      transport,
    tibrv_u32           numBytes);
/* Network transports with embedded license tickets */

extern tibrv_status
tibrvTransport_CreateLicensed(
    tibrvTransport*     transport,
    const char*         service,
    const char*         network,
    const char*         daemonStr,
    const char*         license_ticket);

extern tibrv_status
tibrvTransport_RequestReliability(
    tibrvTransport      transport,
    tibrv_f64           reliability);

#ifdef  __cplusplus
}
#endif

#endif /* _INCLUDED_tibrvtransport_h */

