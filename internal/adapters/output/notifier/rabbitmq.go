// Package notifier handles notification implementations
package notifier

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/core/ports"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	topicPrefix string = "grp_"
)

type subscriptionInfo struct {
	queueName string
}

// RabbitMQNotifier implements RabbitMQ based notifications
type RabbitMQNotifier[T ports.NotifierCodec] struct {
	codec         T
	registry      ports.Registry
	subscriptions map[string]subscriptionInfo
	conn          *amqp.Connection
	ch            *amqp.Channel
	mu            sync.Mutex
}

// NewRabbitMQNotifierWithLocal creates a new instance for Local connection to RMQ
func NewRabbitMQNotifierWithLocal[T ports.NotifierCodec](url string, r ports.Registry) (*RabbitMQNotifier[T], error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQNotifier[T]{
		registry:      r,
		mu:            sync.Mutex{},
		conn:          conn,
		ch:            ch,
		subscriptions: map[string]subscriptionInfo{},
	}, nil
}

// NewRabbitMQNotifierWithTLS creates a new instance for RMQ connection using TLS
func NewRabbitMQNotifierWithTLS[T ports.NotifierCodec](url string, tls *tls.Config, r ports.Registry) (*RabbitMQNotifier[T], error) {
	conn, err := amqp.DialTLS(url, tls)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQNotifier[T]{
		registry:      r,
		mu:            sync.Mutex{},
		conn:          conn,
		ch:            ch,
		subscriptions: map[string]subscriptionInfo{},
	}, nil
}

// Broadcast message
func (n *RabbitMQNotifier[T]) Broadcast(m domain.Notication) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	t := internalTopic(m.To)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if _, ok := n.subscriptions[t]; ok {
		body, err := n.codec.Encode(m)

		if err != nil {
			return fmt.Errorf("failed to parse json: %s", err)
		}
		// TODO: PMS: check options later
		err = n.ch.PublishWithContext(ctx,
			t, // exchange
			"",
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		if err != nil {
			return fmt.Errorf("failed to publish to topici '%s': '%s", m.To, err)
		}
	}

	return nil
}

// Subscribe to notifications
func (n *RabbitMQNotifier[T]) Subscribe(ep ports.Endpoint) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	room := ep.Room()

	err := n.registry.Register(ep)
	if err != nil {
		return fmt.Errorf("failed to register endpoit: %v", err)
	}
	t := internalTopic(room)

	if _, ok := n.subscriptions[t]; !ok {
		queueName, err := n.initTopic(t)
		if err != nil {
			return fmt.Errorf("failed to init topic '%s': %s", room, err)
		}

		n.subscriptions[t] = subscriptionInfo{
			queueName: queueName,
		}

		err = n.createSubscription(queueName, ep)
		if err != nil {
			return fmt.Errorf("failed to failed to create subscription to topic '%s': %s", room, err)
		}
	}
	return nil
}

// Unsubscribe unsubscribes the endpoint to receive notificatoins
func (n *RabbitMQNotifier[T]) Unsubscribe(ep ports.Endpoint) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.registry.DeRegister(ep)
}

func (n *RabbitMQNotifier[T]) createSubscription(queue string, r ports.Receiver) error {

	msgs, err := n.ch.Consume(
		queue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if rr := recover(); rr != nil {
				r.Recover()
			}
		}()

		// TODO: terminate this using a channel otherwise it will be a go routine leak
		for d := range msgs {

			m := domain.Notication{}

			err := n.codec.Decode(d.Body, &m)
			if err != nil {
				// TODO add logger
				continue
			}
			n.registry.Notify(m)
		}
	}()
	return nil
}
func (n *RabbitMQNotifier[T]) initTopic(t string) (string, error) {
	// create the exchange
	err := n.ch.ExchangeDeclare(
		t,        // name
		"fanout", // type
		false,    // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return "", err
	}

	q, err := n.ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return "", err
	}

	err = n.ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		t,      // exchange
		false,
		nil)
	if err != nil {
		return "", err
	}
	return q.Name, nil
}

func internalTopic(dst string) string {
	return fmt.Sprintf("%s%s", topicPrefix, dst)
}
