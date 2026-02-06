package middleware

import (
	"context"
	"net/http"
	"strings"
	"todo-api/internal/domain"
	"todo-api/internal/utils"
)

func AuthMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				utils.ParseError(w, domain.ErrUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := utils.ValidateJWT(token, jwtSecret)
			if err != nil {
				utils.ParseError(w, domain.ErrUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
