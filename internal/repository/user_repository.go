package repository

import (
	"errors"
	"github.com/LiteMove/light-stack/internal/model"
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
		Where("tenant_id = ? AND email = ?", tenantID, email).
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
