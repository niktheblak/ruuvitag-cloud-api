package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad measurement")
		return
	}
	err = s.Writer.Write(r.Context(), sd)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not write measurement")
		return
	}
	fmt.Fprint(w, "OK")
}
