package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	Port             string
	OpenLibraryURL   string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://library:library@localhost:5432/library?sslmode=disable"),
		Port:           getEnv("PORT", "8080"),
		OpenLibraryURL: getEnv("OPEN_LIBRARY_BASE_URL", "https://openlibrary.org"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
