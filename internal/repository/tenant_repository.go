package repository

import (
	"errors"

	"github.com/LiteMove/light-stack/internal/model"

	"gorm.io/gorm"
)

// TenantRepository 租户数据访问接口
type TenantRepository interface {
	// 创建租户
	Create(tenant *model.Tenant) error
	// 根据ID获取租户
	GetByID(id uint64) (*model.Tenant, error)
	// 根据域名获取租户
	GetByDomain(domain string) (*model.Tenant, error)
	// 更新租户
	Update(tenant *model.Tenant) error
	// 删除租户
	Delete(id uint64) error
	// 检查域名是否存在
	DomainExists(domain string) (bool, error)
	// 获取租户列表（分页）
	GetList(page, pageSize int, keyword string, status int) ([]*model.Tenant, int64, error)
	// 更新租户状态
	UpdateStatus(id uint64, status int) error
	// 检查租户名称是否存在
	NameExists(name string) (bool, error)
	// 检查租户是否有用户
	HasUsers(id uint64) (bool, error)
	// 获取启用的租户列表
	GetSelectList() ([]*model.Tenant, error)
}

// tenantRepository 租户数据访问实现
type tenantRepository struct {
	db *gorm.DB
}

// NewTenantRepository 创建租户数据访问实例
func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{
		db: db,
	}
}

// Create 创建租户
func (r *tenantRepository) Create(tenant *model.Tenant) error {
	return r.db.Create(tenant).Error
}

// GetByID 根据ID获取租户
func (r *tenantRepository) GetByID(id uint64) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.First(&tenant, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}
	return &tenant, nil
}

// GetByDomain 根据域名获取租户
func (r *tenantRepository) GetByDomain(domain string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.Where("domain = ?", domain).First(&tenant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant not found")
		}
		return nil, err
	}
	return &tenant, nil
}

// Update 更新租户
func (r *tenantRepository) Update(tenant *model.Tenant) error {
	return r.db.Save(tenant).Error
}

// Delete 删除租户（软删除）
func (r *tenantRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Tenant{}, id).Error
}

// DomainExists 检查域名是否存在
func (r *tenantRepository) DomainExists(domain string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Tenant{}).
		Where("domain = ?", domain).
		Count(&count).Error
	return count > 0, err
}

// GetList 获取租户列表（分页）
func (r *tenantRepository) GetList(page, pageSize int, keyword string, status int) ([]*model.Tenant, int64, error) {
	var tenants []*model.Tenant
	var total int64

	query := r.db.Model(&model.Tenant{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR domain LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

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
		Order("created_at DESC").
		Find(&tenants).Error
	if err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

// UpdateStatus 更新租户状态
func (r *tenantRepository) UpdateStatus(id uint64, status int) error {
	return r.db.Model(&model.Tenant{}).Where("id = ?", id).Update("status", status).Error
}

// NameExists 检查租户名称是否存在
func (r *tenantRepository) NameExists(name string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Tenant{}).
		Where("name = ?", name).
		Count(&count).Error
	return count > 0, err
}

// HasUsers 检查租户是否有用户
func (r *tenantRepository) HasUsers(id uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).
		Where("tenant_id = ?", id).
		Count(&count).Error
	return count > 0, err
}

func (r *tenantRepository) GetSelectList() ([]*model.Tenant, error) {
	var tenants []*model.Tenant
	err := r.db.Find(&tenants).Error
	if err != nil {
		return nil, err
	}
	return tenants, nil
}
