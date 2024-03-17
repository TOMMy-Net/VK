package middleware

import (
	"net/http"

	"github.com/TOMMy-Net/VK/db"
	"github.com/TOMMy-Net/VK/services"
)

func TokenAuthWare(storage *db.Storage, auth *services.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			if req.URL.Path == "/api/auth" || req.Method == http.MethodGet{
				next.ServeHTTP(w, req)
				return
			}
			b, err := auth.AuthUser(token)
			if b != nil && err == nil {
				if b["role"].(float64) == 1 {
					next.ServeHTTP(w, req)
					return
				} else if b["role"].(float64) == 0 && req.Method == http.MethodGet {
					next.ServeHTTP(w, req)
					return
				}
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

		})
	}
}
