package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置根结构
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	AI       AIConfig       `mapstructure:"ai"`
	OAuth    OAuthConfig    `mapstructure:"oauth"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug / release
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN            string `mapstructure:"dsn"`
	MaxOpenConns   int    `mapstructure:"max_open_conns"`
	MaxIdleConns   int    `mapstructure:"max_idle_conns"`
	MaxLifetimeSec int    `mapstructure:"max_lifetime_sec"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr         string        `mapstructure:"addr"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	PoolSize     int           `mapstructure:"pool_size"`
	MinIdleConns int           `mapstructure:"min_idle_conns"`
	MaxRetries   int           `mapstructure:"max_retries"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	AccessTTL  int    `mapstructure:"access_ttl_min"`  // minutes
	RefreshTTL int    `mapstructure:"refresh_ttl_hour"` // hours
}

// AIConfig AI 服务配置
type AIConfig struct {
	DeepSeekKey string  `mapstructure:"deepseek_key"`
	FallbackKey string  `mapstructure:"fallback_key"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
}

// OAuthConfig OAuth2 配置
type OAuthConfig struct {
	Apple  AppleOAuthConfig  `mapstructure:"apple"`
	Google GoogleOAuthConfig `mapstructure:"google"`
}

// AppleOAuthConfig Apple Sign In 配置
type AppleOAuthConfig struct {
	ClientID   string `mapstructure:"client_id"`
	TeamID     string `mapstructure:"team_id"`
	KeyID      string `mapstructure:"key_id"`
	PrivateKey string `mapstructure:"private_key"`
}

// GoogleOAuthConfig Google Sign In 配置
type GoogleOAuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// Load 加载配置
// 优先级：环境变量 > .env 文件 > config.yaml > 默认值
func Load() (*Config, error) {
	v := viper.New()

	// 1. 设置默认值
	setDefaults(v)

	// 2. 读取 .env 文件（如果存在）
	v.SetConfigFile(".env")
	_ = v.ReadInConfig() // 忽略错误，.env 是可选的

	// 3. 读取 config.yaml
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	_ = v.ReadInConfig() // 忽略错误，config.yaml 是可选的

	// 4. 环境变量覆盖（最高优先级）
	// 格式：LAGOM_SERVER_PORT, LAGOM_DATABASE_DSN
	v.SetEnvPrefix("LAGOM")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 5. 解析到结构体
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// 6. 从独立环境变量读取敏感配置（不通过 viper 前缀）
	readSensitiveEnv(&cfg)

	return &cfg, nil
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	// Server
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")

	// Database
	v.SetDefault("database.dsn", "postgres://lagom:lagom@localhost:5432/lagom?sslmode=disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.max_lifetime_sec", 3600)

	// Redis
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.pool_size", 20)
	v.SetDefault("redis.min_idle_conns", 5)
	v.SetDefault("redis.max_retries", 3)
	v.SetDefault("redis.dial_timeout", "5s")
	v.SetDefault("redis.read_timeout", "3s")
	v.SetDefault("redis.write_timeout", "3s")

	// JWT
	v.SetDefault("jwt.secret", "lagom-dev-secret-change-in-production")
	v.SetDefault("jwt.access_ttl_min", 15)
	v.SetDefault("jwt.refresh_ttl_hour", 168)

	// AI
	v.SetDefault("ai.max_tokens", 300)
	v.SetDefault("ai.temperature", 0.7)
}

// readSensitiveEnv 从独立环境变量读取敏感配置
func readSensitiveEnv(cfg *Config) {
	// AI Keys
	if v := viper.GetString("DEEPSEEK_API_KEY"); v != "" {
		cfg.AI.DeepSeekKey = v
	}
	if v := viper.GetString("OPENAI_API_KEY"); v != "" {
		cfg.AI.FallbackKey = v
	}

	// OAuth Secrets
	if v := viper.GetString("GOOGLE_CLIENT_SECRET"); v != "" {
		cfg.OAuth.Google.ClientSecret = v
	}
	if v := viper.GetString("APPLE_PRIVATE_KEY"); v != "" {
		cfg.OAuth.Apple.PrivateKey = v
	}

	// Database password from dedicated env (extract from DSN or separate)
	if v := viper.GetString("DB_PASSWORD"); v != "" {
		cfg.Database.DSN = strings.ReplaceAll(cfg.Database.DSN, ":lagom@", ":"+v+"@")
	}
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

// IsDebug 是否调试模式
func (c *Config) IsDebug() bool {
	return c.Server.Mode == "debug"
}

// IsRelease 是否生产模式
func (c *Config) IsRelease() bool {
	return c.Server.Mode == "release"
}
