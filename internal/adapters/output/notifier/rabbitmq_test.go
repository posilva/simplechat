package notifier

import (
	"crypto/tls"
	"testing"

	"github.com/posilva/simplechat/internal/core/domain"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

const (
	localURL    string = "amqp://guest:guest@localhost:5672/"
	localURLSSL string = "amqps://guest:guest@localhost:5671/"
)

func TestNewRabbitMQNotifierWithLocal(t *testing.T) {
	got, err := NewRabbitMQNotifierWithLocal(localURL)

	expectedType := &RabbitMQNotifier{}
	assert.NoError(t, err, "expected not an error")
	assert.NotNil(t, got, "expected notifier to be not nil")
	assert.IsType(t, expectedType, got, "expected RabbitMQNotifier ")
	assert.NotNil(t, got.connection, "connection should be non nil")
	assert.NotNil(t, got.channel, "channel should be non nil")
}

func TestNewRabbitMQNotifierWithTLS(t *testing.T) {

	got, err := NewRabbitMQNotifierWithTLS(localURLSSL, &tls.Config{InsecureSkipVerify: true})

	expectedType := &RabbitMQNotifier{}
	assert.NoError(t, err, "expected not an error")
	assert.NotNil(t, got, "expected notifier to be not nil")
	assert.IsType(t, expectedType, got, "expected RabbitMQNotifier ")
	assert.NotNil(t, got.connection, "connection should be non nil")
	assert.NotNil(t, got.channel, "channel should be non nil")

}

func TestRabbitMQNotifier_Broadcast(t *testing.T) {
	type fields struct {
		connection *amqp.Connection
		channel    *amqp.Channel
	}
	type args struct {
		m domain.ModeratedMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &RabbitMQNotifier{
				connection: tt.fields.connection,
				channel:    tt.fields.channel,
			}
			if err := n.Broadcast(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("RabbitMQNotifier.Broadcast() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRabbitMQNotifier_Register(t *testing.T) {
	type fields struct {
		connection *amqp.Connection
		channel    *amqp.Channel
	}
	type args struct {
		id    string
		topic string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &RabbitMQNotifier{
				connection: tt.fields.connection,
				channel:    tt.fields.channel,
			}
			if err := n.Register(tt.args.id, tt.args.topic); (err != nil) != tt.wantErr {
				t.Errorf("RabbitMQNotifier.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRabbitMQNotifier_DeRegister(t *testing.T) {
	type fields struct {
		connection *amqp.Connection
		channel    *amqp.Channel
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &RabbitMQNotifier{
				connection: tt.fields.connection,
				channel:    tt.fields.channel,
			}
			if err := n.DeRegister(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("RabbitMQNotifier.DeRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
