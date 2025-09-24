package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/LiteMove/light-stack/pkg/jwt"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
)

// ResponseMiddleware 响应中间件
func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 处理panic恢复
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			response.InternalServerError(c, err.Error())
			return
		}
	}
}

// RequestLogMiddleware 请求日志中间件
func RequestLogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
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

// SuperAdminAuthMiddleware 超级管理员权限中间件
func SuperAdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin := c.GetBool("is_super_admin")

		if !isAdmin {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalJWTAuthMiddleware 可选的JWT认证中间件（不强制要求认证）
func OptionalJWTAuthMiddleware() gin.HandlerFunc {
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
