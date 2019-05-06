package measurement

import "time"

const Collection = "measurements"

type Measurement struct {
	Name        string    `json:"name" firestore:"name"`
	MAC         string    `json:"mac" firestore:"mac"`
	Timestamp   time.Time `json:"ts" firestore:"ts"`
	Temperature float64   `json:"temperature" firestore:"temperature"`
	Humidity    float64   `json:"humidity" firestore:"humidity"`
	Pressure    float64   `json:"pressure" firestore:"pressure"`
	ID          string    `json:"id" firestore:"-"`
}
