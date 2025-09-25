package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	AIAPIKey     string
	AIBaseURL    string
	ASRAPIKey    string
	TTSAPIKey    string
	VisionAPIKey string
	JWTSecret    string
	Environment  string
}

func Load() *Config {
	// 加载.env文件（如果存在）
	_ = godotenv.Load()

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		AIAPIKey:     getEnv("AI_API_KEY", ""),
		AIBaseURL:    getEnv("AI_BASE_URL", ""),
		ASRAPIKey:    getEnv("ASR_API_KEY", ""),
		TTSAPIKey:    getEnv("TTS_API_KEY", ""),
		VisionAPIKey: getEnv("VISION_API_KEY", ""),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
