package middleware

import (
	"github.com/LiteMove/light-stack/pkg/permission"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// CheckPermission 权限验证中间件 - 检查用户是否有指定权限码
// 用法: middleware.CheckPermission("system:user:create", "system:user:update")
func CheckPermission(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查是否是超级管理员 - 超级管理员拥有所有权限
		if c.GetBool("is_super_admin") {
			c.Next()
			return
		}

		// 检查权限 - 只要有其中一个权限即可
		if permission.Cache.HasAnyPermission(userID, codes...) {
			c.Next()
			return
		}

		response.Forbidden(c, "没有操作权限")
		c.Abort()
	}
}

// CheckAllPermissions 检查用户是否拥有所有指定权限 - 必须拥有全部权限
// 用法: middleware.CheckAllPermissions("system:user:read", "system:user:list")
func CheckAllPermissions(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查是否是超级管理员 - 超级管理员拥有所有权限
		if c.GetBool("is_super_admin") {
			c.Next()
			return
		}

		// 逐一检查是否拥有全部权限
		for _, code := range codes {
			if !permission.Cache.HasAnyPermission(userID, code) {
				response.Forbidden(c, "权限不足")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// CheckRole 角色验证中间件 - 检查用户是否有指定角色
// 用法: middleware.CheckRole("admin", "super_admin")
func CheckRole(roleCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查角色 - 只要有其中一个角色即可
		if permission.Cache.HasAnyRole(userID, roleCodes...) {
			c.Next()
			return
		}

		response.Forbidden(c, "没有角色权限")
		c.Abort()
	}
}

// CheckAllRoles 检查用户是否拥有所有指定角色 - 必须拥有全部角色
// 用法: middleware.CheckAllRoles("admin", "editor")
func CheckAllRoles(roleCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 逐一检查是否拥有全部角色
		for _, roleCode := range roleCodes {
			if !permission.Cache.HasAnyRole(userID, roleCode) {
				response.Forbidden(c, "角色权限不足")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// CheckPermissionOrRole 权限或角色验证中间件 - 满足权限或角色任一条件即可
// 用法: middleware.CheckPermissionOrRole([]string{"system:user:create"}, []string{"admin"})
func CheckPermissionOrRole(permissions []string, roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查是否是超级管理员 - 超级管理员拥有所有权限
		if c.GetBool("is_super_admin") {
			c.Next()
			return
		}

		// 先检查权限
		if len(permissions) > 0 && permission.Cache.HasAnyPermission(userID, permissions...) {
			c.Next()
			return
		}

		// 再检查角色
		if len(roles) > 0 && permission.Cache.HasAnyRole(userID, roles...) {
			c.Next()
			return
		}

		response.Forbidden(c, "权限不足")
		c.Abort()
	}
}

// CheckPermissionAndRole 权限和角色验证中间件 - 必须同时满足权限和角色条件
// 用法: middleware.CheckPermissionAndRole([]string{"system:user:create"}, []string{"admin"})
func CheckPermissionAndRole(permissions []string, roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查是否是超级管理员 - 超级管理员拥有所有权限
		if c.GetBool("is_super_admin") {
			c.Next()
			return
		}

		// 必须同时满足权限和角色
		hasPermission := len(permissions) == 0 || permission.Cache.HasAnyPermission(userID, permissions...)
		hasRole := len(roles) == 0 || permission.Cache.HasAnyRole(userID, roles...)

		if hasPermission && hasRole {
			c.Next()
			return
		}

		response.Forbidden(c, "权限和角色验证失败")
		c.Abort()
	}
}

// SuperAdmin 超级管理员验证中间件 - 快捷方式
func SuperAdmin() gin.HandlerFunc {
	return CheckRole("super_admin")
}

// Admin 管理员验证中间件 - 快捷方式（包含超级管理员）
func Admin() gin.HandlerFunc {
	return CheckRole("admin", "super_admin")
}

// Owner 资源拥有者验证中间件 - 验证用户是否是资源的拥有者
// 用法: middleware.Owner("user_id") - 检查路径参数中的user_id是否等于当前用户ID
func Owner(paramKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查是否是超级管理员 - 超级管理员拥有所有权限
		if c.GetBool("is_super_admin") {
			c.Next()
			return
		}

		// 获取路径参数中的资源拥有者ID
		resourceUserID := c.GetUint64(paramKey)
		if resourceUserID == 0 {
			response.BadRequest(c, "无效的资源ID")
			c.Abort()
			return
		}

		// 检查是否是资源拥有者
		if userID == resourceUserID {
			c.Next()
			return
		}

		response.Forbidden(c, "只能操作自己的资源")
		c.Abort()
	}
}

// OwnerOrRole 资源拥有者或指定角色验证中间件
// 用法: middleware.OwnerOrRole("user_id", "admin", "moderator")
func OwnerOrRole(paramKey string, roleCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetUint64("userId")
		if userID == 0 {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查是否是超级管理员 - 超级管理员拥有所有权限
		if c.GetBool("is_super_admin") {
			c.Next()
			return
		}

		// 获取路径参数中的资源拥有者ID
		resourceUserID := c.GetUint64(paramKey)
		if resourceUserID == 0 {
			response.BadRequest(c, "无效的资源ID")
			c.Abort()
			return
		}

		// 检查是否是资源拥有者
		if userID == resourceUserID {
			c.Next()
			return
		}

		// 检查是否有指定角色
		if permission.Cache.HasAnyRole(userID, roleCodes...) {
			c.Next()
			return
		}

		response.Forbidden(c, "权限不足")
		c.Abort()
	}
}
