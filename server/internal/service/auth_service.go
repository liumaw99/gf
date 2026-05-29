package service

import "context"

// AuthService 认证服务
type AuthService struct{}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{}
}

// AppleSignIn Apple 登录
func (s *AuthService) AppleSignIn(ctx context.Context, idToken string) (string, error) {
	// TODO: 验证 Apple ID Token，生成 JWT
	return "", nil
}

// GoogleSignIn Google 登录
func (s *AuthService) GoogleSignIn(ctx context.Context, idToken string) (string, error) {
	// TODO: 验证 Google ID Token，生成 JWT
	return "", nil
}
