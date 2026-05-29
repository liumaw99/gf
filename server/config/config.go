package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	AI       AIConfig
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port int
	Mode string // debug / release
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN            string
	MaxOpenConns   int
	MaxIdleConns   int
	MaxLifetimeSec int
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string
	AccessTTL  int // minutes
	RefreshTTL int // hours
}

// AIConfig AI 服务配置
type AIConfig struct {
	DeepSeekKey string
	FallbackKey string
	MaxTokens   int
	Temperature float64
}

// Load 从环境变量加载配置
func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: getEnvInt("SERVER_PORT", 8080),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			DSN:            getEnv("DATABASE_DSN", "postgres://lagom:lagom@localhost:5432/lagom?sslmode=disable"),
			MaxOpenConns:   getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:   getEnvInt("DB_MAX_IDLE_CONNS", 5),
			MaxLifetimeSec: getEnvInt("DB_MAX_LIFETIME_SEC", 3600),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "lagom-dev-secret-change-in-production"),
			AccessTTL:  getEnvInt("JWT_ACCESS_TTL_MIN", 15),
			RefreshTTL: getEnvInt("JWT_REFRESH_TTL_HOUR", 168), // 7 days
		},
		AI: AIConfig{
			DeepSeekKey: getEnv("DEEPSEEK_API_KEY", ""),
			FallbackKey: getEnv("OPENAI_API_KEY", ""),
			MaxTokens:   getEnvInt("AI_MAX_TOKENS", 300),
			Temperature: getEnvFloat("AI_TEMPERATURE", 0.7),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return defaultValue
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Server.Mode == "release" {
		if c.AI.DeepSeekKey == "" {
			return fmt.Errorf("DEEPSEEK_API_KEY is required in release mode")
		}
		if c.JWT.Secret == "" || c.JWT.Secret == "lagom-dev-secret-change-in-production" {
			return fmt.Errorf("JWT_SECRET must be set to a secure value in release mode")
		}
	}
	return nil
}
