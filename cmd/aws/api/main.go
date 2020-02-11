package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/measurement"
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
		table = "ruuvitag"
	}
	svc = NewService(dyndb, table)
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, fmt.Errorf("no context")
	}
	if lc.Identity.CognitoIdentityID == "" {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusForbidden}, nil
	}
	name := request.PathParameters["name"]
	if name == "" {
		return events.APIGatewayProxyResponse{Body: "Name must be specified", StatusCode: http.StatusBadRequest}, nil
	}
	from, to, err := server.ParseTimeRange(request.QueryStringParameters["from"], request.QueryStringParameters["to"])
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Invalid time range", StatusCode: http.StatusBadRequest}, nil
	}
	limit := server.ParseLimit(request.QueryStringParameters["limit"])
	measurements, err := svc.ListMeasurements(ctx, name, from, to, limit)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "Failed to query measurements", StatusCode: http.StatusInternalServerError}, err
	}
	body, _ := json.Marshal(measurements)
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: http.StatusOK}, nil
}

func main() {
	log.Println("Starting service")
	lambda.Start(HandleRequest)
}
