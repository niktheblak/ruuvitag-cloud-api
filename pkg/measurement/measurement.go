package measurement

import "time"

const Kind = "Measurement"

type Measurement struct {
	Name        string    `json:"name" datastore:"name"`
	MAC         string    `json:"mac" datastore:"mac"`
	Timestamp   time.Time `json:"ts" datastore:"ts"`
	Temperature float64   `json:"temperature" datastore:"temperature,noindex"`
	Humidity    float64   `json:"humidity" datastore:"humidity,noindex"`
	Pressure    float64   `json:"pressure" datastore:"pressure,noindex"`
	ID          int64     `json:"id" datastore:"-"`
}
