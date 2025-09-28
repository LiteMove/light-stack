package service

import (
	"errors"
	repository2 "github.com/LiteMove/light-stack/internal/modules/system/repository"

	"github.com/LiteMove/light-stack/internal/modules/system/model"
	"github.com/LiteMove/light-stack/pkg/logger"
)

// roleService 角色服务实现
type roleService struct {
	roleRepo repository2.RoleRepository
	userRepo repository2.UserRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleRepo repository2.RoleRepository, userRepo repository2.UserRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	TenantID    uint64 `json:"tenant_id"`
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Code        string `json:"code" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"max=255"`
	SortOrder   int    `json:"sortOrder"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
	Status      int    `json:"status" validate:"oneof=1 2"`
	SortOrder   int    `json:"sortOrder"`
}

// RoleService 角色服务
type RoleService interface {
	// 创建角色
	Create(req *CreateRoleRequest) (*model.RoleProfile, error)
	// 更新角色
	Update(id uint64, req *UpdateRoleRequest) (*model.RoleProfile, error)
	// 删除角色
	Delete(id uint64) error
	// 获取角色信息
	GetByID(id uint64) (*model.RoleProfile, error)
	// 获取角色列表
	GetList(page, pageSize int, status int) ([]*model.Role, int64, error)
	// 为用户分配角色
	AssignRolesToUser(userID uint64, roleIDs []uint64) error
	// 移除用户角色
	RemoveUserRoles(userID uint64, roleIDs []uint64) error
	// 获取所有启用的角色
	GetEnabledRoles(isSuper bool) ([]*model.Role, error)
}

// 角色服务实现

// Create 创建角色
func (s *roleService) Create(req *CreateRoleRequest) (*model.RoleProfile, error) {
	// 检查角色编码是否已存在
	exists, err := s.roleRepo.CodeExists(req.Code)
	if err != nil {
		logger.Error("Failed to check role code existence:", err)
		return nil, errors.New("创建失败")
	}
	if exists {
		return nil, errors.New("角色编码已存在")
	}

	// 创建角色
	role := &model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1, // 默认启用
		SortOrder:   req.SortOrder,
	}

	if err := s.roleRepo.Create(role); err != nil {
		logger.Error("Failed to create role:", err)
		return nil, errors.New("创建失败")
	}

	logger.WithField("roleId", role.ID).Info("Role created successfully")
	profile := role.ToProfile()
	return &profile, nil
}

// Update 更新角色
func (s *roleService) Update(id uint64, req *UpdateRoleRequest) (*model.RoleProfile, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	// 更新角色信息
	role.Name = req.Name
	role.Description = req.Description
	role.Status = req.Status
	role.SortOrder = req.SortOrder

	if err := s.roleRepo.Update(role); err != nil {
		logger.WithField("roleId", id).Error("Failed to update role:", err)
		return nil, errors.New("更新失败")
	}

	logger.WithField("roleId", id).Info("Role updated successfully")
	profile := role.ToProfile()
	return &profile, nil
}

// Delete 删除角色
func (s *roleService) Delete(id uint64) error {
	// 检查角色是否还有用户在使用
	count, err := s.roleRepo.GetRoleUserCount(id)
	if err != nil {
		return errors.New("删除失败")
	}
	if count > 0 {
		return errors.New("该角色还有用户在使用，无法删除")
	}

	if err := s.roleRepo.Delete(id); err != nil {
		logger.WithField("roleId", id).Error("Failed to delete role:", err)
		return errors.New("删除失败")
	}

	logger.WithField("roleId", id).Info("Role deleted successfully")
	return nil
}

// GetByID 获取角色信息
func (s *roleService) GetByID(id uint64) (*model.RoleProfile, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	profile := role.ToProfile()
	return &profile, nil
}

// GetList 获取角色列表
func (s *roleService) GetList(page, pageSize int, status int) ([]*model.Role, int64, error) {
	return s.roleRepo.GetList(page, pageSize, status)
}

// AssignRolesToUser 为用户分配角色
func (s *roleService) AssignRolesToUser(userID uint64, roleIDs []uint64) error {
	return s.roleRepo.UpdateUserRoles(userID, roleIDs)
}

// RemoveUserRoles 移除用户角色
func (s *roleService) RemoveUserRoles(userID uint64, roleIDs []uint64) error {
	return s.roleRepo.RemoveUserRoles(userID, roleIDs)
}

// GetEnabledRoles 获取所有启用的角色
func (s *roleService) GetEnabledRoles(isSuper bool) ([]*model.Role, error) {
	return s.roleRepo.GetEnabledRoles(isSuper)
}
