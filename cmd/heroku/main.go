package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement/postgres"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("$DATABASE_URL must be set")
	}
	ctx := context.Background()
	writer, err := postgres.New(ctx, dbUrl, "measurements")
	if err != nil {
		log.Fatal(err)
	}
	router := httprouter.New()
	srv := &Server{
		Writer: writer,
	}
	router.POST("/receive", srv.Receive)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
