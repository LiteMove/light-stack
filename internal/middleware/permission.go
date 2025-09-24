package middleware

import (
	"github.com/LiteMove/light-stack/pkg/permission"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// CheckPermission 权限验证中间件
func CheckPermission(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查权限
		if permission.Cache.HasAnyPermission(userID, codes...) {
			c.Next()
			return
		}

		response.Forbidden(c, "没有操作权限")
		c.Abort()
	}
}

// CheckRole 角色验证中间件
func CheckRole(roleCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查角色
		if permission.Cache.HasAnyRole(userID, roleCodes...) {
			c.Next()
			return
		}

		response.Forbidden(c, "没有角色权限")
		c.Abort()
	}
}

// CheckPermissionOrRole 权限或角色验证中间件（满足其一即可）
func CheckPermissionOrRole(permissions []string, roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查权限
		if len(permissions) > 0 && permission.Cache.HasAnyPermission(userID, permissions...) {
			c.Next()
			return
		}

		// 检查角色
		if len(roles) > 0 && permission.Cache.HasAnyRole(userID, roles...) {
			c.Next()
			return
		}

		response.Forbidden(c, "没有权限")
		c.Abort()
	}
}
