package model

import "time"

// RoleMenus 角色菜单权限关联模型
type RoleMenus struct {
	Id        uint64    `json:"id" gorm:"primarykey"`
	RoleId    uint64    `json:"roleId" gorm:"not null;uniqueIndex:uk_role_menu;index:idx_role_id" validate:"required"`
	MenuId    uint64    `json:"menuId" gorm:"not null;uniqueIndex:uk_role_menu;index:idx_menu_id" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
}

// TableName 指定表名
func (RoleMenus) TableName() string {
	return "role_menus"
}
