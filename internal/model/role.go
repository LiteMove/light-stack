package model

import (
	"time"
	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	TenantID    uint64         `json:"tenant_id" gorm:"not null;default:0;index" validate:"min=0"`
	Name        string         `json:"name" gorm:"not null;size:100" validate:"required,min=1,max=100"`
	Code        string         `json:"code" gorm:"not null;size:50;uniqueIndex:uk_tenant_code" validate:"required,min=1,max=50"`
	Description string         `json:"description" gorm:"size:255" validate:"max=255"`
	Status      int            `json:"status" gorm:"not null;default:1" validate:"required,oneof=1 2"`
	IsSystem    bool           `json:"is_system" gorm:"not null;default:false"`
	SortOrder   int            `json:"sort_order" gorm:"not null;default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
	Menus       []Menu       `json:"menus,omitempty" gorm:"many2many:role_menus;"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// RoleProfile 角色资料（简化版本）
type RoleProfile struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	IsSystem    bool      `json:"is_system"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToProfile 转换为角色资料
func (r *Role) ToProfile() RoleProfile {
	return RoleProfile{
		ID:          r.ID,
		Name:        r.Name,
		Code:        r.Code,
		Description: r.Description,
		Status:      r.Status,
		IsSystem:    r.IsSystem,
		SortOrder:   r.SortOrder,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// BeforeCreate 创建前的钩子
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.Status == 0 {
		r.Status = 1 // 默认启用
	}
	return nil
}

// Permission 权限模型
type Permission struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null;size:100" validate:"required,max=100"`
	Code        string         `json:"code" gorm:"not null;size:100;uniqueIndex" validate:"required,max=100"`
	Type        string         `json:"type" gorm:"not null;default:api;size:20" validate:"required,oneof=api menu"`
	Resource    string         `json:"resource" gorm:"size:255" validate:"max=255"`
	Action      string         `json:"action" gorm:"size:50" validate:"max=50"`
	Description string         `json:"description" gorm:"size:255" validate:"max=255"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// Menu 菜单模型
type Menu struct {
	ID             uint           `json:"id" gorm:"primarykey"`
	ParentID       uint           `json:"parent_id" gorm:"not null;default:0;index"`
	Name           string         `json:"name" gorm:"not null;size:100" validate:"required,max=100"`
	Path           string         `json:"path" gorm:"size:255" validate:"max=255"`
	Component      string         `json:"component" gorm:"size:255" validate:"max=255"`
	Icon           string         `json:"icon" gorm:"size:100" validate:"max=100"`
	SortOrder      int            `json:"sort_order" gorm:"not null;default:0"`
	IsHidden       bool           `json:"is_hidden" gorm:"not null;default:false"`
	IsSystem       bool           `json:"is_system" gorm:"not null;default:false"`
	PermissionCode string         `json:"permission_code" gorm:"size:100;index" validate:"max=100"`
	Meta           string         `json:"meta" gorm:"type:json"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_menus;"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}