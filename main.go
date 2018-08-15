package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
)

type Service struct {
	ctx context.Context
	ds  *datastore.Client
}

func (s *Service) GetMeasurements(name string, start, end time.Time) (measurements []Measurement, err error) {
	if !start.IsZero() && !end.IsZero() && start.After(end) {
		return nil, fmt.Errorf("start timestamp cannot after end timestamp")
	}
	filters := make(map[string]interface{})
	filters["name"] = name
	if !start.IsZero() {
		filters["ts >="] = start
	}
	if !end.IsZero() {
		filters["ts <"] = end
	}
	query := datastore.NewQuery("Measurement")
	for k, v := range filters {
		query = query.Filter(k, v)
	}
	query = query.Order("-ts")
	keys, err := s.ds.GetAll(s.ctx, query, &measurements)
	if err != nil {
		return
	}
	for i, key := range keys {
		measurements[i].id = key.ID
	}
	return
}

type Measurement struct {
	Name        string    `datastore:"name"`
	MAC         string    `datastore:"mac"`
	Timestamp   time.Time `datastore:"ts"`
	Temperature float64   `datastore:"temperature"`
	Humidity    float64   `datastore:"humidity"`
	Pressure    float64   `datastore:"pressure"`
	id          int64
}

func main() {
	ctx := context.Background()
	// Read project ID from environment variable DATASTORE_PROJECT_ID
	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		log.Fatalf("Error while creating client: %v", err)
	}
	service := Service{
		ctx: ctx,
		ds:  client,
	}
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.HandleFunc("/measurements", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		name := query.Get("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Query parameter name must be specified"))
			return
		}
		var start, end time.Time
		if query.Get("start") != "" {
			start, _ = time.Parse("2006-01-02", query.Get("start"))
		}
		if query.Get("end") != "" {
			end, _ = time.Parse("2006-01-02", query.Get("end"))
		}
		if !start.IsZero() && !end.IsZero() && start == end {
			end = end.AddDate(0, 0, 1)
		}
		measurements, err := service.GetMeasurements(name, start, end)
		if err != nil {
			log.Printf("Error while querying measurements: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error while querying measurements"))
			return
		}
		body, _ := json.Marshal(measurements)
		w.Write(body)
	})
	http.ListenAndServe(":8080", nil)
	appengine.Main()
}
