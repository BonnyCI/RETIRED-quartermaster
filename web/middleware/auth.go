package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/web/endpoints/auth"
)

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header["Authorization"][0][7:]

		valid, claims := auth.ValidateToken(token)
		if !valid {
			jww.ERROR.Printf("Authentication Faild with token: %s", token)
			http.Error(w, "Authentication Failed!", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "Claims", claims)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func AdminMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := ctx.Value("Claims").(auth.UserClaims)

		if !ok {
			http.Error(w, "Context doesn't contain authentication Claims.", http.StatusBadRequest)
		}

		if !claims.Admin {
			http.Error(w, "Not member of Admin group", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func SelfOrAdminMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		user := params["user"]

		ctx := r.Context()
		claims, ok := ctx.Value("Claims").(auth.UserClaims)

		if !ok {
			http.Error(w, "Context doesn't contain authentication Claims.", http.StatusBadRequest)
		}

		if !claims.Admin || claims.Username != user {
			http.Error(w, "Not member of Admin group, or the modifing their own status.", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}

var AuthAndAdmin = []func(http.HandlerFunc) http.HandlerFunc{AuthMiddleware, AdminMiddleware}
var AuthAndSelfOrAdmin = []func(http.HandlerFunc) http.HandlerFunc{AuthMiddleware, SelfOrAdminMiddleware}
