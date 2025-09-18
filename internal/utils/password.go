package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost bcrypt默认代价
	DefaultCost = bcrypt.DefaultCost
	// MinCost bcrypt最小代价
	MinCost = bcrypt.MinCost
	// MaxCost bcrypt最大代价
	MaxCost = bcrypt.MaxCost
)

// HashPassword 使用bcrypt对密码进行哈希加密
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// VerifyPassword 验证密码是否正确
func VerifyPassword(hashedPassword, password string) bool {
	if hashedPassword == "" || password == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// HashPasswordWithCost 使用指定代价对密码进行哈希加密
func HashPasswordWithCost(password string, cost int) (string, error) {
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	if cost < MinCost || cost > MaxCost {
		return "", fmt.Errorf("invalid cost: must be between %d and %d", MinCost, MaxCost)
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be positive")
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// GenerateResetToken 生成密码重置令牌
func GenerateResetToken() (string, error) {
	return GenerateRandomString(32)
}

// SecureCompare 安全比较两个字符串（防止时序攻击）
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// ValidatePasswordStrength 验证密码强度
func ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	if len(password) > 128 {
		return errors.New("password must be less than 128 characters long")
	}

	hasLower := false
	hasDigit := false

	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	var requirements []string
	if !hasLower {
		requirements = append(requirements, "lowercase letter")
	}
	if !hasDigit {
		requirements = append(requirements, "number")
	}

	if len(requirements) > 0 {
		return fmt.Errorf("password must contain at least one %s", strings.Join(requirements, " and "))
	}

	return nil
}

// IsPasswordExpired 检查密码是否过期（可根据业务需求实现）
func IsPasswordExpired(lastChanged int64, expiryDays int) bool {
	if expiryDays <= 0 {
		return false // 不设置过期时间
	}

	// 这里可以根据实际需求实现密码过期逻辑
	// 例如：检查lastChanged时间戳与当前时间的差值
	return false
}
