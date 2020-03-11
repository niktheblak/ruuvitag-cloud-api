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
		badRequest(w)
		return
	}
	if sd.Addr == "" {
		badRequest(w)
		return
	}
	if sd.Name == "" {
		badRequest(w)
		return
	}
	if sd.Timestamp.IsZero() {
		sd.Timestamp = time.Now()
	}
	err = s.Writer.Write(r.Context(), sd)
	if err != nil {
		log.Print(err)
		response(w, http.StatusInternalServerError, "Cloud not write measurement")
		return
	}
	fmt.Fprint(w, "OK")
}

func badRequest(w http.ResponseWriter) {
	response(w, http.StatusBadRequest, "Bad request")
}

func response(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprint(w, message)
	return
}
