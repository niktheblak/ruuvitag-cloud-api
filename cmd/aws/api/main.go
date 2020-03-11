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
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/api"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement/aws"
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
		return BadRequest("Name must be specified"), nil
	}
	from, to, err := api.ParseTimeRange(request.QueryStringParameters["from"], request.QueryStringParameters["to"])
	if err != nil {
		return BadRequest("Invalid time range"), nil
	}
	limit := api.ParseLimit(request.QueryStringParameters["limit"])
	measurements, err := svc.ListMeasurements(ctx, name, from, to, limit)
	if err != nil {
		return InternalServerError("Failed to query measurements"), err
	}
	body, _ := json.Marshal(measurements)
	return OKResponse(string(body)), nil
}

func main() {
	log.Println("Starting service")
	lambda.Start(HandleRequest)
}
