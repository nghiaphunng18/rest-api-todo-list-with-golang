package utils

import (
	"errors"
	"time"
	"todo-api/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var jwtConfig config.JWTConfig

type ContextKey string
const UserIDKey ContextKey = "userID"

func SetJWTConfig(cfg config.JWTConfig) {
	if cfg.Expiration <= 0 {
		cfg.Expiration = time.Hour
	}
	jwtConfig = cfg
}

func GenerateJWT(userID uint, jwtSecret string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":      time.Now().Add(jwtConfig.Expiration).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string, jwtSecret string) (uint, error) {
	// parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil 
	})

	if err != nil {
		return 0, err // token invalid or time expired
	}

    // extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims or token not valid")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id claim missing or invalid")
	}

	return uint(userIDFloat), nil
}