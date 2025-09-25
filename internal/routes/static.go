package routes

import (
	"github.com/LiteMove/light-stack/internal/config"
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/gin-gonic/gin"
)

// StaticRoutes 静态文件路由模块
type StaticRoutes struct {
	BaseRouteModule
}

// NewStaticRoutes 创建静态文件路由模块
func NewStaticRoutes() RouteModule {
	return &StaticRoutes{
		BaseRouteModule: BaseRouteModule{
			Name:   "static",
			Prefix: "/static",
		},
	}
}

// RegisterRoutes 注册静态文件路由
func (r *StaticRoutes) RegisterRoutes(router *gin.RouterGroup, c container.Container) {
	// 静态文件服务 - 使用配置文件中的base_url
	cfg := config.Get()
	baseURL := cfg.File.BaseURL
	if baseURL == "" {
		baseURL = "/api/static"
	}

	// 公开文件 - 无需认证
	r.registerPublicFiles(router, baseURL)

	// 私有文件 - 需要认证和权限验证
	r.registerPrivateFiles(router, baseURL, c)
}

// registerPublicFiles 注册公开文件路由
func (r *StaticRoutes) registerPublicFiles(router *gin.RouterGroup, baseURL string) {
	// 直接在引擎上注册静态文件服务
	// 注意：这里需要在主路由注册器中处理，因为Static方法需要在gin.Engine上调用
}

// registerPrivateFiles 注册私有文件路由
func (r *StaticRoutes) registerPrivateFiles(router *gin.RouterGroup, baseURL string, c container.Container) {
	privateFiles := router.Group(baseURL + "/private")
	privateFiles.Use(middleware.JWTAuthMiddleware())
	privateFiles.Use(middleware.TenantMiddleware(c.TenantService()))
	{
		// 私有文件访问需要验证用户是否有权限访问该租户的文件
		privateFiles.StaticFS("/", gin.Dir("./uploads/private", false))
	}
}

// RegisterStaticFiles 在gin.Engine上注册静态文件服务
func RegisterStaticFiles(engine *gin.Engine, c container.Container) {
	// 静态文件服务 - 使用配置文件中的base_url
	cfg := config.Get()
	baseURL := cfg.File.BaseURL
	if baseURL == "" {
		baseURL = "/api/static"
	}

	// 公开文件 - 无需认证
	engine.Static(baseURL+"/public", "./uploads/public")
}
