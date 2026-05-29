package service

import (
	"context"
	"strings"

	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/ent"
	"github.com/lagom/lagom-server/internal/pkg/ai"
)

// AIService AI 对话服务
type AIService struct {
	cfg       *config.Config
	client    *ent.Client
	aiClient  *ai.DeepSeekClient
}

// NewAIService 创建 AI 服务
func NewAIService(cfg *config.Config, client *ent.Client, aiClient *ai.DeepSeekClient) *AIService {
	return &AIService{cfg: cfg, client: client, aiClient: aiClient}
}

// Chat 流式对话
func (s *AIService) Chat(ctx context.Context, userID string, characterID *string, userMessage string) (<-chan string, error) {
	ch := make(chan string, 10)
	close(ch)
	return ch, nil
}

// BuildSystemPrompt 构建 System Prompt
func (s *AIService) BuildSystemPrompt(ctx context.Context, userID string, characterID *string) string {
	var prompt strings.Builder
	prompt.WriteString("你是 Lagom 的 AI 伙伴，帮助用户养成自律习惯。\n")
	prompt.WriteString("回复要温暖、简短（不超过100字），善用emoji增加温度。\n")
	prompt.WriteString("失败时给予共情和鼓励，绝不批评。\n\n")
	return prompt.String()
}
