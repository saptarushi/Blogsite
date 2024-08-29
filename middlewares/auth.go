package middlewares

import (
	"Blogsite/utils"
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type key int

const (
	UserIDKey key = iota
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		userIDStr, err := utils.ParseJWT(tokenString)
		if err != nil {
			log.Printf("Error parsing JWT: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 64)
		if err != nil {
			log.Printf("Error converting user ID: %v", err)
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		log.Printf("Extracted user ID: %v (type: %T)", userID, userID)
		ctx := context.WithValue(r.Context(), UserIDKey, uint(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
