package globals

import (
	"github.com/LiteMove/light-stack/internal/controller"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/database"
	"gorm.io/gorm"
)

// 全局服务实例
var (
	// Repository 层
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	menuRepo   repository.MenuRepository
	tenantRepo repository.TenantRepository
	fileRepo   *repository.FileRepository

	// Service 层
	authSvc    service.AuthService
	userSvc    service.UserService
	roleSvc    service.RoleService
	menuSvc    service.MenuService
	tenantSvc  service.TenantService
	fileSvc    *service.FileService
	profileSvc service.ProfileService

	// Controller 层
	authCtrl    *controller.AuthController
	userCtrl    *controller.UserController
	roleCtrl    *controller.RoleController
	menuCtrl    *controller.MenuController
	tenantCtrl  *controller.TenantController
	fileCtrl    *controller.FileController
	profileCtrl *controller.ProfileController
)

// Init 初始化所有服务
func Init() {
	db := database.GetDB()
	if db == nil {
		panic("database not initialized")
	}

	initRepositories(db)
	initServices()
	initControllers()
}

func initRepositories(db *gorm.DB) {
	userRepo = repository.NewUserRepository(db)
	roleRepo = repository.NewRoleRepository(db)
	menuRepo = repository.NewMenuRepository(db)
	tenantRepo = repository.NewTenantRepository(db)
	fileRepo = repository.NewFileRepository(db)
}

func initServices() {
	authSvc = service.NewAuthService(userRepo, roleRepo, menuRepo)
	userSvc = service.NewUserService(userRepo, roleRepo)
	roleSvc = service.NewRoleService(roleRepo, userRepo)
	menuSvc = service.NewMenuService(menuRepo, roleRepo)
	tenantSvc = service.NewTenantService(tenantRepo, userRepo)
	profileSvc = service.NewProfileService(userRepo, roleRepo, tenantRepo)
	fileSvc = service.NewFileService(fileRepo, tenantSvc)
}

func initControllers() {
	authCtrl = controller.NewAuthController(authSvc, roleSvc, menuSvc)
	userCtrl = controller.NewUserController(userSvc)
	roleCtrl = controller.NewRoleController(roleSvc)
	menuCtrl = controller.NewMenuController(menuSvc)
	tenantCtrl = controller.NewTenantController(tenantSvc)
	fileCtrl = controller.NewFileController(fileSvc)
	profileCtrl = controller.NewProfileController(profileSvc)
}

// === Service 获取函数 ===
func AuthSvc() service.AuthService       { return authSvc }
func UserSvc() service.UserService       { return userSvc }
func RoleSvc() service.RoleService       { return roleSvc }
func MenuSvc() service.MenuService       { return menuSvc }
func TenantSvc() service.TenantService   { return tenantSvc }
func FileSvc() *service.FileService      { return fileSvc }
func ProfileSvc() service.ProfileService { return profileSvc }

// === Controller 获取函数 ===
func AuthCtrl() *controller.AuthController       { return authCtrl }
func UserCtrl() *controller.UserController       { return userCtrl }
func RoleCtrl() *controller.RoleController       { return roleCtrl }
func MenuCtrl() *controller.MenuController       { return menuCtrl }
func TenantCtrl() *controller.TenantController   { return tenantCtrl }
func FileCtrl() *controller.FileController       { return fileCtrl }
func ProfileCtrl() *controller.ProfileController { return profileCtrl }

// === 权限检查函数 ===
func CheckUserRole(userID uint64, roleCode string) bool {
	// 超级管理员拥有所有权限
	if roleCode == "super_admin" {
		roles, err := userSvc.GetUserRoles(userID)
		if err != nil {
			return false
		}
		for _, role := range roles {
			if role.Code == "super_admin" {
				return true
			}
		}
		return false
	}

	// 检查是否有指定角色或超级管理员权限
	roles, err := userSvc.GetUserRoles(userID)
	if err != nil {
		return false
	}

	for _, role := range roles {
		if role.Code == "super_admin" || role.Code == roleCode {
			return true
		}
	}
	return false
}
