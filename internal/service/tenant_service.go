package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/repository"
)

// TenantService 租户服务接口
type TenantService interface {
	// 基础CRUD操作
	CreateTenant(tenant *model.Tenant) error
	GetTenant(id uint64) (*model.Tenant, error)
	GetTenantByDomain(domain string) (*model.Tenant, error)
	UpdateTenant(tenant *model.Tenant) error
	DeleteTenant(id uint64) error

	// 查询操作
	GetTenantList(page, pageSize int, keyword string, status int) ([]*model.Tenant, int64, error)
	CheckDomainExists(domain string) (bool, error)
	CheckNameExists(name string) (bool, error)

	// 状态操作
	UpdateTenantStatus(id uint64, status int) error
	CheckTenantExpired(tenant *model.Tenant) bool

	// 租户验证
	ValidateTenant(domain string) (*model.Tenant, error)
	GetSelectList() ([]*model.Tenant, error)
}

// tenantService 租户服务实现
type tenantService struct {
	tenantRepo repository.TenantRepository
	userRepo   repository.UserRepository
}

func (s *tenantService) GetSelectList() ([]*model.Tenant, error) {
	return s.tenantRepo.GetSelectList()
}

// NewTenantService 创建租户服务
func NewTenantService(tenantRepo repository.TenantRepository, userRepo repository.UserRepository) TenantService {
	return &tenantService{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
	}
}

// CreateTenant 创建租户
func (s *tenantService) CreateTenant(tenant *model.Tenant) error {
	// 检查租户名称是否已存在
	if tenant.Name != "" {
		exists, err := s.tenantRepo.NameExists(tenant.Name)
		if err != nil {
			return fmt.Errorf("检查租户名称是否存在失败: %w", err)
		}
		if exists {
			return errors.New("租户名称已存在")
		}
	}

	// 检查域名是否已存在（如果提供了域名）
	if tenant.Domain != "" {
		exists, err := s.tenantRepo.DomainExists(tenant.Domain)
		if err != nil {
			return fmt.Errorf("检查域名是否存在失败: %w", err)
		}
		if exists {
			return errors.New("域名已存在")
		}
	}

	// 设置默认值
	if tenant.Status == 0 {
		tenant.Status = 1 // 默认启用
	}

	// 创建租户
	if err := s.tenantRepo.Create(tenant); err != nil {
		return fmt.Errorf("创建租户失败: %w", err)
	}

	return nil
}

// GetTenant 获取租户
func (s *tenantService) GetTenant(id uint64) (*model.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取租户失败: %w", err)
	}
	return tenant, nil
}

// GetTenantByDomain 根据域名获取租户
func (s *tenantService) GetTenantByDomain(domain string) (*model.Tenant, error) {
	tenant, err := s.tenantRepo.GetByDomain(domain)
	if err != nil {
		return nil, fmt.Errorf("根据域名获取租户失败: %w", err)
	}
	return tenant, nil
}

// UpdateTenant 更新租户
func (s *tenantService) UpdateTenant(tenant *model.Tenant) error {
	// 获取原租户信息
	existingTenant, err := s.tenantRepo.GetByID(tenant.ID)
	if err != nil {
		return fmt.Errorf("获取原租户信息失败: %w", err)
	}

	// 如果租户名称发生变化，检查新名称是否已存在
	if tenant.Name != existingTenant.Name {
		exists, err := s.tenantRepo.NameExists(tenant.Name)
		if err != nil {
			return fmt.Errorf("检查租户名称是否存在失败: %w", err)
		}
		if exists {
			return errors.New("租户名称已存在")
		}
	}

	// 如果域名发生变化，检查新域名是否已存在
	if tenant.Domain != "" && tenant.Domain != existingTenant.Domain {
		exists, err := s.tenantRepo.DomainExists(tenant.Domain)
		if err != nil {
			return fmt.Errorf("检查域名是否存在失败: %w", err)
		}
		if exists {
			return errors.New("域名已存在")
		}
	}

	// 更新租户
	if err := s.tenantRepo.Update(tenant); err != nil {
		return fmt.Errorf("更新租户失败: %w", err)
	}

	return nil
}

// DeleteTenant 删除租户
func (s *tenantService) DeleteTenant(id uint64) error {
	// 检查租户是否存在
	tenant, err := s.tenantRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("获取租户信息失败: %w", err)
	}

	// 不允许删除系统租户（ID为0）
	if tenant.ID == 0 {
		return errors.New("不允许删除系统租户")
	}

	// 检查租户下是否还有用户
	hasUsers, err := s.tenantRepo.HasUsers(id)
	if err != nil {
		return fmt.Errorf("检查租户用户失败: %w", err)
	}
	if hasUsers {
		return errors.New("租户下还有用户，无法删除，请先删除所有用户")
	}

	// 删除租户
	if err := s.tenantRepo.Delete(id); err != nil {
		return fmt.Errorf("删除租户失败: %w", err)
	}

	return nil
}

// GetTenantList 获取租户列表
func (s *tenantService) GetTenantList(page, pageSize int, keyword string, status int) ([]*model.Tenant, int64, error) {
	tenants, total, err := s.tenantRepo.GetList(page, pageSize, keyword, status)
	if err != nil {
		return nil, 0, fmt.Errorf("获取租户列表失败: %w", err)
	}

	return tenants, total, nil
}

// CheckDomainExists 检查域名是否存在
func (s *tenantService) CheckDomainExists(domain string) (bool, error) {
	exists, err := s.tenantRepo.DomainExists(domain)
	if err != nil {
		return false, fmt.Errorf("检查域名是否存在失败: %w", err)
	}
	return exists, nil
}

// CheckNameExists 检查租户名称是否存在
func (s *tenantService) CheckNameExists(name string) (bool, error) {
	exists, err := s.tenantRepo.NameExists(name)
	if err != nil {
		return false, fmt.Errorf("检查租户名称是否存在失败: %w", err)
	}
	return exists, nil
}

// UpdateTenantStatus 更新租户状态
func (s *tenantService) UpdateTenantStatus(id uint64, status int) error {
	// 检查租户是否存在
	tenant, err := s.tenantRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("获取租户信息失败: %w", err)
	}

	// 不允许禁用系统租户
	if tenant.ID == 0 && status != 1 {
		return errors.New("不允许修改系统租户状态")
	}

	// 更新状态
	if err := s.tenantRepo.UpdateStatus(id, status); err != nil {
		return fmt.Errorf("更新租户状态失败: %w", err)
	}

	return nil
}

// CheckTenantExpired 检查租户是否已过期
func (s *tenantService) CheckTenantExpired(tenant *model.Tenant) bool {
	if tenant.ExpiredAt == nil {
		return false
	}
	return tenant.ExpiredAt.Before(time.Now())
}

// ValidateTenant 验证租户
func (s *tenantService) ValidateTenant(domain string) (*model.Tenant, error) {
	// 根据域名获取租户
	tenant, err := s.tenantRepo.GetByDomain(domain)
	if err != nil {
		return nil, errors.New("租户不存在")
	}

	// 检查租户状态
	if !tenant.IsActive() {
		return nil, errors.New("租户已被禁用")
	}

	// 检查租户是否过期
	if tenant.IsExpired() {
		return nil, errors.New("租户已过期")
	}

	return tenant, nil
}
