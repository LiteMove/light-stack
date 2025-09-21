package model

import (
	"time"
)

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint64    `json:"id" gorm:"primarykey"`
	UserID    uint64    `json:"userId" gorm:"not null;uniqueIndex:uk_user_role;index:idx_user_id" validate:"required"`
	RoleID    uint64    `json:"roleId" gorm:"not null;uniqueIndex:uk_user_role;index:idx_role_id" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RoleMenuPermission 角色菜单权限关联模型
type RoleMenuPermission struct {
	ID        uint64    `json:"id" gorm:"primarykey"`
	RoleID    uint64    `json:"roleId" gorm:"not null;uniqueIndex:uk_role_menu;index:idx_role_id" validate:"required"`
	MenuID    uint64    `json:"menuId" gorm:"not null;uniqueIndex:uk_role_menu;index:idx_menu_id" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName 指定表名
func (RoleMenuPermission) TableName() string {
	return "role_menu_permissions"
}

// UserRoleInfo 用户角色信息（包含角色详情）
type UserRoleInfo struct {
	UserID      uint64    `json:"userId"`
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	Email       string    `json:"email"`
	Roles       []Role    `json:"roles"`
	RoleCodes   []string  `json:"role_codes"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
}

// RoleWithUsers 角色及其用户信息
type RoleWithUsers struct {
	Role
	UserCount int    `json:"user_count"`
	Users     []User `json:"users,omitempty"`
}
