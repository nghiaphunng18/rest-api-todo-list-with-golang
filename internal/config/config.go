package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	
	JWTSecret string

	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

type JWTConfig struct {
	Expiration time.Duration
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(".env file not found") 
	}

	cfg := &Config{
		ServerPort: os.Getenv("SERVER_PORT"),
		// auth
		JWTSecret: os.Getenv("JWT_SECRET"),
		// connect to database
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
	}
    
    //default port
    if cfg.ServerPort == "" {
        cfg.ServerPort = "8080"
    }

	return cfg, nil
}

func LoadJWTConfig() JWTConfig {
	expiration := time.Hour

	if v := os.Getenv("JWT_EXPIRATION_HOURS"); v != "" {
		if hours, err := strconv.Atoi(v); err == nil && hours > 0 {
			expiration = time.Duration(hours) * time.Hour
		}
	}

	return JWTConfig{
		Expiration: expiration,
	}
}