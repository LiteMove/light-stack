package routes

import (
	"github.com/LiteMove/light-stack/internal/shared/globals"
	"github.com/LiteMove/light-stack/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterGeneratorRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	// 代码生成器路由（需要认证）
	generator := v1.Group("/generator")
	generator.Use(middleware.Auth())
	{
		generator.GET("/tables", globals.GeneratorCtrl().GetTableList)
		generator.GET("/tables/:table/columns", globals.GeneratorCtrl().GetTableColumns)
		generator.POST("/generate", globals.GeneratorCtrl().GenerateCode)
		generator.GET("/config", globals.GenConfigCtrl().GetConfigList)
		generator.POST("/config", globals.GenConfigCtrl().CreateConfig)
		generator.PUT("/config/:id", globals.GenConfigCtrl().UpdateConfig)
		generator.DELETE("/config/:id", globals.GenConfigCtrl().DeleteConfig)
	}
}
