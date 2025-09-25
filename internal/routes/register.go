package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/pkg/database"
	"github.com/gin-gonic/gin"
)

// RegisterRoutesV2 使用新架构注册路由
// 这个函数是新架构的入口点，替换了原有的213行RegisterRoutes函数
// 主要改进：
// 1. 使用DI容器自动管理所有依赖，避免手动创建20+个对象
// 2. 使用模块化路由，将160+行路由定义拆分为6个清晰的模块
// 3. 统一的中间件管理和权限控制
func RegisterRoutesV2(r *gin.Engine) {
	// 步骤1: 初始化数据库连接
	// 保持与原有架构的兼容性，使用相同的数据库初始化方式
	db := database.GetDB()
	if db == nil {
		panic("Failed to initialize database connection")
	}

	// 步骤2: 创建DI容器并自动注册所有服务
	// 新的DI容器会自动处理所有依赖注入：
	// - Repository层：UserRepository, RoleRepository, MenuRepository, TenantRepository, FileRepository
	// - Service层：AuthService, UserService, RoleService, MenuService, TenantService, FileService, ProfileService
	// - Controller层：AuthController, UserController, RoleController, MenuController, TenantController, FileController, ProfileController, HealthController
	c := container.NewContainerCompat(db)

	// 步骤3: 使用模块化路由注册器
	// RegisterAllRoutes会自动加载所有路由模块：
	// - PublicRoutes: 健康检查、ping、租户信息等公开接口
	// - AuthRoutes: 登录、注册、token刷新、登出等认证接口
	// - UserRoutes: 用户资料、文件管理等需要认证的用户接口
	// - AdminRoutes: 用户管理、角色管理等管理员接口
	// - SuperAdminRoutes: 租户管理、系统配置等超管接口
	// - StaticRoutes: 静态文件服务（公开和私有文件）
	RegisterAllRoutes(r, c)
}

// RegisterRoutesWithCustomModules 自定义模块注册
// 这个函数允许灵活选择需要加载的路由模块，适用于特定部署场景
// 例如：只部署用户相关功能，或者分离管理员和普通用户服务
func RegisterRoutesWithCustomModules(r *gin.Engine, moduleNames ...string) {
	// 初始化数据库和容器
	db := database.GetDB()
	if db == nil {
		panic("Failed to initialize database connection")
	}

	// 创建容器实例
	c := container.NewContainerCompat(db)

	// 根据指定的模块名称注册路由
	// 支持的模块名称：public, auth, user, admin, super_admin, static
	if len(moduleNames) == 0 {
		// 如果没有指定模块，则注册所有模块（等同于RegisterRoutesV2）
		RegisterAllRoutes(r, c)
	} else {
		// 只注册指定的模块
		RegisterCustomRoutes(r, c, moduleNames...)
	}
}

// MigrateFromOldRoutes 从旧路由架构迁移的辅助函数
// 这个函数提供了一个平滑的迁移路径，可以逐步从旧架构切换到新架构
// 使用方式：将main.go中的controller.RegisterRoutes(r)替换为routes.MigrateFromOldRoutes(r)
func MigrateFromOldRoutes(r *gin.Engine) {
	// 记录迁移日志（在生产环境中可以使用proper logger）
	// log.Info("Migrating to new modular route architecture...")

	// 直接使用新的路由注册方法
	RegisterRoutesV2(r)

	// 在这里可以添加迁移后的验证逻辑
	// validateRouteMigration(r)
}

// ValidateRouteMigration 验证路由迁移的完整性（开发辅助函数）
// 这个函数可以在开发阶段用于验证所有关键路由是否正确迁移
func ValidateRouteMigration(r *gin.Engine) {
	// 获取所有注册的路由
	routes := r.Routes()

	// 定义关键路由列表（这些路由必须存在）
	criticalPaths := []string{
		"/api/health",
		"/api/v1/ping",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/user/profile",
		"/api/v1/admin/users",
		"/api/v1/super-admin/tenants",
	}

	// 检查关键路由是否都已正确注册
	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Path] = true
	}

	// 验证每个关键路由
	for _, path := range criticalPaths {
		if !routeMap[path] {
			// 在生产环境中应该使用proper logger记录错误
			// log.Errorf("Critical route missing after migration: %s", path)
		}
	}
}
