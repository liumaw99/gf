package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/internal/pkg/logger"
)

// App 应用程序
type App struct {
	router *gin.Engine
	server *http.Server
	cfg    *config.Config
	log    *slog.Logger
}

// New 创建应用实例
func New() (*App, error) {
	// 加载配置
	cfg, err := config.Load()

	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	// 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	log := logger.New(logger.Config{Mode: cfg.Server.Mode})
	logger.SetDefault(log)
	log.Info("config loaded", "server_mode", cfg.Server.Mode, "server_port", cfg.Server.Port)

	// 创建路由
	router := setupRouter(log)

	return &App{
		router: router,
		cfg:    cfg,
		log:    log,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
			Handler: router,
		},
	}, nil
}

// Run 启动服务
func (a *App) Run() error {
	a.log.Info("server starting", "addr", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown 优雅关闭
func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
