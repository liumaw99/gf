package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/internal/pkg/response"
)

// JWTAuth JWT 认证中间件
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format")
			return
		}

		token := parts[1]

		// TODO: 验证 JWT
		_ = token
		_ = secret

		// 临时跳过验证，开发阶段
		c.Set("user_id", "dev-user-id")
		c.Next()
	}
}
