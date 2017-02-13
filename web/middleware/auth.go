package middleware

import (
	"net/http"

	"github.com/goji/httpauth"

	"github.com/bonnyci/quartermaster/lib"
)

var authDefaults = httpauth.AuthOptions{
	Realm:    "quartermaster",
	AuthFunc: lib.AuthToken,
}

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpauth.BasicAuth(authDefaults)(h).ServeHTTP(w, r)
	})
}

func AdminMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aU, _, _ := r.BasicAuth()
		if !lib.IsAdmin(aU) {
			http.Error(w, "Not member of Admin group", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

var AuthAndAdmin = []func(http.HandlerFunc) http.HandlerFunc{AuthMiddleware, AdminMiddleware}
