package service

import "context"

// CharacterService 名人角色服务
type CharacterService struct{}

// NewCharacterService 创建角色服务
func NewCharacterService() *CharacterService {
	return &CharacterService{}
}

// RecommendCharacter 推荐角色
func (s *CharacterService) RecommendCharacter(ctx context.Context, answers []string) (string, float64, error) {
	// TODO: 基于答案匹配角色
	return "haruki_murakami_v1", 0.85, nil
}
