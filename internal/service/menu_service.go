package service

import (
	"errors"
	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/repository"
	"gorm.io/gorm"
)

// MenuService 菜单服务接口
type MenuService interface {
	// 基础CRUD操作
	CreateMenu(menu *model.Menu) error
	GetMenu(id uint64) (*model.Menu, error)
	UpdateMenu(menu *model.Menu) error
	DeleteMenu(id uint64) error

	// 查询操作
	GetMenuList(page, pageSize int, name string, status int) ([]model.MenuProfile, int64, error)
	GetMenuTree() ([]model.MenuTreeNode, error)
	GetUserMenuTree(userID uint64) ([]model.MenuTreeNode, error)
	GetRoleMenus(roleID uint64) ([]model.MenuProfile, error)

	// 状态操作
	UpdateMenuStatus(id uint64, status int) error
	BatchUpdateMenuStatus(ids []uint64, status int) error

	// 权限相关
	AssignMenusToRole(roleID uint64, menuIDs []uint64) error
	GetMenuPermissions(userID uint64) ([]string, error)
	CheckMenuPermission(userID uint64, menuCode string) (bool, error)
}

// menuService 菜单服务实现
type menuService struct {
	menuRepo repository.MenuRepository
	roleRepo repository.RoleRepository
}

// NewMenuService 创建菜单服务
func NewMenuService(menuRepo repository.MenuRepository, roleRepo repository.RoleRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
		roleRepo: roleRepo,
	}
}

// CreateMenu 创建菜单
func (s *menuService) CreateMenu(menu *model.Menu) error {
	// 检查菜单代码是否已存在
	existMenu, err := s.menuRepo.GetByCode(menu.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existMenu != nil {
		return errors.New("菜单代码已存在")
	}

	// 如果有父菜单，检查父菜单是否存在
	if menu.ParentID != 0 {
		parent, err := s.menuRepo.GetByID(menu.ParentID)
		if err != nil {
			return errors.New("父菜单不存在")
		}
		if parent.Status != 1 {
			return errors.New("父菜单已禁用")
		}
	}

	return s.menuRepo.Create(menu)
}

// GetMenu 获取菜单
func (s *menuService) GetMenu(id uint64) (*model.Menu, error) {
	return s.menuRepo.GetByID(id)
}

// UpdateMenu 更新菜单
func (s *menuService) UpdateMenu(menu *model.Menu) error {
	// 检查菜单是否存在
	existMenu, err := s.menuRepo.GetByID(menu.ID)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 如果修改了代码，检查新代码是否已存在
	if existMenu.Code != menu.Code {
		codeMenu, err := s.menuRepo.GetByCode(menu.Code)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if codeMenu != nil {
			return errors.New("菜单代码已存在")
		}
	}

	// 检查父菜单有效性
	if menu.ParentID != 0 {
		// 不能将自己设为父菜单
		if menu.ParentID == menu.ID {
			return errors.New("不能将自己设为父菜单")
		}

		// 检查父菜单是否存在
		parent, err := s.menuRepo.GetByID(menu.ParentID)
		if err != nil {
			return errors.New("父菜单不存在")
		}
		if parent.Status != 1 {
			return errors.New("父菜单已禁用")
		}

		// 检查是否形成循环引用
		if s.hasCircularReference(menu.ID, menu.ParentID) {
			return errors.New("不能形成循环引用")
		}
	}

	return s.menuRepo.Update(menu)
}

// DeleteMenu 删除菜单
func (s *menuService) DeleteMenu(id uint64) error {
	// 检查菜单是否存在
	menu, err := s.menuRepo.GetByID(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 检查是否为系统菜单
	if menu.IsSystem {
		return errors.New("系统菜单不能删除")
	}

	// 检查是否有子菜单
	isParent, err := s.menuRepo.IsParent(id)
	if err != nil {
		return err
	}
	if isParent {
		return errors.New("存在子菜单，不能删除")
	}

	return s.menuRepo.Delete(id)
}

// GetMenuList 获取菜单列表
func (s *menuService) GetMenuList(page, pageSize int, name string, status int) ([]model.MenuProfile, int64, error) {
	offset := (page - 1) * pageSize
	menus, total, err := s.menuRepo.GetList(offset, pageSize, name, status)
	if err != nil {
		return nil, 0, err
	}

	profiles := make([]model.MenuProfile, len(menus))
	for i, menu := range menus {
		profiles[i] = menu.ToProfile()
	}

	return profiles, total, nil
}

// GetMenuTree 获取菜单树
func (s *menuService) GetMenuTree() ([]model.MenuTreeNode, error) {
	menus, err := s.menuRepo.GetTree()
	if err != nil {
		return nil, err
	}

	return s.buildMenuTree(menus, 0), nil
}

// GetUserMenuTree 获取用户菜单树
func (s *menuService) GetUserMenuTree(userID uint64) ([]model.MenuTreeNode, error) {
	menus, err := s.menuRepo.GetUserMenus(userID)
	if err != nil {
		return nil, err
	}

	return s.buildMenuTree(menus, 0), nil
}

// GetRoleMenus 获取角色菜单
func (s *menuService) GetRoleMenus(roleID uint64) ([]model.MenuProfile, error) {
	menus, err := s.menuRepo.GetRoleMenus(roleID)
	if err != nil {
		return nil, err
	}

	profiles := make([]model.MenuProfile, len(menus))
	for i, menu := range menus {
		profiles[i] = menu.ToProfile()
	}

	return profiles, nil
}

// UpdateMenuStatus 更新菜单状态
func (s *menuService) UpdateMenuStatus(id uint64, status int) error {
	// 检查菜单是否存在
	menu, err := s.menuRepo.GetByID(id)
	if err != nil {
		return errors.New("菜单不存在")
	}

	// 检查是否为系统菜单
	if menu.IsSystem && status != 1 {
		return errors.New("系统菜单不能禁用")
	}

	return s.menuRepo.UpdateStatus(id, status)
}

// BatchUpdateMenuStatus 批量更新菜单状态
func (s *menuService) BatchUpdateMenuStatus(ids []uint64, status int) error {
	// 检查系统菜单
	for _, id := range ids {
		menu, err := s.menuRepo.GetByID(id)
		if err != nil {
			continue
		}
		if menu.IsSystem && status != 1 {
			return errors.New("系统菜单不能禁用")
		}
	}

	return s.menuRepo.BatchUpdateStatus(ids, status)
}

// AssignMenusToRole 为角色分配菜单
func (s *menuService) AssignMenusToRole(roleID uint64, menuIDs []uint64) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return errors.New("角色不存在")
	}

	// 检查菜单是否存在
	for _, menuID := range menuIDs {
		_, err := s.menuRepo.GetByID(menuID)
		if err != nil {
			return errors.New("菜单不存在")
		}
	}

	// 调用repository层的方法
	return s.menuRepo.AssignMenusToRole(roleID, menuIDs)
}

