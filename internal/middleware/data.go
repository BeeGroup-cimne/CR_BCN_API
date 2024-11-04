package middleware

import (
	"context"
	"go-api/internal/helpers/jwtHelpers"
	"go-api/internal/models/enums/jwtType"
	"net/http"
	"strings"
)

func DataMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Malformed Authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := jwtHelpers.DecodeToken(tokenString)

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)

		if claims.Type != jwtType.Access {
			http.Error(w, "Invalid token type", http.StatusUnauthorized)
			return
		}

		// TODO: MIDDLEWARE CONTENT FOR DATA ENDPOINTS

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
