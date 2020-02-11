package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/measurement"
)

type Server struct {
	meas measurement.Service
}

func NewServer(meas measurement.Service) *Server {
	return &Server{
		meas: meas,
	}
}

func (s *Server) GetMeasurementHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	name := ps.ByName("name")
	ts, err := time.Parse(time.RFC3339, ps.ByName("ts"))
	if err != nil {
		http.Error(w, "Invalid timestamp", http.StatusBadRequest)
		return
	}
	m, err := s.meas.GetMeasurement(ctx, name, ts)
	switch err {
	case nil:
	case measurement.ErrNotFound:
		http.Error(w, "Measurement with given ID not found", http.StatusNotFound)
		return
	default:
		log.Printf("Error while querying measurement: %v", err)
		http.Error(w, "Error while querying measurement", http.StatusInternalServerError)
		return
	}
	writeJSON(w, m)
}

func (s *Server) ListMeasurementsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	ctx := r.Context()
	query := r.URL.Query()
	from, to, err := ParseTimeRange(query.Get("from"), query.Get("to"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := ParseLimit(query.Get("limit"))
	measurements, err := s.meas.ListMeasurements(ctx, name, from, to, limit)
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
