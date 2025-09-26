package globals

import (
	"github.com/LiteMove/light-stack/internal/controller"
	"github.com/LiteMove/light-stack/internal/generator"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/database"
	"gorm.io/gorm"
)

// 全局服务实例
var (
	// Repository 层
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	menuRepo       repository.MenuRepository
	tenantRepo     repository.TenantRepository
	fileRepo       *repository.FileRepository
	dictRepo       repository.DictRepository
	dbAnalyzerRepo *repository.DBAnalyzerRepository
	genConfigRepo  *repository.GenConfigRepository

	// Generator 层
	templateEngine *generator.TemplateEngine
	codeGenerator  *generator.CodeGenerator
	filePackager   *generator.FilePackager

	// Service 层
	authSvc       service.AuthService
	userSvc       service.UserService
	roleSvc       service.RoleService
	menuSvc       service.MenuService
	tenantSvc     service.TenantService
	fileSvc       *service.FileService
	profileSvc    service.ProfileService
	dashboardSvc  service.DashboardService
	dictSvc       service.DictService
	dbAnalyzerSvc *service.DBAnalyzerService
	genConfigSvc  *service.GenConfigService

	// Controller 层
	authCtrl      *controller.AuthController
	userCtrl      *controller.UserController
	roleCtrl      *controller.RoleController
	menuCtrl      *controller.MenuController
	tenantCtrl    *controller.TenantController
	fileCtrl      *controller.FileController
	profileCtrl   *controller.ProfileController
	dashboardCtrl *controller.DashboardController
	dictCtrl      *controller.DictController
	generatorCtrl *controller.GeneratorController
	genConfigCtrl *controller.GenConfigController
)

// Init 初始化所有服务
func Init() {
	db := database.GetDB()
	if db == nil {
		panic("database not initialized")
	}

	initRepositories(db)
	initGenerators()
	initServices()
	initControllers()
}

func initRepositories(db *gorm.DB) {
	userRepo = repository.NewUserRepository(db)
	roleRepo = repository.NewRoleRepository(db)
	menuRepo = repository.NewMenuRepository(db)
	tenantRepo = repository.NewTenantRepository(db)
	fileRepo = repository.NewFileRepository(db)
	dictRepo = repository.NewDictRepository(db)
	dbAnalyzerRepo = repository.NewDBAnalyzerRepository(db)
	genConfigRepo = repository.NewGenConfigRepository(db)
}

func initGenerators() {
	templateEngine = generator.NewTemplateEngine()
	if err := templateEngine.LoadTemplates("templates"); err != nil {
		panic("Failed to load templates: " + err.Error())
	}
	codeGenerator = generator.NewCodeGenerator(templateEngine)
	filePackager = generator.NewFilePackager("generated")
}

func initServices() {
	authSvc = service.NewAuthService(userRepo, roleRepo, menuRepo)
	userSvc = service.NewUserService(userRepo, roleRepo)
	roleSvc = service.NewRoleService(roleRepo, userRepo)
	menuSvc = service.NewMenuService(menuRepo, roleRepo)
	tenantSvc = service.NewTenantService(tenantRepo, userRepo)
	profileSvc = service.NewProfileService(userRepo, roleRepo, tenantRepo)
	fileSvc = service.NewFileService(fileRepo, tenantSvc)
	dashboardSvc = service.NewDashboardService(userRepo, tenantRepo, fileRepo)
	dictSvc = service.NewDictService(dictRepo)
	dbAnalyzerSvc = service.NewDBAnalyzerService(dbAnalyzerRepo, database.GetDB())
	genConfigSvc = service.NewGenConfigService(genConfigRepo, dbAnalyzerSvc)

}

func initControllers() {
	authCtrl = controller.NewAuthController(authSvc, roleSvc, menuSvc)
	userCtrl = controller.NewUserController(userSvc)
	roleCtrl = controller.NewRoleController(roleSvc)
	menuCtrl = controller.NewMenuController(menuSvc)
	tenantCtrl = controller.NewTenantController(tenantSvc)
	fileCtrl = controller.NewFileController(fileSvc)
	profileCtrl = controller.NewProfileController(profileSvc)
	dashboardCtrl = controller.NewDashboardController(dashboardSvc)
	dictCtrl = controller.NewDictController(dictSvc)
	generatorCtrl = controller.NewGeneratorController(dbAnalyzerSvc, genConfigSvc, codeGenerator, filePackager)
	genConfigCtrl = controller.NewGenConfigController(genConfigSvc)
}

// === Service 获取函数 ===
func AuthSvc() service.AuthService           { return authSvc }
func UserSvc() service.UserService           { return userSvc }
func RoleSvc() service.RoleService           { return roleSvc }
func MenuSvc() service.MenuService           { return menuSvc }
func TenantSvc() service.TenantService       { return tenantSvc }
func FileSvc() *service.FileService          { return fileSvc }
func ProfileSvc() service.ProfileService     { return profileSvc }
func DashboardSvc() service.DashboardService { return dashboardSvc }
func DictSvc() service.DictService           { return dictSvc }

// Generator 获取函数
func TemplateEngine() *generator.TemplateEngine { return templateEngine }
func CodeGenerator() *generator.CodeGenerator   { return codeGenerator }
func FilePackager() *generator.FilePackager     { return filePackager }

// === Controller 获取函数 ===
func AuthCtrl() *controller.AuthController           { return authCtrl }
func UserCtrl() *controller.UserController           { return userCtrl }
func RoleCtrl() *controller.RoleController           { return roleCtrl }
func MenuCtrl() *controller.MenuController           { return menuCtrl }
func TenantCtrl() *controller.TenantController       { return tenantCtrl }
func FileCtrl() *controller.FileController           { return fileCtrl }
func ProfileCtrl() *controller.ProfileController     { return profileCtrl }
func DashboardCtrl() *controller.DashboardController { return dashboardCtrl }
func DictCtrl() *controller.DictController           { return dictCtrl }
func GeneratorCtrl() *controller.GeneratorController { return generatorCtrl }
func GenConfigCtrl() *controller.GenConfigController { return genConfigCtrl }

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
