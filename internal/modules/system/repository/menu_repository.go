package repository

import (
	"github.com/LiteMove/light-stack/internal/modules/system/model"
	"gorm.io/gorm"
)

// MenuRepository 菜单数据访问层接口
type MenuRepository interface {
	// 基础CRUD操作
	Create(menu *model.Menu) error
	GetByID(id uint64) (*model.Menu, error)
	GetByCode(code string) (*model.Menu, error)
	Update(menu *model.Menu) error
	Delete(id uint64) error

	// 查询操作
	GetList(offset, limit int, name string, status int) ([]model.Menu, int64, error)
	GetAll() ([]model.Menu, error)
	GetByParentID(parentID uint64) ([]model.Menu, error)
	GetTree() ([]model.Menu, error)
	GetUserMenus(userID uint64) ([]model.Menu, error)
	GetUserPermissions(userID uint64) ([]string, error)
	GetRoleMenus(roleID uint64) ([]model.Menu, error)

	// 状态操作
	UpdateStatus(id uint64, status int) error

	// 树形结构操作
	GetChildren(parentID uint64) ([]model.Menu, error)
	GetParents(id uint64) ([]model.Menu, error)
	IsParent(id uint64) (bool, error)
	UpdateSortOrder(id uint64, sortOrder int) error

	// 角色菜单关联操作
	AssignMenusToRole(roleID uint64, menuIDs []uint64) error
	RemoveMenusFromRole(roleID uint64, menuIDs []uint64) error
	ClearRoleMenus(roleID uint64) error
}

// menuRepository 菜单数据访问层实现
type menuRepository struct {
	db *gorm.DB
}

// NewMenuRepository 创建菜单数据访问层
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

