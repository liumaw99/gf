package middleware

import "github.com/gin-gonic/gin"

// RateLimiter 限流中间件（占位实现）
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 基于 Redis 实现令牌桶限流
		c.Next()
	}
}
