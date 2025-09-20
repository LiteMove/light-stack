package controller

import (
	"strconv"

	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// MenuController 菜单控制器
type MenuController struct {
	menuService service.MenuService
	validator   *validator.Validate
}

// NewMenuController 创建菜单控制器
func NewMenuController(menuService service.MenuService) *MenuController {
	return &MenuController{
		menuService: menuService,
		validator:   validator.New(),
	}
}

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	ParentID  uint64 `json:"parent_id" validate:""`
	Name      string `json:"name" validate:"required,min=1,max=100"`
	Code      string `json:"code" validate:"required,min=1,max=100"`
	Type      string `json:"type" validate:"required,oneof=directory menu permission"`
	Path      string `json:"path" validate:"max=255"`
	Component string `json:"component" validate:"max=255"`
	Icon      string `json:"icon" validate:"max=100"`
	Resource  string `json:"resource" validate:"max=255"`
	Action    string `json:"action" validate:"max=50"`
	SortOrder int    `json:"sort_order"`
	IsHidden  bool   `json:"is_hidden"`
	Status    int    `json:"status" validate:"required,oneof=1 2"`
	Meta      string `json:"meta"`
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	ParentID  uint64 `json:"parent_id"`
	Name      string `json:"name" validate:"required,min=1,max=100"`
	Code      string `json:"code" validate:"required,min=1,max=100"`
	Type      string `json:"type" validate:"required,oneof=directory menu permission"`
	Path      string `json:"path" validate:"max=255"`
	Component string `json:"component" validate:"max=255"`
	Icon      string `json:"icon" validate:"max=100"`
	Resource  string `json:"resource" validate:"max=255"`
	Action    string `json:"action" validate:"max=50"`
	SortOrder int    `json:"sort_order"`
	IsHidden  bool   `json:"is_hidden"`
	Status    int    `json:"status" validate:"required,oneof=1 2"`
	Meta      string `json:"meta"`
}

