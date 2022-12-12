package tibco

/*
#include <stdio.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"sync"
)

type TibCmListener struct {
	CmOption    *TibCmOption
	lock        sync.Mutex
	infoChan    chan<- string
	errChan     chan<- error
	cmTransport *CmTransport
	events      []*CmEvent
	messagePool *sync.Pool
}

func NewTibCmListener(opt *TibCmOption, infoChan chan<- string, errChan chan<- error) (*TibCmListener, error) {
	if err := tibrvOpen(); err != nil {
		return nil, err
	}

	transport, err := NewTransport(opt.TibOption)
	if err != nil {
		return nil, err
	}

	cmTransport, err := NewCmTransport(transport, opt)
	if err != nil {
		return nil, err
	}

	listener := &TibCmListener{
		CmOption:    opt,
		lock:        sync.Mutex{},
		infoChan:    infoChan,
		errChan:     errChan,
		cmTransport: cmTransport,
		events:      make([]*CmEvent, 0),
		messagePool: func() *sync.Pool {
			if opt.PooledMessage {
				return &sync.Pool{
					New: func() any {
						msg, err := NewMessage()
						if err != nil {
							errChan <- err

							return nil
						}

						return msg
					},
				}
			}

			return nil
		}(),
	}

	return listener, nil
}

func (l *TibCmListener) Close() error {
	if err := l.cmTransport.Close(); err != nil {
		return err
	}

	for i := range l.events {
		if err := l.events[i].Close(false); err != nil {
			return err
		}
	}

	if err := tibrvClose(); err != nil {
		return err
	}

	return nil
}

func (l *TibCmListener) Listen(subjectName string, cb TibrvCmEventCallback) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.cmTransport == nil {
		return errors.New("nil cmTransport")
	}

	listener, err := NewCmListener(nil, l.cmTransport, subjectName, cb)
	if err != nil {
		return err
	}

	l.events = append(l.events, listener)

	go func() {
		q, err := listener.GetQueue()
		if err != nil {
			l.errChan <- fmt.Errorf("1 %w", err)

			return
		}

		for {
			if err := q.Dispatch(); err != nil {
				l.errChan <- fmt.Errorf("2 %w", err)

				break
			}
		}
	}()

	return nil
}

func (l *TibCmListener) AddListener(cmName, subjedtName string) error {
	return l.cmTransport.AddListener(cmName, subjedtName)
}
