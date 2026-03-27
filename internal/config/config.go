package config

import (
	"os"
	"strconv"
	"strings"

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
	FrontendBaseURL   string
	StaticMountPrefix string
	StaticLocalDir    string

	AdminLoginUsername string
	AdminLoginPassword string
	AdminJWTSecret     string
	AdminTokenExpireHours int
}

func Load() Config {
	loadEnvFiles()

	publicBaseURL := getEnv("PUBLIC_BASE_URL", "http://localhost:8080")

	return Config{
		AppEnv:            getEnv("APP_ENV", "dev"),
		AppName:           getEnv("APP_NAME", "pm-backend"),
		Port:              getEnv("APP_PORT", "8080"),
		APIPrefix:         getEnv("API_PREFIX", "/api/v1"),
		AllowOrigins:      getEnv("ALLOW_ORIGINS", "*"),
		UseMock:           getEnv("USE_MOCK", "true") == "true",
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		PublicBaseURL:     publicBaseURL,
		FrontendBaseURL:   getEnv("FRONTEND_BASE_URL", publicBaseURL),
		StaticMountPrefix: getEnv("STATIC_MOUNT_PREFIX", "/static"),
		StaticLocalDir:    getEnv("STATIC_LOCAL_DIR", "./public"),
		AdminLoginUsername: getEnv("ADMIN_LOGIN_USERNAME", "admin"),
		AdminLoginPassword: getEnv("ADMIN_LOGIN_PASSWORD", "123456"),
		AdminJWTSecret: getEnv("ADMIN_JWT_SECRET", "pm-admin-dev-secret"),
		AdminTokenExpireHours: getEnvInt("ADMIN_TOKEN_EXPIRE_HOURS", 72),
	}
}

func loadEnvFiles() {
	mode := normalizeEnvMode(os.Getenv("APP_ENV"))
	files := []string{
		".env",
		".env." + mode,
		".env.local",
		".env." + mode + ".local",
	}
	_ = godotenv.Overload(files...)
}

func normalizeEnvMode(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "prod", "production":
		return "production"
	default:
		return "development"
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
