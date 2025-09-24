package model

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单权限模型
type Menu struct {
	ID        uint64         `json:"id" gorm:"primarykey"`
	ParentID  uint64         `json:"parentId" gorm:"not null;default:0;index"`
	Name      string         `json:"name" gorm:"not null;size:100" validate:"required,max=100"`
	Code      string         `json:"code" gorm:"size:100;index" validate:"required_if_permission,permission_code,max=100"`
	Type      string         `json:"type" gorm:"not null;default:menu;size:20" validate:"required,oneof=directory menu permission"`
	Path      string         `json:"path" gorm:"size:255" validate:"max=255"`
	Component string         `json:"component" gorm:"size:255" validate:"max=255"`
	Icon      string         `json:"icon" gorm:"size:100" validate:"max=100"`
	SortOrder int            `json:"sortOrder" gorm:"not null;default:0"`
	IsHidden  bool           `json:"isHidden" gorm:"not null;default:false"`
	Status    int            `json:"status" gorm:"not null;default:1" validate:"required,oneof=1 2"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Roles    []Role `json:"roles,omitempty" gorm:"many2many:role_menus;"`
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
	SortOrder int           `json:"sortOrder"`
	IsHidden  bool          `json:"isHidden"`
	Status    int           `json:"status"`
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
		SortOrder: m.SortOrder,
		IsHidden:  m.IsHidden,
		Status:    m.Status,
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
