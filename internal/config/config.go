package config

import (
	"os"
)

type Config struct {
	Port           string
	TessDataPrefix string
	StaticDir      string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "3000"),
		TessDataPrefix: getEnv("TESSDATA_PREFIX", "/opt/homebrew/share/tessdata"),
		StaticDir:      getEnv("STATIC_DIR", "static"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
