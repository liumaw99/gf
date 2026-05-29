package service

import "context"

// HabitService 习惯服务
type HabitService struct{}

// NewHabitService 创建习惯服务
func NewHabitService() *HabitService {
	return &HabitService{}
}

// GeneratePlan 生成习惯计划
func (s *HabitService) GeneratePlan(ctx context.Context, goal string) error {
	// TODO: 调用 AI 生成 21 天渐进计划
	return nil
}

// CalculateHSI 计算习惯强度指数
func (s *HabitService) CalculateHSI(completionRate, qualityScore, moodScore float64) float64 {
	return completionRate*0.4 + qualityScore*0.3 + moodScore*0.3
}
