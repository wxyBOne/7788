// Package config 提供应用程序配置管理
package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用程序配置结构
type Config struct {
	Port         string // 服务器端口
	DatabaseURL  string // 数据库连接URL
	AIAPIKey     string // AI服务API密钥
	AIBaseURL    string // AI服务基础URL
	AIModel      string // AI模型名称
	ASRAPIKey    string // 语音识别API密钥
	TTSAPIKey    string // 语音合成API密钥
	VisionAPIKey string // 视觉识别API密钥
	JWTSecret    string // JWT密钥
	Environment  string // 运行环境
}

// Load 加载应用程序配置
func Load() *Config {
	// 加载.env文件（如果存在）
	_ = godotenv.Load()

	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		AIAPIKey:     getEnv("AI_API_KEY", ""),
		AIBaseURL:    getEnv("AI_BASE_URL", ""),
		AIModel:      getEnv("AI_MODEL", "qwen3-max"),
		ASRAPIKey:    getEnv("ASR_API_KEY", ""),
		TTSAPIKey:    getEnv("TTS_API_KEY", ""),
		VisionAPIKey: getEnv("VISION_API_KEY", ""),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量并转换为布尔值
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
