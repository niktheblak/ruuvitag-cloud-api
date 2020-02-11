package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/measurement"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/server"
	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
)

type Query struct {
	Name  string `json:"name"`
	From  string `json:"from"`
	To    string `json:"to"`
	Limit int    `json:"limit"`
}

var (
	dyndb *dynamodb.DynamoDB
	svc   measurement.Service
)

func init() {
	log.Println("Creating session")
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	dyndb = dynamodb.New(sess)
	table := os.Getenv("TABLE")
	if table == "" {
		table = "ruuvitag"
	}
	svc = NewService(dyndb, table)
}

func HandleRequest(ctx context.Context, q Query) ([]sensor.Data, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("no context")
	}
	if lc.Identity.CognitoIdentityID == "" {
		return nil, fmt.Errorf("unauthorized")
	}
	if q.Name == "" {
		return nil, fmt.Errorf("name must be specified")
	}
	from, to, err := server.ParseTimeRange(q.From, q.To)
	if err != nil {
		return nil, err
	}
	if q.Limit <= 0 {
		q.Limit = 20
	}
	return svc.ListMeasurements(ctx, q.Name, from, to, q.Limit)
}

func main() {
	log.Println("Starting service")
	lambda.Start(HandleRequest)
}
