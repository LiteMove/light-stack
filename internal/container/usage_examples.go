package container

import (
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/database"
)

// 新的简化DI容器使用示例

// 示例1: 使用全局默认容器（最简单）
func SimpleUsage() {
	// 初始化
	db := database.GetDB()
	RegisterDefaults(db)

	// 直接获取服务，无需任何工厂方法
	userService := GetService[service.UserService]()
	authService := MustGetService[service.AuthService]()

	// 使用服务
	_ = userService
	_ = authService
}

// 示例2: 使用局部容器实例
func LocalContainerUsage() {
	db := database.GetDB()

	// 创建容器并自动注册所有服务
	container := NewDI()
	container.RegisterDefaults(db)

	// 获取服务
	userService := Get[service.UserService](container)
	authService := MustGet[service.AuthService](container)

	_ = userService
	_ = authService
}

// 示例3: 保持向后兼容性
func BackwardCompatibleUsage() {
	db := database.GetDB()

	// 使用原有的接口，但内部使用新的DI容器
	container := NewContainerCompat(db)

	// 原有代码无需修改
	userController := container.UserController()
	authService := container.AuthService()

	_ = userController
	_ = authService
}

// 对比：旧版本 vs 新版本

// 旧版本：需要写大量工厂方法
/*
func (c *container) UserService() service.UserService {
	c.mu.RLock()
	if c.userService != nil {
		c.mu.RUnlock()
		return c.userService
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.userService == nil {
		c.userService = service.NewUserService(
			c.UserRepository(),
			c.RoleRepository(),
		)
	}
	return c.userService
}
*/

// 新版本：只需要一次注册
/*
container.RegisterSingleton(func(c *DIContainer) interface{} {
	return service.NewUserService(
		Get[repository.UserRepository](c),
		Get[repository.RoleRepository](c),
	)
}, (*service.UserService)(nil))
*/
