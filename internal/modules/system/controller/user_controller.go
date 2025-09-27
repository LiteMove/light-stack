package controller

import (
	"strconv"

	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserController 用户控制器
type UserController struct {
	userService service.UserService
	validator   *validator.Validate
}

// NewUserController 创建用户控制器
func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
		validator:   validator.New(),
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Nickname string `json:"nickname" validate:"required,min=1,max=100"`
	Email    string `json:"email" validate:"omitempty,email,max=255"`
	Phone    string `json:"phone" validate:"omitempty,max=20"`
	Avatar   string `json:"avatar" validate:"omitempty,max=255"`
	Password string `json:"password" validate:"omitempty,min=6,max=50"`
	Status   int    `json:"status" validate:"required,oneof=1 2"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Nickname string `json:"nickname" validate:"required,min=1,max=100"`
	Email    string `json:"email" validate:"omitempty,email,max=255"`
	Phone    string `json:"phone" validate:"omitempty,max=20"`
	Avatar   string `json:"avatar" validate:"omitempty,max=255"`
	Status   int    `json:"status" validate:"required,oneof=1 2"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status" validate:"oneof=0 1 2"`
	RoleID   uint64 `form:"roleId"`
}

// UpdateUserStatusRequest 更新用户状态请求
type UpdateUserStatusRequest struct {
	Status int `json:"status" validate:"required,oneof=1 2"`
}

// BatchUpdateUserStatusRequest 批量更新用户状态请求
type BatchUpdateUserStatusRequest struct {
	IDs    []uint64 `json:"ids" validate:"required,min=1"`
	Status int      `json:"status" validate:"required,oneof=1 2"`
}

// ChangeUserPasswordRequest 修改用户密码请求
type ChangeUserPasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,min=6,max=50"`
	NewPassword string `json:"newPassword" validate:"required,min=6,max=50"`
}

// AssignUserRolesRequest 分配用户角色请求
type AssignUserRolesRequest struct {
	RoleIDs []uint64 `json:"roleIds" validate:"required"`
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.Error(ctx, 400, "参数验证失败: "+err.Error())
		return
	}

	// 从上下文获取租户ID
	tenantID, exists := middleware.GetTenantIDFromContext(ctx)
	if !exists {
		tenantID = uint64(1) // 默认系统租户
	}

	// 创建用户对象
	user := &model.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Password: req.Password,
		Status:   req.Status,
		IsSystem: false,
	}

	// 处理可选字段
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	// 设置头像
	user.Avatar = req.Avatar

	user.TenantID = tenantID

	// 调用服务创建用户
	if err := c.userService.CreateUser(user); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
	})
}

// GetUsers 获取用户列表
func (c *UserController) GetUsers(ctx *gin.Context) {
	var req UserListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, 400, "请求参数格式错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.Error(ctx, 400, "参数验证失败: "+err.Error())
		return
	}

	// 从上下文获取租户ID
	tenantID, exists := middleware.GetTenantIDFromContext(ctx)
	if !exists {
		tenantID = uint64(1) // 默认系统租户
	}

	// 调用服务获取用户列表
	users, total, err := c.userService.GetUserList(tenantID, req.Page, req.PageSize, req.Keyword, req.Status, req.RoleID)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	// 返回分页数据
	response.Success(ctx, gin.H{
		"list":      users,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetUser 获取用户详情
func (c *UserController) GetUser(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, 400, "用户ID格式错误")
		return
	}

	// 调用服务获取用户
	user, err := c.userService.GetUserWithRoles(id)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, user)
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, 400, "用户ID格式错误")
		return
	}

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.Error(ctx, 400, "参数验证失败: "+err.Error())
		return
	}

	// 获取原用户信息
	existingUser, err := c.userService.GetUser(id)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	// 不可禁用超级管理员/系统用户
	if id == model.SuperAdminId || existingUser.IsSystem {
		response.BadRequest(ctx, "不可禁用超级管理员/系统用户")
		return
	}

	// 更新用户信息
	existingUser.Username = req.Username
	existingUser.Nickname = req.Nickname
	existingUser.Status = req.Status

	// 处理可选字段
	if req.Email != "" {
		existingUser.Email = &req.Email
	} else {
		existingUser.Email = nil
	}
	if req.Phone != "" {
		existingUser.Phone = &req.Phone
	} else {
		existingUser.Phone = nil
	}
	// 更新头像
	existingUser.Avatar = req.Avatar

	// 调用服务更新用户
	if err := c.userService.UpdateUser(existingUser); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"id":       existingUser.ID,
		"username": existingUser.Username,
		"nickname": existingUser.Nickname,
	})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, 400, "用户ID格式错误")
		return
	}

	// 调用服务删除用户
	if err := c.userService.DeleteUser(id); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "删除成功",
	})
}

// UpdateUserStatus 更新用户状态
func (c *UserController) UpdateUserStatus(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, 400, "用户ID格式错误")
		return
	}

	var req UpdateUserStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.Error(ctx, 400, "参数验证失败: "+err.Error())
		return
	}

	// 调用服务更新状态
	if err := c.userService.UpdateUserStatus(id, req.Status); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "状态更新成功",
	})
}

// BatchUpdateUserStatus 批量更新用户状态
func (c *UserController) BatchUpdateUserStatus(ctx *gin.Context) {
	var req BatchUpdateUserStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.Error(ctx, 400, "参数验证失败: "+err.Error())
		return
	}

	// 调用服务批量更新状态
	if err := c.userService.BatchUpdateUserStatus(req.IDs, req.Status); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "批量状态更新成功",
	})
}

// ResetPassword 重置密码
func (c *UserController) ResetPassword(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "用户ID格式错误")
		return
	}

	// 调用服务重置密码
	newPassword, err := c.userService.ResetPassword(id)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message":     "密码重置成功",
		"newPassword": newPassword,
	})
}

// AssignUserRoles 为用户分配角色
func (c *UserController) AssignUserRoles(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "用户ID格式错误")
		return
	}

	var req AssignUserRolesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	if id == model.SuperAdminUserId {
		response.BadRequest(ctx, "超级管理员用户不允许分配角色")
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 调用服务分配角色
	if err := c.userService.AssignUserRoles(id, req.RoleIDs); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "角色分配成功",
	})
}

// GetUserRoles 获取用户角色
func (c *UserController) GetUserRoles(ctx *gin.Context) {
	// 获取用户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, 400, "用户ID格式错误")
		return
	}

	// 调用服务获取用户角色
	roles, err := c.userService.GetUserRoles(id)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, roles)
}
