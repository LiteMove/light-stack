package container

import (
	"github.com/LiteMove/light-stack/internal/controller"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/service"
)

// Container 依赖注入容器接口
type Container interface {
	// Repository 工厂方法 - 接口类型不使用指针
	UserRepository() repository.UserRepository
	RoleRepository() repository.RoleRepository
	MenuRepository() repository.MenuRepository
	TenantRepository() repository.TenantRepository
	FileRepository() *repository.FileRepository

	// Service 工厂方法 - 接口类型不使用指针
	AuthService() service.AuthService
	UserService() service.UserService
	RoleService() service.RoleService
	MenuService() service.MenuService
	TenantService() service.TenantService
	FileService() *service.FileService
	ProfileService() service.ProfileService

	// Controller 工厂方法 - 结构体使用指针
	AuthController() *controller.AuthController
	UserController() *controller.UserController
	RoleController() *controller.RoleController
	MenuController() *controller.MenuController
	TenantController() *controller.TenantController
	FileController() *controller.FileController
	ProfileController() *controller.ProfileController
	HealthController() *controller.HealthController
}
