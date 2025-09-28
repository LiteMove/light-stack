package routes

import (
	"github.com/LiteMove/light-stack/internal/shared/globals"
	"github.com/LiteMove/light-stack/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	// 认证相关路由（无需认证）
	auth := v1.Group("/auth")
	{
		auth.POST("/login", globals.AuthCtrl().Login)
		auth.POST("/register", globals.AuthCtrl().Register)
		auth.POST("/refresh", globals.AuthCtrl().RefreshToken)
		auth.POST("/logout", globals.AuthCtrl().Logout)
	}

	// 用户档案路由（需要认证）
	profile := v1.Group("/profile")
	profile.Use(middleware.Auth())
	{
		profile.GET("", globals.AuthCtrl().GetProfile)
		profile.PUT("", globals.AuthCtrl().UpdateProfile)
		profile.PUT("/password", globals.AuthCtrl().ChangePassword)
		// 用户查看和修改租户配置
		profile.Use(middleware.TenantMiddleware(globals.TenantSvc()))
		{
			profile.GET("/tenant-config", globals.ProfileCtrl().GetTenantConfig)    // 获取所在租户配置
			profile.PUT("/tenant-config", globals.ProfileCtrl().UpdateTenantConfig) // 更新所在租户配置
		}
	}
}
