package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/measurement"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/measurement/aws"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/server"
)

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
		log.Fatal("Environment variable TABLE must be specified")
	}
	svc = aws.NewService(dyndb, table)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := request.PathParameters["name"]
	if name == "" {
		return badRequest("Name must be specified"), nil
	}
	from, to, err := server.ParseTimeRange(request.QueryStringParameters["from"], request.QueryStringParameters["to"])
	if err != nil {
		return badRequest("Invalid time range"), nil
	}
	limit := server.ParseLimit(request.QueryStringParameters["limit"])
	measurements, err := svc.ListMeasurements(ctx, name, from, to, limit)
	if err != nil {
		return internalServerError("Failed to query measurements"), err
	}
	body, _ := json.Marshal(measurements)
	return okResponse(string(body)), nil
}

func main() {
	log.Println("Starting service")
	lambda.Start(HandleRequest)
}
