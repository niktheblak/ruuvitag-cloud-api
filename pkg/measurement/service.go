package measurement

import (
	"context"
	"time"

	"github.com/niktheblak/ruuvitag-cloud-api/pkg/sensor"
)

type Service interface {
	GetMeasurement(ctx context.Context, name string, ts time.Time) (sensor.Data, error)
	ListMeasurements(ctx context.Context, name string, from, to time.Time, limit int) ([]sensor.Data, error)
	Write(ctx context.Context, sd sensor.Data) error
}
