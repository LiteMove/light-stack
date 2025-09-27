package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/utils"
)

// UserService 用户服务接口
type UserService interface {
	// 基础CRUD操作
	CreateUser(user *model.User) error
	GetUser(id uint64) (*model.User, error)
	GetUserWithRoles(id uint64) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint64) error

	// 查询操作
	GetUserList(tenantID uint64, page, pageSize int, keyword string, status int, roleID uint64) ([]*model.User, int64, error)
	GetUserByUsername(tenantID uint64, username string) (*model.User, error)
	GetUserByEmail(tenantID uint64, email string) (*model.User, error)

	// 状态操作
	UpdateUserStatus(id uint64, status int) error
	BatchUpdateUserStatus(ids []uint64, status int) error

	// 密码相关
	ChangePassword(id uint64, oldPassword, newPassword string) error
	ResetPassword(id uint64) (string, error)

	// 用户验证
	ValidateUser(tenantID uint64, username, password string) (*model.User, error)
	CheckUsernameExists(tenantID uint64, username string) (bool, error)
	CheckEmailExists(tenantID uint64, email string) (bool, error)

	// 角色管理
	AssignUserRoles(userID uint64, roleIDs []uint64) error
	RemoveUserRoles(userID uint64, roleIDs []uint64) error
	GetUserRoles(userID uint64) ([]*model.Role, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) UserService {
	return &userService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// CreateUser 创建用户
func (s *userService) CreateUser(user *model.User) error {
	// 检查用户名是否已存在
	exists, err := s.userRepo.UsernameExists(user.TenantID, user.Username)
	if err != nil {
		return fmt.Errorf("检查用户名是否存在失败: %w", err)
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在（如果提供了邮箱）
	if user.Email != nil && *user.Email != "" {
		exists, err := s.userRepo.EmailExists(user.TenantID, *user.Email)
		if err != nil {
			return fmt.Errorf("检查邮箱是否存在失败: %w", err)
		}
		if exists {
			return errors.New("邮箱已存在")
		}
	}

	// 检查手机号是否已存在（如果提供了手机号）
	if user.Phone != nil && *user.Phone != "" {
		exists, err := s.userRepo.PhoneExists(user.TenantID, *user.Phone)
		if err != nil {
			return fmt.Errorf("检查手机号是否存在失败: %w", err)
		}
		if exists {
			return errors.New("手机号已存在")
		}
	}

	// 如果没有设置密码，生成默认密码
	if user.Password == "" {
		user.Password = "123456" // 默认密码
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}
	user.Password = hashedPassword

	// 设置默认值
	if user.Status == 0 {
		user.Status = 1 // 默认启用
	}

	// 创建用户
	if err := s.userRepo.Create(user); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// GetUser 获取用户
func (s *userService) GetUser(id uint64) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return user, nil
}

// GetUserWithRoles 获取用户（包含角色）
func (s *userService) GetUserWithRoles(id uint64) (*model.User, error) {
	user, err := s.userRepo.GetByIDWithRoles(id)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}
	return user, nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(user *model.User) error {
	// 获取原用户信息
	existingUser, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return fmt.Errorf("获取原用户信息失败: %w", err)
	}

	// 如果用户名发生变化，检查新用户名是否已存在
	if user.Username != existingUser.Username {
		exists, err := s.userRepo.UsernameExists(user.TenantID, user.Username)
		if err != nil {
			return fmt.Errorf("检查用户名是否存在失败: %w", err)
		}
		if exists {
			return errors.New("用户名已存在")
		}
	}

	// 如果邮箱发生变化，检查新邮箱是否已存在
	if user.Email != nil && *user.Email != "" {
		// 检查邮箱是否真的发生了变化
		existingEmailValue := ""
		if existingUser.Email != nil {
			existingEmailValue = *existingUser.Email
		}
		if *user.Email != existingEmailValue {
			exists, err := s.userRepo.EmailExists(user.TenantID, *user.Email)
			if err != nil {
				return fmt.Errorf("检查邮箱是否存在失败: %w", err)
			}
			if exists {
				return errors.New("邮箱已存在")
			}
		}
	}

	// 如果手机号发生变化，检查新手机号是否已存在
	if user.Phone != nil && *user.Phone != "" {
		// 检查手机号是否真的发生了变化
		existingPhoneValue := ""
		if existingUser.Phone != nil {
			existingPhoneValue = *existingUser.Phone
		}
		if *user.Phone != existingPhoneValue {
			exists, err := s.userRepo.PhoneExists(user.TenantID, *user.Phone)
			if err != nil {
				return fmt.Errorf("检查手机号是否存在失败: %w", err)
			}
			if exists {
				return errors.New("手机号已存在")
			}
		}
	}

	// 不允许修改密码（使用专门的修改密码方法）
	user.Password = existingUser.Password

	// 更新用户
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("更新用户失败: %w", err)
	}

	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint64) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 不允许删除系统用户
	if user.IsSystem {
		return errors.New("不允许删除系统用户")
	}

	// 删除用户
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(tenantID uint64, page, pageSize int, keyword string, status int, roleID uint64) ([]*model.User, int64, error) {
	// TODO: 目前repository层的GetList方法不支持关键词和角色筛选
	// 这里先使用基础的分页查询，后续需要扩展repository方法
	users, total, err := s.userRepo.GetList(tenantID, page, pageSize, status)
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 简单的关键词过滤（在内存中进行，生产环境应该在数据库层面实现）
	if keyword != "" {
		keyword = strings.ToLower(keyword)
		var filteredUsers []*model.User
		for _, user := range users {
			emailStr := ""
			if user.Email != nil {
				emailStr = *user.Email
			}
			if strings.Contains(strings.ToLower(user.Username), keyword) ||
				strings.Contains(strings.ToLower(user.Nickname), keyword) ||
				strings.Contains(strings.ToLower(emailStr), keyword) {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
		total = int64(len(users))
	}

	return users, total, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(tenantID uint64, username string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(tenantID, username)
	if err != nil {
		return nil, fmt.Errorf("根据用户名获取用户失败: %w", err)
	}
	return user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *userService) GetUserByEmail(tenantID uint64, email string) (*model.User, error) {
	user, err := s.userRepo.GetByEmail(tenantID, email)
	if err != nil {
		return nil, fmt.Errorf("根据邮箱获取用户失败: %w", err)
	}
	return user, nil
}

// UpdateUserStatus 更新用户状态
func (s *userService) UpdateUserStatus(id uint64, status int) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 不允许禁用系统用户
	if user.IsSystem && status == 2 {
		return errors.New("不允许禁用系统用户")
	}

	// 不允许禁用超级管理员
	if user.ID == model.SuperAdminUserId && status == 2 {
		return errors.New("不允许禁用超级管理员")
	}

	// 更新状态
	if err := s.userRepo.UpdateStatus(id, status); err != nil {
		return fmt.Errorf("更新用户状态失败: %w", err)
	}

	return nil
}

// BatchUpdateUserStatus 批量更新用户状态
func (s *userService) BatchUpdateUserStatus(ids []uint64, status int) error {
	for _, id := range ids {
		// 检查每个用户
		user, err := s.userRepo.GetByID(id)
		if err != nil {
			continue // 跳过不存在的用户
		}

		// 不允许禁用系统用户
		if user.IsSystem && status == 2 {
			continue // 跳过系统用户
		}

		// 不允许禁用超级管理员
		if user.ID == model.SuperAdminUserId && status == 2 {
			continue // 跳过超级管理员
		}

		// 更新状态
		if err := s.userRepo.UpdateStatus(id, status); err != nil {
			// 记录错误但继续处理其他用户
			continue
		}
	}

	return nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(id uint64, oldPassword, newPassword string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 验证原密码
	if !utils.VerifyPassword(user.Password, oldPassword) {
		return errors.New("原密码不正确")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(id, hashedPassword); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(id uint64) (string, error) {
	// 生成新密码（6位随机数字）
	newPassword := "123456" // 简单示例，生产环境应该生成随机密码

	// 加密密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(id, hashedPassword); err != nil {
		return "", fmt.Errorf("重置密码失败: %w", err)
	}

	return newPassword, nil
}

// ValidateUser 验证用户
func (s *userService) ValidateUser(tenantID uint64, username, password string) (*model.User, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsernameWithRoles(tenantID, username)
	if err != nil {
		return nil, errors.New("用户名或密码不正确")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if !utils.VerifyPassword(user.Password, password) {
		return nil, errors.New("用户名或密码不正确")
	}

	return user, nil
}

// CheckUsernameExists 检查用户名是否存在
func (s *userService) CheckUsernameExists(tenantID uint64, username string) (bool, error) {
	exists, err := s.userRepo.UsernameExists(tenantID, username)
	if err != nil {
		return false, fmt.Errorf("检查用户名是否存在失败: %w", err)
	}
	return exists, nil
}

// CheckEmailExists 检查邮箱是否存在
func (s *userService) CheckEmailExists(tenantID uint64, email string) (bool, error) {
	exists, err := s.userRepo.EmailExists(tenantID, email)
	if err != nil {
		return false, fmt.Errorf("检查邮箱是否存在失败: %w", err)
	}
	return exists, nil
}

// AssignUserRoles 为用户分配角色
func (s *userService) AssignUserRoles(userID uint64, roleIDs []uint64) error {
	// 获取用户信息
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 验证角色是否存在
	for _, roleID := range roleIDs {
		_, err := s.roleRepo.GetByID(roleID)
		if err != nil {
			return fmt.Errorf("角色ID %d 不存在", roleID)
		}
	}

	// 分配角色
	if err := s.userRepo.BatchAssignRoles(userID, roleIDs); err != nil {
		return fmt.Errorf("分配角色失败: %w", err)
	}

	return nil
}

// RemoveUserRoles 移除用户角色
func (s *userService) RemoveUserRoles(userID uint64, roleIDs []uint64) error {
	// 获取用户信息
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 验证角色是否存在
	for _, roleID := range roleIDs {
		_, err := s.roleRepo.GetByID(roleID)
		if err != nil {
			return fmt.Errorf("角色ID %d 不存在", roleID)
		}
	}

	// 移除角色
	if err := s.userRepo.BatchRemoveRoles(userID, roleIDs); err != nil {
		return fmt.Errorf("移除角色失败: %w", err)
	}

	return nil
}

// GetUserRoles 获取用户角色
func (s *userService) GetUserRoles(userID uint64) ([]*model.Role, error) {
	// 获取带角色信息的用户
	user, err := s.userRepo.GetByIDWithRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 转换类型
	var roles []*model.Role
	for i := range user.Roles {
		roles = append(roles, &user.Roles[i])
	}
	return roles, nil
}
