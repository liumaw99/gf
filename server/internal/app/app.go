package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/ent"
	"github.com/lagom/lagom-server/internal/handler"
	"github.com/lagom/lagom-server/internal/pkg/redis"
)

// App 应用程序
type App struct {
	router *gin.Engine
	server *http.Server
	cfg    *config.Config
	client *ent.Client
	redis  *redis.Client
}

// New 创建应用实例（由 Wire 注入依赖）
func New(
	cfg *config.Config,
	client *ent.Client,
	rdb *redis.Client,
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

	router := setupRouter(cfg, authHandler, chatHandler, habitHandler, characterHandler)

	return &App{
		router: router,
		cfg:    cfg,
		client: client,
		redis:  rdb,
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
	if a.redis != nil {
		a.redis.Close()
	}
	if a.client != nil {
		a.client.Close()
	}
	return a.server.Shutdown(ctx)
}
