package measurement

import (
	"context"
	"errors"
	"time"

	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
)

var ErrNotFound = errors.New("measurement with given ID not found")

type Service interface {
	GetMeasurement(ctx context.Context, id string) (sensor.Data, error)
	ListMeasurements(ctx context.Context, name string, from, to time.Time, limit int) ([]sensor.Data, error)
}
