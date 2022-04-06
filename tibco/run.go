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
func (t *tibcoMq) Run(ctx context.Context, configID string, initConfig mq.ConfigFunc) error {
	conf, err := initConfig(configID)
	if err != nil {
		return fmt.Errorf("tibco init config error, Error: %w", err)
	}

	if conf == nil {
		return fmt.Errorf("nil config")
	}

	t.service = conf.GetString(TibcoService)
	t.network = conf.GetString(TibcoNetwork)
	t.daemon = conf.GetStringSlice(TibcoDaemon)
	t.fieldName = conf.GetString(TibcoFieldName)
	t.targetSubjectName = conf.GetString(TibcoSubjectName)

	if err := t.init(t.service, t.network, t.daemon); err != nil {
		return fmt.Errorf("tibco init error, Error: %w", err)
	}

	return nil
}

func (t *tibcoMq) Send(ctx context.Context, c chan<- mq.MqResponse, msg []mq.MqMessage) {
	for _, m := range msg {
		err := t.makeMsg(t.targetSubjectName, t.fieldName, string(m.Msg))

		if err != nil {
			c <- mq.MqResponse{
				Msg: nil,
				Err: fmt.Errorf("tibco send error, %w", err),
			}

			return
		}

		if m.IsEvent {
			t.send()
		} else {
			re, err := t.sendRequest(ctx)
			if err != nil {
				c <- mq.MqResponse{
					Msg: nil,
					Err: fmt.Errorf("%s, %s\n", err, m.Msg),
				}

				continue
			}

			c <- mq.MqResponse{
				Msg: []byte(re),
				Err: err,
			}
		}
	}
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
			errChan <- err

			return
		}

		defer func(){
			_re.Destroy()
		}()

		s, err := _re.GetString(t.fieldName, 0)
		if err != nil {
			errChan <- err

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
		return err
	}

	if err := t.message.SetSendSubject(targetSubjectName); err != nil {
		return err
	}

	if err := t.message.AddString(fieldName, msg, 0); err != nil {
		return err
	}

	return nil
}

func (t *tibcoMq) init(service, network string, daemons []string) error {
	if err := tibrvOpen(); err != nil {
		return err
	}

	t.transport = NewTransport()

	if err := t.transport.create(
		service,
		network,
		daemons,
	); err != nil {
		return err
	}

	t.message = NewMessage()

	if err := t.message.Create(); err != nil {
		return err
	}

	return nil
}