// Create 创建菜单
func (r *menuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

// GetByID 根据ID获取菜单
func (r *menuRepository) GetByID(id uint64) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// GetByCode 根据代码获取菜单
func (r *menuRepository) GetByCode(code string) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.Where("code = ?", code).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// Update 更新菜单
func (r *menuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

// Delete 删除菜单
func (r *menuRepository) Delete(id uint64) error {
	return r.db.Delete(&model.Menu{}, id).Error
}

// GetList 获取菜单列表
func (r *menuRepository) GetList(offset, limit int, name string, status int) ([]model.Menu, int64, error) {
	var menus []model.Menu
	var total int64

	query := r.db.Model(&model.Menu{})

	// 条件过滤
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	err := query.Order("sort_order ASC, id ASC").
		Offset(offset).
		Limit(limit).
		Find(&menus).Error

	return menus, total, err
}

// GetAll 获取所有菜单
func (r *menuRepository) GetAll() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Order("sort_order ASC, id ASC").Find(&menus).Error
	return menus, err
}

// GetByParentID 根据父级ID获取菜单
func (r *menuRepository) GetByParentID(parentID uint64) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("parent_id = ?", parentID).
		Order("sort_order ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

// GetTree 获取菜单树
func (r *menuRepository) GetTree() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("status = ?", 1).
		Order("parent_id ASC, sort_order ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

// GetUserMenus 获取用户菜单
func (r *menuRepository) GetUserMenus(userID uint64) ([]model.Menu, error) {
	// 先检查用户是否为超级管理员
	var isSuper bool
	err := r.db.Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.code = 'super_admin' AND roles.status = 1", userID).
		Select("COUNT(*)").
		Row().Scan(&isSuper)
	if err != nil {
		return nil, err
	}

	var menus []model.Menu

	// 如果是超级管理员，返回所有菜单和目录（不包括权限类型）
	if isSuper {
		err := r.db.Where("status = ? AND type IN ('directory', 'menu')", 1).
			Order("sort_order ASC, id ASC").
			Find(&menus).Error
		return menus, err
	}

	// 普通用户通过角色权限查询
	err = r.db.Table("menus").
		Select("DISTINCT menus.*").
		Joins("JOIN role_menus ON menus.id = role_menus.menu_id").
		Joins("JOIN user_roles ON role_menus.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND menus.status = ? AND menus.type IN ('directory', 'menu')", userID, 1).
		Order("menus.sort_order ASC, menus.id ASC").
		Find(&menus).Error
	return menus, err
}

// GetUserPermissions 获取用户权限代码
func (r *menuRepository) GetUserPermissions(userID uint64) ([]string, error) {
	// 先检查用户是否为超级管理员
	var isSuper bool
	err := r.db.Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.code = 'super_admin' AND roles.status = 1", userID).
		Select("COUNT(*)").
		Row().Scan(&isSuper)
	if err != nil {
		return nil, err
	}

	var permissions []string

	// 如果是超级管理员，返回所有权限代码（type='permission'的code）
	if isSuper {
		err := r.db.Model(&model.Menu{}).
			Where("status = ? AND type = 'permission' AND code != ''", 1).
			Pluck("code", &permissions).Error
		return permissions, err
	}

	// 普通用户通过角色权限查询
	err = r.db.Table("menus").
		Select("DISTINCT menus.code").
		Joins("JOIN role_menus ON menus.id = role_menus.menu_id").
		Joins("JOIN user_roles ON role_menus.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND menus.status = ? AND menus.code != ''", userID, 1).
		Pluck("code", &permissions).Error
	return permissions, err
}

// GetRoleMenus 获取角色菜单
func (r *menuRepository) GetRoleMenus(roleID uint64) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Table("menus").
		Select("menus.*").
		Joins("JOIN role_menus ON menus.id = role_menus.menu_id").
		Where("role_menus.role_id = ? AND menus.status = ?", roleID, 1).
		Order("menus.sort_order ASC, menus.id ASC").
		Find(&menus).Error
	return menus, err
}

// UpdateStatus 更新菜单状态
func (r *menuRepository) UpdateStatus(id uint64, status int) error {
	return r.db.Model(&model.Menu{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// GetChildren 获取子菜单
func (r *menuRepository) GetChildren(parentID uint64) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("parent_id = ?", parentID).
		Order("sort_order ASC, id ASC").
		Find(&menus).Error
	return menus, err
}

// GetParents 获取父级菜单链
func (r *menuRepository) GetParents(id uint64) ([]model.Menu, error) {
	var parents []model.Menu
	var currentID = id

	for currentID != 0 {
		var menu model.Menu
		err := r.db.Where("id = ?", currentID).First(&menu).Error
		if err != nil {
			break
		}
		if menu.ParentID != 0 {
			parents = append(parents, menu)
		}
		currentID = menu.ParentID
	}

	// 反转切片，让顶级父节点在前面
	for i, j := 0, len(parents)-1; i < j; i, j = i+1, j-1 {
		parents[i], parents[j] = parents[j], parents[i]
	}

	return parents, nil
}

// IsParent 检查是否为父节点
func (r *menuRepository) IsParent(id uint64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).
		Where("parent_id = ?", id).
		Count(&count).Error
	return count > 0, err
}

// UpdateSortOrder 更新排序
func (r *menuRepository) UpdateSortOrder(id uint64, sortOrder int) error {
	return r.db.Model(&model.Menu{}).
		Where("id = ?", id).
		Update("sort_order", sortOrder).Error
}

// AssignMenusToRole 为角色分配菜单
func (r *menuRepository) AssignMenusToRole(roleID uint64, menuIDs []uint64) error {
	// 先清除现有关联
	if err := r.ClearRoleMenus(roleID); err != nil {
		return err
	}

	// 如果没有菜单ID，直接返回
	if len(menuIDs) == 0 {
		return nil
	}

	var roleMenu []model.RoleMenus
	for _, menuID := range menuIDs {
		roleMenu = append(roleMenu, model.RoleMenus{
			RoleId: roleID,
			MenuId: menuID,
		})
	}

	return r.db.Table("role_menus").Create(&roleMenu).Error
}

// RemoveMenusFromRole 从角色中移除菜单
func (r *menuRepository) RemoveMenusFromRole(roleID uint64, menuIDs []uint64) error {
	return r.db.Table("role_menus").
		Where("role_id = ? AND menu_id IN ?", roleID, menuIDs).
		Delete(nil).Error
}

// ClearRoleMenus 清除角色的所有菜单
func (r *menuRepository) ClearRoleMenus(roleID uint64) error {
	return r.db.Table("role_menus").
		Where("role_id = ?", roleID).
		Delete(nil).Error
}
