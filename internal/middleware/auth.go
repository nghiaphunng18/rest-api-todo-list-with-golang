package middleware

import (
	"context"
	"net/http"
	"strings"
	"todo-api/internal/utils"
)

func AuthMiddleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			// extract the authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Authorization header required")
				return
			}
			
			if !strings.HasPrefix(authHeader, "Bearer ") {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid token format (must be Bearer <token>)")
				return
			}
			
			// exact token string
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			
			// validate the JWT token
			userID, err := utils.ValidateJWT(tokenString, jwtSecret)
			if err != nil {
				utils.ErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}
			
			// add the user_id to context
			ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}