package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/niktheblak/ruuvitag-cloud-api/aws/pkg/service"
	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/sensor"
)

var (
	writer *service.Service
)

func init() {
	table := os.Getenv("TABLE")
	if table == "" {
		log.Fatal("$TABLE must be specified")
	}
	slog.Info("Creating session")
	var err error
	writer, err = service.New(table)
	if err != nil {
		log.Fatalln(err)
	}
}

func HandleRequest(ctx context.Context, sd sensor.Data) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "Received measurement", slog.String("addr", sd.Addr), slog.Any("measurement", sd))
	return writer.Write(ctx, sd)
}

func main() {
	slog.Info("Starting writer")
	lambda.Start(HandleRequest)
}
