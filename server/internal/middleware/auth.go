package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/internal/pkg/jwt"
	"github.com/lagom/lagom-server/internal/pkg/response"
)

// JWTAuth JWT 认证中间件
func JWTAuth(mgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := mgr.Validate(token)
		if err != nil {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
