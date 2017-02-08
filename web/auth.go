package web

import (
	"net/http"

	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/lib"
)

func AdminMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aU, _, _ := r.BasicAuth()
		if !lib.IsAdmin(aU) {
			return
		}
		h.ServeHTTP(w, r)
	})
}
