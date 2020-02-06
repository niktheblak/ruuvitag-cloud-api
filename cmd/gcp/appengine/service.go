package main

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	fs "github.com/niktheblak/ruuvitag-cloud-api/cmd/gcp/firestore"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/measurement"
	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const Collection = "measurements"

type Service struct {
	client *firestore.Client
}

func NewService(client *firestore.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) GetMeasurement(ctx context.Context, id string) (sensor.Data, error) {
	r, err := s.client.Collection(Collection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		err = measurement.ErrNotFound
	}
	if err != nil {
		return sensor.Data{}, nil
	}
	var fm fs.Measurement
	err = r.DataTo(&fm)
	return toSensorData(fm), nil
}

func (s *Service) ListMeasurements(ctx context.Context, name string, from, to time.Time, limit int) (measurements []sensor.Data, err error) {
	coll := s.client.Collection(Collection)
	query := coll.OrderBy("ts", firestore.Desc).Where("name", "==", name)
	if !from.IsZero() {
		query = query.Where("ts", ">=", from)
	}
	if !to.IsZero() {
		query = query.Where("ts", "<", to)
		if to.Sub(from) <= time.Hour*24 {
			// Don't limit number of results if the query is for less than one day
			limit = 0
		}
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	docs := query.Documents(ctx)
	defer docs.Stop()
	var doc *firestore.DocumentSnapshot
	for doc, err = docs.Next(); err == nil; doc, err = docs.Next() {
		var fm fs.Measurement
		err = doc.DataTo(&fm)
		if err != nil {
			return
		}
		measurements = append(measurements, toSensorData(fm))
	}
	if err == iterator.Done {
		err = nil
	}
	return
}

func toSensorData(fm fs.Measurement) sensor.Data {
	return sensor.Data{
		Name:          fm.Name,
		Addr:          fm.MAC,
		Timestamp:     fm.Timestamp,
		Temperature:   fm.Temperature,
		Humidity:      fm.Humidity,
		Pressure:      fm.Pressure,
		Battery:       fm.Battery,
		AccelerationX: fm.AccelerationX,
		AccelerationY: fm.AccelerationY,
		AccelerationZ: fm.AccelerationZ,
	}
}
