package notifier

import (
	"crypto/tls"
	"fmt"

	"github.com/posilva/simplechat/internal/core/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQNotifier struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

// NewRabbitMQNotifierWithLocal creates a new instance for Local connection to RMQ
func NewRabbitMQNotifierWithLocal(url string) (*RabbitMQNotifier, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQNotifier{
		connection: conn,
		channel:    ch,
	}, nil
}

// NewRabbitMQNotifierWithTLS creates a new instance for RMQ connection using TLS
func NewRabbitMQNotifierWithTLS(url string, tls *tls.Config) (*RabbitMQNotifier, error) {
	conn, err := amqp.DialTLS(url, tls)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQNotifier{
		connection: conn,
		channel:    ch,
	}, nil
}

// Broadcast message
func (n *RabbitMQNotifier) Broadcast(m domain.ModeratedMessage) error {
	dest := m.To
	topic := fmt.Sprintf("chat_ep_%s", dest)
	_ = topic
	return nil
}

// Register
func (n *RabbitMQNotifier) Register(id string, topic string) error {
	return nil
}

// DeRegister
func (n *RabbitMQNotifier) DeRegister(id string) error {
	return nil
}
