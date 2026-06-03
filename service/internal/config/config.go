package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string
	Env     string
	DBURL   string

	JWTSecret         string
	JWTRefreshSecret  string
	JWTAccessExpire   time.Duration
	JWTRefreshExpire  time.Duration

	S3Endpoint  string
	S3Bucket    string
	S3AccessKey string
	S3SecretKey string
	S3Region    string

	LLMProvider string
	LLMEndpoint string
	LLMModel    string
	OpenAIKey   string
	AnthropicKey string

	FrontendURL string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:   getEnv("PORT", "8080"),
		Env:    getEnv("ENV", "development"),
		DBURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/quikram?sslmode=disable"),

		JWTSecret:        getEnv("JWT_ACCESS_SECRET", "dev-access-secret"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "dev-refresh-secret"),
		JWTAccessExpire:  getDuration("JWT_ACCESS_EXPIRE", 15*time.Minute),
		JWTRefreshExpire: getDuration("JWT_REFRESH_EXPIRE", 720*time.Hour),

		S3Endpoint:  getEnv("S3_ENDPOINT", ""),
		S3Bucket:    getEnv("S3_BUCKET", "quikram-files"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey: getEnv("S3_SECRET_KEY", ""),
		S3Region:    getEnv("S3_REGION", "auto"),

		LLMProvider:  getEnv("LLM_PROVIDER", "openai"),
		LLMEndpoint:  getEnv("LLM_ENDPOINT", "https://api.openai.com/v1"),
		LLMModel:     getEnv("LLM_MODEL", ""),
		OpenAIKey:    getEnv("OPENAI_API_KEY", ""),
		AnthropicKey: getEnv("ANTHROPIC_API_KEY", ""),

		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
