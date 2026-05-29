package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []any{},
	})
}

// Create 创建习惯
func (h *HabitHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "not implemented",
	})
}

// Complete 标记完成
func (h *HabitHandler) Complete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "not implemented",
	})
}
