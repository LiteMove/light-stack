package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// PublicRoutes 公开接口路由模块
type PublicRoutes struct {
	BaseRouteModule
}

// NewPublicRoutes 创建公开接口路由模块
func NewPublicRoutes() RouteModule {
	return &PublicRoutes{
		BaseRouteModule: BaseRouteModule{
			Name:   "public",
			Prefix: "",
		},
	}
}

// RegisterRoutes 注册公开接口路由
func (r *PublicRoutes) RegisterRoutes(router *gin.RouterGroup, c container.Container) {
	// 健康检查
	router.GET("/health", r.healthCheck)

	// 租户信息接口（公开，无需认证）
	router.GET("/tenant/info", c.TenantController().GetTenantByDomain)

	// V1 版本API
	v1 := router.Group("/v1")
	{
		// 基础路由组
		v1.GET("/ping", r.ping)
	}
}

// healthCheck 健康检查
func (r *PublicRoutes) healthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ok",
		"service": "light-stack",
	})
}

// ping 基础连接测试
func (r *PublicRoutes) ping(c *gin.Context) {
	response.Success(c, gin.H{
		"message": "pong",
	})
}
