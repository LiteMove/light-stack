package permission

import (
	repository2 "github.com/LiteMove/light-stack/internal/modules/system/repository"
	"github.com/LiteMove/light-stack/pkg/logger"
)

func LoadUserPermissions(userID uint64, menuRepo repository2.MenuRepository) error {
	permissions, err := menuRepo.GetUserPermissions(userID)
	if err != nil {
		logger.WithField("userId", userID).Error("Failed to load user permissions:", err)
		return err
	}

	Cache.LoadUserPermissions(userID, permissions)
	logger.WithField("userId", userID).Debug("User permissions loaded successfully")
	return nil
}

func LoadUserRoles(userID uint64, roleRepo repository2.RoleRepository) error {
	roles, err := roleRepo.GetUserRoles(userID)
	if err != nil {
		logger.WithField("userId", userID).Error("Failed to load user roles:", err)
		return err
	}

	roleCodes := make([]string, len(roles))
	for i, role := range roles {
		roleCodes[i] = role.Code
	}

	Cache.LoadUserRoles(userID, roleCodes)
	logger.WithField("userId", userID).Debug("User roles loaded successfully")
	return nil
}

func LoadUserData(userID uint64, menuRepo repository2.MenuRepository, roleRepo repository2.RoleRepository) error {
	// 加载权限
	if err := LoadUserPermissions(userID, menuRepo); err != nil {
		return err
	}

	// 加载角色
	if err := LoadUserRoles(userID, roleRepo); err != nil {
		return err
	}

	return nil
}

func ClearUserPermissions(userID uint64) {
	Cache.ClearUserPermissions(userID)
	logger.WithField("userId", userID).Debug("User permissions cleared")
}
