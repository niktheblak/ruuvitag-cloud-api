package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/api"
	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/measurement"
	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/sensor"
)

type Server struct {
	svc    measurement.Service
	logger *slog.Logger
}

func NewServer(svc measurement.Service, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	return &Server{
		svc:    svc,
		logger: logger,
	}
}

func (s *Server) Receive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)
	var sd sensor.Data
	err := dec.Decode(&sd)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Invalid measurement", slog.Any("error", err))
		http.Error(w, "Invalid measurement", http.StatusBadRequest)
		return
	}
	err = validate(sd)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Measurement failed validation", slog.Any("error", err))
		http.Error(w, "Invalid measurement", http.StatusBadRequest)
		return
	}
	if sd.Timestamp.IsZero() {
		sd.Timestamp = time.Now()
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Received measurement", slog.Any("measurement", sd))
	err = s.svc.Write(r.Context(), sd)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error while writing measurement", slog.Any("error", err))
		http.Error(w, "Cloud not write measurement", http.StatusInternalServerError)
		return
	}
	if _, err := fmt.Fprint(w, "OK"); err != nil {
		panic(err)
	}
}

func (s *Server) GetMeasurements(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	query := r.URL.Query()
	from, to, err := api.ParseTimeRange(query.Get("from"), query.Get("to"))
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Invalid time range", slog.String("from", query.Get("from")), slog.String("to", query.Get("to")), slog.Any("error", err))
		http.Error(w, "Invalid time range", http.StatusBadRequest)
		return
	}
	limit := api.ParseLimit(query.Get("limit"))
	var measurements []sensor.Data
	if query.Get("ts") != "" {
		ts, err := time.Parse(time.RFC3339Nano, query.Get("ts"))
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Invalid timestamp", slog.String("timestamp", query.Get("ts")), slog.Any("error", err))
			http.Error(w, "Invalid timestamp", http.StatusBadRequest)
			return
		}
		m, err := s.svc.GetMeasurement(r.Context(), name, ts)
		measurements = append(measurements, m)
	} else {
		measurements, err = s.svc.ListMeasurements(r.Context(), name, from, to, limit)
	}
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error while reading measurements", slog.Any("error", err))
		http.Error(w, "Error while reading measurements", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if len(measurements) > 1 {
		err = enc.Encode(map[string]interface{}{
			"measurements": measurements,
		})
	} else {
		err = enc.Encode(measurements[0])
	}
	if err != nil {
		panic(err)
	}
}

func validate(sd sensor.Data) error {
	if sd.Addr == "" {
		return fmt.Errorf("empty address")
	}
	if sd.Name == "" {
		return fmt.Errorf("empty name")
	}
	if sd.Temperature == 0 && sd.Humidity == 0 && sd.Pressure == 0 {
		return fmt.Errorf("all main readings are zero")
	}
	return nil
}
