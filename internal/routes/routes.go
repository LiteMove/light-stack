package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RouteManager 路由管理器
type RouteManager struct {
	modules   []RouteModule
	container container.Container
}

// NewRouteManager 创建路由管理器
func NewRouteManager(c container.Container) *RouteManager {
	return &RouteManager{
		modules:   make([]RouteModule, 0),
		container: c,
	}
}

// RegisterModule 注册路由模块
func (rm *RouteManager) RegisterModule(module RouteModule) {
	rm.modules = append(rm.modules, module)
}

// RegisterModules 批量注册路由模块
func (rm *RouteManager) RegisterModules(modules ...RouteModule) {
	for _, module := range modules {
		rm.RegisterModule(module)
	}
}

// RegisterRoutes 注册所有路由模块到gin引擎
func (rm *RouteManager) RegisterRoutes(engine *gin.Engine) {
	// API 分组
	api := engine.Group("/api")
	// 可选登录中间件，不强制登录
	api.Use(middleware.OptionalJWTAuthMiddleware())
	// 租户中间件，根据域名识别租户
	api.Use(middleware.TenantMiddleware(rm.container.TenantService()))

	// 注册所有模块路由
	for _, module := range rm.modules {
		module.RegisterRoutes(api, rm.container)
	}

	// 静态文件服务（直接在引擎上注册）
	RegisterStaticFiles(engine, rm.container)
}

// RegisterAllRoutes 注册所有默认路由模块
func RegisterAllRoutes(engine *gin.Engine, c container.Container) {
	// 创建路由管理器
	manager := NewRouteManager(c)

	// 注册所有模块
	manager.RegisterModules(
		NewPublicRoutes(),     // 公开接口
		NewAuthRoutes(),       // 认证相关
		NewUserRoutes(),       // 用户相关
		NewAdminRoutes(),      // 管理员
		NewSuperAdminRoutes(), // 超级管理员
		NewStaticRoutes(),     // 静态文件
	)

	// 注册路由到引擎
	manager.RegisterRoutes(engine)
}

// RegisterRoutesWithContainer 便捷函数：使用新的DI容器注册路由
func RegisterRoutesWithContainer(engine *gin.Engine) {
	// 使用全局默认容器
	c := container.NewContainerCompat(nil) // 这里需要传入数据库实例
	RegisterAllRoutes(engine, c)
}

// GetAvailableModules 获取所有可用的路由模块
func GetAvailableModules() []RouteModule {
	return []RouteModule{
		NewPublicRoutes(),
		NewAuthRoutes(),
		NewUserRoutes(),
		NewAdminRoutes(),
		NewSuperAdminRoutes(),
		NewStaticRoutes(),
	}
}

// RegisterCustomRoutes 自定义路由注册（用于灵活配置）
func RegisterCustomRoutes(engine *gin.Engine, c container.Container, moduleNames ...string) {
	manager := NewRouteManager(c)
	allModules := GetAvailableModules()

	// 根据名称筛选模块
	if len(moduleNames) == 0 {
		// 如果没有指定模块，注册所有模块
		manager.RegisterModules(allModules...)
	} else {
		// 只注册指定的模块
		for _, name := range moduleNames {
			for _, module := range allModules {
				if module.GetName() == name {
					manager.RegisterModule(module)
					break
				}
			}
		}
	}

	// 注册路由到引擎
	manager.RegisterRoutes(engine)
}
