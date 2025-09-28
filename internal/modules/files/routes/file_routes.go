package routes

import (
	"github.com/LiteMove/light-stack/internal/shared/globals"
	"github.com/LiteMove/light-stack/internal/shared/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	// 文件管理路由（需要认证）
	files := v1.Group("/files")
	files.Use(middleware.Auth())
	{
		files.GET("", globals.FileCtrl().GetAllFiles)       // 获取所有文件列表
		files.GET("/:id", globals.FileCtrl().GetFile)       // 获取文件信息
		files.DELETE("/:id", globals.FileCtrl().DeleteFile) // 删除文件
	}

}
