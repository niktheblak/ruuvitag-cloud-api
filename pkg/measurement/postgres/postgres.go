package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
)

type postgresWriter struct {
	db         *sql.DB
	insertStmt *sql.Stmt
}

func New(ctx context.Context, connStr, table string) (measurement.Writer, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	insertStmt, err := db.PrepareContext(ctx, fmt.Sprintf(`
INSERT INTO %s (
  mac,
  name,
  ts,
  temperature,
  humidity,
  pressure,
  acceleration_x,
  acceleration_y,
  acceleration_z,
  movement_counter,
  battery
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, table))
	if err != nil {
		return nil, err
	}
	return &postgresWriter{
		db:         db,
		insertStmt: insertStmt,
	}, nil
}

func (p *postgresWriter) Write(ctx context.Context, data sensor.Data) error {
	_, err := p.insertStmt.ExecContext(ctx, data.Addr, data.Name, data.Timestamp, data.Temperature, data.Humidity, data.Pressure, data.AccelerationX, data.AccelerationY, data.AccelerationZ, data.MovementCounter, data.Battery)
	return err
}

func (p *postgresWriter) Close() error {
	p.insertStmt.Close()
	return p.db.Close()
}
