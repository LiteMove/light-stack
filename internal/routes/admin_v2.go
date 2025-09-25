package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/gin-gonic/gin"
)

// AdminRoutesV2 管理员路由模块（使用统一权限中间件）
// 这个版本展示了如何使用新的统一权限中间件替换原有的权限检查
type AdminRoutesV2 struct {
	BaseRouteModule
}

// NewAdminRoutesV2 创建升级版管理员路由模块
func NewAdminRoutesV2() RouteModule {
	return &AdminRoutesV2{
		BaseRouteModule: BaseRouteModule{
			Name:   "admin_v2",
			Prefix: "/v1/admin",
		},
	}
}

// RegisterRoutes 注册管理员路由（使用统一权限中间件）
func (r *AdminRoutesV2) RegisterRoutes(router *gin.RouterGroup, c container.Container) {
	// 管理员路由基础组（需要JWT认证）
	admin := router.Group("/v1/admin")
	admin.Use(middleware.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		// 用户管理模块
		r.registerUserManagementV2(admin, c)

		// 角色管理模块
		r.registerRoleManagementV2(admin, c)

		// 文件管理模块
		r.registerFileManagementV2(admin, c)
	}
}

// registerUserManagementV2 注册用户管理路由（使用统一权限中间件）
func (r *AdminRoutesV2) registerUserManagementV2(admin *gin.RouterGroup, c container.Container) {
	users := admin.Group("/users")
	{
		// 创建用户 - 使用统一权限配置
		users.POST("",
			middleware.UnifiedAuth(middleware.UserManage.WithCustomMessage("需要用户创建权限")),
			c.UserController().CreateUser)

		// 获取用户列表 - 使用预定义权限配置
		users.GET("",
			middleware.UnifiedAuth(middleware.UserView),
			c.UserController().GetUsers)

		// 获取用户详情 - 组合权限和角色验证
		users.GET("/:id",
			middleware.UnifiedAuth(middleware.AuthConfig{
				RequireAuth:  true,
				Permissions:  []string{"system:user:detail", "system:user:list"}, // 支持多个权限
				Roles:        []string{"admin", "user_manager", "viewer"},        // 支持多个角色
				Logic:        middleware.AuthLogicOR,                             // 权限或角色任一满足即可
				ErrorMessage: "需要用户查看权限或管理员角色",
			}),
			c.UserController().GetUser)

		// 更新用户 - 使用AND逻辑，必须同时拥有权限和角色
		users.PUT("/:id",
			middleware.UnifiedAuth(middleware.AuthConfig{
				RequireAuth:  true,
				Permissions:  []string{"system:user:update"},
				Roles:        []string{"admin", "user_manager"},
				Logic:        middleware.AuthLogicAND, // 必须同时满足权限和角色
				ErrorMessage: "需要同时拥有用户更新权限和管理员角色",
			}),
			c.UserController().UpdateUser)

		// 删除用户 - 高权限操作，仅限特定角色
		users.DELETE("/:id",
			middleware.UnifiedAuth(middleware.AuthConfig{
				RequireAuth:  true,
				Permissions:  []string{"system:user:delete"},
				Roles:        []string{"admin"}, // 仅限管理员
				Logic:        middleware.AuthLogicAND,
				ErrorMessage: "仅管理员可删除用户",
			}),
			c.UserController().DeleteUser)

		// 批量操作 - 使用链式配置
		users.PUT("/batch/status",
			middleware.UnifiedAuth(
				middleware.UserManage.
					WithLogic(middleware.AuthLogicAND).
					WithCustomMessage("批量操作需要完整用户管理权限"),
			),
			c.UserController().BatchUpdateUserStatus)

		// 角色分配 - 特殊权限要求
		users.PUT("/:id/roles",
			middleware.UnifiedAuth(middleware.AuthConfig{
				RequireAuth:  true,
				Permissions:  []string{"system:user:role:assign"},
				Roles:        []string{"admin", "hr_manager"},
				Logic:        middleware.AuthLogicOR,
				ErrorMessage: "需要角色分配权限或HR管理员角色",
			}),
			c.UserController().AssignUserRoles)
	}
}

// registerRoleManagementV2 注册角色管理路由（演示不同权限策略）
func (r *AdminRoutesV2) registerRoleManagementV2(admin *gin.RouterGroup, c container.Container) {
	roles := admin.Group("/roles")
	{
		// 获取角色选择列表 - 较松的权限要求
		roles.GET("/select-list",
			middleware.UnifiedAuth(middleware.AuthConfig{
				RequireAuth: true,
				// 不设置具体权限，任何已认证用户都可访问
				ErrorMessage: "需要登录后访问",
			}),
			c.RoleController().GetEnabledRoles)

		// 角色详细管理 - 需要特定权限
		roles.GET("",
			middleware.UnifiedAuth(middleware.RoleManage),
			c.RoleController().GetRoles)
	}
}

// registerFileManagementV2 注册文件管理路由（演示文件权限控制）
func (r *AdminRoutesV2) registerFileManagementV2(admin *gin.RouterGroup, c container.Container) {
	files := admin.Group("/files")
	{
		// 查看所有文件 - 使用预定义配置
		files.GET("",
			middleware.UnifiedAuth(middleware.FileManage),
			c.FileController().GetAllFiles)

		// 删除文件 - 高权限操作
		files.DELETE("/:id",
			middleware.UnifiedAuth(middleware.AuthConfig{
				RequireAuth:  true,
				Permissions:  []string{"system:file:delete", "system:file:manage"},
				Roles:        []string{"admin", "file_manager"},
				Logic:        middleware.AuthLogicAND, // 必须同时有权限和角色
				ErrorMessage: "文件删除需要管理权限和相应角色",
			}),
			c.FileController().DeleteFile)
	}
}
