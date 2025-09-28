package controller

import (
	"strconv"
	"strings"

	"github.com/LiteMove/light-stack/internal/shared/utils"

	"github.com/LiteMove/light-stack/internal/modules/system/model"
	"github.com/LiteMove/light-stack/internal/modules/system/service"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TenantController 租户控制器
type TenantController struct {
	tenantService service.TenantService
	validator     *validator.Validate
}

// NewTenantController 创建租户控制器
func NewTenantController(tenantService service.TenantService) *TenantController {
	return &TenantController{
		tenantService: tenantService,
		validator:     validator.New(),
	}
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=100"`
	Domain    string `json:"domain" validate:"omitempty,max=100"`
	Status    int    `json:"status" validate:"required,oneof=1 2 3 4"`
	ExpiredAt string `json:"expiredAt" validate:"omitempty"`
	Config    string `json:"config" validate:"omitempty"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=100"`
	Domain    string `json:"domain" validate:"omitempty,max=100"`
	Status    int    `json:"status" validate:"required,oneof=1 2 3 4"`
	ExpiredAt string `json:"expiredAt" validate:"omitempty"`
	Config    string `json:"config" validate:"omitempty"`
}

// TenantListRequest 租户列表请求
type TenantListRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status" validate:"oneof=0 1 2 3 4"`
}

// UpdateTenantStatusRequest 更新租户状态请求
type UpdateTenantStatusRequest struct {
	Status int `json:"status" validate:"required,oneof=1 2 3 4"`
}

// CreateTenant 创建租户
func (c *TenantController) CreateTenant(ctx *gin.Context) {
	var req CreateTenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 创建租户对象
	tenant := &model.Tenant{
		Name:   req.Name,
		Domain: req.Domain,
		Status: req.Status,
		Config: req.Config,
	}

	// 处理过期时间
	if req.ExpiredAt != "" {
		time, err := utils.ParseToTime(req.ExpiredAt)
		if err != nil {
			response.BadRequest(ctx, "无法解析过期时间: "+err.Error())
			return
		}
		tenant.ExpiredAt = time
	}

	// 调用服务创建租户
	if err := c.tenantService.CreateTenant(tenant); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"id":     tenant.ID,
		"name":   tenant.Name,
		"domain": tenant.Domain,
	})
}

// GetSelectList 获取下拉租户列表
func (c *TenantController) GetSelectList(ctx *gin.Context) {
	tenants, err := c.tenantService.GetSelectList()
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.Success(ctx, tenants)
}

// GetTenants 获取租户列表
func (c *TenantController) GetTenants(ctx *gin.Context) {
	var req TenantListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
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
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 调用服务获取租户列表
	tenants, total, err := c.tenantService.GetTenantList(req.Page, req.PageSize, req.Keyword, req.Status)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 返回分页数据
	response.Success(ctx, gin.H{
		"list":      tenants,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetTenant 获取租户详情
func (c *TenantController) GetTenant(ctx *gin.Context) {
	// 获取租户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "租户ID格式错误")
		return
	}

	// 调用服务获取租户
	tenant, err := c.tenantService.GetTenant(id)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, tenant)
}

// UpdateTenant 更新租户
func (c *TenantController) UpdateTenant(ctx *gin.Context) {
	// 获取租户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "租户ID格式错误")
		return
	}

	var req UpdateTenantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 获取原租户信息
	existingTenant, err := c.tenantService.GetTenant(id)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	// 更新租户信息
	existingTenant.Name = req.Name
	existingTenant.Domain = req.Domain
	existingTenant.Status = req.Status
	existingTenant.Config = req.Config

	// 处理过期时间
	if req.ExpiredAt != "" {
		time, err := utils.ParseToTime(req.ExpiredAt)
		if err != nil {
			response.BadRequest(ctx, "过期时间格式错误")
			return
		}
		existingTenant.ExpiredAt = time
	}

	// 调用服务更新租户
	if err := c.tenantService.UpdateTenant(existingTenant); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"id":     existingTenant.ID,
		"name":   existingTenant.Name,
		"domain": existingTenant.Domain,
	})
}

// DeleteTenant 删除租户
func (c *TenantController) DeleteTenant(ctx *gin.Context) {
	// 获取租户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "租户ID格式错误")
		return
	}
	// 禁止删除系统租户
	if id == model.SystemTenantId {
		response.BadRequest(ctx, "禁止删除系统租户")
		return
	}
	// 调用服务删除租户
	if err := c.tenantService.DeleteTenant(id); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "删除成功",
	})
}

// UpdateTenantStatus 更新租户状态
func (c *TenantController) UpdateTenantStatus(ctx *gin.Context) {
	// 获取租户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "租户ID格式错误")
		return
	}

	var req UpdateTenantStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}
	// 禁止禁用系统租户
	if id == model.SystemTenantId && req.Status != model.TenantStatusActive {
		response.BadRequest(ctx, "禁止禁用系统租户")
		return
	}

	// 调用服务更新状态
	if err := c.tenantService.UpdateTenantStatus(id, req.Status); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "状态更新成功",
	})
}

