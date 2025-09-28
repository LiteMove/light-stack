package controller

import (
	"strconv"
	"strings"

	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/permission"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
	roleService service.RoleService
	menuService service.MenuService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService service.AuthService,
	roleService service.RoleService,
	menuService service.MenuService) *AuthController {
	return &AuthController{
		authService: authService,
		roleService: roleService,
		menuService: menuService,
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数格式错误")
		return
	}
	// 从上下文获取租户ID
	tenantID, exists := middleware.GetTenantIDFromContext(ctx)
	if !exists {
		tenantID = uint64(1) // 默认系统租户
	}

	tokenResp, err := c.authService.Login(tenantID, &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, tokenResp)
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数格式错误")
		return
	}

	// 从上下文获取租户ID
	tenantID, exists := middleware.GetTenantIDFromContext(ctx)
	if !exists {
		tenantID = uint64(1) // 默认系统租户
	}
	user, err := c.authService.Register(tenantID, &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, user)
}

// RefreshToken 刷新token
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	// 从请求头中获取token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		response.Unauthorized(ctx, "缺少Authorization头")
		return
	}

	// 提取token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		response.Unauthorized(ctx, "无效的Authorization格式")
		return
	}

	tokenResp, err := c.authService.RefreshToken(tokenString)
	if err != nil {
		response.Unauthorized(ctx, err.Error())
		return
	}

	response.Success(ctx, tokenResp)
}

// GetProfile 获取用户信息
func (c *AuthController) GetProfile(ctx *gin.Context) {
	// 从上下文中获取用户ID（由JWT中间件设置）
	userId := ctx.GetUint64("userId")
	if userId == 0 {
		response.Unauthorized(ctx, "未授权")
		return
	}

	profile, err := c.authService.GetUserProfile(userId)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	profile.Menus, err = c.menuService.GetUserMenuTree(userId)
	if err != nil {
		response.BadRequest(ctx, "获取用户菜单树失败")
		return
	}

	// 从权限缓存中获取权限列表
	perms, exists := permission.Cache.GetUserPermissions(userId)
	if exists {
		permList := make([]string, 0, len(perms))
		for p := range perms {
			permList = append(permList, p)
		}
		profile.Permissions = permList
	} else {
		// 如果缓存中没有，从数据库获取并缓存
		profile.Permissions, err = c.menuService.GetMenuPermissions(userId)
		if err != nil {
			response.BadRequest(ctx, "获取用户权限失败")
			return
		}
		// 加载到缓存
		permission.Cache.LoadUserPermissions(userId, profile.Permissions)
	}

	response.Success(ctx, profile)
}

// UpdateProfile 更新用户信息
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userId")
	if !exists {
		response.Unauthorized(ctx, "未授权")
		return
	}

	var req service.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数格式错误")
		return
	}

	profile, err := c.authService.UpdateUserProfile(userID.(uint64), &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, profile)
}

// ChangePassword 修改密码
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userId")
	if !exists {
		response.Unauthorized(ctx, "未授权")
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数格式错误")
		return
	}

	err := c.authService.ChangePassword(userID.(uint64), req.OldPassword, req.NewPassword)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{"message": "密码修改成功"})
}

// AssignUserRoles 为用户分配角色
func (c *AuthController) AssignUserRoles(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的用户ID")
		return
	}

	var req AssignRolesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数格式错误")
		return
	}

	err = c.authService.AssignUserRoles(uint64(userID), req.RoleIDs)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{"message": "角色分配成功"})
}

// GetUserRoles 获取用户角色
func (c *AuthController) GetUserRoles(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的用户ID")
		return
	}

	roles, err := c.authService.GetUserRoles(uint64(userID))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, roles)
}

// Logout 用户登出
func (c *AuthController) Logout(ctx *gin.Context) {
	// 获取用户ID并清除权限缓存
	userID := ctx.GetUint64("userId")
	if userID != 0 {
		permission.ClearUserPermissions(userID)
	}

	// 对于JWT，登出通常由客户端删除token即可
	// 如果需要服务端记录登出状态，可以将token加入黑名单
	response.Success(ctx, gin.H{"message": "登出成功"})
}
