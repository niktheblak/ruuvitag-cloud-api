package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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
	ctx := appengine.NewContext(r)
	appID := getAppID(ctx)
	client, err := datastore.NewClient(ctx, appID)
	if err != nil {
		log.Errorf(ctx, "Error while creating datastore client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while creating datastore client"))
		return
	}
	defer client.Close()
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ID parameter must be numeric"))
		return
	}
	m, err := NewService(ctx, client).GetMeasurement(id)
	switch err {
	case nil:
	case datastore.ErrNoSuchEntity:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Measurement with given ID not found"))
		return
	default:
		log.Errorf(ctx, "Error while querying measurement %v: %v", id, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while querying measurement"))
		return
	}
	enc := json.NewEncoder(w)
	err = enc.Encode(m)
	if err != nil {
		log.Errorf(ctx, "Error while writing response: %v", err)
	}
}

func ListMeasurementsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	ctx := appengine.NewContext(r)
	appID := getAppID(ctx)
	client, err := datastore.NewClient(ctx, appID)
	if err != nil {
		log.Errorf(ctx, "Error while creating datastore client: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while creating datastore client"))
		return
	}
	defer client.Close()
	query := r.URL.Query()
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
	measurements, err := NewService(ctx, client).ListMeasurements(name, from, to, int(limit))
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

func getAppID(ctx context.Context) string {
	id := appengine.AppID(ctx)
	if appengine.IsDevAppServer() || id == "None" {
		return "ruuvitag-212713"
	}
	return id
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	m := r.PathPrefix("/measurements").Subrouter()
	m.HandleFunc("/{name}", ListMeasurementsHandler)
	m.HandleFunc("/{name}/{id:[0-9]+}", GetMeasurementHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
	appengine.Main()
}
