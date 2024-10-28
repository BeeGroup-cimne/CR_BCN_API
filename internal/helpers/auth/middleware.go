package auth

import (
	"context"
	"go-api/internal/helpers/jwtHelpers"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
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

		//user, err := dbOperations.FetchNeo4jUser(claims.Email)
		//if err != nil {
		//	http.Error(w, "User not found", http.StatusUnauthorized)
		//	return
		//}
		//ctx := context.WithValue(r.Context(), "user", user)

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
