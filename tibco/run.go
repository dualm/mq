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

	"github.com/dualm/common"
	"github.com/dualm/mq"
)

const (
	TIBCO_MAX_WAITING_TIME = 100
)

type tibcoMq struct {
	*TibOption
	infoChan  chan<- string
	errChan   chan<- error
	transport *Transport
	message   *Message
}

type TibOption struct {
	FieldName string
	Service   string
	Network   string
	Daemon    []string
}

func New(opt *TibOption, infoChan chan<- string, errChan chan<- error) *tibcoMq {
	return &tibcoMq{
		TibOption: opt,
		infoChan:  infoChan,
		errChan:   errChan,
	}
}

// Run implements mq.Mq
func (t *tibcoMq) Run(ctx context.Context) error {
	if err := t.init(); err != nil {
		return fmt.Errorf("tibco init error, Error: %w", err)
	}

	return nil
}

func (t *tibcoMq) Send(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage, targetSubjectName string) <-chan struct{} {
	c := make(chan struct{})

	go func() {
		for _, m := range msg {
			if m.Msg == nil {
				if responseChan != nil {
					common.SendInNewRT(ctx, mq.MqResponse{}, responseChan)
				}

				continue
			}

			err := t.makeMsg(targetSubjectName, t.FieldName, string(m.Msg))

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

func (t *tibcoMq) Close(ctx context.Context) {
	t.message.Destroy()
	t.transport.Destroy()
}

func (t *tibcoMq) send() error {
	return t.transport.send(t.message)
}

func (t *tibcoMq) sendRequest(ctx context.Context) (string, error) {
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

func (t *tibcoMq) makeMsg(targetSubjectName, fieldName string, msg string) error {
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

func (t *tibcoMq) init() error {
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
