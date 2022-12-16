package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/niktheblak/ruuvitag-cloud-api/pkg/api"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/sensor"
)

type Server struct {
	svc measurement.Service
}

func NewServer(svc measurement.Service) *Server {
	return &Server{
		svc: svc,
	}
}

func (s *Server) Receive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)
	var sd sensor.Data
	err := dec.Decode(&sd)
	if err != nil {
		badRequest(w, "Invalid measurement: bad JSON")
		return
	}
	err = validate(sd)
	if err != nil {
		badRequest(w, "Invalid measurement: "+err.Error())
		return
	}
	if sd.Timestamp.IsZero() {
		sd.Timestamp = time.Now()
	}
	err = s.svc.Write(r.Context(), sd)
	if err != nil {
		log.Print(err)
		response(w, http.StatusInternalServerError, "Cloud not write measurement")
		return
	}
	if _, err := fmt.Fprint(w, "OK"); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) GetMeasurements(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	query := r.URL.Query()
	from, to, err := api.ParseTimeRange(query.Get("from"), query.Get("to"))
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	limit := api.ParseLimit(query.Get("limit"))
	var measurements []sensor.Data
	if query.Get("ts") != "" {
		ts, err := time.Parse(time.RFC3339Nano, query.Get("ts"))
		if err != nil {
			badRequest(w, err.Error())
			return
		}
		m, err := s.svc.GetMeasurement(r.Context(), name, ts)
		measurements = append(measurements, m)
	} else {
		measurements, err = s.svc.ListMeasurements(r.Context(), name, from, to, limit)
	}
	if err != nil {
		internalServerError(w, err.Error())
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
		internalServerError(w, err.Error())
		return
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

func badRequest(w http.ResponseWriter, message string) {
	response(w, http.StatusBadRequest, message)
}

func internalServerError(w http.ResponseWriter, message string) {
	response(w, http.StatusInternalServerError, message)
}

func response(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	if _, err := fmt.Fprint(w, message); err != nil {
		log.Fatal(err)
	}
	return
}
