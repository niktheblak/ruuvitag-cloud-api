package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
)

type dynamoDBWriter struct {
	sess  *session.Session
	db    *dynamodb.DynamoDB
	table string
}

func NewDynamoDBWriter(table string) measurement.Writer {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	db := dynamodb.New(sess)
	return &dynamoDBWriter{
		sess:  sess,
		db:    db,
		table: table,
	}
}

func (r *dynamoDBWriter) Write(ctx context.Context, sd sensor.Data) error {
	item, err := dynamodbattribute.MarshalMap(sd)
	if err != nil {
		return err
	}
	_, err = r.db.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(r.table),
	})
	return err
}
