package model

import (
	"gorm.io/gorm"
	"time"
)

// Role 角色模型
type Role struct {
	TenantBaseModel

	Name        string `json:"name" gorm:"not null;size:100" validate:"required,min=1,max=100"`
	Code        string `json:"code" gorm:"not null;size:50;uniqueIndex" validate:"required,min=1,max=50"`
	Description string `json:"description" gorm:"size:255" validate:"max=255"`
	Status      int    `json:"status" gorm:"not null;default:1" validate:"required,oneof=1 2"`
	IsSystem    bool   `json:"is_system" gorm:"not null;default:false"`
	SortOrder   int    `json:"sort_order" gorm:"not null;default:0"`

	// 关联关系
	Users  []User `json:"users,omitempty" gorm:"many2many:user_roles;"`
	Menus  []Menu `json:"menus,omitempty" gorm:"many2many:role_menu_permissions;"`
	Tenant Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
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
	ID          uint64         `json:"id" gorm:"primarykey"`
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

// Menu 菜单权限模型
type Menu struct {
	ID        uint64         `json:"id" gorm:"primarykey"`
	ParentID  uint64         `json:"parent_id" gorm:"not null;default:0;index"`
	Name      string         `json:"name" gorm:"not null;size:100" validate:"required,max=100"`
	Code      string         `json:"code" gorm:"not null;size:100;uniqueIndex" validate:"required,max=100"`
	Type      string         `json:"type" gorm:"not null;default:menu;size:20" validate:"required,oneof=directory menu permission"`
	Path      string         `json:"path" gorm:"size:255" validate:"max=255"`
	Component string         `json:"component" gorm:"size:255" validate:"max=255"`
	Icon      string         `json:"icon" gorm:"size:100" validate:"max=100"`
	Resource  string         `json:"resource" gorm:"size:255" validate:"max=255"`
	Action    string         `json:"action" gorm:"size:50" validate:"max=50"`
	SortOrder int            `json:"sort_order" gorm:"not null;default:0"`
	IsHidden  bool           `json:"is_hidden" gorm:"not null;default:false"`
	IsSystem  bool           `json:"is_system" gorm:"not null;default:false"`
	Status    int            `json:"status" gorm:"not null;default:1" validate:"required,oneof=1 2"`
	Meta      string         `json:"meta" gorm:"type:json"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Roles    []Role `json:"roles,omitempty" gorm:"many2many:role_menu_permissions;"`
	Children []Menu `json:"children,omitempty" gorm:"-"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// MenuProfile 菜单资料（简化版本）
type MenuProfile struct {
	ID        uint64        `json:"id"`
	ParentID  uint64        `json:"parent_id"`
	Name      string        `json:"name"`
	Code      string        `json:"code"`
	Type      string        `json:"type"`
	Path      string        `json:"path"`
	Component string        `json:"component"`
	Icon      string        `json:"icon"`
	Resource  string        `json:"resource"`
	Action    string        `json:"action"`
	SortOrder int           `json:"sort_order"`
	IsHidden  bool          `json:"is_hidden"`
	IsSystem  bool          `json:"is_system"`
	Status    int           `json:"status"`
	Meta      string        `json:"meta"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Children  []MenuProfile `json:"children,omitempty"`
}

// ToProfile 转换为菜单资料
func (m *Menu) ToProfile() MenuProfile {
	return MenuProfile{
		ID:        m.ID,
		ParentID:  m.ParentID,
		Name:      m.Name,
		Code:      m.Code,
		Type:      m.Type,
		Path:      m.Path,
		Component: m.Component,
		Icon:      m.Icon,
		Resource:  m.Resource,
		Action:    m.Action,
		SortOrder: m.SortOrder,
		IsHidden:  m.IsHidden,
		IsSystem:  m.IsSystem,
		Status:    m.Status,
		Meta:      m.Meta,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// MenuTreeNode 菜单树节点
type MenuTreeNode struct {
	MenuProfile
	Children []MenuTreeNode `json:"children,omitempty"`
}

// ToTreeNode 转换为菜单树节点
func (m *Menu) ToTreeNode() MenuTreeNode {
	return MenuTreeNode{
		MenuProfile: m.ToProfile(),
		Children:    make([]MenuTreeNode, 0),
	}
}

// BeforeCreate 创建前的钩子
func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.Status == 0 {
		m.Status = 1 // 默认启用
	}
	return nil
}
