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
	"strings"

	"github.com/dualm/common"
	"github.com/dualm/mq"
)

const (
	TIBCO_MAX_WAITING_TIME = 100
)

type tibcoMq struct {
	infoChan          chan<- string
	errChan           chan<- error
	targetSubjectName string
	fieldName         string
	service           string
	network           string
	daemon            []string
	transport         *Transport
	message           *Message
}

func New(infoChan chan<- string, errChan chan<- error) mq.Mq {
	return &tibcoMq{
		infoChan: infoChan,
		errChan:  errChan,
	}
}

// Run implements mq.Mq
func (t *tibcoMq) Run(ctx context.Context, initConfig mq.ConfigFunc, configID string, keys string) (map[string]string, error) {
	conf, err := initConfig(configID)
	if err != nil {
		return nil, fmt.Errorf("tibco init config error, Error: %w", err)
	}

	if conf == nil {
		return nil, fmt.Errorf("nil config")
	}

	t.service = common.GetString(conf, keys, TibcoService)
	t.network = common.GetString(conf, keys, TibcoNetwork)
	t.daemon = common.GetStringSlice(conf, keys, TibcoDaemon)
	t.fieldName = common.GetString(conf, keys, TibcoFieldName)
	t.targetSubjectName = common.GetString(conf, keys, TibcoSubjectName)

	if err := t.init(t.service, t.network, t.daemon); err != nil {
		return nil, fmt.Errorf("tibco init error, Error: %w", err)
	}

	return map[string]string{
		"Service":           t.service,
		"Network":           t.network,
		"Daemon":            strings.Join(t.daemon, ","),
		"FieldName":         t.fieldName,
		"TargetSubjectName": t.targetSubjectName,
	}, nil
}

func (t *tibcoMq) Send(ctx context.Context, responseChan chan<- mq.MqResponse, msg []mq.MqMessage) <-chan struct{} {
	for _, m := range msg {
		if m.Msg == nil {
			if responseChan != nil {
				responseChan <- mq.MqResponse{}
			}

			continue
		}

		err := t.makeMsg(t.targetSubjectName, t.fieldName, string(m.Msg))

		if err != nil {
			responseChan <- mq.MqResponse{
				Msg: nil,
				Err: fmt.Errorf("make tibco message error, %w", err),
			}

			return nil
		}

		if m.IsEvent {
			t.send()

			if responseChan != nil {
				responseChan <- mq.MqResponse{}
			}
		} else {
			re, err := t.sendRequest(ctx)
			if err != nil {
				responseChan <- mq.MqResponse{
					Msg: nil,
					Err: fmt.Errorf("%s, %s\n", err, m.Msg),
				}

				continue
			}

			responseChan <- mq.MqResponse{
				Msg: []byte(re),
				Err: err,
			}
		}
	}

	return nil
}

func (t *tibcoMq) Close(ctx context.Context) {
	t.message.Destroy()
	t.transport.destroy()
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

		s, err := _re.GetString(t.fieldName, 0)
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

func (t *tibcoMq) init(service, network string, daemons []string) error {
	if err := tibrvOpen(); err != nil {
		return fmt.Errorf("tibrv open error, Error: %w", err)
	}

	t.transport = NewTransport()

	if err := t.transport.create(
		service,
		network,
		daemons,
	); err != nil {
		return fmt.Errorf("creating transport error, Error: %w", err)
	}

	t.message = NewMessage()

	if err := t.message.Create(); err != nil {
		return fmt.Errorf("create message error, Error: %w", err)
	}

	return nil
}
