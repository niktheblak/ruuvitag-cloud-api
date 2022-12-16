//go:build !aws

package aws

import (
	"context"
	"time"

	"github.com/niktheblak/ruuvitag-cloud-api/pkg/errs"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/sensor"
)

type DummyService struct {
}

func New(table string) (*DummyService, error) {
	return &DummyService{}, nil
}

func (s *DummyService) GetMeasurement(ctx context.Context, name string, ts time.Time) (sensor.Data, error) {
	return sensor.Data{}, errs.ErrNotImplemented
}

func (s *DummyService) ListMeasurements(ctx context.Context, name string, from, to time.Time, limit int) ([]sensor.Data, error) {
	return nil, errs.ErrNotImplemented
}

func (s *DummyService) Write(ctx context.Context, sd sensor.Data) error {
	return errs.ErrNotImplemented
}
