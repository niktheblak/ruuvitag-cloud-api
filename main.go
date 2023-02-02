package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/niktheblak/ruuvitag-cloud-api/pkg/auth"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement/postgres"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/middleware"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://127.0.0.1/ruuvitag"
	}
	table := os.Getenv("DATABASE_TABLE")
	if table == "" {
		table = "measurements"
	}
	ctx := context.Background()
	svc, err := postgres.New(ctx, connStr, table)
	if err != nil {
		log.Fatal(err)
	}
	tokens := strings.Split(os.Getenv("ALLOWED_TOKENS"), ",")
	var authenticator auth.Authenticator
	if len(tokens) > 0 {
		log.Printf("Authorized %d tokens", len(tokens))
		authenticator = auth.Static(tokens...)
	} else {
		log.Println("Allowing all tokens")
		authenticator = auth.AlwaysAllow()
	}
	router := httprouter.New()
	srv := server.NewServer(svc)
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "OK")
	})
	router.GET("/measurements/:name", middleware.Authenticator(srv.GetMeasurements, authenticator))
	router.POST("/receive", middleware.Authenticator(srv.Receive, authenticator))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
