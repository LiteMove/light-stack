package routes

import (
	"github.com/LiteMove/light-stack/internal/shared/globals"
	"github.com/LiteMove/light-stack/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterSystemRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	// 需要认证的系统管理路由
	admin := v1.Group("/admin")
	admin.Use(middleware.Auth())
	{
		// 用户管理
		users := admin.Group("/users")
		{
			users.POST("", middleware.CheckPermission("system:user:create"), globals.UserCtrl().CreateUser)                        // 创建用户
			users.GET("", middleware.CheckPermission("system:user:list"), globals.UserCtrl().GetUsers)                             // 获取用户列表
			users.GET("/:id", middleware.CheckPermission("system:user:detail"), globals.UserCtrl().GetUser)                        // 获取用户详情
			users.PUT("/:id", middleware.CheckPermission("system:user:update"), globals.UserCtrl().UpdateUser)                     // 更新用户
			users.DELETE("/:id", middleware.CheckPermission("system:user:delete"), globals.UserCtrl().DeleteUser)                  // 删除用户
			users.PUT("/:id/status", middleware.CheckPermission("system:user:update"), globals.UserCtrl().UpdateUserStatus)        // 更新用户状态
			users.PUT("/batch/status", middleware.CheckPermission("system:user:update"), globals.UserCtrl().BatchUpdateUserStatus) // 批量更新用户状态
			users.POST("/:id/reset-password", middleware.CheckPermission("system:user:reset"), globals.UserCtrl().ResetPassword)   // 重置密码
			users.PUT("/:id/roles", middleware.CheckPermission("system:user:role:assign"), globals.UserCtrl().AssignUserRoles)     // 为用户分配角色
			users.GET("/:id/roles", middleware.CheckPermission("system:user:role:assign"), globals.UserCtrl().GetUserRoles)        // 获取用户角色
		}

		// 角色管理
		roles := admin.Group("/roles")
		{
			roles.POST("", globals.RoleCtrl().CreateRole)                 // 创建角色
			roles.GET("/select-list", globals.RoleCtrl().GetEnabledRoles) // 获取下拉角色列表
			roles.GET("", globals.RoleCtrl().GetRoles)                    // 获取角色列表
			roles.GET("/:id", globals.RoleCtrl().GetRole)                 // 获取角色详情
			roles.PUT("/:id", globals.RoleCtrl().UpdateRole)              // 更新角色
			roles.DELETE("/:id", globals.RoleCtrl().DeleteRole)           // 删除角色
			roles.GET("/:id/menus", globals.MenuCtrl().GetRoleMenus)      // 获取角色菜单
			roles.PUT("/:id/menus", globals.MenuCtrl().AssignMenusToRole) // 为角色分配菜单
		}

		// 菜单管理
		menus := admin.Group("/menus")
		{
			menus.POST("", globals.MenuCtrl().CreateMenu)                 // 创建菜单
			menus.GET("", globals.MenuCtrl().GetMenus)                    // 获取菜单列表
			menus.GET("/tree", globals.MenuCtrl().GetMenuTree)            // 获取菜单树
			menus.GET("/:id", globals.MenuCtrl().GetMenu)                 // 获取菜单详情
			menus.PUT("/:id", globals.MenuCtrl().UpdateMenu)              // 更新菜单
			menus.DELETE("/:id", globals.MenuCtrl().DeleteMenu)           // 删除菜单
			menus.PUT("/:id/status", globals.MenuCtrl().UpdateMenuStatus) // 更新菜单状态
		}

		// 租户管理
		tenants := admin.Group("/tenants")
		{
			tenants.POST("", globals.TenantCtrl().CreateTenant) // 创建租户
			tenants.GET("", globals.TenantCtrl().GetTenants)
			tenants.GET("/list", globals.TenantCtrl().GetSelectList)            // 获取租户列表
			tenants.GET("/:id", globals.TenantCtrl().GetTenant)                 // 获取租户详情
			tenants.PUT("/:id", globals.TenantCtrl().UpdateTenant)              // 更新租户
			tenants.DELETE("/:id", globals.TenantCtrl().DeleteTenant)           // 删除租户
			tenants.PUT("/:id/status", globals.TenantCtrl().UpdateTenantStatus) // 更新租户状态
			tenants.GET("/check-domain", globals.TenantCtrl().CheckDomain)      // 检查域名可用性
			tenants.GET("/check-name", globals.TenantCtrl().CheckName)          // 检查名称可用性
			tenants.GET("/:id/config", globals.TenantCtrl().GetTenantConfig)    // 获取租户配置
			tenants.PUT("/:id/config", globals.TenantCtrl().UpdateTenantConfig) // 更新租户配置
		}

		// 字典管理
		dicts := admin.Group("/dicts")
		{
			// 字典类型管理
			dictTypes := dicts.Group("/types")
			{
				dictTypes.POST("", globals.DictCtrl().CreateType)       // 创建字典类型
				dictTypes.GET("", globals.DictCtrl().GetTypeList)       // 获取字典类型列表
				dictTypes.GET("/:id", globals.DictCtrl().GetType)       // 获取字典类型详情
				dictTypes.PUT("/:id", globals.DictCtrl().UpdateType)    // 更新字典类型
				dictTypes.DELETE("/:id", globals.DictCtrl().DeleteType) // 删除字典类型
			}

			// 字典数据管理
			dictData := dicts.Group("/data")
			{
				dictData.POST("", globals.DictCtrl().CreateData)                        // 创建字典数据
				dictData.GET("/type/:type", globals.DictCtrl().GetDataList)             // 获取字典数据列表
				dictData.GET("/:id", globals.DictCtrl().GetData)                        // 获取字典数据详情
				dictData.PUT("/:id", globals.DictCtrl().UpdateData)                     // 更新字典数据
				dictData.DELETE("/:id", globals.DictCtrl().DeleteData)                  // 删除字典数据
				dictData.PUT("/batch/status", globals.DictCtrl().BatchUpdateDataStatus) // 批量更新状态
				dictData.DELETE("/batch", globals.DictCtrl().BatchDeleteData)           // 批量删除
			}
		}
	}
}
