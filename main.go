package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
)

const Kind = "Measurement"

type Measurement struct {
	Name        string    `json:"name" datastore:"name"`
	MAC         string    `json:"mac" datastore:"mac"`
	Timestamp   time.Time `json:"ts" datastore:"ts"`
	Temperature float64   `json:"temperature" datastore:"temperature"`
	Humidity    float64   `json:"humidity" datastore:"humidity"`
	Pressure    float64   `json:"pressure" datastore:"pressure"`
	ID          int64     `json:"id" datastore:"-"`
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

func (s *Service) GetMeasurement(id int64) (measurement Measurement, err error) {
	key := datastore.IDKey(Kind, id, nil)
	err = s.client.Get(s.ctx, key, &measurement)
	measurement.ID = id
	return
}

func (s *Service) ListMeasurements(name string, from, to time.Time, limit int) (measurements []Measurement, err error) {
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		return nil, fmt.Errorf("from timestamp cannot be after to timestamp")
	}
	filters := make(map[string]interface{})
	filters["name ="] = name
	if !from.IsZero() {
		filters["ts >="] = from
	}
	if !to.IsZero() {
		filters["ts <"] = to
	}
	query := datastore.NewQuery(Kind)
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
		measurements[i].ID = key.ID
	}
	return
}

func GetMeasurementHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		log.Printf("Error while creating datastore client: %v", err)
		http.Error(w, "Error while creating datastore client", http.StatusInternalServerError)
		return
	}
	defer client.Close()
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "ID parameter must be numeric", http.StatusBadRequest)
		return
	}
	service := NewService(ctx, client)
	measurement, err := service.GetMeasurement(id)
	switch err {
	case nil:
	case datastore.ErrNoSuchEntity:
		http.Error(w, "Measurement with given ID not found", http.StatusNotFound)
		return
	default:
		log.Printf("Error while querying measurement %v: %v", id, err)
		http.Error(w, "Error while querying measurement", http.StatusInternalServerError)
		return
	}
	writeJSON(ctx, w, measurement)
}

func ListMeasurementsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		log.Printf("Error while creating datastore client: %v", err)
		http.Error(w, "Error while creating datastore client", http.StatusInternalServerError)
		return
	}
	defer client.Close()
	query := r.URL.Query()
	from, to := parseTimeRange(query)
	limit := parseLimit(query)
	service := NewService(ctx, client)
	measurements, err := service.ListMeasurements(name, from, to, limit)
	if err != nil {
		log.Printf("Error while querying measurements: %v", err)
		http.Error(w, "Error while querying measurement", http.StatusInternalServerError)
		return
	}
	writeJSON(ctx, w, measurements)
}

func writeJSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	err := enc.Encode(v)
	if err != nil {
		log.Printf("Error while writing response: %v", err)
	}
}

func parseTimeRange(query url.Values) (from time.Time, to time.Time) {
	if query.Get("from") != "" {
		from, _ = time.Parse("2006-01-02", query.Get("from"))
	}
	if query.Get("to") != "" {
		to, _ = time.Parse("2006-01-02", query.Get("to"))
	}
	if !from.IsZero() && !to.IsZero() && from == to {
		to = to.AddDate(0, 0, 1)
	}
	return
}

func parseLimit(query url.Values) int {
	var limit int64
	if query.Get("limit") != "" {
		limit, _ = strconv.ParseInt(query.Get("limit"), 10, 64)
	}
	if limit <= 0 {
		limit = 20
	}
	return int(limit)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	r.HandleFunc("/_ah/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	m := r.PathPrefix("/measurements").Subrouter()
	m.HandleFunc("/{name}", ListMeasurementsHandler)
	m.HandleFunc("/{name}/{id:[0-9]+}", GetMeasurementHandler)
	http.Handle("/", r)
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
