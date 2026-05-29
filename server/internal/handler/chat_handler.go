package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/internal/service"
)

// ChatHandler 对话处理器
type ChatHandler struct {
	svc *service.AIService
}

// NewChatHandler 创建对话处理器
func NewChatHandler(svc *service.AIService) *ChatHandler {
	return &ChatHandler{svc: svc}
}

// Chat 流式对话（SSE）
func (h *ChatHandler) Chat(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.SSEvent("delta", "{\"content\":\"嗨\"}")
	c.SSEvent("done", "")
}

// GetHistory 获取对话历史
func (h *ChatHandler) GetHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []any{},
	})
}
