package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/pkg/database"
	"github.com/gin-gonic/gin"
)

// 使用示例

// 示例1: 使用新的模块化路由（推荐）
func ExampleModularRoutes() {
	// 创建gin引擎
	engine := gin.New()

	// 初始化容器
	db := database.GetDB()
	c := container.NewContainerCompat(db)

	// 注册所有路由模块
	RegisterAllRoutes(engine, c)

	// 启动服务器
	engine.Run(":8080")
}

// 示例2: 自定义路由模块组合
func ExampleCustomRoutes() {
	engine := gin.New()
	db := database.GetDB()
	c := container.NewContainerCompat(db)

	// 只注册特定模块
	RegisterCustomRoutes(engine, c, "public", "auth", "user")

	engine.Run(":8080")
}

// 示例3: 手动管理路由模块
func ExampleManualRoutes() {
	engine := gin.New()
	db := database.GetDB()
	c := container.NewContainerCompat(db)

	// 创建路由管理器
	manager := NewRouteManager(c)

	// 手动选择和配置模块
	manager.RegisterModule(NewPublicRoutes())
	manager.RegisterModule(NewAuthRoutes())
	manager.RegisterModule(NewUserRoutes())

	// 注册路由
	manager.RegisterRoutes(engine)

	engine.Run(":8080")
}

// CustomRoutes 自定义路由模块示例
type CustomRoutes struct {
	BaseRouteModule
}

// NewCustomRoutes 创建自定义路由模块
func NewCustomRoutes() RouteModule {
	return &CustomRoutes{
		BaseRouteModule: BaseRouteModule{
			Name:   "custom",
			Prefix: "/v1/custom",
		},
	}
}

// RegisterRoutes 注册自定义路由
func (r *CustomRoutes) RegisterRoutes(router *gin.RouterGroup, c container.Container) {
	custom := router.Group("/v1/custom")
	{
		custom.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "custom route"})
		})
	}
}

// 示例4: 添加自定义路由模块
func ExampleCustomModule() {
	// 使用自定义模块
	engine := gin.New()
	db := database.GetDB()
	c := container.NewContainerCompat(db)

	manager := NewRouteManager(c)
	manager.RegisterModule(NewCustomRoutes())
	manager.RegisterRoutes(engine)

	engine.Run(":8080")
}
