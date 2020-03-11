package measurement

import (
	"context"

	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
)

type Writer interface {
	Write(ctx context.Context, sd sensor.Data) error
}
