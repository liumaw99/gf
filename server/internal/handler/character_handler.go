package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/internal/pkg/response"
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
	response.Success(c, []any{})
}

// GetDetail 获取角色详情
func (h *CharacterHandler) GetDetail(c *gin.Context) {
	response.Success(c, gin.H{})
}

// Recommend 角色推荐
func (h *CharacterHandler) Recommend(c *gin.Context) {
	response.Success(c, gin.H{
		"recommended_character_id": "haruki_murakami_v1",
		"match_score":              0.85,
	})
}

// GetUserCharacters 获取用户角色配置
func (h *CharacterHandler) GetUserCharacters(c *gin.Context) {
	response.Success(c, []any{})
}

// SelectCharacter 选择角色
func (h *CharacterHandler) SelectCharacter(c *gin.Context) {
	response.NotImplemented(c, "")
}

// SwitchCharacter 切换角色
func (h *CharacterHandler) SwitchCharacter(c *gin.Context) {
	response.NotImplemented(c, "")
}
