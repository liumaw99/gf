package app

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/ent"
	"github.com/lagom/lagom-server/internal/handler"
	"github.com/lagom/lagom-server/internal/middleware"
	"github.com/lagom/lagom-server/internal/pkg/redis"
	"github.com/lagom/lagom-server/internal/pkg/response"
)

// setupRouter 配置完整路由和中间件链
// 中间件顺序（Gin 执行顺序）：Recovery → Logger → RateLimiter → CORS → (JWTAuth)
func setupRouter(
	cfg *config.Config,
	client *ent.Client,
	rdb *redis.Client,
	authHandler *handler.AuthHandler,
	chatHandler *handler.ChatHandler,
	habitHandler *handler.HabitHandler,
	characterHandler *handler.CharacterHandler,
) *gin.Engine {
	router := gin.New()

	// 1. Recovery — 必须最前，捕获所有 panic，防止服务崩溃
	router.Use(gin.Recovery())

	// 2. 结构化请求日志
	router.Use(middleware.Logger())

	// 3. 限流（公开路由层）— 基于 Redis 的占位实现
	router.Use(middleware.RateLimiter())

	// 4. CORS — 处理跨域和 OPTIONS 预检
	router.Use(middleware.CORS())

	// 5. 安全响应头中间件
	router.Use(middleware.SecurityHeaders())

	// --- 公开路由 ---
	router.GET("/health", healthHandler)
	router.GET("/ready", readinessHandler(client, rdb))

	api := router.Group("/api/v1")
	{
		// 公开 API
		api.GET("/ping", func(c *gin.Context) {
			response.Success(c, gin.H{
				"message": "pong",
				"time":    time.Now().Unix(),
			})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/apple", authHandler.AppleSignIn)
			auth.POST("/google", authHandler.GoogleSignIn)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			authorized.POST("/chat", chatHandler.Chat)
			authorized.GET("/chat/history", chatHandler.GetHistory)

			authorized.GET("/habits", habitHandler.List)
			authorized.POST("/habits", habitHandler.Create)
			authorized.POST("/habits/:id/complete", habitHandler.Complete)

			authorized.GET("/characters", characterHandler.List)
			authorized.GET("/characters/:id", characterHandler.GetDetail)
			authorized.POST("/characters/recommend", characterHandler.Recommend)
			authorized.GET("/user/characters", characterHandler.GetUserCharacters)
			authorized.POST("/user/characters", characterHandler.SelectCharacter)
			authorized.POST("/user/characters/:id/switch", characterHandler.SwitchCharacter)
		}
	}

	// 404 处理
	router.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "endpoint not found")
	})

	return router
}

// healthHandler 存活检查 — 只要进程在运行即返回健康
func healthHandler(c *gin.Context) {
	response.Success(c, gin.H{
		"status":    "healthy",
		"service":   "lagom-server",
		"timestamp": time.Now().Unix(),
	})
}

// readinessHandler 就绪检查 — 验证下游依赖（DB、Redis）是否可用
func readinessHandler(client *ent.Client, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		checks := gin.H{}
		allOK := true

		// 检查数据库
		if client != nil {
			_, err := client.User.Query().Limit(1).All(ctx)
			if err != nil {
				checks["database"] = gin.H{"status": "error", "error": err.Error()}
				allOK = false
			} else {
				checks["database"] = "ok"
			}
		} else {
			checks["database"] = gin.H{"status": "skipped", "reason": "client not initialized"}
		}

		// 检查 Redis
		if rdb != nil {
			if err := rdb.Ping(ctx).Err(); err != nil {
				checks["redis"] = gin.H{"status": "error", "error": err.Error()}
				allOK = false
			} else {
				checks["redis"] = "ok"
			}
		} else {
			checks["redis"] = gin.H{"status": "skipped", "reason": "client not initialized"}
		}

		if !allOK {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "not_ready",
				"service": "lagom-server",
				"checks":  checks,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ready",
			"service": "lagom-server",
			"checks":  checks,
		})
	}
}
