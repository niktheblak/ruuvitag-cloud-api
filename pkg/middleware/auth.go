package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type UsersAndPasswordHashes map[string][]byte

func BasicAuth(h httprouter.Handle, users UsersAndPasswordHashes) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()
		if !hasAuth {
			auth(w)
			return
		}
		hashed, ok := users[user]
		if !ok {
			auth(w)
			return
		}
		err := bcrypt.CompareHashAndPassword(hashed, []byte(password))
		if err != nil {
			auth(w)
			return
		}
		h(w, r, ps)

	}
}

func auth(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}
