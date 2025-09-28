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
		files.GET("", globals.FileCtrl().GetAllFiles)                // 获取所有文件列表
		files.GET("/:id", globals.FileCtrl().GetFile)                // 获取文件信息
		files.GET("/:id/private", globals.FileCtrl().GetPrivateFile) // 获取私有文件内容
		files.DELETE("/:id", globals.FileCtrl().DeleteFile)          // 删除文件
		files.POST("/upload", globals.FileCtrl().UploadFile)         // 上传文件
		files.GET("/user", globals.FileCtrl().GetUserFiles)          // 获取用户文件列表
	}

}
