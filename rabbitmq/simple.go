package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleProducer interface {
	Send(ctx context.Context, mandatory bool, immediate bool, msg []byte) error // todo add comment
	Close() error
}

type SimpleConsumer interface {
	Recv() (chan []byte, error) // todo add comment
	Close() error
}

// todo add Confirm mode
type simple struct {
	dialOption     *DialOption
	queueOption    *QueueOption
	consumerOption *ConsumerOption
	errorChan      chan error
	sessionChan    chan chan *Session
	session        *Session
}

func NewSimpleProducer(dialOption *DialOption, queueOption *QueueOption) (SimpleProducer, error) {
	s := &simple{
		dialOption:  dialOption,
		queueOption: queueOption,
	}

	var err error
	s.session, err = simpleSession(dialOption, queueOption)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func NewSimpleConsumer(dialOption *DialOption, queueOption *QueueOption, consumerOption *ConsumerOption, errorChan chan error) (SimpleConsumer, error) {
	s := &simple{
		dialOption:     dialOption,
		queueOption:    queueOption,
		errorChan:      errorChan,
		consumerOption: consumerOption,
	}

	var err error
	s.session, err = simpleSession(dialOption, queueOption)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *simple) Send(ctx context.Context, mandatory, immediate bool, msg []byte) error {
	ch := s.session.Channel

	return ch.PublishWithContext(ctx, "", s.session.Queue.Name, mandatory, immediate, amqp.Publishing{
		Body: msg,
	})
}

func (s *simple) Recv() (chan []byte, error) {
	ch := s.session.Channel

	msgs, err := consumeDeclare(ch, "", s.consumerOption)
	if err != nil {
		return nil, err
	}

	re := make(chan []byte)

	go func() {
		for i := range msgs {
			go func(delivery amqp.Delivery) {
				if !s.consumerOption.AutoAck {
					err := delivery.Ack(false)
					if err != nil {
						s.errorChan <- err

						return
					}
				}

				re <- delivery.Body
			}(i)
		}
	}()

	return re, nil
}

func (s *simple) Close() error {
	err := s.session.Channel.Close()
	if err != nil {
		return err
	}

	err = s.session.Connection.Close()
	if err != nil {
		return err
	}

	return nil
}

func simpleSession(dialOption *DialOption, queueOption *QueueOption) (*Session, error) {
	conn, err := amqp.DialConfig(dialOption.url(), amqp.Config{
		Vhost: dialOption.VHost,
	})
	if err != nil {
		return nil, fmt.Errorf("rabbitmq/simple cannot (re)dial: %v: %q", err, dialOption.url())
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq/simple cannot create channel: %v", err)
	}

	queue, err := queueDeclare(ch, queueOption)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq/simple cannot create channel: %v", err)
	}

	s := &Session{
		Connection: conn,
		Channel:    ch,
		Queue:      queue,
	}

	return s, nil
}
