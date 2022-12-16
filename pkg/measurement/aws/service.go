//go:build aws

package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/niktheblak/ruuvitag-cloud-api/pkg/errs"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/sensor"
)

type key struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"ts"`
}

type Service struct {
	client *dynamodb.DynamoDB
	table  string
}

func New(table string) (*Service, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	client := dynamodb.New(sess)
	return &Service{
		client: client,
		table:  table,
	}, nil
}

func (s *Service) GetMeasurement(ctx context.Context, name string, ts time.Time) (sd sensor.Data, err error) {
	k := key{
		Name:      name,
		Timestamp: ts,
	}
	kattrs, err := dynamodbattribute.MarshalMap(k)
	if err != nil {
		return
	}
	res, err := s.client.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.table),
		Key:       kattrs,
	})
	if err != nil {
		return
	}
	if res == nil {
		err = errs.ErrNotFound
		return
	}
	err = dynamodbattribute.UnmarshalMap(res.Item, &sd)
	return
}

func (s *Service) ListMeasurements(ctx context.Context, name string, from, to time.Time, limit int) (measurements []sensor.Data, err error) {
	res, err := s.client.QueryWithContext(ctx, &dynamodb.QueryInput{
		ExpressionAttributeNames: map[string]*string{
			"#ts":   aws.String("ts"),
			"#name": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(name),
			},
			":from": {
				S: aws.String(from.Format(time.RFC3339)),
			},
			":to": {
				S: aws.String(to.Format(time.RFC3339)),
			},
		},
		KeyConditionExpression: aws.String("#name = :name AND #ts BETWEEN :from AND :to"),
		Limit:                  aws.Int64(int64(limit)),
		ProjectionExpression:   nil,
		ScanIndexForward:       nil,
		Select:                 nil,
		TableName:              aws.String(s.table),
	})
	if err != nil {
		return
	}
	for _, v := range res.Items {
		var sd sensor.Data
		err = dynamodbattribute.UnmarshalMap(v, &sd)
		if err != nil {
			return
		}
		measurements = append(measurements, sd)
	}
	return
}

func (s *Service) Write(ctx context.Context, sd sensor.Data) error {
	item, err := dynamodbattribute.MarshalMap(sd)
	if err != nil {
		return err
	}
	_, err = s.client.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(s.table),
	})
	return err
}
