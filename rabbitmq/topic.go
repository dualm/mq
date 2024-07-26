package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
)

type TopicProducer struct {
	// todo add confirm option with error channel
	sendExchangeOption *ExchangeOption // 关联的发送exchange信息

	recvReplyChan sync.Map

	sendSession      *Session // 发送Request和Report的Session
	recvReplySession *Session // 收到Reply的Session

	recvReplyQueue amqp.Queue // 提供永久ReplyQueue供消费发送Reply消息
}

type TopicConsumer struct {
	// todo add confirm option with error channel
	recvExchangeOption  *ExchangeOption // 关联的exchange信息
	replyExchangeOption *ExchangeOption // 关联的发送Reply的exchange信息

	consumerOption   *ConsumerOption // 消费的配置
	recvSession      *Session        // 接收消息的Session
	sendReplySession *Session        // 发送Reply的Session
	errorChan        chan error
}

func NewTopicProducer(dialOption *DialOption, sendExchangeOption, recvExchangeOption *ExchangeOption, recvQueueOption *QueueOption, recvReplyRoutingKey string) (*TopicProducer, error) {
	t := &TopicProducer{
		sendExchangeOption: sendExchangeOption,
	}

	var err error

	t.sendSession, err = topicSession(dialOption, sendExchangeOption)
	if err != nil {
		return nil, err
	}

	t.recvReplySession, err = topicSession(dialOption, recvExchangeOption)
	if err != nil {
		return nil, err
	}

	t.recvReplyQueue, err = t.recvReplySession.Channel.QueueDeclare(recvQueueOption.Name, recvQueueOption.Durable, recvQueueOption.AutoDelete, recvQueueOption.Exclusive, recvQueueOption.NoWait, nil)
	if err != nil {
		return nil, err
	}

	err = t.recvReplySession.Channel.QueueBind(t.recvReplyQueue.Name, recvReplyRoutingKey, recvExchangeOption.Name, recvQueueOption.NoWait, nil)
	if err != nil {
		return nil, fmt.Errorf("bind receive queue error, %w", err)
	}

	err = t.recvReply()
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TopicProducer) Close() error {
	switch {
	case t.sendSession != nil:
		err := t.sendSession.Channel.Close()
		if err != nil {
			return err
		}

		fallthrough
	case t.recvReplySession != nil:
		err := t.recvReplySession.Channel.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TopicProducer) SendReport(ctx context.Context, routingKey string, mandatory, immediate bool, msg []byte) error {
	return t.sendSession.Channel.PublishWithContext(
		ctx,
		t.sendExchangeOption.Name,
		routingKey,
		mandatory,
		immediate,
		amqp.Publishing{
			Body: msg,
		},
	)
}

func (t *TopicProducer) SendRequest(ctx context.Context, routingKey string, mandatory, immediate bool, msg []byte) ([]byte, error) {
	_correlationId := newCorrelationId()
	c := make(chan amqp.Delivery)
	t.recvReplyChan.Store(_correlationId, c)

	err := t.sendSession.Channel.PublishWithContext(
		ctx,
		t.sendExchangeOption.Name,
		routingKey,
		mandatory,
		immediate,
		amqp.Publishing{
			Body:          msg,
			ReplyTo:       t.recvReplyQueue.Name,
			CorrelationId: _correlationId,
		},
	)
	if err != nil {
		return nil, err
	}

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("get reply error, %w", ctx.Err())
		case delivery := <-c:
			if delivery.CorrelationId != _correlationId {
				return nil, fmt.Errorf("get wrong reply, want: %s, got: %s", _correlationId, delivery.CorrelationId)
			}

			return delivery.Body, nil
		}
	}
}

func (t *TopicProducer) recvReply() error {
	replies, err := t.recvReplySession.Channel.Consume(
		t.recvReplyQueue.Name,
		"",
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("consume reply error, %w", err)
	}

	go func() {
		for i := range replies {
			go func(delivery amqp.Delivery) {
				_value, ok := t.recvReplyChan.Load(delivery.CorrelationId)
				if ok {
					_c := _value.(chan amqp.Delivery)
					_c <- delivery
				}
			}(i)
		}
	}()

	return nil
}

func NewTopicConsumer(dialOption *DialOption, recvExchangeOption, replyExchangeOption *ExchangeOption, consumerOption *ConsumerOption, errChan chan error) (*TopicConsumer, error) {
	t := &TopicConsumer{
		recvExchangeOption:  recvExchangeOption,
		replyExchangeOption: replyExchangeOption,
		consumerOption:      consumerOption,
		errorChan:           errChan,
	}

	var err error
	t.recvSession, err = topicSession(dialOption, recvExchangeOption)
	if err != nil {
		return nil, err
	}

	t.sendReplySession, err = topicSession(dialOption, recvExchangeOption)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *TopicConsumer) Close() error {
	switch {
	case t.recvSession != nil:
		err := t.recvSession.Channel.Close()
		if err != nil {
			return err
		}

		fallthrough
	case t.sendReplySession != nil:
		err := t.sendReplySession.Channel.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TopicConsumer) Recv(queueOption *QueueOption, routingKey string) (chan *Receive, error) {
	_q, err := queueDeclare(t.recvSession.Channel, queueOption)
	if err != nil {
		return nil, fmt.Errorf("declare receive queue error, %w", err)
	}

	err = t.recvSession.Channel.QueueBind(_q.Name, routingKey, t.recvExchangeOption.Name, false, nil)
	if err != nil {
		return nil, fmt.Errorf("bind receive queue error %w", err)
	}

	re := make(chan *Receive, 8)

	_msgs, err := consumeDeclare(t.recvSession.Channel, _q.Name, t.consumerOption)
	if err != nil {
		return nil, err
	}

	go func() {
		for i := range _msgs {
			if !t.consumerOption.AutoAck {
				if err := i.Ack(false); err != nil {
					t.errorChan <- err

					continue
				}
			}

			re <- &Receive{
				Body:          i.Body,
				CorrelationId: i.CorrelationId,
				ReplyTo:       i.ReplyTo,
			}
		}
	}()

	return re, nil
}

func (t *TopicConsumer) SendReply(ctx context.Context, received *Receive, mandatory, immediate bool, msg []byte) error {
	return t.sendReplySession.Channel.PublishWithContext(
		ctx,
		t.replyExchangeOption.Name,
		received.ReplyTo,
		mandatory,
		immediate,
		amqp.Publishing{
			CorrelationId: received.CorrelationId,
			Body:          msg,
		},
	)
}

func topicSession(dialOption *DialOption, exchangeOption *ExchangeOption) (*Session, error) {
	if exchangeOption.Kind != TopicKind {
		return nil, fmt.Errorf("wrong exchange kind: %s, should be topic", exchangeOption.Kind)
	}

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

	err = exchangeDeclare(ch, exchangeOption)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq/simple cannot create channel: %v", err)
	}

	s := &Session{
		Connection: conn,
		Channel:    ch,
	}

	return s, nil
}
