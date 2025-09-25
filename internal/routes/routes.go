package routes

import (
	"github.com/LiteMove/light-stack/internal/config"
	"github.com/LiteMove/light-stack/internal/globals"
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// 初始化所有服务
	globals.Init()

	// API 分组
	api := r.Group("/api")

	// 全局中间件
	api.Use(middleware.OptionalAuth())                        // 可选认证
	api.Use(middleware.TenantMiddleware(globals.TenantSvc())) // 租户处理

	// 注册各模块路由
	registerPublicRoutes(api)
	registerAuthRoutes(api)
	registerUserRoutes(api)
	registerAdminRoutes(api)
	registerSuperAdminRoutes(api)
	registerStaticRoutes(r)
}

// 公开接口（无需认证）
func registerPublicRoutes(api *gin.RouterGroup) {
	// 健康检查
	api.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":  "ok",
			"service": "light-stack",
		})
	})

	// 租户信息接口（公开，无需认证）
	api.GET("/tenant/info", globals.TenantCtrl().GetTenantByDomain)
}

// 认证相关接口
func registerAuthRoutes(api *gin.RouterGroup) {
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
			auth.POST("/login", globals.AuthCtrl().Login)          // 用户登录
			auth.POST("/register", globals.AuthCtrl().Register)    // 用户注册
			auth.POST("/refresh", globals.AuthCtrl().RefreshToken) // 刷新token
			auth.POST("/logout", globals.AuthCtrl().Logout)        // 用户登出
		}
	}
}

// 用户相关接口（需要认证）
func registerUserRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		// 用户相关路由（需要认证）
		user := v1.Group("/user")
		user.Use(middleware.Auth()) // 应用JWT认证中间件
		{
			user.GET("/profile", globals.AuthCtrl().GetProfile)      // 获取用户信息
			user.PUT("/profile", globals.AuthCtrl().UpdateProfile)   // 更新用户信息
			user.PUT("/password", globals.AuthCtrl().ChangePassword) // 修改密码
			user.GET("/:id/roles", globals.AuthCtrl().GetUserRoles)  // 获取用户角色
		}

		// 用户查看和修改租户配置
		user.Use(middleware.TenantMiddleware(globals.TenantSvc()))
		{
			user.GET("/tenant-config", globals.ProfileCtrl().GetTenantConfig)    // 获取所在租户配置
			user.PUT("/tenant-config", globals.ProfileCtrl().UpdateTenantConfig) // 更新所在租户配置
		}

		// 文件相关路由（需要认证）
		files := v1.Group("/files")
		files.Use(middleware.Auth()) // 应用JWT认证中间件
		{
			files.POST("/upload", globals.FileCtrl().UploadFile)         // 文件上传
			files.GET("", globals.FileCtrl().GetAllFiles)                // 获取文件列表（按租户）
			files.GET("/list", globals.FileCtrl().GetAllFiles)           // 获取文件列表（兼容旧接口）
			files.GET("/:id", globals.FileCtrl().GetFile)                // 获取文件信息（包含access_url用于下载/预览）
			files.GET("/:id/private", globals.FileCtrl().GetPrivateFile) // 获取私有文件内容
			files.DELETE("/:id", globals.FileCtrl().DeleteFile)          // 删除文件
		}

		// 仪表盘相关路由（需要认证）
		dashboard := v1.Group("/dashboard")
		dashboard.Use(middleware.Auth()) // 应用JWT认证中间件
		{
			dashboard.GET("/stats", globals.DashboardCtrl().GetDashboardStats) // 获取仪表盘统计数据
		}
	}
}

// 管理员接口
func registerAdminRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		// 管理员相关路由
		admin := v1.Group("/admin")
		admin.Use(middleware.Auth()) // 应用JWT认证中间件
		{
			// 用户角色管理
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
				roles.GET("/select-list", globals.RoleCtrl().GetEnabledRoles) // 获取角色列表 - 管理员角色
			}

			// 文件管理
			files := admin.Group("/files")
			{
				files.GET("", globals.FileCtrl().GetAllFiles)       // 获取所有文件列表
				files.GET("/:id", globals.FileCtrl().GetFile)       // 获取文件信息
				files.DELETE("/:id", globals.FileCtrl().DeleteFile) // 删除文件
			}
		}
	}
}

// 超级管理员接口
func registerSuperAdminRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		// 超管相关路由（需要超级管理员权限）
		superAdmin := v1.Group("/super-admin")
		superAdmin.Use(middleware.Auth())       // 应用JWT认证中间件
		superAdmin.Use(middleware.SuperAdmin()) // 超级管理员权限中间件
		{
			// 租户管理
			tenants := superAdmin.Group("/tenants")
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

			// 系统信息（超管专用）
			system := superAdmin.Group("/system")
			{
				system.GET("/info", globals.DashboardCtrl().GetSystemInfo) // 获取系统信息
			}

			// 角色管理
			roles := superAdmin.Group("/roles")
			{
				roles.POST("", globals.RoleCtrl().CreateRole)                 // 创建角色
				roles.GET("", globals.RoleCtrl().GetRoles)                    // 获取角色列表
				roles.GET("/:id", globals.RoleCtrl().GetRole)                 // 获取角色详情
				roles.PUT("/:id", globals.RoleCtrl().UpdateRole)              // 更新角色
				roles.DELETE("/:id", globals.RoleCtrl().DeleteRole)           // 删除角色
				roles.GET("/:id/menus", globals.MenuCtrl().GetRoleMenus)      // 获取角色菜单
				roles.PUT("/:id/menus", globals.MenuCtrl().AssignMenusToRole) // 为角色分配菜单
			}

			// 菜单管理
			menus := superAdmin.Group("/menus")
			{
				menus.POST("", globals.MenuCtrl().CreateMenu)                 // 创建菜单
				menus.GET("", globals.MenuCtrl().GetMenus)                    // 获取菜单列表
				menus.GET("/tree", globals.MenuCtrl().GetMenuTree)            // 获取菜单树
				menus.GET("/:id", globals.MenuCtrl().GetMenu)                 // 获取菜单详情
				menus.PUT("/:id", globals.MenuCtrl().UpdateMenu)              // 更新菜单
				menus.DELETE("/:id", globals.MenuCtrl().DeleteMenu)           // 删除菜单
				menus.PUT("/:id/status", globals.MenuCtrl().UpdateMenuStatus) // 更新菜单状态
			}
		}
	}
}

// 静态文件服务
func registerStaticRoutes(r *gin.Engine) {
	// 静态文件服务 - 使用配置文件中的base_url
	cfg := config.Get()
	baseURL := cfg.File.BaseURL
	if baseURL == "" {
		baseURL = "/api/static"
	}

	// 公开文件 - 无需认证
	r.Static(baseURL+"/public", "./uploads/public")

	// 私有文件 - 需要认证和权限验证
	privateFiles := r.Group(baseURL + "/private")
	privateFiles.Use(middleware.Auth())
	privateFiles.Use(middleware.TenantMiddleware(globals.TenantSvc()))
	{
		// 私有文件访问需要验证用户是否有权限访问该租户的文件
		privateFiles.StaticFS("/", gin.Dir("./uploads/private", false))
	}
}
