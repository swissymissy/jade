package middleware

import (
	"log"
	"net/http"

	"github.com/swissymissy/jade/internal/auth"
)

// middleware to check for auth
func AuthRequired(next http.HandlerFunc, jwtSecret string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from header
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("Missing or Invalid auth header: %s\n", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// validate token
		adminID, err := auth.ValidateJWT(token, jwtSecret)
		if err != nil {
			log.Printf("Invalid token: %s\n", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// token is valid, log and let the request go though
		log.Printf("Authenticated admin: %s\n", adminID)
		next.ServeHTTP(w, r)
	})
}