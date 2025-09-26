package main

import (
	"log"

	"github.com/LiteMove/light-stack/internal/config"
	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/utils"
	"github.com/LiteMove/light-stack/pkg/database"
	"github.com/LiteMove/light-stack/pkg/logger"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatal("Failed to initialize config:", err)
	}

	// 初始化日志
	logger.Init()

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	db := database.GetDB()

	// 自动迁移数据库表
	logger.Info("Starting database migration...")

	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.UserRole{},
		&model.RoleMenus{},
		&model.File{},
		&model.GenTableConfig{},
		&model.GenTableColumn{},
		&model.GenHistory{},
	)

	if err != nil {
		logger.Fatal("Failed to migrate database:", err)
	}

	logger.Info("Database migration completed successfully!")

	// 创建基础数据
	createBasicData()

	// 创建文件管理菜单
	createFileManagementMenus()

	logger.Info("Migration process finished!")
}

// createBasicData 创建基础数据
func createBasicData() {
	// 创建基础角色
	createBasicRoles()

	// 创建默认管理员用户
	createDefaultAdmin()

	// 分配管理员角色
	assignAdminRoles()
}

// createBasicRoles 创建基础角色
func createBasicRoles() {
	db := database.GetDB()

	roles := []model.Role{
		{
			Name:        "超级管理员",
			Code:        "super_admin",
			Description: "拥有系统所有权限",
			Status:      1,
			IsSystem:    true,
			SortOrder:   1,
		},
		{
			Name:        "租户管理员",
			Code:        "tenant_admin",
			Description: "租户管理员，管理本租户下的用户和角色",
			Status:      1,
			IsSystem:    true,
			SortOrder:   2,
		},
		{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户，只能查看和操作自己的信息",
			Status:      1,
			IsSystem:    true,
			SortOrder:   3,
		},
	}

	for _, role := range roles {
		var existingRole model.Role
		result := db.Where("code = ?", role.Code).First(&existingRole)
		if result.Error != nil {
			// 角色不存在，创建新角色
			if err := db.Create(&role).Error; err != nil {
				logger.Error("Failed to create role:", role.Code, err)
			} else {
				logger.Info("Created role:", role.Code)
			}
		} else {
			logger.Info("Role already exists:", role.Code)
		}
	}
}

// createDefaultAdmin 创建默认管理员用户
func createDefaultAdmin() {
	db := database.GetDB()

	// 检查是否已存在管理员用户
	var count int64
	db.Model(&model.User{}).Where("username = ? AND tenant_id = ?", "admin", 0).Count(&count)

	if count == 0 {
		// 加密默认密码
		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			logger.Error("Failed to hash default admin password:", err)
			return
		}

		// 创建默认管理员
		email := "admin@lightstack.com"
		admin := model.User{
			Username: "admin",
			Password: hashedPassword,
			Nickname: "系统管理员",
			Email:    &email,
			Status:   1,
			IsSystem: true,
		}

		if err := db.Create(&admin).Error; err != nil {
			logger.Error("Failed to create default admin user:", err)
		} else {
			logger.Info("Default admin user created successfully")
			logger.Info("Username: admin, Password: admin123")
			logger.Warn("Please change the default password after first login!")
		}
	} else {
		logger.Info("Admin user already exists, skipping creation")
	}
}

// assignAdminRoles 分配管理员角色
func assignAdminRoles() {
	db := database.GetDB()

	// 获取管理员用户
	var admin model.User
	if err := db.Where("username = ? AND tenant_id = ?", "admin", 0).First(&admin).Error; err != nil {
		logger.Error("Failed to find admin user:", err)
		return
	}

	// 获取超级管理员角色
	var superAdminRole model.Role
	if err := db.Where("code = ? AND tenant_id = ?", "super_admin", 0).First(&superAdminRole).Error; err != nil {
		logger.Error("Failed to find super_admin role:", err)
		return
	}

	// 检查是否已分配角色
	var userRoleCount int64
	db.Model(&model.UserRole{}).Where("user_id = ? AND role_id = ?", admin.ID, superAdminRole.ID).Count(&userRoleCount)

	if userRoleCount == 0 {
		// 为管理员分配超级管理员角色
		userRole := model.UserRole{
			UserID: admin.ID,
			RoleID: superAdminRole.ID,
		}

		if err := db.Create(&userRole).Error; err != nil {
			logger.Error("Failed to assign super_admin role to admin:", err)
		} else {
			logger.Info("Assigned super_admin role to admin user")
		}
	} else {
		logger.Info("Admin user already has super_admin role")
	}
}

// createFileManagementMenus 创建文件管理菜单
func createFileManagementMenus() {
	db := database.GetDB()

	menus := []model.Menu{
		{
			ParentID:  0,
			Name:      "文件管理",
			Code:      "file_management",
			Type:      "menu",
			Path:      "/files",
			Component: "system/files",
			Icon:      "FolderOpened",
			SortOrder: 500,
			IsHidden:  false,
			Status:    1,
		},
	}

	for _, menu := range menus {
		var existingMenu model.Menu
		result := db.Where("code = ?", menu.Code).First(&existingMenu)
		if result.Error != nil {
			// 菜单不存在，创建新菜单
			if err := db.Create(&menu).Error; err != nil {
				logger.Error("Failed to create menu:", menu.Code, err)
			} else {
				logger.Info("Created menu:", menu.Code)

				// 为超级管理员和租户管理员角色分配文件管理权限
				assignFileMenuToRoles(menu.ID)
			}
		} else {
			logger.Info("Menu already exists:", menu.Code)
		}
	}
}

// assignFileMenuToRoles 为角色分配文件管理菜单权限
func assignFileMenuToRoles(menuID uint64) {
	db := database.GetDB()

	// 获取需要分配权限的角色
	roleCodes := []string{"super_admin", "tenant_admin"}

	for _, roleCode := range roleCodes {
		var role model.Role
		if err := db.Where("code = ?", roleCode).First(&role).Error; err != nil {
			logger.Error("Failed to find role:", roleCode, err)
			continue
		}

		// 检查是否已分配权限
		var count int64
		db.Model(&model.RoleMenus{}).Where("role_id = ? AND menu_id = ?", role.ID, menuID).Count(&count)

		if count == 0 {
			// 为角色分配菜单权限
			roleMenu := model.RoleMenus{
				RoleId: role.ID,
				MenuId: menuID,
			}

			if err := db.Create(&roleMenu).Error; err != nil {
				logger.Error("Failed to assign menu to role:", roleCode, err)
			} else {
				logger.Info("Assigned file management menu to role:", roleCode)
			}
		}
	}
}
