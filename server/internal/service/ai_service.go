package service

import "context"

// AIService AI 对话服务
type AIService struct{}

// NewAIService 创建 AI 服务
func NewAIService() *AIService {
	return &AIService{}
}

// Chat 流式对话
func (s *AIService) Chat(ctx context.Context, userID string, characterID *string, message string) (<-chan string, error) {
	ch := make(chan string, 10)
	close(ch)
	return ch, nil
}

// BuildSystemPrompt 构建 System Prompt
func (s *AIService) BuildSystemPrompt(ctx context.Context, userID string, characterID *string) string {
	// TODO: 根据用户画像和角色档案组装 Prompt
	return "你是 Lagom 的 AI 伙伴，帮助用户养成自律习惯。"
}
