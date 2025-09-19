package main

import (
	"fmt"
	"log"

	"github.com/LiteMove/light-stack/internal/config"
	"github.com/LiteMove/light-stack/internal/controller"
	"github.com/LiteMove/light-stack/internal/middleware"
	"github.com/LiteMove/light-stack/internal/utils"
	"github.com/LiteMove/light-stack/pkg/cache"
	"github.com/LiteMove/light-stack/pkg/database"
	"github.com/LiteMove/light-stack/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化日志
	logger.Init()

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化 Redis
	if err := cache.Init(); err != nil {
		log.Fatal("Failed to initialize cache:", err)
	}

	// 设置 Gin 模式
	if config.Get().App.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 引擎
	r := gin.New()

	// 设置中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	r.Use(middleware.ResponseMiddleware())

	// 注册路由
	controller.RegisterRoutes(r)

	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Hashed Password: %s\n", hashedPassword)

	// 启动服务器
	port := config.Get().Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
