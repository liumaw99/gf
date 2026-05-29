package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HabitHandler 习惯处理器
type HabitHandler struct {
	// TODO: inject habit service
}

// NewHabitHandler 创建习惯处理器
func NewHabitHandler() *HabitHandler {
	return &HabitHandler{}
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
