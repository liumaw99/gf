package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/internal/handler"
	"github.com/lagom/lagom-server/internal/middleware"
)

// setupRouter 配置完整路由
func setupRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	chatHandler *handler.ChatHandler,
	habitHandler *handler.HabitHandler,
	characterHandler *handler.CharacterHandler,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger())
	router.Use(middleware.CORS())

	router.GET("/health", healthHandler)
	router.GET("/ready", readinessHandler)

	api := router.Group("/api/v1")
	{
		// 公开路由
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
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

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "endpoint not found",
		})
	})

	return router
}

// healthHandler 健康检查
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "lagom-server",
	})
}

// readinessHandler 就绪检查
func readinessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"service": "lagom-server",
		"checks": gin.H{
			"database": "ok",
			"redis":    "ok",
		},
	})
}

// requestLogger 请求日志中间件
func requestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s | %3d | %13v | %15s | %-7s %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
		)
	})
}
