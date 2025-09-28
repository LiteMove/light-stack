package routes

import (
	"github.com/LiteMove/light-stack/internal/shared/globals"
	"github.com/LiteMove/light-stack/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterGeneratorRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	// 代码生成器路由（需要认证）
	generator := v1.Group("/gen")
	generator.Use(middleware.Auth())
	{
		// 数据库表分析
		generator.GET("/tables", globals.GeneratorCtrl().GetTableList)                       // 获取数据库表列表
		generator.GET("/tables/:tableName", globals.GeneratorCtrl().GetTableInfo)            // 获取表结构信息
		generator.GET("/tables/:tableName/columns", globals.GeneratorCtrl().GetTableColumns) // 获取表字段信息

		// 获取系统菜单树（用于选择父级菜单）
		generator.GET("/menus/tree", globals.GeneratorCtrl().GetSystemMenus) // 获取系统菜单树

		// 生成配置管理
		generator.POST("/configs", globals.GenConfigCtrl().CreateConfig)                         // 创建生成配置
		generator.GET("/configs", globals.GenConfigCtrl().GetConfigList)                         // 获取配置列表
		generator.GET("/configs/:id", globals.GenConfigCtrl().GetConfig)                         // 获取配置详情
		generator.PUT("/configs/:id", globals.GenConfigCtrl().UpdateConfig)                      // 更新配置
		generator.DELETE("/configs/:id", globals.GenConfigCtrl().DeleteConfig)                   // 删除配置
		generator.POST("/configs/import/:tableName", globals.GenConfigCtrl().ImportTableConfig)  // 导入表配置
		generator.GET("/configs/table/:tableName", globals.GenConfigCtrl().GetConfigByTableName) // 根据表名获取配置

		// 代码生成
		generator.GET("/preview/:configId", globals.GeneratorCtrl().PreviewCode)   // 临时预览接口
		generator.POST("/generate", globals.GeneratorCtrl().GenerateCode)          // 生成代码
		generator.GET("/download/:taskId", globals.GeneratorCtrl().DownloadCode)   // 下载代码包
		generator.GET("/templates", globals.GeneratorCtrl().GetAvailableTemplates) // 获取可用模板

		// 生成历史
		generator.GET("/history", globals.GeneratorCtrl().GetHistory)           // 获取生成历史
		generator.DELETE("/history/:id", globals.GeneratorCtrl().DeleteHistory) // 删除生成历史
	}
}