// MenuListRequest 菜单列表请求
type MenuListRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"page_size" validate:"min=1,max=100"`
	Name     string `form:"name"`
	Status   int    `form:"status" validate:"oneof=0 1 2"`
}

// UpdateStatusRequest 更新状态请求
type UpdateStatusRequest struct {
	Status int `json:"status" validate:"required,oneof=1 2"`
}

// BatchUpdateStatusRequest 批量更新状态请求
type BatchUpdateStatusRequest struct {
	IDs    []uint64 `json:"ids" validate:"required,min=1"`
	Status int      `json:"status" validate:"required,oneof=1 2"`
}

// AssignMenusRequest 分配菜单请求
type AssignMenusRequest struct {
	MenuIDs []uint64 `json:"menu_ids" validate:"required"`
}

// CreateMenu 创建菜单
func (mc *MenuController) CreateMenu(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数格式错误")
		return
	}

	if err := mc.validator.Struct(&req); err != nil {
		response.BadRequest(c, "参数验证失败")
		return
	}

	menu := &model.Menu{
		ParentID:  req.ParentID,
		Name:      req.Name,
		Code:      req.Code,
		Type:      req.Type,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Resource:  req.Resource,
		Action:    req.Action,
		SortOrder: req.SortOrder,
		IsHidden:  req.IsHidden,
		Status:    req.Status,
		Meta:      req.Meta,
	}

	if err := mc.menuService.CreateMenu(menu); err != nil {
		response.BadRequest(c, "创建菜单失败")
		return
	}

	response.Success(c, menu.ToProfile())
}

// GetMenus 获取菜单列表
func (mc *MenuController) GetMenus(c *gin.Context) {
	var req MenuListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数格式错误")
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	if err := mc.validator.Struct(&req); err != nil {
		response.BadRequest(c, "参数验证失败")
		return
	}

	menus, total, err := mc.menuService.GetMenuList(req.Page, req.PageSize, req.Name, req.Status)
	if err != nil {
		response.BadRequest(c, "获取菜单列表失败")
		return
	}

	response.Success(c, gin.H{
		"list":      menus,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetMenu 获取菜单详情
func (mc *MenuController) GetMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	menu, err := mc.menuService.GetMenu(id)
	if err != nil {
		response.NotFound(c, "菜单不存在")
		return
	}

	response.Success(c, menu.ToProfile())
}

// UpdateMenu 更新菜单
func (mc *MenuController) UpdateMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数格式错误")
		return
	}

	if err := mc.validator.Struct(&req); err != nil {
		response.BadRequest(c, "参数验证失败")
		return
	}

	menu := &model.Menu{
		ID:        id,
		ParentID:  req.ParentID,
		Name:      req.Name,
		Code:      req.Code,
		Type:      req.Type,
		Path:      req.Path,
		Component: req.Component,
		Icon:      req.Icon,
		Resource:  req.Resource,
		Action:    req.Action,
		SortOrder: req.SortOrder,
		IsHidden:  req.IsHidden,
		Status:    req.Status,
		Meta:      req.Meta,
	}

	if err := mc.menuService.UpdateMenu(menu); err != nil {
		response.BadRequest(c, "更新菜单失败")
		return
	}

	response.Success(c, menu.ToProfile())
}

// DeleteMenu 删除菜单
func (mc *MenuController) DeleteMenu(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	if err := mc.menuService.DeleteMenu(id); err != nil {
		response.BadRequest(c, "删除菜单失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// GetMenuTree 获取菜单树
func (mc *MenuController) GetMenuTree(c *gin.Context) {
	tree, err := mc.menuService.GetMenuTree()
	if err != nil {
		response.InternalServerError(c, "获取菜单树失败")
		return
	}

	response.Success(c, tree)
}

// GetUserMenuTree 获取用户菜单树
func (mc *MenuController) GetUserMenuTree(c *gin.Context) {
	// 从上下文获取用户ID（由认证中间件注入）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "用户ID不存在")
		return
	}

	tree, err := mc.menuService.GetUserMenuTree(userID.(uint64))
	if err != nil {
		response.BadRequest(c, "获取用户菜单树失败")
		return
	}

	response.Success(c, tree)
}

// GetRoleMenus 获取角色菜单
func (mc *MenuController) GetRoleMenus(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	menus, err := mc.menuService.GetRoleMenus(roleID)
	if err != nil {
		response.BadRequest(c, "获取角色菜单失败")
		return
	}

	response.Success(c, menus)
}

// UpdateMenuStatus 更新菜单状态
func (mc *MenuController) UpdateMenuStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的菜单ID")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数格式错误")
		return
	}

	if err := mc.validator.Struct(&req); err != nil {
		response.BadRequest(c, "参数验证失败")
		return
	}

	if err := mc.menuService.UpdateMenuStatus(id, req.Status); err != nil {
		response.BadRequest(c, "更新菜单状态失败")
		return
	}

	response.Success(c, gin.H{"message": "更新状态成功"})
}

// BatchUpdateMenuStatus 批量更新菜单状态
func (mc *MenuController) BatchUpdateMenuStatus(c *gin.Context) {
	var req BatchUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数格式错误")
		return
	}

	if err := mc.validator.Struct(&req); err != nil {
		response.BadRequest(c, "参数验证失败")
		return
	}

	if err := mc.menuService.BatchUpdateMenuStatus(req.IDs, req.Status); err != nil {
		response.BadRequest(c, "批量更新菜单状态失败")
		return
	}

	response.Success(c, gin.H{"message": "批量更新状态成功"})
}

// AssignMenusToRole 为角色分配菜单
func (mc *MenuController) AssignMenusToRole(c *gin.Context) {
	idStr := c.Param("id")
	roleID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的角色ID")
		return
	}

	var req AssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数格式错误")
		return
	}

	if err := mc.validator.Struct(&req); err != nil {
		response.BadRequest(c, "参数验证失败")
		return
	}

	if err := mc.menuService.AssignMenusToRole(roleID, req.MenuIDs); err != nil {
		response.BadRequest(c, "分配菜单失败")
		return
	}

	response.Success(c, gin.H{"message": "分配菜单成功"})
}

// GetMenuPermissions 获取用户菜单权限
func (mc *MenuController) GetMenuPermissions(c *gin.Context) {
	// 从上下文获取用户ID（由认证中间件注入）
	userId := c.GetUint64("user_id")
	if userId == 0 {
		response.Unauthorized(c, "用户ID不存在")
		return
	}

	permissions, err := mc.menuService.GetMenuPermissions(userId)
	if err != nil {
		response.BadRequest(c, "获取菜单权限失败")
		return
	}

	response.Success(c, gin.H{"permissions": permissions})
}
