package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/api"

	"github.com/niktheblak/ruuvitag-cloud-api/aws/pkg/service"
)

var svc *service.Service

func init() {
	table := os.Getenv("TABLE")
	if table == "" {
		log.Fatal("Environment variable TABLE must be specified")
	}
	var err error
	svc, err = service.New(table)
	if err != nil {
		log.Fatal(err)
	}
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
	slog.Info("Starting service")
	lambda.Start(HandleRequest)
}
