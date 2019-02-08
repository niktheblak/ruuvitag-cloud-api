package service

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
)

type Service struct {
	ctx    context.Context
	client *datastore.Client
}

func NewService(ctx context.Context, client *datastore.Client) *Service {
	return &Service{
		ctx:    ctx,
		client: client,
	}
}

func (s *Service) GetMeasurement(id int64) (m measurement.Measurement, err error) {
	key := datastore.IDKey(measurement.Kind, id, nil)
	err = s.client.Get(s.ctx, key, &m)
	m.ID = id
	return
}

func (s *Service) ListMeasurements(name string, from, to time.Time, limit int) (measurements []measurement.Measurement, err error) {
	filters := make(map[string]interface{})
	filters["name ="] = name
	if !from.IsZero() {
		filters["ts >="] = from
	}
	if !to.IsZero() {
		filters["ts <"] = to
		if to.Sub(from) <= time.Hour*24 {
			// Don't limit number of results if the query is for less than one day
			limit = 0
		}
	}
	query := datastore.NewQuery(measurement.Kind)
	for k, v := range filters {
		query = query.Filter(k, v)
	}
	query = query.Order("-ts")
	if limit > 0 {
		query = query.Limit(limit)
	}
	keys, err := s.client.GetAll(s.ctx, query, &measurements)
	if err != nil {
		return
	}
	for i, key := range keys {
		measurements[i].ID = key.ID
	}
	return
}
