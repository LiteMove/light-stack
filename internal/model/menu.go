package model

import (
	"gorm.io/gorm"
	"time"
)

// Menu 菜单权限模型
type Menu struct {
	ID        uint64         `json:"id" gorm:"primarykey"`
	ParentID  uint64         `json:"parentId" gorm:"not null;default:0;index"`
	Name      string         `json:"name" gorm:"not null;size:100" validate:"required,max=100"`
	Code      string         `json:"code" gorm:"not null;size:100;uniqueIndex" validate:"required,max=100"`
	Type      string         `json:"type" gorm:"not null;default:menu;size:20" validate:"required,oneof=directory menu permission"`
	Path      string         `json:"path" gorm:"size:255" validate:"max=255"`
	Component string         `json:"component" gorm:"size:255" validate:"max=255"`
	Icon      string         `json:"icon" gorm:"size:100" validate:"max=100"`
	Resource  string         `json:"resource" gorm:"size:255" validate:"max=255"`
	Action    string         `json:"action" gorm:"size:50" validate:"max=50"`
	SortOrder int            `json:"sortOrder" gorm:"not null;default:0"`
	IsHidden  bool           `json:"isHidden" gorm:"not null;default:false"`
	IsSystem  bool           `json:"isSystem" gorm:"not null;default:false"`
	Status    int            `json:"status" gorm:"not null;default:1" validate:"required,oneof=1 2"`
	Meta      string         `json:"meta" gorm:"type:json"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
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
	ParentID  uint64        `json:"parentId"`
	Name      string        `json:"name"`
	Code      string        `json:"code"`
	Type      string        `json:"type"`
	Path      string        `json:"path"`
	Component string        `json:"component"`
	Icon      string        `json:"icon"`
	Resource  string        `json:"resource"`
	Action    string        `json:"action"`
	SortOrder int           `json:"sortOrder"`
	IsHidden  bool          `json:"isHidden"`
	IsSystem  bool          `json:"isSystem"`
	Status    int           `json:"status"`
	Meta      string        `json:"meta"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
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
