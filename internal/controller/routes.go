package controller

import (
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/database"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Check 健康检查
func (h *HealthController) Check(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ok",
		"service": "light-stack",
	})
}

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine) {
	// 初始化依赖
	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	menuRepo := repository.NewMenuRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	authService := service.NewAuthService(userRepo, roleRepo)
	userService := service.NewUserService(userRepo, roleRepo)
	roleService := service.NewRoleService(roleRepo, userRepo)
	menuService := service.NewMenuService(menuRepo, roleRepo)
	tenantService := service.NewTenantService(tenantRepo, userRepo)
	authController := NewAuthController(authService, roleService)
	userController := NewUserController(userService)
	roleController := NewRoleController(roleService)
	menuController := NewMenuController(menuService)
	tenantController := NewTenantController(tenantService)
	healthController := NewHealthController()

	// API 分组
	api := r.Group("/api")
	api.Use(middleware.TenantMiddleware(tenantService))
	{
		// 健康检查
		api.GET("/health", healthController.Check)

		// V1 版本API
		v1 := api.Group("/v1")
		{
			// 基础路由组
			v1.GET("/ping", func(c *gin.Context) {
				response.Success(c, gin.H{
					"message": "pong",
				})
			})

			// 认证相关路由（无需认证）
			auth := v1.Group("/auth")
			{
				auth.POST("/login", authController.Login)          // 用户登录
				auth.POST("/register", authController.Register)    // 用户注册
				auth.POST("/refresh", authController.RefreshToken) // 刷新token
				auth.POST("/logout", authController.Logout)        // 用户登出
			}

			// 用户相关路由（需要认证）
			user := v1.Group("/user")
			user.Use(middleware.JWTAuthMiddleware()) // 应用JWT认证中间件
			{
				user.GET("/profile", authController.GetProfile)      // 获取用户信息
				user.PUT("/profile", authController.UpdateProfile)   // 更新用户信息
				user.PUT("/password", authController.ChangePassword) // 修改密码
				user.GET("/:id/roles", authController.GetUserRoles)  // 获取用户角色

				// 用户菜单相关
				user.GET("/menus", menuController.GetUserMenuTree)          // 获取用户菜单树
				user.GET("/permissions", menuController.GetMenuPermissions) // 获取用户菜单权限
			}

			// 管理员相关路由（需要管理员权限）
			admin := v1.Group("/admin")
			admin.Use(middleware.JWTAuthMiddleware())   // 应用JWT认证中间件
			admin.Use(middleware.AdminAuthMiddleware()) // 应用管理员权限中间件
			{
				// 租户管理
				tenants := admin.Group("/tenants")
				{
					tenants.POST("", tenantController.CreateTenant)                 // 创建租户
					tenants.GET("", tenantController.GetTenants)                    // 获取租户列表
					tenants.GET("/:id", tenantController.GetTenant)                 // 获取租户详情
					tenants.PUT("/:id", tenantController.UpdateTenant)              // 更新租户
					tenants.DELETE("/:id", tenantController.DeleteTenant)           // 删除租户
					tenants.PUT("/:id/status", tenantController.UpdateTenantStatus) // 更新租户状态
					tenants.GET("/check-domain", tenantController.CheckDomain)      // 检查域名可用性
					tenants.GET("/check-name", tenantController.CheckName)          // 检查名称可用性
				}

				// 角色管理
				roles := admin.Group("/roles")
				{
					roles.POST("", roleController.CreateRole)                 // 创建角色
					roles.GET("", roleController.GetRoles)                    // 获取角色列表
					roles.GET("/:id", roleController.GetRole)                 // 获取角色详情
					roles.PUT("/:id", roleController.UpdateRole)              // 更新角色
					roles.DELETE("/:id", roleController.DeleteRole)           // 删除角色
					roles.GET("/:id/menus", menuController.GetRoleMenus)      // 获取角色菜单
					roles.PUT("/:id/menus", menuController.AssignMenusToRole) // 为角色分配菜单
				}

				// 菜单管理
				menus := admin.Group("/menus")
				{
					menus.POST("", menuController.CreateMenu)                        // 创建菜单
					menus.GET("", menuController.GetMenus)                           // 获取菜单列表
					menus.GET("/tree", menuController.GetMenuTree)                   // 获取菜单树
					menus.GET("/:id", menuController.GetMenu)                        // 获取菜单详情
					menus.PUT("/:id", menuController.UpdateMenu)                     // 更新菜单
					menus.DELETE("/:id", menuController.DeleteMenu)                  // 删除菜单
					menus.PUT("/:id/status", menuController.UpdateMenuStatus)        // 更新菜单状态
					menus.PUT("/batch/status", menuController.BatchUpdateMenuStatus) // 批量更新菜单状态
				}

				// 用户角色管理
				users := admin.Group("/users")
				{
					users.POST("", userController.CreateUser)                        // 创建用户
					users.GET("", userController.GetUsers)                           // 获取用户列表
					users.GET("/:id", userController.GetUser)                        // 获取用户详情
					users.PUT("/:id", userController.UpdateUser)                     // 更新用户
					users.DELETE("/:id", userController.DeleteUser)                  // 删除用户
					users.PUT("/:id/status", userController.UpdateUserStatus)        // 更新用户状态
					users.PUT("/batch/status", userController.BatchUpdateUserStatus) // 批量更新用户状态
					users.POST("/:id/reset-password", userController.ResetPassword)  // 重置密码
					users.PUT("/:id/roles", userController.AssignUserRoles)          // 为用户分配角色
					users.GET("/:id/roles", userController.GetUserRoles)             // 获取用户角色
				}
			}
		}
	}
}
