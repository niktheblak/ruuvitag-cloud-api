package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/auth"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/measurement"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/middleware"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/server"
)

func main() {
	ctx := context.Background()
	var err error
	client, err := firestore.NewClient(ctx, firestore.DetectProjectID)
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
	meas := measurement.NewService(client)
	srv := server.NewServer(meas)
	authenticator := auth.NewFirestoreAuthenticator(client, "users")
	router.GET("/measurements/:name", middleware.Authenticator(srv.ListMeasurementsHandler, authenticator))
	router.GET("/measurements/:name/:id", middleware.Authenticator(srv.GetMeasurementHandler, authenticator))
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
