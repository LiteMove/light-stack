package controller

import (
	"strconv"
	"strings"

	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
	roleService service.RoleService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService service.AuthService, roleService service.RoleService) *AuthController {
	return &AuthController{
		authService: authService,
		roleService: roleService,
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数格式错误")
		return
	}

	// 默认租户ID为0（系统租户）
	if req.TenantID == 0 {
		req.TenantID = 0
	}

	loginResp, err := c.authService.Login(&req)
	if err != nil {
		response.Error(ctx, 401, err.Error())
		return
	}

	response.Success(ctx, loginResp)
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数格式错误")
		return
	}

	// 默认租户ID为0（系统租户）
	if req.TenantID == 0 {
		req.TenantID = 0
	}

	user, err := c.authService.Register(&req)
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
	userID, exists := ctx.Get("user_id")
	if !exists {
		response.Unauthorized(ctx, "未授权")
		return
	}

	profile, err := c.authService.GetUserProfile(userID.(uint64))
	if err != nil {
		response.Error(ctx, 404, err.Error())
		return
	}

	response.Success(ctx, profile)
}

// UpdateProfile 更新用户信息
func (c *AuthController) UpdateProfile(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("user_id")
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
	userID, exists := ctx.Get("user_id")
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
	// 对于JWT，登出通常由客户端删除token即可
	// 如果需要服务端记录登出状态，可以将token加入黑名单
	response.Success(ctx, gin.H{"message": "登出成功"})
}
