package auth

import (
	"context"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username     string `firestore:"username"`
	PasswordHash string `firestore:"password_hash"`
}

type FirestoreAuthenticator struct {
	client     *firestore.Client
	collection string
}

func NewFirestoreAuthenticator(client *firestore.Client, collection string) Authenticator {
	return &FirestoreAuthenticator{
		client:     client,
		collection: collection,
	}
}

func (a *FirestoreAuthenticator) Authenticate(ctx context.Context, username, password string) error {
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
