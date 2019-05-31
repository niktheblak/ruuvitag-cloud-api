package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/measurement"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/middleware"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/server"
)

func readUsers(ctx context.Context, client *firestore.Client) (middleware.UsersAndPasswordHashes, error) {
	coll := client.Collection("users")
	iter := coll.Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return nil, err
	}
	type User struct {
		Username     string `firestore:"username"`
		PasswordHash string `firestore:"password_hash"`
	}
	m := make(map[string][]byte)
	for _, doc := range docs {
		var user User
		if err := doc.DataTo(&user); err != nil {
			return nil, err
		}
		m[user.Username] = []byte(user.PasswordHash)
	}
	return middleware.UsersAndPasswordHashes(m), nil
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
	users, err := readUsers(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	if len(users) == 0 {
		log.Printf("Warning: no users in database!")
	}
	meas := measurement.NewService(client)
	srv := server.NewServer(meas)
	router.GET("/measurements/:name", middleware.BasicAuth(srv.ListMeasurementsHandler, users))
	router.GET("/measurements/:name/:id", middleware.BasicAuth(srv.GetMeasurementHandler, users))
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
