package middleware

import (
	"errors"
	"strings"

	"github.com/LiteMove/light-stack/pkg/jwt"
	"github.com/LiteMove/light-stack/pkg/permission"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// UnifiedAuth 统一权限中间件
func UnifiedAuth(config AuthConfig) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 1. 检查是否需要认证
		if !config.RequireAuth {
			c.Next()
			return
		}

		// 2. JWT认证检查
		userID, err := validateJWT(c)
		if err != nil {
			errorMsg := "未登录或登录已过期"
			if config.ErrorMessage != "" {
				errorMsg = config.ErrorMessage
			}
			response.Unauthorized(c, errorMsg)
			c.Abort()
			return
		}

		// 3. 超管权限检查
		if config.RequireSuperAdmin {
			if !isSuperAdmin(c) {
				errorMsg := "需要超级管理员权限"
				if config.ErrorMessage != "" {
					errorMsg = config.ErrorMessage
				}
				response.Forbidden(c, errorMsg)
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// 4. 普通权限检查
		if !checkPermissions(userID, config, c) {
			errorMsg := "权限不足"
			if config.ErrorMessage != "" {
				errorMsg = config.ErrorMessage
			}
			response.Forbidden(c, errorMsg)
			c.Abort()
			return
		}

		c.Next()
	})
}

// validateJWT 验证JWT并返回用户ID
func validateJWT(c *gin.Context) (uint64, error) {
	// 从请求头中获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, errors.New("token missing")
	}

	// 检查Bearer格式
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return 0, errors.New("invalid token format")
	}

	// 提取token
	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
	if tokenString == "" {
		return 0, errors.New("empty token")
	}

	// 解析token
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}

	// 将用户信息存储到上下文中
	c.Set("userId", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("user_roles", claims.Roles)

	// 设置超管标识
	for _, role := range claims.Roles {
		if role == "super_admin" {
			c.Set("is_super_admin", true)
			break
		}
	}

	return claims.UserID, nil
}

// isSuperAdmin 检查是否为超级管理员
func isSuperAdmin(c *gin.Context) bool {
	return c.GetBool("is_super_admin")
}

// checkPermissions 检查权限
func checkPermissions(userID uint64, config AuthConfig, c *gin.Context) bool {
	// 如果没有权限和角色要求，直接通过
	if len(config.Permissions) == 0 && len(config.Roles) == 0 {
		return true
	}

	hasPermission := false
	hasRole := false

	// 检查权限
	if len(config.Permissions) > 0 {
		hasPermission = permission.Cache.HasAnyPermission(userID, config.Permissions...)
	}

	// 检查角色
	if len(config.Roles) > 0 {
		hasRole = permission.Cache.HasAnyRole(userID, config.Roles...)
	}

	// 根据逻辑判断
	switch config.Logic {
	case AuthLogicAND:
		// AND逻辑：权限和角色都要满足
		permissionResult := len(config.Permissions) == 0 || hasPermission
		roleResult := len(config.Roles) == 0 || hasRole
		return permissionResult && roleResult
	case AuthLogicOR:
		fallthrough
	default:
		// OR逻辑：权限或角色任一满足
		if len(config.Permissions) > 0 && hasPermission {
			return true
		}
		if len(config.Roles) > 0 && hasRole {
			return true
		}
		// 如果有权限或角色要求，但都不满足
		if len(config.Permissions) > 0 || len(config.Roles) > 0 {
			return false
		}
		return true
	}
}

// SimpleAuth 简单权限中间件（兼容现有代码）
func SimpleAuth(permissions ...string) gin.HandlerFunc {
	config := AuthConfig{
		RequireAuth: true,
		Permissions: permissions,
		Logic:       AuthLogicOR,
	}
	return UnifiedAuth(config)
}

// SimpleRoleAuth 简单角色权限中间件（兼容现有代码）
func SimpleRoleAuth(roles ...string) gin.HandlerFunc {
	config := AuthConfig{
		RequireAuth: true,
		Roles:       roles,
		Logic:       AuthLogicOR,
	}
	return UnifiedAuth(config)
}

// AdminAuth 管理员权限中间件（兼容现有代码）
func AdminAuth() gin.HandlerFunc {
	return UnifiedAuth(RequireSuperAdmin)
}
