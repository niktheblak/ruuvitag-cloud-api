package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
)

type postgresService struct {
	db         *sql.DB
	getStmt    *sql.Stmt
	listStmt   *sql.Stmt
	insertStmt *sql.Stmt
}

func New(ctx context.Context, connStr, table string) (measurement.WriterService, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	getStmt, err := db.PrepareContext(ctx, fmt.Sprintf(`
		SELECT
  			mac,
  			temperature,
  			humidity,
  			pressure,
  			acceleration_x,
  			acceleration_y,
  			acceleration_z,
  			movement_counter,
  			battery
		FROM %s
		WHERE name = $1 AND ts = $2
		LIMIT 1
	`, table))
	if err != nil {
		return nil, err
	}
	listStmt, err := db.PrepareContext(ctx, fmt.Sprintf(`
		SELECT
  			mac,
  			ts,
  			temperature,
  			humidity,
  			pressure,
  			acceleration_x,
  			acceleration_y,
  			acceleration_z,
  			movement_counter,
  			battery
		FROM %s
		WHERE name = $1 AND ts >= $2 AND ts <= $3
		ORDER BY ts DESC
		LIMIT $4
	`, table))
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
	return &postgresService{
		db:         db,
		getStmt:    getStmt,
		listStmt:   listStmt,
		insertStmt: insertStmt,
	}, nil
}

func (p *postgresService) GetMeasurement(ctx context.Context, name string, ts time.Time) (sd sensor.Data, err error) {
	row := p.getStmt.QueryRowContext(ctx, name, ts)
	err = row.Scan(&sd.Addr, &sd.Temperature, &sd.Humidity, &sd.Pressure, &sd.AccelerationX, &sd.AccelerationY, &sd.AccelerationZ, &sd.MovementCounter, &sd.Battery)
	if err == sql.ErrNoRows {
		err = measurement.ErrNotFound
	}
	sd.Name = name
	sd.Timestamp = ts
	return
}

func (p *postgresService) ListMeasurements(ctx context.Context, name string, from, to time.Time, limit int) ([]sensor.Data, error) {
	rows, err := p.listStmt.QueryContext(ctx, name, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var measurements []sensor.Data
	for rows.Next() {
		var sd sensor.Data
		err = rows.Scan(&sd.Addr, &sd.Timestamp, &sd.Temperature, &sd.Humidity, &sd.Pressure, &sd.AccelerationX, &sd.AccelerationY, &sd.AccelerationZ, &sd.MovementCounter, &sd.Battery)
		if err != nil {
			return nil, err
		}
		sd.Name = name
		measurements = append(measurements, sd)
	}
	return measurements, rows.Err()
}

func (p *postgresService) Write(ctx context.Context, data sensor.Data) error {
	_, err := p.insertStmt.ExecContext(ctx, data.Addr, data.Name, data.Timestamp, data.Temperature, data.Humidity, data.Pressure, data.AccelerationX, data.AccelerationY, data.AccelerationZ, data.MovementCounter, data.Battery)
	return err
}

func (p *postgresService) Close() error {
	p.getStmt.Close()
	p.listStmt.Close()
	p.insertStmt.Close()
	return p.db.Close()
}
