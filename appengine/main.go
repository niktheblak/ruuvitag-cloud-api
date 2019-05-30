package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/internal/service"
)

func readUsers() (UsersAndPasswordHashes, error) {
	usersEnv := os.Getenv("USERS")
	if usersEnv == "" {
		return nil, nil
	}
	m := make(map[string][]byte)
	users := strings.Split(usersEnv, ",")
	for _, user := range users {
		tokens := strings.SplitN(user, ":", 2)
		if len(tokens) < 2 {
			return nil, fmt.Errorf("invalid user: %s", user)
		}
		m[tokens[0]] = []byte(tokens[1])
	}
	return UsersAndPasswordHashes(m), nil
}

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
	users, err := readUsers()
	if err != nil {
		log.Fatal(err)
	}
	srv := NewServer(service.NewService(client))
	router.GET("/measurements/:name", BasicAuth(srv.ListMeasurementsHandler, users))
	router.GET("/measurements/:name/:id", BasicAuth(srv.GetMeasurementHandler, users))
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
