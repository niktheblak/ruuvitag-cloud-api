package receiver

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
)

const (
	TimeFormat = "2006-01-02T15:04:05.999999"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type JSONMeasurement struct {
	Name        string  `json:"name"`
	MAC         string  `json:"mac"`
	Timestamp   string  `json:"ts"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Pressure    float64 `json:"pressure"`
}

func (jm *JSONMeasurement) ToMeasurement() (*measurement.Measurement, error) {
	ts, err := time.Parse(TimeFormat, jm.Timestamp)
	if err != nil {
		return nil, err
	}
	ts = ts.UTC()
	return &measurement.Measurement{
		Name:        jm.Name,
		MAC:         jm.MAC,
		Timestamp:   ts,
		Temperature: jm.Temperature,
		Humidity:    jm.Humidity,
		Pressure:    jm.Pressure,
	}, nil
}

func ReceiveMeasurement(ctx context.Context, msg PubSubMessage) error {
	if len(msg.Data) == 0 {
		return errors.New("Message does not contain payload")
	}
	var jm JSONMeasurement
	err := json.Unmarshal(msg.Data, &jm)
	if err != nil {
		log.Printf("Failed to parse measurement JSON: %v. Payload: %s", err, string(msg.Data))
		return err
	}
	dsClient, err := datastore.NewClient(ctx, "")
	if err != nil {
		log.Printf("Failed to create datastore client: %v", err)
		return err
	}
	defer dsClient.Close()
	key := datastore.IncompleteKey(measurement.Kind, nil)
	m, err := jm.ToMeasurement()
	if err != nil {
		log.Printf("Failed to convert measurement to stored entity: %v", err)
		return err
	}
	_, err = dsClient.Put(ctx, key, m)
	if err != nil {
		log.Printf("Failed to store measurement: %v", err)
		return err
	}
	return nil
}
