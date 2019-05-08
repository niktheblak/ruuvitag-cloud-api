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

	"cloud.google.com/go/firestore"
	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/service"
)

var client *firestore.Client

func GetMeasurementHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	id := ps.ByName("id")
	srv := service.NewService(ctx, client)
	m, err := srv.GetMeasurement(id)
	switch err {
	case nil:
	case service.ErrNotFound:
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
	query := r.URL.Query()
	from, to, err := parseTimeRange(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := parseLimit(query)
	srv := service.NewService(ctx, client)
	measurements, err := srv.ListMeasurements(name, from, to, limit)
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

func main() {
	ctx := context.Background()
	var err error
	client, err = firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		log.Fatalf("Error while creating datastore client: %v", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintln(w, "OK")
	})
	router.GET("/_ah/health", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintln(w, "OK")
	})
	router.GET("/measurements/:name", ListMeasurementsHandler)
	router.GET("/measurements/:name/:id", GetMeasurementHandler)
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