// CheckDomain 检查域名是否可用
func (c *TenantController) CheckDomain(ctx *gin.Context) {
	domain := ctx.Query("domain")
	if domain == "" {
		response.BadRequest(ctx, "域名参数不能为空")
		return
	}

	// 调用服务检查域名
	exists, err := c.tenantService.CheckDomainExists(domain)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"domain":    domain,
		"available": !exists,
	})
}

// CheckName 检查租户名称是否可用
func (c *TenantController) CheckName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		response.BadRequest(ctx, "名称参数不能为空")
		return
	}

	// 调用服务检查名称
	exists, err := c.tenantService.CheckNameExists(name)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"name":      name,
		"available": !exists,
	})
}

// TenantConfigRequest 租户配置请求
type TenantConfigRequest struct {
	FileStorage model.FileStorageConfig `json:"fileStorage" validate:"required"`
}

// GetTenantConfig 获取租户配置
func (c *TenantController) GetTenantConfig(ctx *gin.Context) {
	// 获取租户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "租户ID格式错误")
		return
	}

	// 调用服务获取租户配置
	config, err := c.tenantService.GetTenantConfig(id)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, config)
}

// UpdateTenantConfig 更新租户配置
func (c *TenantController) UpdateTenantConfig(ctx *gin.Context) {
	// 获取租户ID
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "租户ID格式错误")
		return
	}

	var req TenantConfigRequest
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
	config := &model.TenantConfig{
		FileStorage: req.FileStorage,
	}

	// 调用服务更新租户配置
	if err := c.tenantService.UpdateTenantConfig(id, config); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"message": "配置更新成功",
	})
}

// TenantDisplayInfo 租户展示信息
type TenantDisplayInfo struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	SystemName  string `json:"systemName"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	Copyright   string `json:"copyright"`
}

// GetTenantByDomain 根据域名获取租户展示信息（公开接口）
func (c *TenantController) GetTenantByDomain(ctx *gin.Context) {
	domain := ctx.Query("domain")
	if domain == "" {
		// 如果没有指定域名，尝试从Host头获取
		domain = ctx.Request.Host
		// 去掉端口号
		if idx := strings.Index(domain, ":"); idx != -1 {
			domain = domain[:idx]
		}
	}

	// 对于localhost或127.0.0.1，获取系统租户配置
	if domain == "localhost" || domain == "127.0.0.1" || domain == "" {
		// 获取系统租户（ID=1）的配置
		tenant, err := c.tenantService.GetTenant(model.SystemTenantId)
		if err != nil {
			// 如果获取失败，返回默认配置
			displayInfo := TenantDisplayInfo{
				ID:          1,
				Name:        "系统租户",
				SystemName:  "Light Stack",
				Logo:        "",
				Description: "轻量级管理系统",
				Copyright:   "© 2024 Light Stack. All rights reserved.",
			}
			response.Success(ctx, displayInfo)
			return
		}

		// 解析租户配置
		config, err := tenant.GetConfig()
		if err != nil {
			response.InternalServerError(ctx, "获取租户配置失败")
			return
		}

		// 构建展示信息
		displayInfo := TenantDisplayInfo{
			ID:          tenant.ID,
			Name:        tenant.Name,
			SystemName:  config.SystemName,
			Logo:        config.Logo,
			Description: config.Description,
			Copyright:   config.Copyright,
		}

		// 如果没有配置系统名称，使用租户名称
		if displayInfo.SystemName == "" {
			displayInfo.SystemName = tenant.Name
		}

		// 如果没有配置版权信息，使用默认值
		if displayInfo.Copyright == "" {
			displayInfo.Copyright = "© 2024 " + displayInfo.SystemName + ". All rights reserved."
		}

		response.Success(ctx, displayInfo)
		return
	}

	// 调用服务获取租户信息
	tenant, err := c.tenantService.GetTenantByDomain(domain)
	if err != nil {
		// 如果租户不存在，返回默认配置
		displayInfo := TenantDisplayInfo{
			ID:          1,
			Name:        "系统租户",
			SystemName:  "Light Stack",
			Logo:        "",
			Description: "轻量级管理系统",
			Copyright:   "© 2024 Light Stack. All rights reserved.",
		}
		response.Success(ctx, displayInfo)
		return
	}

	// 检查租户状态
	if !tenant.IsActive() {
		response.BadRequest(ctx, "租户已被禁用")
		return
	}

	// 检查租户是否过期
	if tenant.IsExpired() {
		response.BadRequest(ctx, "租户已过期")
		return
	}

	// 解析租户配置
	config, err := tenant.GetConfig()
	if err != nil {
		response.InternalServerError(ctx, "获取租户配置失败")
		return
	}

	// 构建展示信息
	displayInfo := TenantDisplayInfo{
		ID:          tenant.ID,
		Name:        tenant.Name,
		SystemName:  config.SystemName,
		Logo:        config.Logo,
		Description: config.Description,
		Copyright:   config.Copyright,
	}

	// 如果没有配置系统名称，使用租户名称
	if displayInfo.SystemName == "" {
		displayInfo.SystemName = tenant.Name
	}

	// 如果没有配置版权信息，使用默认值
	if displayInfo.Copyright == "" {
		displayInfo.Copyright = "© 2024 " + displayInfo.SystemName + ". All rights reserved."
	}

	response.Success(ctx, displayInfo)
}
