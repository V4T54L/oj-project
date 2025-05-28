package middleware

import (
	"context"
	"net/http"
	authmodule "online-judge/internal/auth_module"
)

type contextKey string

const UserIDKey = contextKey("user_id")

func JWTAuthMiddleware(auth *authmodule.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// authHeader := r.Header.Get("Authorization")
			// if authHeader == "" {
			// 	http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			// 	return
			// }

			// parts := strings.SplitN(authHeader, " ", 2)
			// if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			// 	http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
			// 	return
			// }

			// userID, err := auth.Validate(parts[1])
			// if err != nil {
			// 	http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			// 	return
			// }

			userID := 1

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
