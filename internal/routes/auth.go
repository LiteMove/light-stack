package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/gin-gonic/gin"
)

// AuthRoutes 认证相关路由模块
type AuthRoutes struct {
	BaseRouteModule
}

// NewAuthRoutes 创建认证路由模块
func NewAuthRoutes() RouteModule {
	return &AuthRoutes{
		BaseRouteModule: BaseRouteModule{
			Name:   "auth",
			Prefix: "/v1/auth",
		},
	}
}

// RegisterRoutes 注册认证相关路由
func (r *AuthRoutes) RegisterRoutes(router *gin.RouterGroup, c container.Container) {
	auth := router.Group("/v1/auth")
	{
		// 用户认证接口（无需认证）
		auth.POST("/login", c.AuthController().Login)          // 用户登录
		auth.POST("/register", c.AuthController().Register)    // 用户注册
		auth.POST("/refresh", c.AuthController().RefreshToken) // 刷新token
		auth.POST("/logout", c.AuthController().Logout)        // 用户登出
	}
}
