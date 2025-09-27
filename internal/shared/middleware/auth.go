package middleware

import (
	"strings"

	"github.com/LiteMove/light-stack/pkg/jwt"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// Auth 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未登录或登陆已过期!")
			c.Abort()
			return
		}

		// 检查Bearer格式
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Unauthorized(c, "未登录或登陆已过期!")
			c.Abort()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, "未登录或登陆已过期!")
			c.Abort()
			return
		}

		// 解析token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "未登录或登陆已过期!")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_roles", claims.Roles)
		for _, role := range claims.Roles {
			if role == "super_admin" {
				c.Set("is_super_admin", true)
				break
			}
		}

		c.Next()
	}
}

// OptionalAuth 可选认证中间件
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Next()
			return
		}

		// 检查Bearer格式
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.Next()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			c.Next()
			return
		}

		// 解析token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_roles", claims.Roles)
		for _, role := range claims.Roles {
			if role == "super_admin" {
				c.Set("is_super_admin", true)
				break
			}
		}
		c.Next()
	}
}

// GetUserIDFromContext 从上下文获取用户ID的辅助函数
func GetUserIDFromContext(c *gin.Context) uint64 {
	userID := c.GetUint64("userId")
	return userID
}
