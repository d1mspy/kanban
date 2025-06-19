package config

import (
	"log"
	"os"
	"sync"
)

type Config struct {
	PostgresURI string
	Host        string
	JWTSecret   []byte
}

var (
	config *Config
	once   sync.Once
)

func Load() {
	once.Do(func ()  {
		pg := os.Getenv("POSTGRES")
		if pg == "" {
			log.Fatal("POSTGRES env is required")
		}

		host := os.Getenv("HOST")
		if host == "" {
			log.Fatal("HOST env is required")
		}

		jwtKey := os.Getenv("JWT_SECRET")
		if jwtKey == "" {
			log.Fatal("JWT_SECRET env is required")
		}

		config = &Config{
			PostgresURI: pg,
			Host: host,
			JWTSecret: []byte(jwtKey),
		}
	})
}

func Get() *Config {
	if config == nil {
		log.Fatal("config not loaded: call config.Load() first")
	}
	return config
}