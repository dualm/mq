package tibco

//#cgo CFLAGS: -I./tibrv
//#cgo LDFLAGS: -L${SRCDIR}/tibrv -ltibrv
/*
#include <stdio.h>
*/
import "C"

var tibTransport *Transport
var tibmessage *Message

func TibInit(service, network string, daemons []string) error {
	if err := tibrvOpen(); err != nil {
		return err
	}

	transport := NewTransport()

	if err := transport.create(
		service,
		network,
		daemons,
	); err != nil {
		return err
	}

	tibTransport = transport

	tibmessage = NewMessage()

	return tibmessage.Create()
}

func TibDestroy() {
	tibmessage.Destroy()
	tibrvClose()
}

func tibSend(targetSubjectName, fieldName string, msg string) error {
	err := tibmessage.Reset()
	if err != nil {
		return err
	}

	tibmessage.SetSendSubject(targetSubjectName)
	tibmessage.AddString(
		fieldName,
		msg, 0)

	return tibTransport.send(tibmessage)
}

func tibSendRequest(targetSubjectName, fieldName string, msg string) (string, error) {
	err := tibmessage.Reset()
	if err != nil {
		return "", err
	}

	if err := tibmessage.SetSendSubject(targetSubjectName); err != nil {
		return "", err
	}

	if err := tibmessage.AddString(
		fieldName,
		msg, 0); err != nil {
		return "", err
	}

	re, err := tibTransport.sendRequest(tibmessage, 20)
	if err != nil {
		return "", err
	}

	return re.GetString(
		fieldName,
		0,
	)
}
