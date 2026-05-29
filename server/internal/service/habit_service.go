package service

import (
	"context"

	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/ent"
)

// HabitService 习惯服务
type HabitService struct {
	cfg    *config.Config
	client *ent.Client
}

// NewHabitService 创建习惯服务
func NewHabitService(cfg *config.Config, client *ent.Client) *HabitService {
	return &HabitService{cfg: cfg, client: client}
}

// GeneratePlan 生成习惯计划
func (s *HabitService) GeneratePlan(ctx context.Context, goal string) error {
	return nil
}

// CalculateHSI 计算习惯强度指数
func (s *HabitService) CalculateHSI(completionRate, qualityScore, moodScore float64) float64 {
	return completionRate*0.4 + qualityScore*0.3 + moodScore*0.3
}
