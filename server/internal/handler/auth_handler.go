package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/internal/service"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	svc *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// AppleSignIn Apple 登录
func (h *AuthHandler) AppleSignIn(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "not implemented",
	})
}

// GoogleSignIn Google 登录
func (h *AuthHandler) GoogleSignIn(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "not implemented",
	})
}

// RefreshToken 刷新 Token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "not implemented",
	})
}
