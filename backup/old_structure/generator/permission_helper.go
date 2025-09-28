package generator

import (
	"fmt"
	"strings"

	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/utils"
)

// PermissionHelper 权限辅助工具
type PermissionHelper struct{}

// NewPermissionHelper 创建权限辅助工具
func NewPermissionHelper() *PermissionHelper {
	return &PermissionHelper{}
}

// GeneratePermissions 生成权限字符串列表
func (h *PermissionHelper) GeneratePermissions(moduleName, businessName string) []string {
	return utils.GeneratePermissions(moduleName, businessName)
}

// GeneratePermissionWithCustomOps 生成自定义操作的权限字符串
func (h *PermissionHelper) GeneratePermissionWithCustomOps(moduleName, businessName string, operations []string) []string {
	permissions := make([]string, len(operations))
	for i, op := range operations {
		permissions[i] = h.FormatPermission(moduleName, businessName, op)
	}
	return permissions
}

// FormatPermission 格式化权限字符串
func (h *PermissionHelper) FormatPermission(moduleName, businessName, operation string) string {
	return fmt.Sprintf("%s:%s:%s",
		strings.ToLower(moduleName),
		strings.ToLower(businessName),
		strings.ToLower(operation))
}

// ParsePermission 解析权限字符串
func (h *PermissionHelper) ParsePermission(permission string) (*Permission, error) {
	parts := strings.Split(permission, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("权限字符串格式错误，应为 '模块:业务:操作' 格式")
	}

	return &Permission{
		Module:    parts[0],
		Business:  parts[1],
		Operation: parts[2],
	}, nil
}

// ValidatePermission 验证权限字符串格式
func (h *PermissionHelper) ValidatePermission(permission string) error {
	if permission == "" {
		return fmt.Errorf("权限字符串不能为空")
	}

	parts := strings.Split(permission, ":")
	if len(parts) != 3 {
		return fmt.Errorf("权限字符串格式错误，应为 '模块:业务:操作' 格式")
	}

	for i, part := range parts {
		if strings.TrimSpace(part) == "" {
			switch i {
			case 0:
				return fmt.Errorf("模块名不能为空")
			case 1:
				return fmt.Errorf("业务名不能为空")
			case 2:
				return fmt.Errorf("操作名不能为空")
			}
		}
	}

	return nil
}

// ValidatePermissions 批量验证权限字符串
func (h *PermissionHelper) ValidatePermissions(permissions []string) error {
	for _, permission := range permissions {
		if err := h.ValidatePermission(permission); err != nil {
			return err
		}
	}
	return nil
}

// GenerateMenuPermissions 生成菜单相关的权限
func (h *PermissionHelper) GenerateMenuPermissions(moduleName, businessName string) map[string]string {
	basePermission := fmt.Sprintf("%s:%s", strings.ToLower(moduleName), strings.ToLower(businessName))

	return map[string]string{
		"menu":   basePermission + ":menu",   // 菜单权限
		"list":   basePermission + ":list",   // 列表权限
		"add":    basePermission + ":add",    // 新增权限
		"edit":   basePermission + ":edit",   // 编辑权限
		"delete": basePermission + ":delete", // 删除权限
		"view":   basePermission + ":view",   // 查看权限
		"export": basePermission + ":export", // 导出权限
		"import": basePermission + ":import", // 导入权限
	}
}

// GeneratePermissionSQL 生成权限SQL语句
func (h *PermissionHelper) GeneratePermissionSQL(data *model.TemplateData) string {
	var sqlBuilder strings.Builder

	sqlBuilder.WriteString("-- 权限配置SQL\n")
	sqlBuilder.WriteString(fmt.Sprintf("-- 模块: %s, 业务: %s\n\n", data.ModuleName, data.BusinessName))

	// 生成权限插入语句
	for i, permission := range data.Permissions {
		perm, _ := h.ParsePermission(permission)
		permissionName := h.getPermissionName(perm.Operation)
		permissionDesc := fmt.Sprintf("%s%s", data.FunctionName, permissionName)

		sqlBuilder.WriteString(fmt.Sprintf(
			"INSERT INTO sys_permissions (permission_code, permission_name, permission_desc, module_name, business_name, operation_type, sort, status, create_time) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', %d, 1, NOW());\n",
			permission,
			permissionName,
			permissionDesc,
			data.ModuleName,
			data.BusinessName,
			perm.Operation,
			i+1,
		))
	}

	return sqlBuilder.String()
}

