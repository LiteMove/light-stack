package container

import (
	"github.com/LiteMove/light-stack/internal/controller"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/service"
	"gorm.io/gorm"
)

// SimpleContainer 简化版容器，保持兼容性
type SimpleContainer struct {
	di *DIContainer
}

// NewSimpleContainer 创建简化版容器
func NewSimpleContainer(db *gorm.DB) *SimpleContainer {
	di := NewDI()
	di.RegisterDefaults(db)

	return &SimpleContainer{
		di: di,
	}
}

// 实现 Container 接口 - 修正类型定义

// Repository 工厂方法
func (c *SimpleContainer) UserRepository() repository.UserRepository {
	return Get[repository.UserRepository](c.di)
}

func (c *SimpleContainer) RoleRepository() repository.RoleRepository {
	return Get[repository.RoleRepository](c.di)
}

func (c *SimpleContainer) MenuRepository() repository.MenuRepository {
	return Get[repository.MenuRepository](c.di)
}

func (c *SimpleContainer) TenantRepository() repository.TenantRepository {
	return Get[repository.TenantRepository](c.di)
}

func (c *SimpleContainer) FileRepository() *repository.FileRepository {
	return Get[*repository.FileRepository](c.di)
}

// Service 工厂方法
func (c *SimpleContainer) AuthService() service.AuthService {
	return Get[service.AuthService](c.di)
}

func (c *SimpleContainer) UserService() service.UserService {
	return Get[service.UserService](c.di)
}

func (c *SimpleContainer) RoleService() service.RoleService {
	return Get[service.RoleService](c.di)
}

func (c *SimpleContainer) MenuService() service.MenuService {
	return Get[service.MenuService](c.di)
}

func (c *SimpleContainer) TenantService() service.TenantService {
	return Get[service.TenantService](c.di)
}

func (c *SimpleContainer) FileService() *service.FileService {
	return Get[*service.FileService](c.di)
}

func (c *SimpleContainer) ProfileService() service.ProfileService {
	return Get[service.ProfileService](c.di)
}

// Controller 工厂方法
func (c *SimpleContainer) AuthController() *controller.AuthController {
	return Get[*controller.AuthController](c.di)
}

func (c *SimpleContainer) UserController() *controller.UserController {
	return Get[*controller.UserController](c.di)
}

func (c *SimpleContainer) RoleController() *controller.RoleController {
	return Get[*controller.RoleController](c.di)
}

func (c *SimpleContainer) MenuController() *controller.MenuController {
	return Get[*controller.MenuController](c.di)
}

func (c *SimpleContainer) TenantController() *controller.TenantController {
	return Get[*controller.TenantController](c.di)
}

func (c *SimpleContainer) FileController() *controller.FileController {
	return Get[*controller.FileController](c.di)
}

func (c *SimpleContainer) ProfileController() *controller.ProfileController {
	return Get[*controller.ProfileController](c.di)
}

func (c *SimpleContainer) HealthController() *controller.HealthController {
	return Get[*controller.HealthController](c.di)
}

// NewContainerCompat 工厂函数，保持向后兼容性
func NewContainerCompat(db *gorm.DB) Container {
	return NewSimpleContainer(db)
}
