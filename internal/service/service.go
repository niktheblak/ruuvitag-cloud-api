package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
)

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

func (s *Service) GetMeasurement(id int64) (m measurement.Measurement, err error) {
	key := datastore.IDKey(measurement.Kind, id, nil)
	err = s.client.Get(s.ctx, key, &m)
	m.ID = id
	return
}

func (s *Service) ListMeasurements(name string, from, to time.Time, limit int) (measurements []measurement.Measurement, err error) {
	filters := make(map[string]interface{})
	filters["name ="] = name
	if !from.IsZero() {
		filters["ts >="] = from
	}
	if !to.IsZero() {
		filters["ts <"] = to
		if to.Sub(from) <= time.Hour*24 {
			// Don't limit number of results if the query is for less than one day
			limit = 0
		}
	}
	query := datastore.NewQuery(measurement.Kind)
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

func GetMeasurementHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		log.Fatalf("Error while creating datastore client: %v", err)
		http.Error(w, "Error while creating datastore client", http.StatusInternalServerError)
		return
	}
	defer client.Close()
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, "ID parameter must be numeric", http.StatusBadRequest)
		return
	}
	service := NewService(ctx, client)
	m, err := service.GetMeasurement(id)
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
	writeJSON(w, m)
}

func ListMeasurementsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	ctx := r.Context()
	client, err := datastore.NewClient(ctx, "")
	if err != nil {
		log.Fatalf("Error while creating datastore client: %v", err)
		http.Error(w, "Error while creating datastore client", http.StatusInternalServerError)
		return
	}
	defer client.Close()
	query := r.URL.Query()
	from, to, err := parseTimeRange(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := parseLimit(query)
	service := NewService(ctx, client)
	measurements, err := service.ListMeasurements(name, from, to, limit)
	if err != nil {
		log.Printf("Error while querying measurements: %v", err)
		http.Error(w, "Error while querying measurement", http.StatusInternalServerError)
		return
	}
	writeJSON(w, measurements)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	err := enc.Encode(v)
	if err != nil {
		log.Printf("Error while writing response: %v", err)
	}
}

func parseTimeRange(query url.Values) (from time.Time, to time.Time, err error) {
	if query.Get("from") != "" {
		from, err = time.Parse("2006-01-02", query.Get("from"))
	}
	if err != nil {
		return
	}
	if query.Get("to") != "" {
		to, err = time.Parse("2006-01-02", query.Get("to"))
	}
	if err != nil {
		return
	}
	if !from.IsZero() && !to.IsZero() && from == to {
		to = to.AddDate(0, 0, 1)
	}
	if to.IsZero() || to.After(time.Now()) {
		to = time.Now().UTC()
	}
	if from.After(to) {
		err = fmt.Errorf("from timestamp cannot be after to timestamp")
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
