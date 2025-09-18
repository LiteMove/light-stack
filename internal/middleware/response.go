package middleware

import (
	"fmt"
	"github.com/LiteMove/light-stack/pkg/jwt"
	"github.com/LiteMove/light-stack/pkg/response"
	"strings"
	"time"

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
			response.Unauthorized(c, "缺少Authorization头")
			c.Abort()
			return
		}

		// 检查Bearer格式
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			response.Unauthorized(c, "无效的Authorization格式")
			c.Abort()
			return
		}

		// 提取token
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			response.Unauthorized(c, "Token不能为空")
			c.Abort()
			return
		}

		// 解析token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "无效的token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminAuthMiddleware 管理员权限中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Unauthorized(c, "未授权")
			c.Abort()
			return
		}

		// 检查是否为管理员
		if userRole.(string) != "super_admin" {
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
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
