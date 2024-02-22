package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/auth"
	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/middleware"
	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/server"
	"github.com/niktheblak/ruuvitag-cloud-api/postgres/pkg/service"
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
	svc, err := service.New(ctx, connStr, table)
	if err != nil {
		log.Fatal(err)
	}
	tokens := strings.Split(os.Getenv("ALLOWED_TOKENS"), ",")
	var authenticator auth.Authenticator
	if len(tokens) > 0 {
		slog.Info("Authorized tokens", "count", len(tokens))
		authenticator = auth.Static(tokens...)
	} else {
		slog.Info("Allowing all tokens")
		authenticator = auth.AlwaysAllow()
	}
	router := httprouter.New()
	srv := server.NewServer(svc, slog.Default())
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if _, err := fmt.Fprint(w, "OK"); err != nil {
			log.Fatal(err)
		}
	})
	router.GET("/measurements/:name", middleware.Authenticator(srv.GetMeasurements, authenticator))
	router.POST("/receive", middleware.Authenticator(srv.Receive, authenticator))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
