package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SuperAdminRoutes 超级管理员路由模块
type SuperAdminRoutes struct {
	BaseRouteModule
}

// NewSuperAdminRoutes 创建超级管理员路由模块
func NewSuperAdminRoutes() RouteModule {
	return &SuperAdminRoutes{
		BaseRouteModule: BaseRouteModule{
			Name:   "super_admin",
			Prefix: "/v1/super-admin",
		},
	}
}

// RegisterRoutes 注册超级管理员路由
func (r *SuperAdminRoutes) RegisterRoutes(router *gin.RouterGroup, c container.Container) {
	// 超管相关路由（需要超级管理员权限）
	superAdmin := router.Group("/v1/super-admin")
	superAdmin.Use(middleware.JWTAuthMiddleware())        // 应用JWT认证中间件
	superAdmin.Use(middleware.SuperAdminAuthMiddleware()) // 应用超级管理员权限中间件
	{
		// 租户管理
		r.registerTenantManagement(superAdmin, c)

		// 角色管理
		r.registerRoleManagement(superAdmin, c)

		// 菜单管理
		r.registerMenuManagement(superAdmin, c)
	}
}

// registerTenantManagement 注册租户管理路由
func (r *SuperAdminRoutes) registerTenantManagement(superAdmin *gin.RouterGroup, c container.Container) {
	tenants := superAdmin.Group("/tenants")
	{
		tenants.POST("", c.TenantController().CreateTenant)                 // 创建租户
		tenants.GET("", c.TenantController().GetTenants)                    // 获取租户列表
		tenants.GET("/list", c.TenantController().GetSelectList)            // 获取租户选择列表
		tenants.GET("/:id", c.TenantController().GetTenant)                 // 获取租户详情
		tenants.PUT("/:id", c.TenantController().UpdateTenant)              // 更新租户
		tenants.DELETE("/:id", c.TenantController().DeleteTenant)           // 删除租户
		tenants.PUT("/:id/status", c.TenantController().UpdateTenantStatus) // 更新租户状态
		tenants.GET("/check-domain", c.TenantController().CheckDomain)      // 检查域名可用性
		tenants.GET("/check-name", c.TenantController().CheckName)          // 检查名称可用性
		tenants.GET("/:id/config", c.TenantController().GetTenantConfig)    // 获取租户配置
		tenants.PUT("/:id/config", c.TenantController().UpdateTenantConfig) // 更新租户配置
	}
}

// registerRoleManagement 注册角色管理路由
func (r *SuperAdminRoutes) registerRoleManagement(superAdmin *gin.RouterGroup, c container.Container) {
	roles := superAdmin.Group("/roles")
	{
		roles.POST("", c.RoleController().CreateRole)                 // 创建角色
		roles.GET("", c.RoleController().GetRoles)                    // 获取角色列表
		roles.GET("/:id", c.RoleController().GetRole)                 // 获取角色详情
		roles.PUT("/:id", c.RoleController().UpdateRole)              // 更新角色
		roles.DELETE("/:id", c.RoleController().DeleteRole)           // 删除角色
		roles.GET("/:id/menus", c.MenuController().GetRoleMenus)      // 获取角色菜单
		roles.PUT("/:id/menus", c.MenuController().AssignMenusToRole) // 为角色分配菜单
	}
}

// registerMenuManagement 注册菜单管理路由
func (r *SuperAdminRoutes) registerMenuManagement(superAdmin *gin.RouterGroup, c container.Container) {
	menus := superAdmin.Group("/menus")
	{
		menus.POST("", c.MenuController().CreateMenu)                 // 创建菜单
		menus.GET("", c.MenuController().GetMenus)                    // 获取菜单列表
		menus.GET("/tree", c.MenuController().GetMenuTree)            // 获取菜单树
		menus.GET("/:id", c.MenuController().GetMenu)                 // 获取菜单详情
		menus.PUT("/:id", c.MenuController().UpdateMenu)              // 更新菜单
		menus.DELETE("/:id", c.MenuController().DeleteMenu)           // 删除菜单
		menus.PUT("/:id/status", c.MenuController().UpdateMenuStatus) // 更新菜单状态
	}
}
