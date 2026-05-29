package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// setupRouter 配置路由
func setupRouter(log *slog.Logger) *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(gin.Recovery())
	router.Use(requestLogger(log))
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
		"status":  "healthy",
		"service": "lagom-server",
	})
}

// readinessHandler 就绪检查
func readinessHandler(c *gin.Context) {
	// TODO: 检查数据库、Redis 等依赖
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
func requestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		attrs := []any{
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP(),
			"latency", time.Since(start),
			"bytes", c.Writer.Size(),
		}

		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		if c.Writer.Status() >= http.StatusInternalServerError {
			log.Error("request completed", attrs...)
			return
		}

		log.Info("request completed", attrs...)
	}
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
