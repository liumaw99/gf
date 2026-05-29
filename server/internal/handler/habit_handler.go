package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/internal/pkg/response"
	"github.com/lagom/lagom-server/internal/service"
)

// HabitHandler 习惯处理器
type HabitHandler struct {
	svc *service.HabitService
}

// NewHabitHandler 创建习惯处理器
func NewHabitHandler(svc *service.HabitService) *HabitHandler {
	return &HabitHandler{svc: svc}
}

// List 获取习惯列表
func (h *HabitHandler) List(c *gin.Context) {
	response.Success(c, []any{})
}

// Create 创建习惯
func (h *HabitHandler) Create(c *gin.Context) {
	response.NotImplemented(c, "")
}

// Complete 标记完成
func (h *HabitHandler) Complete(c *gin.Context) {
	response.NotImplemented(c, "")
}
