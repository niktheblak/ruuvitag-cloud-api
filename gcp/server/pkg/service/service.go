package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/sensor"
)

type Service struct {
	client  *bigquery.Client
	dataset string
	table   string
}

func New(projectID, dataset, table string) (*Service, error) {
	client, err := bigquery.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, err
	}
	return &Service{
		client:  client,
		dataset: dataset,
		table:   table,
	}, nil
}

func (s *Service) GetMeasurement(ctx context.Context, name string, ts time.Time) (sd sensor.Data, err error) {
	// TODO
	q := s.client.Query("")
	q.AllowLargeResults = false
	it, err := q.Read(ctx)
	if err != nil {
		return
	}
	err = it.Next(&sd)
	if err != nil {
		return
	}
	return
}

func (s *Service) ListMeasurements(ctx context.Context, name string, from, to time.Time, limit int) (measurements []sensor.Data, err error) {
	// TODO
	q := s.client.Query("")
	q.AllowLargeResults = false
	it, err := q.Read(ctx)
	if err != nil {
		return
	}
	for {
		var sd sensor.Data
		err = it.Next(&sd)
		if err != nil {
			break
		}
		measurements = append(measurements, sd)
	}
	if errors.Is(err, iterator.Done) {
		err = nil
	}
	return
}

func (s *Service) Write(ctx context.Context, sd sensor.Data) error {
	return fmt.Errorf("not implemented")
}

func (s *Service) Close() error {
	return s.client.Close()
}
