package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/ent"
	"github.com/lagom/lagom-server/internal/handler"
)

// App 应用程序
type App struct {
	router *gin.Engine
	server *http.Server
	cfg    *config.Config
	client *ent.Client
}

// New 创建应用实例（手动组装，用于测试）
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := setupRouter(cfg)

	return &App{
		router: router,
		cfg:    cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: router,
		},
	}, nil
}

// NewWireApp 创建应用实例（Wire 注入）
func NewWireApp(
	cfg *config.Config,
	client *ent.Client,
	authHandler *handler.AuthHandler,
	chatHandler *handler.ChatHandler,
	habitHandler *handler.HabitHandler,
	characterHandler *handler.CharacterHandler,
) *App {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := setupWireRouter(cfg, authHandler, chatHandler, habitHandler, characterHandler)

	return &App{
		router: router,
		cfg:    cfg,
		client: client,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: router,
		},
	}
}

// Run 启动服务
func (a *App) Run() error {
	return a.server.ListenAndServe()
}

// Shutdown 优雅关闭
func (a *App) Shutdown(ctx context.Context) error {
	if a.client != nil {
		a.client.Close()
	}
	return a.server.Shutdown(ctx)
}
