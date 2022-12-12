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

type TibCmSender struct {
	CmOption    *TibCmOption
	lock        sync.Mutex
	infoChan    chan<- string
	errChan     chan<- error
	cmTransport *CmTransport
	message     *Message
	messagePool *sync.Pool
}

type TibCmOption struct {
	*TibOption
	// CmName 显示设置时，为可重用模式, 能够实现程序不可用及重启后的持续消息传输, 但同一时间同一个网络上需要全局唯一. 若缺省则系统自动生成
	// 全局唯一ID, 但此时为不可重用模式.
	CmName string
	// 同一时间, 全局唯一
	LedgerName string
	RequestOld bool
	SyncLedger bool
	RelayAgent string
	LimitTime  float64
}

func NewTibCmSender(opt *TibCmOption, infoChan chan<- string, errChan chan<- error) *TibCmSender {
	return &TibCmSender{
		lock:     sync.Mutex{},
		CmOption: opt,
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

func (t *TibCmSender) Run() error {
	if err := t.init(); err != nil {
		return fmt.Errorf("tibco init error, Error: %w", err)
	}

	return nil
}

func (t *TibCmSender) Close() error {
	err := t.message.Close()
	if err != nil {
		return err
	}

	err = t.cmTransport.Close()
	if err != nil {
		return err
	}

	return tibrvClose()
}

func (t *TibCmSender) send() error {
	return t.cmTransport.Send(t.message)
}

// SendRequest 发送消息并在最长timeOut秒内返回。返回时间等于timeOut时报错消息超时
func (t *TibCmSender) SendRequest(msg, field, subject string, timeOut float64) (*Message, error) {
	if field == "" {
		field = t.CmOption.FieldName
	}

	if subject == "" {
		subject = t.CmOption.TargetSubjectName
	}

	err := t.makeMsg(subject, field, msg)
	if err != nil {
		return nil, err
	}

	return t.cmTransport.SendRequest(t.message, timeOut)
}

// SendReport 发送消息
func (t *TibCmSender) SendReport(msg, field, subject string) error {
	if field == "" {
		field = t.CmOption.FieldName
	}

	if subject == "" {
		subject = t.CmOption.TargetSubjectName
	}

	err := t.makeMsg(subject, field, msg)
	if err != nil {
		return err
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	return t.SendReportMessage(t.message)
}

func (t *TibCmSender) SendReportMessage(msg *Message) error {
	return t.cmTransport.Send(msg)
}

// GetPooledMessage 如果启动池化消息，返回一个新的消息。否则返回空
func (t *TibCmSender) GetPooledMessage() (*Message, error) {
	if t.CmOption.PooledMessage {
		msg := t.messagePool.Get().(*Message)
		err := msg.Reset()
		if err != nil {
			return nil, err
		}

		return msg, nil
	}

	return nil, nil
}

func (t *TibCmSender) PutMessage(msg *Message) {
	t.messagePool.Put(msg)
}

func (t *TibCmSender) AddOpaque(msg *Message, field string, b []byte, fieldId uint16) (*Message, error) {
	if err := msg.AddOpaque(field, b, fieldId); err != nil {
		return nil, err
	}

	return msg, nil
}

func (t *TibCmSender) sendRequest(ctx context.Context) (string, error) {
	resultChan := make(chan string)
	errChan := make(chan error)

	go func() {
		_re, err := t.cmTransport.SendRequest(t.message, MaxWaitingTime)
		if err != nil {
			errChan <- fmt.Errorf("send request error, Error: %w", err)

			return
		}

		defer func() {
			_ = _re.Close()
		}()

		s, err := _re.GetString(t.CmOption.FieldName, 0)
		if err != nil {
			errChan <- fmt.Errorf("GetString from response error, Error: %w", err)

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

func (t *TibCmSender) makeMsg(targetSubjectName, fieldName string, msg string) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	err := t.message.Reset()
	if err != nil {
		return fmt.Errorf("reset message error, Error: %w", err)
	}

	if err := t.message.SetSendSubject(targetSubjectName); err != nil {
		return fmt.Errorf("set send subject error, Error: %w", err)
	}

	if err := t.message.SetCMTimeLimit(t.CmOption.LimitTime); err != nil {
		return err
	}

	if err := t.message.AddString(fieldName, msg, 0); err != nil {
		return fmt.Errorf("add string to message error, Error: %w", err)
	}

	return nil
}

func (t *TibCmSender) init() error {
	var err error

	if err := tibrvOpen(); err != nil {
		return fmt.Errorf("tibrv open error, Error: %w", err)
	}

	transport, err := NewTransport(t.CmOption.TibOption)
	if err != nil {
		return err
	}

	t.cmTransport, err = NewCmTransport(transport, t.CmOption)
	if err != nil {
		return fmt.Errorf("creating cmTransport error, Error: %w", err)
	}

	t.message, err = NewMessage()
	if err != nil {
		return fmt.Errorf("create message error, Error: %w", err)
	}

	return nil
}

func (t *TibCmSender) AddListener(cmName, subject string) error {
	if subject == "" {
		subject = t.CmOption.TargetSubjectName
	}

	return t.cmTransport.AddListener(cmName, subject)
}
