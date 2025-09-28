package globals

import (
	analyticsController "github.com/LiteMove/light-stack/internal/modules/analytics/controller"
	analyticsService "github.com/LiteMove/light-stack/internal/modules/analytics/service"
	authController "github.com/LiteMove/light-stack/internal/modules/auth/controller"
	authService "github.com/LiteMove/light-stack/internal/modules/auth/service"
	fileController "github.com/LiteMove/light-stack/internal/modules/files/controller"
	repository3 "github.com/LiteMove/light-stack/internal/modules/files/repository"
	fileService "github.com/LiteMove/light-stack/internal/modules/files/service"
	generatorController "github.com/LiteMove/light-stack/internal/modules/generator/controller"
	generatorEngine "github.com/LiteMove/light-stack/internal/modules/generator/engine"
	repository4 "github.com/LiteMove/light-stack/internal/modules/generator/repository"
	generatorService "github.com/LiteMove/light-stack/internal/modules/generator/service"
	systemController "github.com/LiteMove/light-stack/internal/modules/system/controller"
	repository2 "github.com/LiteMove/light-stack/internal/modules/system/repository"
	systemService "github.com/LiteMove/light-stack/internal/modules/system/service"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/pkg/database"
	"gorm.io/gorm"
)

// 全局服务实例
var (
	// Repository 层
	userRepo       repository2.UserRepository
	roleRepo       repository2.RoleRepository
	menuRepo       repository2.MenuRepository
	tenantRepo     repository2.TenantRepository
	fileRepo       *repository3.FileRepository
	dictRepo       repository2.DictRepository
	dbAnalyzerRepo *repository.DBAnalyzerRepository
	genConfigRepo  *repository4.GenConfigRepository

	// Generator 层
	templateEngine *generatorEngine.TemplateEngine
	codeGenerator  *generatorEngine.CodeGenerator
	filePackager   *generatorEngine.FilePackager

	// Service 层
	authSvc       authService.AuthService
	userSvc       systemService.UserService
	roleSvc       systemService.RoleService
	menuSvc       systemService.MenuService
	tenantSvc     systemService.TenantService
	fileSvc       *fileService.FileService
	profileSvc    authService.ProfileService
	dashboardSvc  analyticsService.DashboardService
	dictSvc       systemService.DictService
	dbAnalyzerSvc *generatorService.DBAnalyzerService
	genConfigSvc  *generatorService.GenConfigService

	// Controller 层
	authCtrl      *authController.AuthController
	userCtrl      *systemController.UserController
	roleCtrl      *systemController.RoleController
	menuCtrl      *systemController.MenuController
	tenantCtrl    *systemController.TenantController
	fileCtrl      *fileController.FileController
	profileCtrl   *authController.ProfileController
	dashboardCtrl *analyticsController.DashboardController
	dictCtrl      *systemController.DictController
	generatorCtrl *generatorController.GeneratorController
	genConfigCtrl *generatorController.GenConfigController
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
	userRepo = repository2.NewUserRepository(db)
	roleRepo = repository2.NewRoleRepository(db)
	menuRepo = repository2.NewMenuRepository(db)
	tenantRepo = repository2.NewTenantRepository(db)
	fileRepo = repository3.NewFileRepository(db)
	dictRepo = repository2.NewDictRepository(db)
	dbAnalyzerRepo = repository.NewDBAnalyzerRepository(db)
	genConfigRepo = repository4.NewGenConfigRepository(db)
}

func initGenerators() {
	templateEngine = generatorEngine.NewTemplateEngine()
	if err := templateEngine.LoadTemplates("templates"); err != nil {
		panic("Failed to load templates: " + err.Error())
	}
	codeGenerator = generatorEngine.NewCodeGenerator(templateEngine)
	filePackager = generatorEngine.NewFilePackager("generated")
}

func initServices() {
	authSvc = authService.NewAuthService(userRepo, roleRepo, menuRepo)
	userSvc = systemService.NewUserService(userRepo, roleRepo)
	roleSvc = systemService.NewRoleService(roleRepo, userRepo)
	menuSvc = systemService.NewMenuService(menuRepo, roleRepo)
	tenantSvc = systemService.NewTenantService(tenantRepo, userRepo)
	profileSvc = authService.NewProfileService(userRepo, roleRepo, tenantRepo)
	fileSvc = fileService.NewFileService(fileRepo, tenantSvc)
	dashboardSvc = analyticsService.NewDashboardService(userRepo, tenantRepo, fileRepo)
	dictSvc = systemService.NewDictService(dictRepo)
	dbAnalyzerSvc = generatorService.NewDBAnalyzerService(dbAnalyzerRepo, database.GetDB())
	genConfigSvc = generatorService.NewGenConfigService(genConfigRepo, dbAnalyzerSvc)

}

func initControllers() {
	authCtrl = authController.NewAuthController(authSvc)
	userCtrl = systemController.NewUserController(userSvc)
	roleCtrl = systemController.NewRoleController(roleSvc)
	menuCtrl = systemController.NewMenuController(menuSvc)
	tenantCtrl = systemController.NewTenantController(tenantSvc)
	fileCtrl = fileController.NewFileController(fileSvc)
	profileCtrl = authController.NewProfileController(profileSvc)
	dashboardCtrl = analyticsController.NewDashboardController(dashboardSvc)
	dictCtrl = systemController.NewDictController(dictSvc)
	generatorCtrl = generatorController.NewGeneratorController(dbAnalyzerSvc, genConfigSvc, codeGenerator, filePackager, menuSvc)
	genConfigCtrl = generatorController.NewGenConfigController(genConfigSvc)
}

// === Service 获取函数 ===
func AuthSvc() authService.AuthService                { return authSvc }
func UserSvc() systemService.UserService              { return userSvc }
func RoleSvc() systemService.RoleService              { return roleSvc }
func MenuSvc() systemService.MenuService              { return menuSvc }
func TenantSvc() systemService.TenantService          { return tenantSvc }
func FileSvc() *fileService.FileService               { return fileSvc }
func ProfileSvc() authService.ProfileService          { return profileSvc }
func DashboardSvc() analyticsService.DashboardService { return dashboardSvc }
func DictSvc() systemService.DictService              { return dictSvc }

// Generator 获取函数
func TemplateEngine() *generatorEngine.TemplateEngine { return templateEngine }
func CodeGenerator() *generatorEngine.CodeGenerator   { return codeGenerator }
func FilePackager() *generatorEngine.FilePackager     { return filePackager }

// === Controller 获取函数 ===
func AuthCtrl() *authController.AuthController                { return authCtrl }
func UserCtrl() *systemController.UserController              { return userCtrl }
func RoleCtrl() *systemController.RoleController              { return roleCtrl }
func MenuCtrl() *systemController.MenuController              { return menuCtrl }
func TenantCtrl() *systemController.TenantController          { return tenantCtrl }
func FileCtrl() *fileController.FileController                { return fileCtrl }
func ProfileCtrl() *authController.ProfileController          { return profileCtrl }
func DashboardCtrl() *analyticsController.DashboardController { return dashboardCtrl }
func DictCtrl() *systemController.DictController              { return dictCtrl }
func GeneratorCtrl() *generatorController.GeneratorController { return generatorCtrl }
func GenConfigCtrl() *generatorController.GenConfigController { return genConfigCtrl }

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
