//go:build wireinject
// +build wireinject

package app

import (
	"context"
	"fmt"

	"github.com/google/wire"
	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/ent"
	"github.com/lagom/lagom-server/internal/handler"
	"github.com/lagom/lagom-server/internal/pkg/ai"
	"github.com/lagom/lagom-server/internal/pkg/jwt"
	"github.com/lagom/lagom-server/internal/pkg/redis"
	"github.com/lagom/lagom-server/internal/service"
	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
)

// ProvideEntClient 提供 ent 数据库客户端
// 开发模式下自动创建 Schema，生产模式需手动执行 migration
func ProvideEntClient(cfg *config.Config) (*ent.Client, error) {
	client, err := ent.Open(dialect.Postgres, cfg.Database.DSN)
	if err != nil {
		return nil, err
	}
	if cfg.Server.Mode != "release" {
		if err := client.Schema.Create(context.Background()); err != nil {
			return nil, fmt.Errorf("auto create schema: %w", err)
		}
	}
	return client, nil
}

// ProvideJWTManager 提供 JWT 管理器
func ProvideJWTManager(cfg *config.Config) *jwt.Manager {
	return jwt.NewManager(cfg.JWT.Secret)
}

// ProvideDeepSeekClient 提供 DeepSeek AI 客户端
func ProvideDeepSeekClient(cfg *config.Config) *ai.DeepSeekClient {
	return ai.NewDeepSeekClient(cfg.AI.DeepSeekKey)
}

// ProvideRedisConfig 提取 Redis 配置
func ProvideRedisConfig(cfg *config.Config) *config.RedisConfig {
	return &cfg.Redis
}

// AppSet 是完整的应用依赖集合
var AppSet = wire.NewSet(
	// Config
	config.Load,

	// Database
	ProvideEntClient,

	// Utilities
	ProvideJWTManager,
	ProvideDeepSeekClient,
	ProvideRedisConfig,
	redis.New,

	// Services
	service.NewAuthService,
	service.NewAIService,
	service.NewHabitService,
	service.NewCharacterService,

	// Handlers
	handler.NewAuthHandler,
	handler.NewChatHandler,
	handler.NewHabitHandler,
	handler.NewCharacterHandler,

	// App
	New,
)

// InitializeApp 初始化应用（由 Wire 生成）
func InitializeApp() (*App, error) {
	wire.Build(AppSet)
	return nil, nil
}
