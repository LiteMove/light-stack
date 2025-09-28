package controller

import (
	"github.com/LiteMove/light-stack/internal/modules/auth/service"
	systemModel "github.com/LiteMove/light-stack/internal/modules/system/model"
	"github.com/LiteMove/light-stack/internal/shared/middleware"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ProfileController 个人中心控制器
type ProfileController struct {
	profileService service.ProfileService
	validator      *validator.Validate
}

// NewProfileController 创建个人中心控制器
func NewProfileController(profileService service.ProfileService) *ProfileController {
	return &ProfileController{
		profileService: profileService,
		validator:      validator.New(),
	}
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" validate:"required,min=1,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"omitempty,len=11"`
	Avatar   string `json:"avatar" validate:"omitempty,max=255"`
}

// ProfileChangePasswordRequest 修改密码请求
type ProfileChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,min=6"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

// UpdateTenantConfigRequest 更新租户配置请求
type UpdateTenantConfigRequest struct {
	SystemName  string                        `json:"systemName"`
	Logo        string                        `json:"logo"`
	Description string                        `json:"description"`
	Copyright   string                        `json:"copyright"`
	FileStorage systemModel.FileStorageConfig `json:"fileStorage" validate:"required"`
}

// GetProfile 获取个人信息
func (c *ProfileController) GetProfile(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userId")
	if !exists {
		response.Unauthorized(ctx, "用户未登录")
		return
	}

	id := userID.(uint64)

	// 调用服务获取用户信息
	profile, err := c.profileService.GetProfile(id)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, profile)
}

// UpdateProfile 更新个人信息
func (c *ProfileController) UpdateProfile(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userId")
	if !exists {
		response.Unauthorized(ctx, "用户未登录")
		return
	}

	id := userID.(uint64)

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 调用服务更新用户信息
	if err := c.profileService.UpdateProfile(id, req.Nickname, req.Email, req.Phone, req.Avatar); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "个人信息更新成功",
	})
}

// ChangePassword 修改密码
func (c *ProfileController) ChangePassword(ctx *gin.Context) {
	// 获取当前用户ID
	userID, exists := ctx.Get("userId")
	if !exists {
		response.Unauthorized(ctx, "用户未登录")
		return
	}

	id := userID.(uint64)

	var req ProfileChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 调用服务修改密码
	if err := c.profileService.ChangePassword(id, req.OldPassword, req.NewPassword); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "密码修改成功",
	})
}

// GetTenantConfig 获取所在租户配置（仅租户管理员）
func (c *ProfileController) GetTenantConfig(ctx *gin.Context) {
	// 获取当前用户ID和租户ID
	userID, exists := ctx.Get("userId")
	if !exists {
		response.Unauthorized(ctx, "用户未登录")
		return
	}

	tenantID, exists := middleware.GetTenantIDFromContext(ctx)
	if !exists {
		response.BadRequest(ctx, "无法获取租户信息")
		return
	}

	uid := userID.(uint64)

	// 检查是否为超级管理员
	isSuperAdmin := ctx.GetBool("is_super_admin")

	// 如果不是超级管理员，检查是否为租户管理员
	if !isSuperAdmin {
		isAdmin, err := c.profileService.IsTenantAdmin(uid, tenantID)
		if err != nil {
			response.InternalServerError(ctx, err.Error())
			return
		}

		if !isAdmin {
			response.Forbidden(ctx, "仅租户管理员可查看租户配置")
			return
		}
	}

	// 获取租户配置
	config, err := c.profileService.GetTenantConfig(tenantID)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, config)
}

// UpdateTenantConfig 更新所在租户配置（仅租户管理员）
func (c *ProfileController) UpdateTenantConfig(ctx *gin.Context) {
	// 获取当前用户ID和租户ID
	userID, exists := ctx.Get("userId")
	if !exists {
		response.Unauthorized(ctx, "用户未登录")
		return
	}

	tenantID, exists := middleware.GetTenantIDFromContext(ctx)
	if !exists {
		response.BadRequest(ctx, "无法获取租户信息")
		return
	}

	uid := userID.(uint64)

	// 检查是否为超级管理员
	isSuperAdmin := ctx.GetBool("is_super_admin")

	// 如果不是超级管理员，检查是否为租户管理员
	if !isSuperAdmin {
		isAdmin, err := c.profileService.IsTenantAdmin(uid, tenantID)
		if err != nil {
			response.InternalServerError(ctx, err.Error())
			return
		}

		if !isAdmin {
			response.Forbidden(ctx, "仅租户管理员可修改租户配置")
			return
		}
	}

	var req UpdateTenantConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 构建租户配置
	config := &systemModel.TenantConfig{
		SystemName:  req.SystemName,
		Logo:        req.Logo,
		Description: req.Description,
		Copyright:   req.Copyright,
		FileStorage: req.FileStorage,
	}

	// 调用服务更新租户配置
	if err := c.profileService.UpdateTenantConfig(tenantID, config); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "租户配置更新成功",
	})
}
