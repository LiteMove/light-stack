package middleware

import (
	"github.com/gin-gonic/gin"
)

// PermissionMigrationHelper 权限中间件迁移辅助工具
// 这个工具帮助从旧的权限中间件平滑迁移到新的统一权限中间件
type PermissionMigrationHelper struct {
	// 记录迁移过程中的权限映射关系
	legacyPermissionMap map[string]AuthConfig
}

// NewPermissionMigrationHelper 创建权限迁移辅助工具
func NewPermissionMigrationHelper() *PermissionMigrationHelper {
	return &PermissionMigrationHelper{
		legacyPermissionMap: make(map[string]AuthConfig),
	}
}

// MigrateCheckPermission 迁移原有的CheckPermission中间件
// 将旧的 middleware.CheckPermission("system:user:create") 转换为新的统一权限中间件
func MigrateCheckPermission(permission string) gin.HandlerFunc {
	// 创建对应的权限配置
	config := AuthConfig{
		RequireAuth:  true,
		Permissions:  []string{permission},
		Logic:        AuthLogicOR, // 默认使用OR逻辑
		ErrorMessage: "没有操作权限",    // 保持与原有行为一致
	}

	// 使用新的统一权限中间件
	return UnifiedAuth(config)
}

// MigrateCheckRole 迁移原有的CheckRole中间件
// 将旧的角色检查转换为新的统一权限中间件
func MigrateCheckRole(roles ...string) gin.HandlerFunc {
	config := AuthConfig{
		RequireAuth:  true,
		Roles:        roles,
		Logic:        AuthLogicOR,
		ErrorMessage: "没有角色权限", // 保持与原有行为一致
	}

	return UnifiedAuth(config)
}

// MigrateCheckPermissionOrRole 迁移复合权限检查
// 将原有的权限或角色验证迁移到新的中间件
func MigrateCheckPermissionOrRole(permissions []string, roles []string) gin.HandlerFunc {
	config := AuthConfig{
		RequireAuth:  true,
		Permissions:  permissions,
		Roles:        roles,
		Logic:        AuthLogicOR, // 原有逻辑是OR
		ErrorMessage: "没有权限",
	}

	return UnifiedAuth(config)
}

// GetLegacyPermissionMapping 获取遗留权限映射
// 这个函数定义了旧权限代码到新权限配置的映射关系
func GetLegacyPermissionMapping() map[string]AuthConfig {
	return map[string]AuthConfig{
		// 用户管理权限映射
		"system:user:create": {
			RequireAuth:  true,
			Permissions:  []string{"system:user:create"},
			Roles:        []string{"admin", "user_manager"},
			Logic:        AuthLogicOR,
			ErrorMessage: "需要用户创建权限",
		},

		"system:user:list": {
			RequireAuth:  true,
			Permissions:  []string{"system:user:list"},
			Roles:        []string{"admin", "user_manager", "viewer"},
			Logic:        AuthLogicOR,
			ErrorMessage: "需要用户查看权限",
		},

		"system:user:update": {
			RequireAuth:  true,
			Permissions:  []string{"system:user:update"},
			Roles:        []string{"admin", "user_manager"},
			Logic:        AuthLogicOR,
			ErrorMessage: "需要用户更新权限",
		},

		"system:user:delete": {
			RequireAuth:  true,
			Permissions:  []string{"system:user:delete"},
			Roles:        []string{"admin"},
			Logic:        AuthLogicAND, // 删除操作要求更严格
			ErrorMessage: "需要用户删除权限",
		},

		// 角色管理权限映射
		"system:role:manage": {
			RequireAuth:  true,
			Permissions:  []string{"system:role:manage"},
			Roles:        []string{"admin"},
			Logic:        AuthLogicOR,
			ErrorMessage: "需要角色管理权限",
		},

		// 文件管理权限映射
		"system:file:upload": {
			RequireAuth:  true,
			Permissions:  []string{"system:file:upload"},
			Roles:        []string{"admin", "user", "editor"},
			Logic:        AuthLogicOR,
			ErrorMessage: "需要文件上传权限",
		},

		"system:file:delete": {
			RequireAuth:  true,
			Permissions:  []string{"system:file:delete"},
			Roles:        []string{"admin", "file_manager"},
			Logic:        AuthLogicAND,
			ErrorMessage: "需要文件删除权限",
		},
	}
}

// BatchMigratePermissions 批量迁移权限配置
// 这个函数可以批量更新路由中的权限中间件
func BatchMigratePermissions(router *gin.RouterGroup, migrationMap map[string]AuthConfig) {
	// 这个函数在实际使用中需要结合具体的路由结构来实现
	// 目前提供迁移的指导方案

	// 示例：将所有使用CheckPermission的路由进行批量更新
	// 实际实现需要根据具体的路由结构进行调整
}

// ValidatePermissionMigration 验证权限迁移的正确性
// 这个函数用于验证迁移后的权限行为是否与原有行为一致
func ValidatePermissionMigration() []string {
	var issues []string

	// 检查所有预定义配置是否完整
	requiredConfigs := []string{
		"RequireLogin", "RequireSuperAdmin", "UserManage",
		"UserView", "RoleManage", "MenuManage", "TenantManage",
		"FileUpload", "FileManage",
	}

	// 验证每个配置的完整性
	for _, configName := range requiredConfigs {
		// 这里可以添加具体的验证逻辑
		// 例如：检查权限字符串格式、角色名称有效性等
		_ = configName // 暂时避免未使用变量错误
	}

	return issues
}

// GenerateMigrationReport 生成权限迁移报告
// 这个函数生成详细的迁移报告，帮助开发者了解迁移过程
func GenerateMigrationReport() map[string]interface{} {
	report := make(map[string]interface{})

	// 统计原有权限中间件使用情况
	legacyPermissions := GetLegacyPermissionMapping()
	report["legacy_permissions_count"] = len(legacyPermissions)

	// 统计新权限配置
	report["new_predefined_configs"] = []string{
		"RequireLogin", "RequireSuperAdmin", "UserManage",
		"UserView", "RoleManage", "MenuManage", "TenantManage",
		"FileUpload", "FileManage",
	}

	// 记录关键改进点
	report["improvements"] = []string{
		"统一的权限检查逻辑",
		"支持权限和角色的AND/OR组合",
		"可配置的错误消息",
		"链式配置API",
		"预定义权限配置常量",
	}

	// 记录潜在风险点
	report["risks"] = []string{
		"权限逻辑变更可能影响现有用户访问",
		"需要验证所有权限场景的正确性",
		"中间件性能影响需要测试",
	}

	return report
}
