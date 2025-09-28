package model

import (
	"gorm.io/gorm"
	"time"
)

// Role 角色模型
type Role struct {
	BaseModel
	Name        string `json:"name" gorm:"not null;size:100" validate:"required,min=1,max=100"`
	Code        string `json:"code" gorm:"not null;size:50;uniqueIndex" validate:"required,min=1,max=50"`
	Description string `json:"description" gorm:"size:255" validate:"max=255"`
	Status      int    `json:"status" gorm:"not null;default:1" validate:"required,oneof=1 2"`
	IsSystem    bool   `json:"isSystem" gorm:"not null;default:false"`
	SortOrder   int    `json:"sortOrder" gorm:"not null;default:0"`

	// 关联关系
	Users []User `json:"users,omitempty" gorm:"many2many:user_roles;"`
	Menus []Menu `json:"menus,omitempty" gorm:"many2many:role_menus;"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// RoleProfile 角色资料（简化版本）
type RoleProfile struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	IsSystem    bool      `json:"isSystem"`
	SortOrder   int       `json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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

var SuperAdminId uint64 = 1 // 超级管理员角色ID

const (
	RoleStatusEnabled  = 1 // 启用
	RoleStatusDisabled = 2 // 禁用
)
