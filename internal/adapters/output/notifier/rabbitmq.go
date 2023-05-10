// Package notifier handles notification implementations
package notifier

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/posilva/simplechat/internal/adapters/output/registry"
	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/core/ports"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/ksuid"
)

const (
	topicPrefix string = "grp_"
)

type destinationSet map[string]struct{}

type notificationInfo struct {
	destinations destinationSet
	name         string
}

// RabbitMQNotifier implements RabbitMQ based notifications
type RabbitMQNotifier struct {
	queueName string
	id2Topic  map[string]string
	conn      *amqp.Connection
	ch        *amqp.Channel
	queue     amqp.Queue

	reg      ports.Registry
	mu       sync.Mutex
	registry map[string]notificationInfo
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
		conn:      conn,
		ch:        ch,
		registry:  make(map[string]notificationInfo),
		id2Topic:  make(map[string]string),
		queueName: ksuid.New().String(),
		reg:       registry.NewInMemoryRegistry(),
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
		conn:      conn,
		ch:        ch,
		registry:  make(map[string]notificationInfo),
		id2Topic:  make(map[string]string),
		queueName: ksuid.New().String(),
		reg:       registry.NewInMemoryRegistry(),
	}, nil
}

// Broadcast message
func (n *RabbitMQNotifier) Broadcast(m domain.ModeratedMessage) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	t := internalTopic(m.To)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if _, ok := n.registry[t]; ok {
		body, err := toJSON(m)
		if err != nil {
			return fmt.Errorf("failed to parse json: %s", err)
		}

		// TODO: PMS: check options later
		err = n.ch.PublishWithContext(ctx,
			t,            // exchange
			n.queue.Name, // routing key
			false,        // mandatory
			false,        // immediate
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
func (n *RabbitMQNotifier) Subscribe(ep ports.Endpoint) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	id := ep.ID()
	room := ep.Room()

	t := internalTopic(room)

	if v, ok := n.id2Topic[id]; ok {
		if strings.Compare(v, t) != 0 {
			return fmt.Errorf("id: '%s' already registered to different topic: '%s'", id, room)
		}
	}

	n.id2Topic[id] = t

	// add the id to existing set
	if v, ok := n.registry[t]; ok {
		v.destinations[id] = struct{}{}
		n.registry[t] = v
	} else {
		// create the registry entry
		n.registry[t] = notificationInfo{
			name:         t,
			destinations: newDestinationSet(id),
		}
		err := n.initTopic(t)
		if err != nil {
			return fmt.Errorf("failed to init topic '%s': %s", room, err)
		}

		err = n.subscribe(ep)
		if err != nil {
			return fmt.Errorf("failed to subscribe to topic '%s': %s", room, err)
		}

	}
	return nil
}

// Unsubscribe notifications
func (n *RabbitMQNotifier) Unsubscribe(ep ports.Endpoint) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if t, ok := n.id2Topic[ep.ID()]; ok {
		if v, ok := n.registry[t]; ok {
			delete(v.destinations, ep.ID())
			n.registry[t] = v
		}
	}
	delete(n.id2Topic, ep.ID())
	return nil
}

func (n *RabbitMQNotifier) subscribe(r ports.Receiver) error {
	msgs, err := n.ch.Consume(
		n.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if rr := recover(); rr != nil {
				log.Printf("recovering from panic: %v", rr)
				r.Recover()
			}
		}()

		// TODO: terminate this using a channel otherwise it will be a go routine leak
		for d := range msgs {
			m := domain.ModeratedMessage{}
			err := json.Unmarshal(d.Body, &m)

			if err != nil {
				log.Printf("failed to parse received moderated message: %v", err)
				continue
			}

			for k := range n.id2Topic {
				if k != m.From {
					r.Receive(m)
				}
			}
		}
	}()
	return nil
}
func (n *RabbitMQNotifier) initTopic(t string) error {
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
		return err
	}

	q, err := n.ch.QueueDeclare(
		n.queueName, // name
		false,       // durable
		true,        // delete when unused
		true,        // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	if err != nil {
		return err
	}

	err = n.ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		t,      // exchange
		false,
		nil)
	if err != nil {
		return err
	}
	n.queue = q
	return nil
}

func internalTopic(dst string) string {
	return fmt.Sprintf("%s%s", topicPrefix, dst)
}

func newDestinationSet(init string) destinationSet {
	d := make(map[string]struct{})
	d[init] = struct{}{}
	return d
}

func toJSON(m domain.ModeratedMessage) ([]byte, error) {
	return json.Marshal(m)
}
