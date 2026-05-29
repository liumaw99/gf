package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash 使用 bcrypt 对密码进行哈希
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

// Verify 验证密码与哈希是否匹配
func Verify(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
