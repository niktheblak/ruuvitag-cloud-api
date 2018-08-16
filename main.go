package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type Measurement struct {
	Name        string    `json:"name" datastore:"name"`
	MAC         string    `json:"mac" datastore:"mac"`
	Timestamp   time.Time `json:"ts" datastore:"ts"`
	Temperature float64   `json:"temperature" datastore:"temperature"`
	Humidity    float64   `json:"humidity" datastore:"humidity"`
	Pressure    float64   `json:"pressure" datastore:"pressure"`
	id          int64
}

type Service struct {
	ctx    context.Context
	client *datastore.Client
}

func NewService(ctx context.Context, client *datastore.Client) *Service {
	return &Service{
		ctx:    ctx,
		client: client,
	}
}

func (s *Service) GetMeasurements(name string, from, to time.Time, limit int) (measurements []Measurement, err error) {
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		return nil, fmt.Errorf("from timestamp cannot after to timestamp")
	}
	filters := make(map[string]interface{})
	filters["name ="] = name
	if !from.IsZero() {
		filters["ts >="] = from
	}
	if !to.IsZero() {
		filters["ts <"] = to
	}
	query := datastore.NewQuery("Measurement")
	for k, v := range filters {
		query = query.Filter(k, v)
	}
	query = query.Order("-ts")
	if limit > 0 {
		query = query.Limit(limit)
	}
	keys, err := s.client.GetAll(s.ctx, query, &measurements)
	if err != nil {
		return
	}
	for i, key := range keys {
		measurements[i].id = key.ID
	}
	return
}

func GetMeasurementsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	appID := appengine.AppID(ctx)
	client, err := datastore.NewClient(ctx, appID)
	if err != nil {
		log.Errorf(ctx, "Error while creating datastore client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while creating datastore client"))
		return
	}
	defer client.Close()
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Query parameter name must be specified"))
		return
	}
	var from, to time.Time
	if query.Get("from") != "" {
		from, _ = time.Parse("2006-01-02", query.Get("from"))
	}
	if query.Get("to") != "" {
		to, _ = time.Parse("2006-01-02", query.Get("to"))
	}
	if !from.IsZero() && !to.IsZero() && from == to {
		to = to.AddDate(0, 0, 1)
	}
	var limit int64
	if query.Get("limit") != "" {
		limit, _ = strconv.ParseInt(query.Get("limit"), 10, 64)
	}
	if limit <= 0 {
		limit = 20
	}
	measurements, err := NewService(ctx, client).GetMeasurements(name, from, to, int(limit))
	if err != nil {
		log.Errorf(ctx, "Error while querying measurements: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while querying measurements"))
		return
	}
	log.Debugf(ctx, "Read %v measurements for RuuviTag %v", len(measurements), name)
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(measurements)
	if err != nil {
		log.Errorf(ctx, "Error while writing response: %v", err)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.HandleFunc("/measurements", GetMeasurementsHandler)
	http.ListenAndServe(":8080", nil)
	appengine.Main()
}
