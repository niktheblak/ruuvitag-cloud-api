package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/sensor"
)

type Service struct {
	db         *sql.DB
	getStmt    *sql.Stmt
	listStmt   *sql.Stmt
	insertStmt *sql.Stmt
}

func New(ctx context.Context, connStr, table string) (*Service, error) {
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
  			battery_voltage
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
  			battery_voltage
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
  			battery_voltage
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, table))
	if err != nil {
		return nil, err
	}
	return &Service{
		db:         db,
		getStmt:    getStmt,
		listStmt:   listStmt,
		insertStmt: insertStmt,
	}, nil
}

func (s *Service) GetMeasurement(ctx context.Context, name string, ts time.Time) (sd sensor.Data, err error) {
	row := s.getStmt.QueryRowContext(ctx, name, ts)
	err = row.Scan(&sd.Addr, &sd.Temperature, &sd.Humidity, &sd.Pressure, &sd.AccelerationX, &sd.AccelerationY, &sd.AccelerationZ, &sd.MovementCounter, &sd.BatteryVoltage)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("not found")
	}
	sd.Name = name
	sd.Timestamp = ts
	return
}

func (s *Service) ListMeasurements(ctx context.Context, name string, from, to time.Time, limit int) ([]sensor.Data, error) {
	rows, err := s.listStmt.QueryContext(ctx, name, from, to, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var measurements []sensor.Data
	for rows.Next() {
		var sd sensor.Data
		err = rows.Scan(&sd.Addr, &sd.Timestamp, &sd.Temperature, &sd.Humidity, &sd.Pressure, &sd.AccelerationX, &sd.AccelerationY, &sd.AccelerationZ, &sd.MovementCounter, &sd.BatteryVoltage)
		if err != nil {
			return nil, err
		}
		sd.Name = name
		measurements = append(measurements, sd)
	}
	return measurements, rows.Err()
}

func (s *Service) Write(ctx context.Context, data sensor.Data) error {
	_, err := s.insertStmt.ExecContext(ctx, data.Addr, data.Name, data.Timestamp, data.Temperature, data.Humidity, data.Pressure, data.AccelerationX, data.AccelerationY, data.AccelerationZ, data.MovementCounter, data.BatteryVoltage)
	return err
}

func (s *Service) Close() error {
	s.getStmt.Close()
	s.listStmt.Close()
	s.insertStmt.Close()
	return s.db.Close()
}
