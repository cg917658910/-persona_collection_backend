package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv            string
	AppName           string
	Port              string
	APIPrefix         string
	AllowOrigins      string
	UseMock           bool
	DatabaseURL       string
	PublicBaseURL     string
	StaticMountPrefix string
	StaticLocalDir    string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:            getEnv("APP_ENV", "dev"),
		AppName:           getEnv("APP_NAME", "pm-backend"),
		Port:              getEnv("APP_PORT", "8080"),
		APIPrefix:         getEnv("API_PREFIX", "/api/v1"),
		AllowOrigins:      getEnv("ALLOW_ORIGINS", "*"),
		UseMock:           getEnv("USE_MOCK", "true") == "true",
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		PublicBaseURL:     getEnv("PUBLIC_BASE_URL", "http://localhost:8080"),
		StaticMountPrefix: getEnv("STATIC_MOUNT_PREFIX", "/static"),
		StaticLocalDir:    getEnv("STATIC_LOCAL_DIR", "./public"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
