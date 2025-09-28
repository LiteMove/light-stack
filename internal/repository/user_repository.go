package repository

import (
	"errors"
	"github.com/LiteMove/light-stack/internal/modules/system/model"
	"time"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问接口
type UserRepository interface {
	// 创建用户
	Create(user *model.User) error
	// 根据ID获取用户
	GetByID(id uint64) (*model.User, error)
	// 根据ID获取用户（包含角色）
	GetByIDWithRoles(id uint64) (*model.User, error)
	// 根据用户名获取用户
	GetByUsername(tenantID uint64, username string) (*model.User, error)
	// 根据用户名获取用户（包含角色）
	GetByUsernameWithRoles(tenantID uint64, username string) (*model.User, error)
	// 根据邮箱获取用户
	GetByEmail(tenantID uint64, email string) (*model.User, error)
	// 更新用户
	Update(user *model.User) error
	// 更新用户登录信息
	UpdateLoginInfo(id uint64, ip string) error
	// 删除用户
	Delete(id uint64) error
	// 检查用户名是否存在
	UsernameExists(tenantID uint64, username string) (bool, error)
	// 检查邮箱是否存在
	EmailExists(tenantID uint64, email string) (bool, error)
	// 检查手机号是否存在
	PhoneExists(tenantID uint64, phone string) (bool, error)
	// 获取用户列表（分页）
	GetList(tenantID uint64, page, pageSize int, status int) ([]*model.User, int64, error)
	// 更新用户状态
	UpdateStatus(id uint64, status int) error
	// 更新密码
	UpdatePassword(id uint64, hashedPassword string) error
	// 记录登录失败
	RecordLoginFailure(id uint64) error
	// 重置登录失败计数
	ResetLoginFailures(id uint64) error
	// 锁定用户
	LockUser(id uint64, lockUntil time.Time) error
	// 为用户分配角色
	AssignRole(userID uint64, roleIDs []uint64) error
	// 移除用户角色
	RemoveRole(userID uint64, roleID uint64) error
	// 批量分配角色
	BatchAssignRoles(userID uint64, roleIDs []uint64) error
	// 批量移除角色
	BatchRemoveRoles(userID uint64, roleIDs []uint64) error
	// 获取用户总数
	GetTotalCount() (int64, error)
	// 根据租户ID获取用户数量
	GetCountByTenantID(tenantID uint64) (int64, error)
}

// userRepository 用户数据访问实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据访问实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// Create 创建用户
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *userRepository) GetByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByIDWithRoles 根据ID获取用户（包含角色）
func (r *userRepository) GetByIDWithRoles(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(tenantID uint64, username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("tenant_id = ? AND username = ?", tenantID, username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsernameWithRoles 根据用户名获取用户（包含角色）
func (r *userRepository) GetByUsernameWithRoles(tenantID uint64, username string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").
		Where("tenant_id = ? AND username = ?", tenantID, username).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(tenantID uint64, email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("tenant_id = ? AND email = ?", tenantID, email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// UpdateLoginInfo 更新用户登录信息
func (r *userRepository) UpdateLoginInfo(id uint64, ip string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"last_login_at":  &now,
		"last_login_ip":  ip,
		"login_failures": 0,
		"locked_until":   nil,
	}
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(id uint64) error {
	return r.db.Delete(&model.User{}, id).Error
}

// UsernameExists 检查用户名是否存在
func (r *userRepository) UsernameExists(tenantID uint64, username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).
		Where("tenant_id = ? AND username = ?", tenantID, username).
		Count(&count).Error
	return count > 0, err
}

// EmailExists 检查邮箱是否存在
func (r *userRepository) EmailExists(tenantID uint64, email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).
		Where("tenant_id = ? AND email = ? AND email != ''", tenantID, email).
		Count(&count).Error
	return count > 0, err
}

// PhoneExists 检查手机号是否存在
func (r *userRepository) PhoneExists(tenantID uint64, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).
		Where("tenant_id = ? AND phone = ? AND phone != ''", tenantID, phone).
		Count(&count).Error
	return count > 0, err
}

// GetList 获取用户列表（分页）
func (r *userRepository) GetList(tenantID uint64, page, pageSize int, status int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{}).Where("tenant_id = ?", tenantID)

	// 状态筛选
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，包含角色信息
	offset := (page - 1) * pageSize
	err := query.Preload("Roles").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateStatus 更新用户状态
func (r *userRepository) UpdateStatus(id uint64, status int) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// UpdatePassword 更新密码
func (r *userRepository) UpdatePassword(id uint64, hashedPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// RecordLoginFailure 记录登录失败
func (r *userRepository) RecordLoginFailure(id uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		Update("login_failures", gorm.Expr("login_failures + 1")).Error
}

// ResetLoginFailures 重置登录失败计数
func (r *userRepository) ResetLoginFailures(id uint64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("login_failures", 0).Error
}

// LockUser 锁定用户
func (r *userRepository) LockUser(id uint64, lockUntil time.Time) error {
	updates := map[string]interface{}{
		"status":       3, // 锁定状态
		"locked_until": &lockUntil,
	}
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// AssignRole 为用户分配角色（单个角色）
func (r *userRepository) AssignRole(userID uint64, roleIDs []uint64) error {
	return r.BatchAssignRoles(userID, roleIDs)
}

// RemoveRole 移除用户角色（单个角色）
func (r *userRepository) RemoveRole(userID uint64, roleID uint64) error {
	return r.BatchRemoveRoles(userID, []uint64{roleID})
}

// BatchAssignRoles 批量为用户分配角色
func (r *userRepository) BatchAssignRoles(userID uint64, roleIDs []uint64) error {
	// 使用事务处理
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除现有的角色关联（如果需要替换）
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 批量插入新的角色关联
		userRoles := make([]model.UserRole, 0, len(roleIDs))
		for _, roleID := range roleIDs {
			// 检查关联是否已存在
			var count int64
			tx.Model(&model.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
			if count == 0 {
				userRoles = append(userRoles, model.UserRole{
					UserID: userID,
					RoleID: roleID,
				})
			}
		}

		// 批量插入
		if len(userRoles) > 0 {
			if err := tx.Create(&userRoles).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// BatchRemoveRoles 批量移除用户角色
func (r *userRepository) BatchRemoveRoles(userID uint64, roleIDs []uint64) error {
	// 使用事务处理
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除指定的角色关联
		return tx.Where("user_id = ? AND role_id IN ?", userID, roleIDs).Delete(&model.UserRole{}).Error
	})
}

// GetTotalCount 获取用户总数
func (r *userRepository) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}

// GetCountByTenantID 根据租户ID获取用户数量
func (r *userRepository) GetCountByTenantID(tenantID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}
