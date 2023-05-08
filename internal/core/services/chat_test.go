package services

import (
	"reflect"
	"testing"
	"time"

	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/core/ports"
)

func TestNewChatService(t *testing.T) {
	type args struct {
		r ports.Repository
		n ports.Notifier
		m ports.Moderator
	}
	tests := []struct {
		name string
		args args
		want *ChatService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewChatService(tt.args.r, tt.args.n, tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChatService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatService_Send(t *testing.T) {
	type fields struct {
		repository ports.Repository
		notifier   ports.Notifier
		moderator  ports.Moderator
	}
	type args struct {
		m domain.Message
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
			c := &ChatService{
				repository: tt.fields.repository,
				notifier:   tt.fields.notifier,
				moderator:  tt.fields.moderator,
			}
			if err := c.Send(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("ChatService.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChatService_History(t *testing.T) {
	type fields struct {
		repository ports.Repository
		notifier   ports.Notifier
		moderator  ports.Moderator
	}
	type args struct {
		dst   string
		since time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*domain.ModeratedMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ChatService{
				repository: tt.fields.repository,
				notifier:   tt.fields.notifier,
				moderator:  tt.fields.moderator,
			}
			got, err := c.History(tt.args.dst, tt.args.since)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatService.History() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChatService.History() = %v, want %v", got, tt.want)
			}
		})
	}
}
