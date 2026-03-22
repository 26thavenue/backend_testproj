package config

import "os"

type Config struct {
	AppAddr string
	DBPath  string
}

func Load() Config {
	return Config{
		AppAddr: envOr("APP_ADDR", ":8080"),
		DBPath:  envOr("DB_PATH", "./data/analytics.db"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