// getPermissionName 获取权限名称
func (h *PermissionHelper) getPermissionName(operation string) string {
	names := map[string]string{
		"list":   "查看",
		"add":    "新增",
		"edit":   "编辑",
		"delete": "删除",
		"view":   "详情",
		"export": "导出",
		"import": "导入",
		"menu":   "菜单",
	}

	if name, exists := names[operation]; exists {
		return name
	}
	return operation
}

// GroupPermissionsByModule 按模块分组权限
func (h *PermissionHelper) GroupPermissionsByModule(permissions []string) map[string][]string {
	groups := make(map[string][]string)

	for _, permission := range permissions {
		if perm, err := h.ParsePermission(permission); err == nil {
			groups[perm.Module] = append(groups[perm.Module], permission)
		}
	}

	return groups
}

// GroupPermissionsByBusiness 按业务分组权限
func (h *PermissionHelper) GroupPermissionsByBusiness(permissions []string) map[string][]string {
	groups := make(map[string][]string)

	for _, permission := range permissions {
		if perm, err := h.ParsePermission(permission); err == nil {
			key := fmt.Sprintf("%s:%s", perm.Module, perm.Business)
			groups[key] = append(groups[key], permission)
		}
	}

	return groups
}

// FilterPermissionsByOperation 按操作过滤权限
func (h *PermissionHelper) FilterPermissionsByOperation(permissions []string, operations []string) []string {
	operationSet := make(map[string]bool)
	for _, op := range operations {
		operationSet[strings.ToLower(op)] = true
	}

	var filtered []string
	for _, permission := range permissions {
		if perm, err := h.ParsePermission(permission); err == nil {
			if operationSet[strings.ToLower(perm.Operation)] {
				filtered = append(filtered, permission)
			}
		}
	}

	return filtered
}

// GetStandardOperations 获取标准操作列表
func (h *PermissionHelper) GetStandardOperations() []string {
	return []string{"list", "add", "edit", "delete", "view"}
}

// GetExtendedOperations 获取扩展操作列表
func (h *PermissionHelper) GetExtendedOperations() []string {
	return []string{"list", "add", "edit", "delete", "view", "export", "import", "audit"}
}

// Permission 权限结构
type Permission struct {
	Module    string `json:"module"`    // 模块
	Business  string `json:"business"`  // 业务
	Operation string `json:"operation"` // 操作
}

// String 权限字符串表示
func (p *Permission) String() string {
	return fmt.Sprintf("%s:%s:%s", p.Module, p.Business, p.Operation)
}

// IsValid 检查权限是否有效
func (p *Permission) IsValid() bool {
	return p.Module != "" && p.Business != "" && p.Operation != ""
}

// PermissionTemplate 权限模板
type PermissionTemplate struct {
	ModuleName   string   `json:"moduleName"`   // 模块名称
	BusinessName string   `json:"businessName"` // 业务名称
	Operations   []string `json:"operations"`   // 操作列表
	Permissions  []string `json:"permissions"`  // 生成的权限列表
	Description  string   `json:"description"`  // 描述
}

// GenerateFromTemplate 从模板生成权限
func (h *PermissionHelper) GenerateFromTemplate(template *PermissionTemplate) *PermissionTemplate {
	template.Permissions = h.GeneratePermissionWithCustomOps(
		template.ModuleName,
		template.BusinessName,
		template.Operations,
	)
	return template
}
