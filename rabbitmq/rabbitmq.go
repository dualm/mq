package rabbitmq

import (
	"fmt"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	UrlFormat = "amqp://%s:%s@%s:%s"
)

const (
	DirectKind  = "direct"
	TopicKind   = "topic"
	FanoutKind  = "fanout"
	HeadersKind = "headers"
)

type Session struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
}

func (s *Session) Close() error {
	if s.Connection == nil {
		return nil
	}

	return s.Connection.Close()
}

type DialOption struct {
	Username string
	Password string
	Host     string
	Port     string
	VHost    string
}

func (o *DialOption) url() string {
	return fmt.Sprintf(UrlFormat, o.Username, o.Password, o.Host, o.Port)
}

type QueueOption struct {
	Name       string
	AutoDelete bool
	Durable    bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type ExchangeOption struct {
	Name string
	Kind string // 每一个exchange都有一个类型, 为"direct", "fanout", "topic"和"headers"其中之一.
	// Durable && !AutoDelete 在服务器重启后以及在没有任何绑定时都为声明状态. 适用于静态场景.
	// !Durable && AutoDelete 服务器重启后不恢复, 没有任何绑定后会被删除. 减少临时通信对于vhost的污染
	// Durable && AutoDelete 服务器重启后不会被移除, 但当没有任何绑定后会被删除. 用于稳定的临时场景.
	// !Durable && !AutoDelete 服务器重启后移除, 但没有任何绑定时仍会保留. 用于绑定间隔时间较长的临时场景.
	AutoDelete bool
	Durable    bool
	NoWait     bool
	Args       amqp.Table
}

type ConsumerOption struct {
	Name      string
	AutoAck   bool // 接收消息后是否自动发送acknowledge. 如果为true, consumer不需要显式ack回复, Server会在回复消息之前优先发送acknowledge.
	Exclusive bool // 表示是否为consumer独占的queue, true: 独占, false: 共享.
	NoWait    bool // 是否等待服务器确认消息发送请求, 如果不等待则立即开始消息的收发.
	Args      amqp.Table
}

// Message 包括Report和Request
type Message struct {
	amqp.Publishing
}

type Receive struct {
	Body          []byte
	CorrelationId string
	ReplyTo       string
}

func queueDeclare(ch *amqp.Channel, option *QueueOption) (amqp.Queue, error) {
	return ch.QueueDeclare(option.Name, option.Durable, option.AutoDelete, option.Exclusive, option.NoWait, option.Args)
}

func consumeDeclare(ch *amqp.Channel, queueName string, option *ConsumerOption) (<-chan amqp.Delivery, error) {
	return ch.Consume(queueName, option.Name, option.AutoAck, option.Exclusive, false, option.NoWait, option.Args)
}

func exchangeDeclare(ch *amqp.Channel, option *ExchangeOption) error {
	return ch.ExchangeDeclare(option.Name, option.Kind, option.Durable, option.AutoDelete, false, option.NoWait, option.Args)
}

func newCorrelationId() string {
	return uuid.New().String()
}
