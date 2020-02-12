package main

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/measurement"
	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
)

type key struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"ts"`
}

type Service struct {
	client *dynamodb.DynamoDB
	table  string
}

func NewService(client *dynamodb.DynamoDB, table string) measurement.Service {
	return &Service{
		client: client,
		table:  table,
	}
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
		err = measurement.ErrNotFound
		return
	}
	err = dynamodbattribute.UnmarshalMap(res.Item, &sd)
	return
}

func (s *Service) ListMeasurements(ctx context.Context, name string, from, to time.Time, limit int) (measurements []sensor.Data, err error) {
	res, err := s.client.QueryWithContext(ctx, &dynamodb.QueryInput{
		ExpressionAttributeNames: map[string]*string{
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
		FilterExpression:       aws.String("#name = :name"),
		KeyConditionExpression: aws.String("ts BETWEEN :from AND :to"),
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
