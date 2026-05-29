package middleware

import (
	"github.com/gin-gonic/gin"
)

// RateLimiter 限流中间件（占位，待基于 Redis 实现令牌桶）
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 基于 Redis 实现令牌桶限流
		// 如果触发限流：
		// response.TooManyRequests(c, "请求过于频繁，请稍后重试")
		// return
		c.Next()
	}
}
