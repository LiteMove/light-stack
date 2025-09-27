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
		files.POST("/upload", globals.FileCtrl().UploadFile)
		files.GET("", globals.FileCtrl().GetAllFiles)
		files.GET("/:id", globals.FileCtrl().GetFile)
		files.DELETE("/:id", globals.FileCtrl().DeleteFile)
		files.GET("/:id/download", globals.FileCtrl().GetPrivateFile)
	}
}
