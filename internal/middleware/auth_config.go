package middleware

// AuthLogic 权限逻辑类型
type AuthLogic int

const (
	// AuthLogicOR 权限或角色任一满足(默认)
	AuthLogicOR AuthLogic = iota
	// AuthLogicAND 权限和角色都要满足
	AuthLogicAND
)

// AuthConfig 统一权限配置
type AuthConfig struct {
	// 是否需要登录
	RequireAuth bool `json:"require_auth"`

	// 需要的权限列表
	Permissions []string `json:"permissions"`

	// 需要的角色列表
	Roles []string `json:"roles"`

	// 是否需要超管权限
	RequireSuperAdmin bool `json:"require_super_admin"`

	// 权限逻辑: AND/OR
	Logic AuthLogic `json:"logic"`

	// 自定义错误消息
	ErrorMessage string `json:"error_message"`
}

// 预定义权限配置常量

var (
	// 无需认证
	NoAuth = AuthConfig{RequireAuth: false}

	// 仅需登录
	RequireLogin = AuthConfig{RequireAuth: true}

	// 需要超管权限
	RequireSuperAdmin = AuthConfig{
		RequireAuth:       true,
		RequireSuperAdmin: true,
		ErrorMessage:      "需要超级管理员权限",
	}

	// 用户管理权限
	UserManage = AuthConfig{
		RequireAuth:  true,
		Permissions:  []string{"system:user:manage", "system:user:create", "system:user:update", "system:user:delete"},
		Roles:        []string{"admin", "user_manager"},
		Logic:        AuthLogicOR,
		ErrorMessage: "需要用户管理权限",
	}

	// 用户查看权限
	UserView = AuthConfig{
		RequireAuth:  true,
		Permissions:  []string{"system:user:list", "system:user:detail"},
		Roles:        []string{"admin", "user_manager", "viewer"},
		Logic:        AuthLogicOR,
		ErrorMessage: "需要用户查看权限",
	}

	// 角色管理权限
	RoleManage = AuthConfig{
		RequireAuth:  true,
		Permissions:  []string{"system:role:manage", "system:role:create", "system:role:update", "system:role:delete"},
		Roles:        []string{"admin"},
		Logic:        AuthLogicOR,
		ErrorMessage: "需要角色管理权限",
	}

	// 菜单管理权限
	MenuManage = AuthConfig{
		RequireAuth:  true,
		Permissions:  []string{"system:menu:manage", "system:menu:create", "system:menu:update", "system:menu:delete"},
		Roles:        []string{"admin"},
		Logic:        AuthLogicOR,
		ErrorMessage: "需要菜单管理权限",
	}

	// 租户管理权限（仅超管）
	TenantManage = AuthConfig{
		RequireAuth:       true,
		RequireSuperAdmin: true,
		ErrorMessage:      "需要超级管理员权限管理租户",
	}

	// 文件上传权限
	FileUpload = AuthConfig{
		RequireAuth:  true,
		Permissions:  []string{"system:file:upload"},
		Roles:        []string{"admin", "user", "editor"},
		Logic:        AuthLogicOR,
		ErrorMessage: "需要文件上传权限",
	}

	// 文件管理权限
	FileManage = AuthConfig{
		RequireAuth:  true,
		Permissions:  []string{"system:file:manage", "system:file:delete"},
		Roles:        []string{"admin", "file_manager"},
		Logic:        AuthLogicOR,
		ErrorMessage: "需要文件管理权限",
	}
)

// WithCustomMessage 设置自定义错误消息
func (ac AuthConfig) WithCustomMessage(message string) AuthConfig {
	ac.ErrorMessage = message
	return ac
}

// WithLogic 设置权限逻辑
func (ac AuthConfig) WithLogic(logic AuthLogic) AuthConfig {
	ac.Logic = logic
	return ac
}

// WithPermissions 添加权限要求
func (ac AuthConfig) WithPermissions(permissions ...string) AuthConfig {
	ac.Permissions = append(ac.Permissions, permissions...)
	return ac
}

// WithRoles 添加角色要求
func (ac AuthConfig) WithRoles(roles ...string) AuthConfig {
	ac.Roles = append(ac.Roles, roles...)
	return ac
}
