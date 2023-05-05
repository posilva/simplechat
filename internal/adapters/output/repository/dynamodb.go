package repository

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/posilva/simplechat/internal/core/domain"
)

type DynamoDBRepository struct {
	client *dynamodb.Client
}

func NewDynamoDBRepository(cfg aws.Config) (*DynamoDBRepository, error) {
	c := dynamodb.NewFromConfig(cfg)
	return &DynamoDBRepository{
		client: c,
	}, nil
}
func (r *DynamoDBRepository) Store(m domain.ModeratedMessage) error {

	return nil
}

func (r *DynamoDBRepository) Fetch(since time.Duration) ([]*domain.ModeratedMessage, error) {
	return nil, nil
}
