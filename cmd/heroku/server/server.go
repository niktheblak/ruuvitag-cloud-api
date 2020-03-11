package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"github.com/niktheblak/ruuvitag-gollector/pkg/sensor"
)

type Server struct {
	Writer measurement.Writer
}

func (s *Server) Receive(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)
	var sd sensor.Data
	err := dec.Decode(&sd)
	if err != nil {
		BadRequest(w, "Invalid measurement: bad JSON")
		return
	}
	err = Validate(sd)
	if err != nil {
		BadRequest(w, "Invalid measurement: "+err.Error())
		return
	}
	if sd.Timestamp.IsZero() {
		sd.Timestamp = time.Now()
	}
	err = s.Writer.Write(r.Context(), sd)
	if err != nil {
		log.Print(err)
		Response(w, http.StatusInternalServerError, "Cloud not write measurement")
		return
	}
	fmt.Fprint(w, "OK")
}

func Validate(sd sensor.Data) error {
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

func BadRequest(w http.ResponseWriter, message string) {
	Response(w, http.StatusBadRequest, message)
}

func Response(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprint(w, message)
	return
}
