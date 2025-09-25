package container

import (
	"github.com/LiteMove/light-stack/internal/controller"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/service"
	"gorm.io/gorm"
)

// RegisterDefaults 自动注册所有默认服务 - 修正类型定义
func (c *DIContainer) RegisterDefaults(db *gorm.DB) {
	// 注册数据库实例
	c.RegisterInstance(db)

	// === Repository 层 - 接口类型不使用指针 ===
	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return repository.NewUserRepository(Get[*gorm.DB](c))
	}, repository.UserRepository(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return repository.NewRoleRepository(Get[*gorm.DB](c))
	}, repository.RoleRepository(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return repository.NewMenuRepository(Get[*gorm.DB](c))
	}, repository.MenuRepository(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return repository.NewTenantRepository(Get[*gorm.DB](c))
	}, repository.TenantRepository(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return repository.NewFileRepository(Get[*gorm.DB](c))
	}, (*repository.FileRepository)(nil))

	// === Service 层 - 接口类型不使用指针 ===
	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return service.NewAuthService(
			MustGet[repository.UserRepository](c),
			MustGet[repository.RoleRepository](c),
			MustGet[repository.MenuRepository](c),
		)
	}, service.AuthService(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return service.NewUserService(
			MustGet[repository.UserRepository](c),
			MustGet[repository.RoleRepository](c),
		)
	}, service.UserService(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return service.NewRoleService(
			MustGet[repository.RoleRepository](c),
			MustGet[repository.UserRepository](c),
		)
	}, service.RoleService(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return service.NewMenuService(
			MustGet[repository.MenuRepository](c),
			MustGet[repository.RoleRepository](c),
		)
	}, service.MenuService(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return service.NewTenantService(
			MustGet[repository.TenantRepository](c),
			MustGet[repository.UserRepository](c),
		)
	}, service.TenantService(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return service.NewFileService(
			MustGet[*repository.FileRepository](c),
			MustGet[service.TenantService](c),
		)
	}, (*service.FileService)(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return service.NewProfileService(
			MustGet[repository.UserRepository](c),
			MustGet[repository.RoleRepository](c),
			MustGet[repository.TenantRepository](c),
		)
	}, service.ProfileService(nil))

	// === Controller 层 ===
	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return controller.NewAuthController(
			MustGet[service.AuthService](c),
			MustGet[service.RoleService](c),
			MustGet[service.MenuService](c),
		)
	}, (*controller.AuthController)(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return controller.NewUserController(
			MustGet[service.UserService](c),
		)
	}, (*controller.UserController)(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return controller.NewRoleController(
			MustGet[service.RoleService](c),
		)
	}, (*controller.RoleController)(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return controller.NewMenuController(
			MustGet[service.MenuService](c),
		)
	}, (*controller.MenuController)(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return controller.NewTenantController(
			MustGet[service.TenantService](c),
		)
	}, (*controller.TenantController)(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return controller.NewFileController(
			MustGet[*service.FileService](c),
		)
	}, (*controller.FileController)(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return controller.NewProfileController(
			MustGet[service.ProfileService](c),
		)
	}, (*controller.ProfileController)(nil))

	c.RegisterSingleton(func(c *DIContainer) interface{} {
		return controller.NewHealthController()
	}, (*controller.HealthController)(nil))
}

// 全局容器实例
var DefaultContainer = NewDI()

// 便捷函数
func RegisterDefaults(db *gorm.DB) {
	DefaultContainer.RegisterDefaults(db)
}

func GetService[T any]() T {
	return Get[T](DefaultContainer)
}

func MustGetService[T any]() T {
	return MustGet[T](DefaultContainer)
}
