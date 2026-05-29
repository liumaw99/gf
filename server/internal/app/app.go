package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/config"
)

// App 应用程序
type App struct {
	router *gin.Engine
	server *http.Server
	cfg    *config.Config
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

	// 创建路由
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

// Run 启动服务
func (a *App) Run() error {
	return a.server.ListenAndServe()
}

// Shutdown 优雅关闭
func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
