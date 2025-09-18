package controller

import (
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/database"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Check 健康检查
func (h *HealthController) Check(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ok",
		"service": "light-stack",
	})
}

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine) {
	// 初始化依赖
	db := database.GetDB()
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	authService := service.NewAuthService(userRepo, roleRepo)
	roleService := service.NewRoleService(roleRepo, userRepo)
	authController := NewAuthController(authService, roleService)
	roleController := NewRoleController(roleService)
	healthController := NewHealthController()

	// API 分组
	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/health", healthController.Check)

		// V1 版本API
		v1 := api.Group("/v1")
		{
			// 基础路由组
			v1.GET("/ping", func(c *gin.Context) {
				response.Success(c, gin.H{
					"message": "pong",
				})
			})

			// 认证相关路由（无需认证）
			auth := v1.Group("/auth")
			{
				auth.POST("/login", authController.Login)       // 用户登录
				auth.POST("/register", authController.Register) // 用户注册
				auth.POST("/refresh", authController.RefreshToken) // 刷新token
				auth.POST("/logout", authController.Logout)     // 用户登出
			}

			// 用户相关路由（需要认证）
			user := v1.Group("/user")
			user.Use(middleware.JWTAuthMiddleware()) // 应用JWT认证中间件
			{
				user.GET("/profile", authController.GetProfile)       // 获取用户信息
				user.PUT("/profile", authController.UpdateProfile)    // 更新用户信息
				user.PUT("/password", authController.ChangePassword)  // 修改密码
				user.GET("/:id/roles", authController.GetUserRoles)   // 获取用户角色
			}

			// 管理员相关路由（需要管理员权限）
			admin := v1.Group("/admin")
			admin.Use(middleware.JWTAuthMiddleware())    // 应用JWT认证中间件
			admin.Use(middleware.AdminAuthMiddleware())  // 应用管理员权限中间件
			{
				// 角色管理
				roles := admin.Group("/roles")
				{
					roles.POST("", roleController.CreateRole)        // 创建角色
					roles.GET("", roleController.GetRoles)           // 获取角色列表
					roles.GET("/:id", roleController.GetRole)        // 获取角色详情
					roles.PUT("/:id", roleController.UpdateRole)     // 更新角色
					roles.DELETE("/:id", roleController.DeleteRole)  // 删除角色
				}

				// 用户角色管理
				users := admin.Group("/users")
				{
					users.PUT("/:id/roles", authController.AssignUserRoles) // 为用户分配角色
				}
			}
		}
	}
}