// GetMenuPermissions 获取用户菜单权限
func (s *menuService) GetMenuPermissions(userID uint64) ([]string, error) {
	menus, err := s.menuRepo.GetUserMenus(userID)
	if err != nil {
		return nil, err
	}

	permissions := make([]string, 0, len(menus))
	for _, menu := range menus {
		if menu.Code != "" {
			permissions = append(permissions, menu.Code)
		}
	}

	return permissions, nil
}

// CheckMenuPermission 检查菜单权限
func (s *menuService) CheckMenuPermission(userID uint64, menuCode string) (bool, error) {
	permissions, err := s.GetMenuPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, permission := range permissions {
		if permission == menuCode {
			return true, nil
		}
	}

	return false, nil
}

// hasCircularReference 检查是否存在循环引用
func (s *menuService) hasCircularReference(childID, parentID uint64) bool {
	// 简单的循环检查逻辑
	visited := make(map[uint64]bool)
	current := parentID

	for current != 0 {
		if visited[current] {
			return true
		}
		if current == childID {
			return true
		}

		visited[current] = true
		parent, err := s.menuRepo.GetByID(current)
		if err != nil {
			break
		}
		current = parent.ParentID
	}

	return false
}

// buildMenuTree 构建菜单树
func (s *menuService) buildMenuTree(menus []model.Menu, parentID uint64) []model.MenuTreeNode {
	var tree []model.MenuTreeNode

	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := menu.ToTreeNode()
			node.Children = s.buildMenuTree(menus, menu.ID)
			tree = append(tree, node)
		}
	}

	return tree
}
