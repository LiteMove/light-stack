package controller

import (
	"strconv"

	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct {
	roleService service.RoleService
}

// NewRoleController 创建角色控制器
func NewRoleController(roleService service.RoleService) *RoleController {
	return &RoleController{
		roleService: roleService,
	}
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

// GetEnabledRoles 获取启用的角色列表
func (c *RoleController) GetEnabledRoles(ctx *gin.Context) {

	isSuperAdmin := ctx.GetBool("is_super_admin")
	roles, err := c.roleService.GetEnabledRoles(isSuperAdmin)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, roles)
}

// GetRoles 获取角色列表
func (c *RoleController) GetRoles(ctx *gin.Context) {
	// 获取查询参数
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

	roles, total, err := c.roleService.GetList(page, pageSize, status)
	if err != nil {
		response.BadRequest(ctx, "获取角色列表失败")
		return
	}

	response.SuccessWithPage(ctx, roles, total, page, pageSize)
}

// 请求结构体定义

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

// AssignRolesRequest 分配角色请求
type AssignRolesRequest struct {
	RoleIDs []uint64 `json:"roleIds" validate:"required"`
}
