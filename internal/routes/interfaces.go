package routes

import "github.com/gin-gonic/gin"

// ModuleInterface 模块接口
type ModuleInterface interface {
	// InitModule 初始化模块
	InitModule() error

	// RegisterRoutes 注册路由
	RegisterRoutes(r *gin.RouterGroup)

	// GetModuleName 获取模块名称
	GetModuleName() string

	// GetModuleVersion 获取模块版本
	GetModuleVersion() string
}
