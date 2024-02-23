package middleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/niktheblak/ruuvitag-cloud-api/common/pkg/auth"
)

func Authenticator(h httprouter.Handle, authenticator auth.Authenticator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			forbidden(w)
			return
		}
		err := authenticator.Authenticate(r.Context(), token)
		if err != nil {
			forbidden(w)
			return
		}
		h(w, r, ps)

	}
}

func forbidden(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}
