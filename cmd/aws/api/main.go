package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/server"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	db := dynamodb.New(sess)
	table := os.Getenv("TABLE")
	if table == "" {
		table = "measurements"
	}
	meas := NewService(db, table)
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
	srv := server.NewServer(meas)
	router.GET("/measurements/:name", srv.ListMeasurementsHandler)
	router.GET("/measurements/:name/:id", srv.GetMeasurementHandler)
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
