package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement/aws"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/sensor"
)

var (
	writer measurement.Writer
)

func init() {
	table := os.Getenv("TABLE")
	if table == "" {
		log.Fatal("Environment variable TABLE must be specified")
	}
	log.Println("Creating session")
	writer = aws.NewDynamoDBWriter(table)
}

func HandleRequest(ctx context.Context, sd sensor.Data) error {
	log.Printf("Received measurement from %v", sd.Addr)
	return writer.Write(ctx, sd)
}

func main() {
	log.Println("Starting writer")
	lambda.Start(HandleRequest)
}
