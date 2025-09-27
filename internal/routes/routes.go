package routes

import (
	"log"

	"github.com/LiteMove/light-stack/internal/config"
	"github.com/LiteMove/light-stack/internal/modules/analytics"
	"github.com/LiteMove/light-stack/internal/modules/auth"
	"github.com/LiteMove/light-stack/internal/modules/files"
	"github.com/LiteMove/light-stack/internal/modules/generator"
	"github.com/LiteMove/light-stack/internal/modules/system"
	"github.com/LiteMove/light-stack/internal/shared/globals"
	"github.com/LiteMove/light-stack/internal/shared/middleware"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有模块路由
func RegisterRoutes(r *gin.Engine) {
	// 初始化所有服务
	globals.Init()

	// 初始化所有模块
	modules := []ModuleInterface{
		auth.NewModule(),
		system.NewModule(),
		files.NewModule(),
		generator.NewModule(),
		analytics.NewModule(),
	}

	// 初始化模块
	for _, module := range modules {
		if err := module.InitModule(); err != nil {
			log.Printf("Failed to initialize module %s: %v", module.GetModuleName(), err)
		}
	}

	// API 分组
	api := r.Group("/api")

	// 全局中间件
	api.Use(middleware.OptionalAuth())                        // 可选认证
	api.Use(middleware.TenantMiddleware(globals.TenantSvc())) // 租户处理

	// 注册公共路由
	registerPublicRoutes(api)

	// 注册各模块路由
	for _, module := range modules {
		module.RegisterRoutes(api)
	}

	// 注册静态文件路由
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

	// V1 版本API
	v1 := api.Group("/v1")
	{
		// 基础路由组
		v1.GET("/ping", func(c *gin.Context) {
			response.Success(c, gin.H{
				"message": "pong",
			})
		})
	}

	// 租户信息接口（公开，无需认证）
	api.GET("/tenant/info", globals.TenantCtrl().GetTenantByDomain)
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
