package middleware

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/julienschmidt/httprouter"
)

func BasicAuth(h httprouter.Handle, client *firestore.Client) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()
		if !hasAuth {
			authResponse(w)
			return
		}
		err := NewFirebaseAuthenticator(client, "users").Authenticate(r.Context(), user, password)
		if err != nil {
			authResponse(w)
			return
		}
		h(w, r, ps)

	}
}

func authResponse(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
