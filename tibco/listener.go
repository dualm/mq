package tibco

/*
#include <stdio.h>
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"sync"
	"unsafe"
)

type TibListener struct {
	*TibOption
	lock          sync.Mutex
	infoChan      chan<- string
	errChan       chan<- error
	transport     *Transport
	events        []*Event
	messagePool   *sync.Pool
	listenSubject string
	cb            func(Event, Message)
}

func NewTibListener(opt *TibOption, cb func(Event, Message), infoChan chan<- string, errChan chan<- error) (*TibListener, error) {
	if err := tibrvOpen(); err != nil {
		return nil, err
	}

	transport, err := NewTransport(opt)
	if err != nil {
		return nil, err
	}

	listener := &TibListener{
		TibOption: opt,
		lock:      sync.Mutex{},
		infoChan:  infoChan,
		errChan:   errChan,
		transport: transport,
		events:    make([]*Event, 0),
		messagePool: &sync.Pool{
			New: func() any {
				msg, err := NewMessage()
				if err != nil {
					errChan <- err

					return nil
				}

				return msg
			},
		},
		cb: cb,
	}

	return listener, nil
}

func (l *TibListener) Close() error {
	if err := l.transport.Close(); err != nil {
		return err
	}

	for i := range l.events {
		if err := l.events[i].Close(); err != nil {
			return err
		}
	}

	if err := tibrvClose(); err != nil {
		return err
	}

	return nil
}

func (l *TibListener) Listen(subjectName string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.transport == nil {
		return errors.New("all transports are nil")
	}

	_h := cgo.NewHandle(l.cb)
	listener, err := NewListener(nil, l.transport, subjectName, unsafe.Pointer(&_h))
	if err != nil {
		return err
	}

	l.events = append(l.events, listener)

	l.listenSubject = subjectName

	go func() {
		q, err := listener.GetQueue()
		if err != nil {
			l.errChan <- err

			return
		}

		for {
			if err := q.Dispatch(); err != nil {
				l.errChan <- err

				break
			}
		}
	}()

	return nil
}

func (l *TibListener) SendReply(msg string, listenSubject, replySubject string) error {
	_msg, err := NewMessage()
	if err != nil {
		return err
	}

	if err = _msg.AddString(l.TibOption.FieldName, msg, 0); err != nil {
		return err
	}

	if err = _msg.SetSendSubject(listenSubject); err != nil {
		return err
	}

	if err = _msg.SetReplySubject(replySubject); err != nil {
		return err
	}

	return l.transport.Send(_msg)
}

func (l *TibListener) SendReplyMessage(msg string, incoming *Message) error {
	_msg, err := NewMessage()
	if err != nil {
		return err
	}

	if err = _msg.AddString(l.TibOption.FieldName, msg, 0); err != nil {
		return err
	}

	return l.transport.SendReply(_msg, incoming)
}
