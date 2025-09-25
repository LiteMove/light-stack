package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/gin-gonic/gin"
)

// UserRoutes 用户相关路由模块
type UserRoutes struct {
	BaseRouteModule
}

// NewUserRoutes 创建用户路由模块
func NewUserRoutes() RouteModule {
	return &UserRoutes{
		BaseRouteModule: BaseRouteModule{
			Name:   "user",
			Prefix: "/v1/user",
		},
	}
}

// RegisterRoutes 注册用户相关路由
func (r *UserRoutes) RegisterRoutes(router *gin.RouterGroup, c container.Container) {
	// 用户相关路由（需要认证）
	user := router.Group("/v1/user")
	user.Use(middleware.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		// 用户个人信息管理
		user.GET("/profile", c.AuthController().GetProfile)      // 获取用户信息
		user.PUT("/profile", c.AuthController().UpdateProfile)   // 更新用户信息
		user.PUT("/password", c.AuthController().ChangePassword) // 修改密码
		user.GET("/:id/roles", c.AuthController().GetUserRoles)  // 获取用户角色
	}

	// 用户租户配置管理（需要租户中间件）
	tenantConfig := user.Group("")
	tenantConfig.Use(middleware.TenantMiddleware(c.TenantService()))
	{
		tenantConfig.GET("/tenant-config", c.ProfileController().GetTenantConfig)    // 获取所在租户配置
		tenantConfig.PUT("/tenant-config", c.ProfileController().UpdateTenantConfig) // 更新所在租户配置
	}

	// 文件相关路由（需要认证）
	files := router.Group("/v1/files")
	files.Use(middleware.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		files.POST("/upload", c.FileController().UploadFile)         // 文件上传
		files.GET("", c.FileController().GetAllFiles)                // 获取文件列表（按租户）
		files.GET("/list", c.FileController().GetAllFiles)           // 获取文件列表（兼容旧接口）
		files.GET("/:id", c.FileController().GetFile)                // 获取文件信息（包含access_url用于下载/预览）
		files.GET("/:id/private", c.FileController().GetPrivateFile) // 获取私有文件内容
		files.DELETE("/:id", c.FileController().DeleteFile)          // 删除文件
	}
}
