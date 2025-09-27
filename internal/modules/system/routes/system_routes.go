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
			users.GET("", globals.UserCtrl().GetUsers)
			users.POST("", globals.UserCtrl().CreateUser)
			users.PUT("/:id", globals.UserCtrl().UpdateUser)
			users.DELETE("/:id", globals.UserCtrl().DeleteUser)
			users.GET("/:id", globals.UserCtrl().GetUser)
		}

		// 角色管理
		roles := admin.Group("/roles")
		{
			roles.GET("", globals.RoleCtrl().GetRoles)
			roles.POST("", globals.RoleCtrl().CreateRole)
			roles.PUT("/:id", globals.RoleCtrl().UpdateRole)
			roles.DELETE("/:id", globals.RoleCtrl().DeleteRole)
		}

		// 菜单管理
		menus := admin.Group("/menus")
		{
			menus.GET("", globals.MenuCtrl().GetMenus)
			menus.POST("", globals.MenuCtrl().CreateMenu)
			menus.PUT("/:id", globals.MenuCtrl().UpdateMenu)
			menus.DELETE("/:id", globals.MenuCtrl().DeleteMenu)
			menus.GET("/tree", globals.MenuCtrl().GetMenuTree)
		}

		// 租户管理
		tenants := admin.Group("/tenants")
		{
			tenants.GET("", globals.TenantCtrl().GetTenants)
			tenants.POST("", globals.TenantCtrl().CreateTenant)
			tenants.PUT("/:id", globals.TenantCtrl().UpdateTenant)
			tenants.DELETE("/:id", globals.TenantCtrl().DeleteTenant)
		}

		// 字典管理
		dicts := admin.Group("/dicts")
		{
			// 字典类型管理
			dictTypes := dicts.Group("/types")
			{
				dictTypes.GET("", globals.DictCtrl().GetTypeList)
				dictTypes.POST("", globals.DictCtrl().CreateType)
				dictTypes.PUT("/:id", globals.DictCtrl().UpdateType)
				dictTypes.DELETE("/:id", globals.DictCtrl().DeleteType)
			}

			// 字典数据管理
			dictData := dicts.Group("/data")
			{
				dictData.GET("/type/:type", globals.DictCtrl().GetDataList)
				dictData.POST("", globals.DictCtrl().CreateData)
				dictData.PUT("/:id", globals.DictCtrl().UpdateData)
				dictData.DELETE("/:id", globals.DictCtrl().DeleteData)
			}
		}
	}
}
