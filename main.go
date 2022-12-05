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
		log.Fatal("$PORT must be set")
	}
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("$DATABASE_URL must be set")
	}
	ctx := context.Background()
	svc, err := postgres.New(ctx, dbUrl, "measurements")
	if err != nil {
		log.Fatal(err)
	}
	tokens := strings.Split(os.Getenv("ALLOWED_TOKENS"), ",")
	if len(tokens) == 0 {
		log.Fatal("No allowed tokens")
	}
	authenticator := &auth.StaticAuthenticator{
		AllowedTokens: tokens,
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