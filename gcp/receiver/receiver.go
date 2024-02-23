package receiver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"

	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/sensor"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

var (
	projectID string
	bqDataset string
	bqTable   string
)

func init() {
	projectID = os.Getenv("GCP_PROJECT")
	if projectID == "" {
		panic("Environment variable GCP_PROJECT is not set")
	}
	bqDataset = os.Getenv("BIGQUERY_DATASET")
	if bqDataset == "" {
		bqDataset = "ruuvitag"
	}
	bqTable = os.Getenv("BIGQUERY_TABLE")
	if bqTable == "" {
		bqTable = "measurements"
	}
	logger.Info("Initializing receiver", "project_id", projectID, "dataset", bqDataset, "table", bqTable)
	functions.HTTP("ReceiveMeasurement", receive)
}

type measurementSaver struct {
	sensor.Data
	bigquery.ValueSaver
}

func (m *measurementSaver) Save() (row map[string]bigquery.Value, insertID string, err error) {
	if m.Addr == "" {
		err = fmt.Errorf("empty value")
		return
	}
	if m.Timestamp.IsZero() {
		err = fmt.Errorf("empty value")
		return
	}
	row = map[string]bigquery.Value{
		"mac":                m.Addr,
		"name":               m.Name,
		"ts":                 m.Timestamp,
		"temperature":        m.Temperature,
		"humidity":           m.Humidity,
		"pressure":           m.Pressure,
		"dew_point":          m.DewPoint,
		"acceleration_x":     m.AccelerationX,
		"acceleration_y":     m.AccelerationY,
		"acceleration_z":     m.AccelerationZ,
		"battery":            m.BatteryVoltage,
		"tx_power":           m.TxPower,
		"movement_counter":   m.MovementCounter,
		"measurement_number": m.MeasurementNumber,
	}
	return
}

func receive(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	if body != nil {
		defer body.Close()
	}
	dec := json.NewDecoder(body)
	var measurement sensor.Data
	err := dec.Decode(&measurement)
	if err != nil {
		logger.Warn("Invalid measurement", "error", err)
		http.Error(w, "Invalid measurement", http.StatusBadRequest)
		return
	}
	logger.Info("Received measurement", "measurement", measurement)
	ctx := r.Context()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		logger.Error("Could not connect to BigQuery", "error", err)
		http.Error(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}
	defer client.Close()
	table := client.Dataset(bqDataset).Table(bqTable)
	saver := measurementSaver{Data: measurement}
	inserter := table.Inserter()
	inserter.IgnoreUnknownValues = true
	inserter.SkipInvalidRows = true
	if err := inserter.Put(ctx, &saver); err != nil {
		logger.Error("Could not write measurement to database", "error", err)
		http.Error(w, "Could not write measurement to database", http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(measurement); err != nil {
		logger.Error("Could not write output", "error", err)
	}
}
