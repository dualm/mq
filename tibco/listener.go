package tibco

/*
#include <stdio.h>
*/
import "C"
import (
	"errors"
	"sync"
)

type TibListener struct {
	*TibOption
	lock        sync.Mutex
	infoChan    chan<- string
	errChan     chan<- error
	transport   *Transport
	events      []*Event
	messagePool *sync.Pool
}

func NewTibListener(opt *TibOption, infoChan chan<- string, errChan chan<- error) (*TibListener, error) {
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

func (l *TibListener) Listen(subjectName string, cb TibrvEventCallback) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.transport == nil {
		return errors.New("all transports are nil")
	}

	listener, err := NewListener(nil, l.transport, subjectName, cb)
	if err != nil {
		return err
	}

	l.events = append(l.events, listener)

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
