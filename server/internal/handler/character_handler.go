package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/internal/service"
)

// CharacterHandler 名人角色处理器
type CharacterHandler struct {
	svc *service.CharacterService
}

// NewCharacterHandler 创建角色处理器
func NewCharacterHandler(svc *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{svc: svc}
}

// List 获取角色列表
func (h *CharacterHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []any{},
	})
}

// GetDetail 获取角色详情
func (h *CharacterHandler) GetDetail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    gin.H{},
	})
}

// Recommend 角色推荐
func (h *CharacterHandler) Recommend(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"recommended_character_id": "haruki_murakami_v1",
			"match_score":              0.85,
		},
	})
}

// GetUserCharacters 获取用户角色配置
func (h *CharacterHandler) GetUserCharacters(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []any{},
	})
}

// SelectCharacter 选择角色
func (h *CharacterHandler) SelectCharacter(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "not implemented",
	})
}

// SwitchCharacter 切换角色
func (h *CharacterHandler) SwitchCharacter(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"code":    501,
		"message": "not implemented",
	})
}
