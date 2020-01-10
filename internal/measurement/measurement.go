package measurement

import "time"

const Collection = "measurements"

type Measurement struct {
	Name          string    `json:"name" firestore:"name"`
	MAC           string    `json:"mac" firestore:"mac"`
	Timestamp     time.Time `json:"ts" firestore:"ts"`
	Temperature   float64   `json:"temperature" firestore:"temperature"`
	Humidity      float64   `json:"humidity" firestore:"humidity"`
	Pressure      float64   `json:"pressure" firestore:"pressure"`
	Battery       int       `json:"battery" firestore:"battery"`
	AccelerationX int       `json:"acceleration_x" firestore:"acceleration_x"`
	AccelerationY int       `json:"acceleration_y" firestore:"acceleration_y"`
	AccelerationZ int       `json:"acceleration_z" firestore:"acceleration_z"`
	ID            string    `json:"id" firestore:"-"`
}
