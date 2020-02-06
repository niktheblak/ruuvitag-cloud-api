package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
)

const table = "ruuvitag"

var dyndb *dynamodb.DynamoDB

func init() {
	log.Println("Creating session")
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	dyndb = dynamodb.New(sess)
}

func HandleRequest(ctx context.Context, sd sensor.Data) error {
	log.Printf("Received measurement from %v", sd.Addr)
	item, err := dynamodbattribute.MarshalMap(sd)
	if err != nil {
		return err
	}
	_, err = dyndb.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(table),
	})
	return err
}

func main() {
	log.Println("Starting receiver")
	lambda.Start(HandleRequest)
}
