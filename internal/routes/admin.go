package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/gin-gonic/gin"
)

// AdminRoutes 管理员路由模块
type AdminRoutes struct {
	BaseRouteModule
}

// NewAdminRoutes 创建管理员路由模块
func NewAdminRoutes() RouteModule {
	return &AdminRoutes{
		BaseRouteModule: BaseRouteModule{
			Name:   "admin",
			Prefix: "/v1/admin",
		},
	}
}

// RegisterRoutes 注册管理员路由
func (r *AdminRoutes) RegisterRoutes(router *gin.RouterGroup, c container.Container) {
	// 管理员路由（需要认证）
	admin := router.Group("/v1/admin")
	admin.Use(middleware.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		// 用户管理
		r.registerUserManagement(admin, c)

		// 角色管理
		r.registerRoleManagement(admin, c)

		// 文件管理
		r.registerFileManagement(admin, c)
	}
}

// registerUserManagement 注册用户管理路由
func (r *AdminRoutes) registerUserManagement(admin *gin.RouterGroup, c container.Container) {
	users := admin.Group("/users")
	{
		users.POST("", middleware.CheckPermission("system:user:create"), c.UserController().CreateUser)                        // 创建用户
		users.GET("", middleware.CheckPermission("system:user:list"), c.UserController().GetUsers)                             // 获取用户列表
		users.GET("/:id", middleware.CheckPermission("system:user:detail"), c.UserController().GetUser)                        // 获取用户详情
		users.PUT("/:id", middleware.CheckPermission("system:user:update"), c.UserController().UpdateUser)                     // 更新用户
		users.DELETE("/:id", middleware.CheckPermission("system:user:delete"), c.UserController().DeleteUser)                  // 删除用户
		users.PUT("/:id/status", middleware.CheckPermission("system:user:update"), c.UserController().UpdateUserStatus)        // 更新用户状态
		users.PUT("/batch/status", middleware.CheckPermission("system:user:update"), c.UserController().BatchUpdateUserStatus) // 批量更新用户状态
		users.POST("/:id/reset-password", middleware.CheckPermission("system:user:update"), c.UserController().ResetPassword)  // 重置密码
		users.PUT("/:id/roles", middleware.CheckPermission("system:user:role:assign"), c.UserController().AssignUserRoles)     // 为用户分配角色
		users.GET("/:id/roles", middleware.CheckPermission("system:user:role:list"), c.UserController().GetUserRoles)          // 获取用户角色
	}
}

// registerRoleManagement 注册角色管理路由
func (r *AdminRoutes) registerRoleManagement(admin *gin.RouterGroup, c container.Container) {
	roles := admin.Group("/roles")
	{
		roles.GET("/select-list", c.RoleController().GetEnabledRoles) // 获取角色列表 - 管理员角色
	}
}

// registerFileManagement 注册文件管理路由
func (r *AdminRoutes) registerFileManagement(admin *gin.RouterGroup, c container.Container) {
	files := admin.Group("/files")
	{
		files.GET("", c.FileController().GetAllFiles)       // 获取所有文件列表
		files.GET("/:id", c.FileController().GetFile)       // 获取文件信息
		files.DELETE("/:id", c.FileController().DeleteFile) // 删除文件
	}
}
