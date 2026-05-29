package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lagom/lagom-server/config"
	"github.com/lagom/lagom-server/ent"
	"github.com/lagom/lagom-server/ent/user"
	"github.com/lagom/lagom-server/internal/pkg/errors"
	"github.com/lagom/lagom-server/internal/pkg/jwt"
	"github.com/lagom/lagom-server/internal/pkg/password"
	"github.com/lagom/lagom-server/internal/pkg/redis"
)

// TokenPair 令牌对
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // access token 有效期（秒）
}

// AuthService 认证服务
type AuthService struct {
	cfg    *config.Config
	client *ent.Client
	jwtMgr *jwt.Manager
	rdb    *redis.Client
}

// NewAuthService 创建认证服务
func NewAuthService(cfg *config.Config, client *ent.Client, jwtMgr *jwt.Manager, rdb *redis.Client) *AuthService {
	return &AuthService{cfg: cfg, client: client, jwtMgr: jwtMgr, rdb: rdb}
}

// Register 邮箱密码注册
func (s *AuthService) Register(ctx context.Context, email, pwd, name string) (*ent.User, *TokenPair, error) {
	// 检查邮箱是否已存在
	exists, err := s.client.User.Query().Where(user.EmailEQ(email)).Exist(ctx)
	if err != nil {
		return nil, nil, errors.WrapInternal(err)
	}
	if exists {
		return nil, nil, errors.ErrEmailExists
	}

	// hash 密码
	hash, err := password.Hash(pwd)
	if err != nil {
		return nil, nil, errors.WrapInternal(err)
	}

	// 创建用户
	builder := s.client.User.Create().
		SetEmail(email).
		SetPasswordHash(hash)
	if name != "" {
		builder.SetName(name)
	}
	u, err := builder.Save(ctx)
	if err != nil {
		return nil, nil, errors.WrapInternal(err)
	}

	// 签发 token
	pair, err := s.issueTokenPair(u.ID.String())
	if err != nil {
		return nil, nil, err
	}

	return u, pair, nil
}

// Login 邮箱密码登录
func (s *AuthService) Login(ctx context.Context, email, pwd string) (*ent.User, *TokenPair, error) {
	u, err := s.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil, errors.ErrInvalidCredentials
		}
		return nil, nil, errors.WrapInternal(err)
	}

	if u.PasswordHash == nil || *u.PasswordHash == "" {
		return nil, nil, errors.ErrInvalidCredentials
	}

	if err := password.Verify(pwd, *u.PasswordHash); err != nil {
		return nil, nil, errors.ErrInvalidCredentials
	}

	pair, err := s.issueTokenPair(u.ID.String())
	if err != nil {
		return nil, nil, err
	}

	return u, pair, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.jwtMgr.Validate(refreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// 检查黑名单
	blacklistKey := fmt.Sprintf("jwt:blacklist:%s", claims.JTI)
	exists, err := s.rdb.Exists(ctx, blacklistKey).Result()
	if err != nil {
		return nil, errors.WrapInternal(err)
	}
	if exists > 0 {
		return nil, errors.ErrTokenRevoked
	}

	// 将旧 refresh token 加入黑名单（TTL = 剩余有效期）
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl > 0 {
		if err := s.rdb.Set(ctx, blacklistKey, "1", ttl).Err(); err != nil {
			return nil, errors.WrapInternal(err)
		}
	}

	return s.issueTokenPair(claims.UserID)
}

// issueTokenPair 签发 access + refresh token
func (s *AuthService) issueTokenPair(userID string) (*TokenPair, error) {
	accessTTL := time.Duration(s.cfg.JWT.AccessTTL) * time.Minute
	refreshTTL := time.Duration(s.cfg.JWT.RefreshTTL) * time.Hour

	accessToken, err := s.jwtMgr.Generate(userID, accessTTL)
	if err != nil {
		return nil, errors.WrapInternal(err)
	}

	refreshToken, err := s.jwtMgr.Generate(userID, refreshTTL)
	if err != nil {
		return nil, errors.WrapInternal(err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTTL.Seconds()),
	}, nil
}

