package service

import (
	"context"

	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/ent"
)

// CharacterService 名人角色服务
type CharacterService struct {
	cfg    *config.Config
	client *ent.Client
}

// NewCharacterService 创建角色服务
func NewCharacterService(cfg *config.Config, client *ent.Client) *CharacterService {
	return &CharacterService{cfg: cfg, client: client}
}

// RecommendCharacter 推荐角色
func (s *CharacterService) RecommendCharacter(ctx context.Context, answers []string) (string, float64, error) {
	return "haruki_murakami_v1", 0.85, nil
}
