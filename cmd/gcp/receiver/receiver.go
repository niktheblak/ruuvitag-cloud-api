package receiver

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	fs "github.com/niktheblak/ruuvitag-cloud-api/cmd/gcp/firestore"
)

const (
	TimeFormat = "2006-01-02T15:04:05.999999"
	Collection = "measurements"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type measurementPayload struct {
	Name          string  `json:"name"`
	MAC           string  `json:"mac"`
	Timestamp     string  `json:"ts"`
	Temperature   float64 `json:"temperature"`
	Humidity      float64 `json:"humidity"`
	Pressure      float64 `json:"pressure"`
	Battery       int     `json:"battery"`
	AccelerationX int     `json:"acceleration_x"`
	AccelerationY int     `json:"acceleration_y"`
	AccelerationZ int     `json:"acceleration_z"`
}

func (jm *measurementPayload) ToFirestoreMeasurement() (fs.Measurement, error) {
	ts, err := time.Parse(time.RFC3339Nano, jm.Timestamp)
	if err != nil {
		ts, err = time.Parse(TimeFormat, jm.Timestamp)
	}
	if err != nil {
		return fs.Measurement{}, err
	}
	ts = ts.UTC()
	return fs.Measurement{
		Name:          jm.Name,
		MAC:           jm.MAC,
		Timestamp:     ts,
		Temperature:   jm.Temperature,
		Humidity:      jm.Humidity,
		Pressure:      jm.Pressure,
		Battery:       jm.Battery,
		AccelerationX: jm.AccelerationX,
		AccelerationY: jm.AccelerationY,
		AccelerationZ: jm.AccelerationZ,
	}, nil
}

func ReceiveMeasurement(ctx context.Context, msg PubSubMessage) error {
	if len(msg.Data) == 0 {
		return errors.New("message does not contain payload")
	}
	var mp measurementPayload
	err := json.Unmarshal(msg.Data, &mp)
	if err != nil {
		log.Printf("Failed to parse measurement JSON: %v. Payload: %s", err, string(msg.Data))
		return err
	}
	client, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		log.Fatalf("Failed to create datastore client: %v", err)
	}
	defer client.Close()
	m, err := mp.ToFirestoreMeasurement()
	if err != nil {
		log.Printf("Failed to convert measurement to stored entity: %v", err)
		return err
	}
	_, _, err = client.Collection(Collection).Add(ctx, m)
	if err != nil {
		log.Printf("Failed to store measurement: %v", err)
		return err
	}
	return nil
}
