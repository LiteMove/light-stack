package routes

import (
	"github.com/LiteMove/light-stack/internal/shared/globals"
	"github.com/LiteMove/light-stack/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAnalyticsRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	// 分析模块路由（需要认证）
	analytics := v1.Group("/analytics")
	analytics.Use(middleware.Auth())
	{
		analytics.GET("/dashboard", globals.DashboardCtrl().GetDashboardStats)
		analytics.GET("/stats", globals.DashboardCtrl().GetSystemInfo)
	}
}
