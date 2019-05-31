package middleware

import (
	"cloud.google.com/go/firestore"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string `firestore:"username"`
	PasswordHash string `firestore:"password_hash"`
}

type Authenticator interface {
	Authenticate(ctx context.Context, username, password string) error
}

type FirebaseAuthenticator struct {
	client     *firestore.Client
	collection string
}

func NewFirebaseAuthenticator(client *firestore.Client, collection string) Authenticator {
	return &FirebaseAuthenticator{
		client:     client,
		collection: collection,
	}
}

func (a *FirebaseAuthenticator) Authenticate(ctx context.Context, username, password string) error {
	iter := a.client.Collection(a.collection).Where("username", "==", username).Limit(1).Documents(ctx)
	defer iter.Stop()
	doc, err := iter.Next()
	if err != nil {
		return err
	}
	var user User
	if err := doc.DataTo(&user); err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}
