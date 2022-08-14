package tibco

//#cgo CFLAGS: -I./tibrv
//#cgo LDFLAGS: -L${SRCDIR}/tibrv -ltibrv
/*
#include <stdio.h>
*/
import "C"
import (
	"context"
	"fmt"
	"sync"

	"github.com/dualm/common"
	"github.com/dualm/mq"
)

const (
	TIBCO_MAX_WAITING_TIME = 100
)

type TibSender struct {
	*TibOption
	lock        sync.Mutex
	infoChan    chan<- string
	errChan     chan<- error
	transport   *Transport
	message     *Message
	messagePool sync.Pool
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

func NewMq(opt *TibOption, infoChan chan<- string, errChan chan<- error) mq.Mq {
	return &TibSender{
		TibOption:   opt,
		lock:        sync.Mutex{},
		infoChan:    infoChan,
		errChan:     errChan,
		transport:   &Transport{},
		message:     &Message{},
		messagePool: sync.Pool{},
	}
}

func NewTibSender(opt *TibOption, infoChan chan<- string, errChan chan<- error) *TibSender {
	return &TibSender{
		lock:      sync.Mutex{},
		TibOption: opt,
		infoChan:  infoChan,
		errChan:   errChan,
		messagePool: func() sync.Pool {
			if opt.PooledMessage {
				return sync.Pool{
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

			return sync.Pool{}
		}(),
	}
}

// Run implements mq.Mq
func (t *TibSender) Run(ctx context.Context) error {
	if err := t.init(); err != nil {
		return fmt.Errorf("tibco init error, Error: %w", err)
	}

	return nil
}

func (t *TibSender) Send(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) <-chan struct{} {
	c := make(chan struct{})

	go func() {
		for _, m := range msg {
			if m.Msg == nil {
				if responseChan != nil {
					common.SendInNewRT(ctx, mq.MqResponse{}, responseChan)
				}

				continue
			}

			err := t.makeMsg(t.TargetSubjectName, t.FieldName, string(m.Msg))

			if err != nil {
				common.SendInNewRT(ctx, mq.MqResponse{Msg: nil, Err: fmt.Errorf("make tibco message error, %w", err)}, responseChan)

				continue
			}

			if m.IsEvent {
				err := t.send()

				if responseChan != nil {
					common.SendInNewRT(
						ctx,
						mq.MqResponse{
							Msg: nil,
							Err: func(e error) error {
								if e != nil {
									return fmt.Errorf("send tibco message error, %w", err)
								}
								return nil
							}(err),
						},
						responseChan,
					)
				}
			} else {
				re, err := t.sendRequest(ctx)
				if err != nil {
					common.SendInNewRT(
						ctx,
						mq.MqResponse{
							Msg: nil,
							Err: fmt.Errorf("%s, %s\n", err, m.Msg),
						},
						responseChan,
					)

					continue
				}

				common.SendInNewRT(
					ctx,
					mq.MqResponse{
						Msg: []byte(re),
						Err: err,
					},
					responseChan,
				)
			}
		}

		c <- struct{}{}
	}()

	return c
}

func (t *TibSender) SendEvents(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage, targetSubjectName string) <-chan struct{} {
	c := make(chan struct{})

	go func() {
		_events := make([]*Message, 0)

		for i, m := range msg {
			if m.Msg == nil {
				if responseChan != nil {
					common.SendInNewRT(ctx, mq.MqResponse{}, responseChan)
				}

				continue
			}

			if !m.IsEvent {
				if responseChan != nil {
					common.SendInNewRT(ctx, mq.MqResponse{
						Err: fmt.Errorf("message is not a event, message index is %d", i),
					}, responseChan)
				}
			}

			if !t.TibOption.PooledMessage {
				err := t.makeMsg(targetSubjectName, t.FieldName, string(m.Msg))

				if err != nil {
					common.SendInNewRT(ctx, mq.MqResponse{Msg: nil, Err: fmt.Errorf("make tibco message error, %w", err)}, responseChan)

					continue
				}
			} else {
				msg, ok := t.messagePool.Get().(*Message)
				if !ok {
					t.errChan <- fmt.Errorf("convert to Message error, message index: %d", i)
				}

				_events = append(_events, msg)
			}
		}

		err := t.transport.Sendv(_events)

		if responseChan != nil {
			common.SendInNewRT(
				ctx,
				mq.MqResponse{
					Msg: nil,
					Err: func(e error) error {
						if e != nil {
							return fmt.Errorf("send tibco message error, %w", err)
						}
						return nil
					}(err),
				},
				responseChan,
			)
		}

		c <- struct{}{}
	}()

	return c
}

func (t *TibSender) Close(ctx context.Context) {
	t.message.Destroy()
	t.transport.Destroy()

	tibrvClose()
}

func (t *TibSender) send() error {
	return t.transport.Send(t.message)
}

func (t *TibSender) sendRequest(ctx context.Context) (string, error) {
	resultChan := make(chan string)
	errChan := make(chan error)

	go func() {
		_re, err := t.transport.sendRequest(t.message, TIBCO_MAX_WAITING_TIME)
		if err != nil {
			errChan <- fmt.Errorf("send request error, Error: %w", err)

			return
		}

		defer func() {
			_re.Destroy()
		}()

		s, err := _re.GetString(t.FieldName, 0)
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

func (t *TibSender) makeMsg(targetSubjectName, fieldName string, msg string) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	err := t.message.Reset()
	if err != nil {
		return fmt.Errorf("reset message error, Error: %w", err)
	}

	if err := t.message.SetSendSubject(targetSubjectName); err != nil {
		return fmt.Errorf("set send subject error, Error: %w", err)
	}

	if err := t.message.AddString(fieldName, msg, 0); err != nil {
		return fmt.Errorf("add string to message error, Error: %w", err)
	}

	return nil
}

func (t *TibSender) init() error {
	var err error

	if err := tibrvOpen(); err != nil {
		return fmt.Errorf("tibrv open error, Error: %w", err)
	}

	t.transport, err = NewTransport(t.Service, t.Network, t.Daemon)
	if err != nil {
		return fmt.Errorf("creating transport error, Error: %w", err)
	}

	t.message, err = NewMessage()
	if err != nil {
		return fmt.Errorf("create message error, Error: %w", err)
	}

	return nil
}
