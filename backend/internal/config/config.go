package config

import (
	"log"
	"os"
)

type Config struct {
	PostgresURI string
	Host        string
	JWTSecret   string
}

func Load() *Config {
	pg := os.Getenv("POSTGRES")
	if pg == "" {
		log.Fatal("POSTGRES env is required")
	}

	host := os.Getenv("HOST")
	if host == "" {
		log.Fatal("HOST env is required")
	}

	jwtKey := os.Getenv("JWT_SECRET")
	if host == "" {
		log.Fatal("JWT_SECRET env is required")
	}

	return &Config{
		PostgresURI: pg,
		Host: host,
		JWTSecret: jwtKey,
	}
}