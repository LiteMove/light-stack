package routes

import (
	"github.com/LiteMove/light-stack/internal/container"
	"github.com/gin-gonic/gin"
)

// RouteRegistrar 路由注册器接口
type RouteRegistrar interface {
	RegisterRoutes(group *gin.RouterGroup, container container.Container)
}

// RouterConfig 路由配置
type RouterConfig struct {
	Container container.Container
	Engine    *gin.Engine
}

// RouteModule 路由模块接口
type RouteModule interface {
	// GetPrefix 获取路由前缀
	GetPrefix() string

	// RegisterRoutes 注册路由
	RegisterRoutes(router *gin.RouterGroup, container container.Container)

	// GetName 获取模块名称
	GetName() string
}

// BaseRouteModule 基础路由模块
type BaseRouteModule struct {
	Name   string
	Prefix string
}

// GetName 获取模块名称
func (b *BaseRouteModule) GetName() string {
	return b.Name
}

// GetPrefix 获取路由前缀
func (b *BaseRouteModule) GetPrefix() string {
	return b.Prefix
}
