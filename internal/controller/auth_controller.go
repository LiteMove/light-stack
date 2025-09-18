package controller

import (
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
	roleService service.RoleService
}

// RoleController 角色控制器
type RoleController struct {
	roleService service.RoleService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService service.AuthService, roleService service.RoleService) *AuthController {
	return &AuthController{
		authService: authService,
		roleService: roleService,
	}
}

// NewRoleController 创建角色控制器
func NewRoleController(roleService service.RoleService) *RoleController {
	return &RoleController{
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

// 角色控制器实现

// CreateRole 创建角色
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req service.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数格式错误")
		return
	}

	// 默认租户ID为0（系统租户）
	if req.TenantID == 0 {
		req.TenantID = 0
	}

	role, err := c.roleService.Create(&req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, role)
}

// UpdateRole 更新角色
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的角色ID")
		return
	}

	var req service.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数格式错误")
		return
	}

	role, err := c.roleService.Update(uint64(roleID), &req)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, role)
}

// DeleteRole 删除角色
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的角色ID")
		return
	}

	err = c.roleService.Delete(uint64(roleID))
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{"message": "角色删除成功"})
}

// GetRole 获取角色信息
func (c *RoleController) GetRole(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		response.BadRequest(ctx, "无效的角色ID")
		return
	}

	role, err := c.roleService.GetByID(uint64(roleID))
	if err != nil {
		response.NotFound(ctx, err.Error())
		return
	}

	response.Success(ctx, role)
}

// GetRoles 获取角色列表
func (c *RoleController) GetRoles(ctx *gin.Context) {
	// 获取查询参数
	tenantIDStr := ctx.DefaultQuery("tenant_id", "0")
	tenantID, _ := strconv.ParseUint(tenantIDStr, 10, 64)

	pageStr := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	statusStr := ctx.DefaultQuery("status", "0")
	status, _ := strconv.Atoi(statusStr)

	roles, total, err := c.roleService.GetList(tenantID, page, pageSize, status)
	if err != nil {
		response.InternalServerError(ctx, "获取角色列表失败")
		return
	}

	response.SuccessWithPage(ctx, roles, total, page, pageSize)
}

// 请求结构体定义

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// AssignRolesRequest 分配角色请求
type AssignRolesRequest struct {
	RoleIDs []uint64 `json:"role_ids" validate:"required"`
}
