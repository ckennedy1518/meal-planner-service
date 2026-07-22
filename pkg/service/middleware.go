package service

import (
	"context"
	"net/http"
	"strings"
)

// We'll store the full auth payload in the request context so handlers can use it later.
type contextKey string

const claimsKey = contextKey("claims")

func AuthMiddleware(auth *AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		payload, err := auth.ValidateAuthHeader(authHeader)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Inject payload into request context.
		ctx := context.WithValue(r.Context(), claimsKey, payload)

		// Pass down to the next handler with the updated context.
		next(w, r.WithContext(ctx))
	}
}

// Helper so handlers can grab the auth payload from context.
func GetClaims(r *http.Request) (map[string]any, bool) {
	payload, ok := r.Context().Value(claimsKey).(map[string]any)
	return payload, ok
}
