package repository

import (
	"errors"

	"github.com/LiteMove/light-stack/internal/model"

	"gorm.io/gorm"
)

// RoleRepository 角色数据访问接口
type RoleRepository interface {
	// 创建角色
	Create(role *model.Role) error
	// 根据ID获取角色
	GetByID(id uint64) (*model.Role, error)
	// 根据编码获取角色
	GetByCode(code string) (*model.Role, error)
	// 更新角色
	Update(role *model.Role) error
	// 删除角色
	Delete(id uint64) error
	// 检查角色编码是否存在
	CodeExists(code string) (bool, error)
	// 获取角色列表（分页）
	GetList(page, pageSize int, status int) ([]*model.Role, int64, error)
	// 获取用户的角色列表
	GetUserRoles(userID uint64) ([]*model.Role, error)
	// 为用户分配角色
	AssignRolesToUser(userID uint64, roleIDs []uint64) error
	// 移除用户角色
	RemoveUserRoles(userID uint64, roleIDs []uint64) error
	// 更新用户角色（先清空再分配）
	UpdateUserRoles(userID uint64, roleIDs []uint64) error
	// 获取角色的用户数量
	GetRoleUserCount(roleID uint64) (int64, error)
	// 获取角色及其用户信息
	GetRoleWithUsers(roleID uint64) (*model.RoleWithUsers, error)
}

// roleRepository 角色数据访问实现
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色数据访问实例
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

// Create 创建角色
func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// GetByID 根据ID获取角色
func (r *roleRepository) GetByID(id uint64) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("Users").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据编码获取角色
func (r *roleRepository) GetByCode(code string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("code = ?", code).Preload("Users").First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

// Update 更新角色
func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色（软删除）
func (r *roleRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Role{}, id).Error
}

// CodeExists 检查角色编码是否存在
func (r *roleRepository) CodeExists(code string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Role{}).
		Where("code = ?", code).
		Count(&count).Error
	return count > 0, err
}

// GetList 获取角色列表（分页）
func (r *roleRepository) GetList(page, pageSize int, status int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := r.db.Model(&model.Role{})

	// 状态筛选
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Order("sort_order ASC, created_at DESC").
		Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// GetUserRoles 获取用户的角色列表
func (r *roleRepository) GetUserRoles(userID uint64) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.status = 1", userID).
		Order("roles.sort_order ASC").
		Find(&roles).Error
	return roles, err
}

// AssignRolesToUser 为用户分配角色
func (r *roleRepository) AssignRolesToUser(userID uint64, roleIDs []uint64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	// 创建用户角色关联记录
	var userRoles []model.UserRole
	for _, roleID := range roleIDs {
		userRoles = append(userRoles, model.UserRole{
			UserID: userID,
			RoleID: roleID,
		})
	}

	return r.db.Create(&userRoles).Error
}

// RemoveUserRoles 移除用户角色
func (r *roleRepository) RemoveUserRoles(userID uint64, roleIDs []uint64) error {
	if len(roleIDs) == 0 {
		return nil
	}

	return r.db.Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&model.UserRole{}).Error
}

// UpdateUserRoles 更新用户角色（先清空再分配）
func (r *roleRepository) UpdateUserRoles(userID uint64, roleIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除现有角色关联
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 分配新角色
		if len(roleIDs) > 0 {
			var userRoles []model.UserRole
			for _, roleID := range roleIDs {
				userRoles = append(userRoles, model.UserRole{
					UserID: userID,
					RoleID: roleID,
				})
			}
			return tx.Create(&userRoles).Error
		}

		return nil
	})
}

// GetRoleUserCount 获取角色的用户数量
func (r *roleRepository) GetRoleUserCount(roleID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.UserRole{}).Where("role_id = ?", roleID).Count(&count).Error
	return count, err
}

// GetRoleWithUsers 获取角色及其用户信息
func (r *roleRepository) GetRoleWithUsers(roleID uint64) (*model.RoleWithUsers, error) {
	var role model.Role
	err := r.db.Preload("Users").First(&role, roleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	roleWithUsers := &model.RoleWithUsers{
		Role:      role,
		UserCount: len(role.Users),
		Users:     role.Users,
	}

	return roleWithUsers, nil
}
