package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/niktheblak/ruuvitag-cloud-api/pkg/auth"
)

func Authenticator(h httprouter.Handle, authenticator auth.Authenticator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()
		if !hasAuth {
			authResponse(w)
			return
		}
		err := authenticator.Authenticate(r.Context(), user, password)
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
