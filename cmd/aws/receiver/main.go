package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement/aws"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/sensor"
)

var (
	writer measurement.Service
)

func init() {
	table := os.Getenv("TABLE")
	if table == "" {
		log.Fatal("$TABLE must be specified")
	}
	log.Println("Creating session")
	var err error
	writer, err = aws.New(table)
	if err != nil {
		log.Fatalln(err)
	}
}

func HandleRequest(ctx context.Context, sd sensor.Data) error {
	slog.Info("Received measurement", "addr", sd.Addr, "measurement", sd)
	return writer.Write(ctx, sd)
}

func main() {
	slog.Info("Starting writer")
	lambda.Start(HandleRequest)
}
