package tibco

/*
#include <stdio.h>
*/
import "C"
import (
	"context"
	"fmt"
	"sync"
)

const (
	MaxWaitingTime = 100
)

type TibSender struct {
	Option      *TibOption
	lock        sync.Mutex
	infoChan    chan<- string
	errChan     chan<- error
	transport   *Transport
	message     *Message
	messagePool *sync.Pool
}

type TibOption struct {
	FieldName         string
	Service           string
	Network           string
	Daemon            []string
	TargetSubjectName string
	SourceSubjectName string
	PooledMessage     bool
}

func NewTibSender(opt *TibOption, infoChan chan<- string, errChan chan<- error) *TibSender {
	return &TibSender{
		lock:     sync.Mutex{},
		Option:   opt,
		infoChan: infoChan,
		errChan:  errChan,
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
}

func (t *TibSender) Run() error {
	if err := t.init(); err != nil {
		return fmt.Errorf("tibco init error, Error: %w", err)
	}

	return nil
}

func (t *TibSender) Close() error {
	err := t.message.Close()
	if err != nil {
		return err
	}

	err = t.transport.Close()
	if err != nil {
		return err
	}

	return tibrvClose()
}

func (t *TibSender) send() error {
	return t.transport.Send(t.message)
}

// SendRequest 发送消息并在最长timeOut秒内返回。返回时间等于timeOut时报错消息超时
func (t *TibSender) SendRequest(msg string, timeOut float64) (*Message, error) {
	err := t.makeMsg(msg, t.message)
	if err != nil {
		return nil, err
	}

	return t.transport.sendRequest(t.message, timeOut)
}

// SendReport 发送消息
func (t *TibSender) SendReport(msg string) error {
	err := t.makeMsg(msg, t.message)
	if err != nil {
		return err
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	return t.transport.Send(t.message)
}

func (t *TibSender) SendReports(messages []string) error {
	tibMessages := make([]*Message, 0, len(messages))

	for i := range messages {
		_m := t.messagePool.Get().(*Message)

		err := _m.Reset()
		if err != nil {
			return err
		}

		if err = t.makeMsg(messages[i], _m); err != nil {
			return err
		}

		tibMessages = append(tibMessages, _m)
	}

	err := t.transport.Sendv(tibMessages)
	if err != nil {
		return err
	}

	for i := range tibMessages {
		t.messagePool.Put(tibMessages[i])
	}

	return err
}

func (t *TibSender) sendRequest(ctx context.Context) (string, error) {
	resultChan := make(chan string)
	errChan := make(chan error)

	go func() {
		_re, err := t.transport.sendRequest(t.message, MaxWaitingTime)
		if err != nil {
			errChan <- fmt.Errorf("send request error, Error: %w", err)

			return
		}

		defer func() {
			_ = _re.Close()
		}()

		s, err := _re.GetString(t.Option.FieldName, 0)
		if err != nil {
			errChan <- fmt.Errorf("get field name from response error, Error: %w", err)

			return
		}

		resultChan <- s
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errChan:
		return "", err
	case result := <-resultChan:
		return result, nil
	}
}

func (t *TibSender) makeMsg(msg string, tibMessage *Message) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	err := tibMessage.Reset()
	if err != nil {
		return fmt.Errorf("reset message error, Error: %w", err)
	}

	if err := tibMessage.SetSendSubject(t.Option.TargetSubjectName); err != nil {
		return fmt.Errorf("set send subject error, Error: %w", err)
	}

	if err := tibMessage.AddString(t.Option.FieldName, msg, 0); err != nil {
		return fmt.Errorf("add string to message error, Error: %w", err)
	}

	return nil
}

func (t *TibSender) init() error {
	var err error

	if err := tibrvOpen(); err != nil {
		return fmt.Errorf("tibrv open error, Error: %w", err)
	}

	t.transport, err = NewTransport(t.Option)
	if err != nil {
		return fmt.Errorf("creating cmTransport error, Error: %w", err)
	}

	t.message, err = NewMessage()
	if err != nil {
		return fmt.Errorf("create message error, Error: %w", err)
	}

	return nil
}
