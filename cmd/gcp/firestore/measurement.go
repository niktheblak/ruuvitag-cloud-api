package firestore

import (
	"time"
)

type Measurement struct {
	Name          string    `firestore:"name"`
	MAC           string    `firestore:"mac"`
	Timestamp     time.Time `firestore:"ts"`
	Temperature   float64   `firestore:"temperature"`
	Humidity      float64   `firestore:"humidity"`
	Pressure      float64   `firestore:"pressure"`
	Battery       int       `firestore:"battery"`
	AccelerationX int       `firestore:"acceleration_x"`
	AccelerationY int       `firestore:"acceleration_y"`
	AccelerationZ int       `firestore:"acceleration_z"`
}
