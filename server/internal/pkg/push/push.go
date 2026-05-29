package push

import "context"

// Service 推送服务（占位实现）
type Service struct{}

// NewService 创建推送服务
func NewService() *Service {
	return &Service{}
}

// Send 发送推送
func (s *Service) Send(ctx context.Context, userID, title, body string) error {
	// TODO: 集成 FCM / APNs
	return nil
}
