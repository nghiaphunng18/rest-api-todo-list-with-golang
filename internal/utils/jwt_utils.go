package utils

import (
	"time"
	"todo-api/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

var jwtConfig config.JWTConfig

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