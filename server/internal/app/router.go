package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/config"
)

// setupRouter 配置路由
func setupRouter(cfg *config.Config) *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(gin.Recovery())
	router.Use(requestLogger())
	router.Use(corsMiddleware())

	// 健康检查
	router.GET("/health", healthHandler)
	router.GET("/ready", readinessHandler)

	// API v1
	api := router.Group("/api/v1")
	{
		// 公开路由
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		// TODO: 认证路由
		// auth := api.Group("/auth")
		// {
		//     auth.POST("/apple", authHandler.AppleSignIn)
		//     auth.POST("/google", authHandler.GoogleSignIn)
		// }

		// TODO: 需要认证的路由
		// authorized := api.Group("")
		// authorized.Use(middleware.JWTAuth(cfg.JWT.Secret))
		// {
		//     authorized.POST("/chat", chatHandler.Chat)
		//     authorized.GET("/habits", habitHandler.List)
		// }
	}

	// 404 处理
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
		"status": "healthy",
		"service": "lagom-server",
	})
}

// readinessHandler 就绪检查
func readinessHandler(c *gin.Context) {
	// TODO: 检查数据库、Redis 等依赖
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
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

// corsMiddleware CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
