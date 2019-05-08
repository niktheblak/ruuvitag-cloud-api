package service

import (
	"context"
	"errors"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrNotFound = errors.New("measurement with given ID not found")

type Service struct {
	ctx    context.Context
	client *firestore.Client
}

func NewService(ctx context.Context, client *firestore.Client) *Service {
	return &Service{
		ctx:    ctx,
		client: client,
	}
}

func (s *Service) GetMeasurement(id string) (m measurement.Measurement, err error) {
	r, err := s.client.Collection(measurement.Collection).Doc(id).Get(s.ctx)
	if status.Code(err) == codes.NotFound {
		err = ErrNotFound
		return
	}
	if err != nil {
		return
	}
	err = r.DataTo(&m)
	if err != nil {
		return
	}
	m.ID = id
	return
}

func (s *Service) ListMeasurements(name string, from, to time.Time, limit int) (measurements []measurement.Measurement, err error) {
	coll := s.client.Collection(measurement.Collection)
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
	docs := query.Documents(s.ctx)
	defer docs.Stop()
	var doc *firestore.DocumentSnapshot
	for doc, err = docs.Next(); err == nil; doc, err = docs.Next() {
		var m measurement.Measurement
		err = doc.DataTo(&m)
		if err != nil {
			return
		}
		m.ID = doc.Ref.ID
		measurements = append(measurements, m)
	}
	if err == iterator.Done {
		err = nil
	}
	return
}
