package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/posilva/simplechat/internal/core/domain"
)

const (
	// formats
	datetimeNoMilliFormat string = "20060102T150405"
	datetimeFormat        string = datetimeNoMilliFormat + ".999999"
	// ddb fields
	hashKeyName string = "pk"
	sortKeyName string = "sk"
	// ddb prefixes
	pkPrefix string = "GRP#"
	skPrefix string = "MSG#"
)

// ChatEntryRecord represents a dynamodb table record
type messageRecord struct {
	PK          string `dynamodbav:"pk" json:"pk"`
	SK          string `dynamodbav:"sk" json:"sk"`
	Source      string `dynamodbav:"src" json:"src"`
	Message     string `dynamodbav:"msg" json:"msg"`
	Filtered    string `dynamodbav:"filtered,omitempty" json:"filtered"`
	FilterLevel string `dynamodbav:"filter_level,omitempty" json:"filter_level"`
}

// DynamoDBRepository implements Repository interface for DynamoDB
type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewDynamoDBRepository creates a new DynamoDB repository
func NewDynamoDBRepository(cfg aws.Config, table string) (*DynamoDBRepository, error) {
	c := dynamodb.NewFromConfig(cfg)
	return &DynamoDBRepository{
		client:    c,
		tableName: table,
	}, nil
}

// Store implementes the Repository interface for DynamoDB
func (r *DynamoDBRepository) Store(m domain.ModeratedMessage) error {
	it := messageRecord{
		PK:      chatPk(m.To),
		SK:      chatSk(),
		Source:  m.From,
		Message: m.Payload,
	}
	if m.Level > 0 {
		it.FilterLevel = strconv.Itoa(int(m.Level))
		it.Filtered = m.FilteredPayload
	}
	item, err := attributevalue.MarshalMap(it)
	if err != nil {
		return fmt.Errorf("failed to marshal item map: %v", err)
	}

	keyEx := expression.AttributeNotExists(expression.Name("pk")).And(expression.AttributeNotExists(expression.Name("sk")))
	expr, err := expression.NewBuilder().WithCondition(keyEx).Build()
	if err != nil {
		return fmt.Errorf("failed to build expression with condition: %v", err)
	}

	input := dynamodb.PutItemInput{
		TableName:                   aws.String(r.tableName),
		Item:                        item,
		ConditionExpression:         expr.KeyCondition(),
		ReturnItemCollectionMetrics: types.ReturnItemCollectionMetricsSize,
		ReturnConsumedCapacity:      types.ReturnConsumedCapacityTotal,
		ReturnValues:                types.ReturnValueAllOld,
	}

	output, err := r.client.PutItem(context.Background(), &input)
	if err != nil {
		return fmt.Errorf("failed to put item: %v", err)
	}
	// TODO: add metrics
	// TODO: maybe return some information here
	_ = output
	return nil
}

// Fetch implementes the Repository interface for DynamoDB
func (r *DynamoDBRepository) Fetch(key string, since time.Duration) ([]*domain.ModeratedMessage, error) {
	now := time.Now().UTC()

	// we increment 10 minutes just to make sure we catch the most recent chats recorded
	latestTs := now.Add(time.Minute * time.Duration(10))
	untilTs := now.Add(-since)

	latest := skPrefix + latestTs.Format(datetimeNoMilliFormat) + ".999999"
	until := skPrefix + untilTs.Format(datetimeNoMilliFormat) + ".000000"

	keyEx := expression.Key(hashKeyName).Equal(expression.Value(chatPk(key)))
	keyEx = keyEx.And(expression.KeyBetween(expression.Key(sortKeyName),
		expression.Value(until),
		expression.Value(latest)),
	)
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()

	if err != nil {
		return nil, fmt.Errorf("failed to build Query condition expression: %v", err)
	}

	input := dynamodb.QueryInput{
		TableName:                 aws.String(r.tableName),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	output, err := r.client.Query(context.Background(), &input)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages delivered to %s: %v", key, err)
	}

	var records []messageRecord

	err = attributevalue.UnmarshalListOfMaps(output.Items, &records)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal list of maps: %v", err)

	}
	var result []*domain.ModeratedMessage
	for _, m := range records {
		result = append(result, toModeratedMessage(m))
	}
	// TODO: add metrics
	// TODO: maybe return some information here
	return result, nil
}

func toModeratedMessage(m messageRecord) *domain.ModeratedMessage {
	to := strings.TrimPrefix(m.PK, pkPrefix)
	id := strings.TrimPrefix(m.SK, skPrefix)

	mm := domain.ModeratedMessage{
		Message: domain.Message{
			Payload: m.Message,
			From:    m.Source,
			To:      to,
		},
		ID: id,
	}

	// TODO: avoid ignoring errors
	v, _ := strconv.Atoi(m.FilterLevel)
	mm.Level = uint(v)
	mm.FilteredPayload = m.Message

	if v > 0 {
		mm.FilteredPayload = m.Filtered
	}

	return &mm
}

func chatPk(d string) string {
	return fmt.Sprintf("%s%s", pkPrefix, d)
}

func chatSk() string {
	return chatSkTs(time.Now().UTC())
}

func chatSkTs(ts time.Time) string {
	return fmt.Sprintf("%s%s", skPrefix, ts.Format(datetimeFormat))
}
