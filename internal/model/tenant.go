package model

import (
	"gorm.io/gorm"
	"time"
)

// Tenant 租户模型
type Tenant struct {
	BaseModel
	Name      string     `json:"name" gorm:"not null;size:100" validate:"required,min=1,max=100"`
	Domain    string     `json:"domain" gorm:"size:100;uniqueIndex:uk_domain" validate:"max=100"`
	Status    int        `json:"status" gorm:"not null;default:1;index" validate:"required,oneof=1 2 3 4"`
	ExpiredAt *time.Time `json:"expired_at" gorm:"index"`
	Config    string     `json:"config" gorm:"type:json"`

	// 关联关系
	Users []User `json:"users,omitempty" gorm:"foreignKey:TenantID"`
}

// TableName 指定表名
func (Tenant) TableName() string {
	return "tenants"
}

// TenantProfile 租户资料（简化版本）
type TenantProfile struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Domain    string     `json:"domain"`
	Status    int        `json:"status"`
	ExpiredAt *time.Time `json:"expired_at"`
	Config    string     `json:"config"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToProfile 转换为租户资料
func (t *Tenant) ToProfile() TenantProfile {
	return TenantProfile{
		ID:        t.ID,
		Name:      t.Name,
		Domain:    t.Domain,
		Status:    t.Status,
		ExpiredAt: t.ExpiredAt,
		Config:    t.Config,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// IsActive 检查租户是否为活跃状态
func (t *Tenant) IsActive() bool {
	return t.Status == 1
}

// IsExpired 检查租户是否已过期
func (t *Tenant) IsExpired() bool {
	if t.ExpiredAt == nil {
		return false
	}
	return t.ExpiredAt.Before(time.Now())
}

// BeforeCreate 创建前的钩子
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if t.Status == 0 {
		t.Status = 1 // 默认启用
	}
	return nil
}
