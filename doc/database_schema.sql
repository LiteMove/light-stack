-- ========================================
-- light-stack 数据库表结构设计
-- ========================================

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ========================================
-- 1. 租户管理表
-- ========================================

-- 租户信息表
CREATE TABLE `tenants` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '租户ID',
  `name` varchar(100) NOT NULL COMMENT '租户名称',
  `domain` varchar(100) DEFAULT NULL COMMENT '租户域名',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '租户状态：1-启用 2-禁用 3-试用 4-过期',
  `expired_at` datetime DEFAULT NULL COMMENT '过期时间',
  `config` json DEFAULT NULL COMMENT '租户配置信息（Logo、主题色等）',
  `max_users` int(11) DEFAULT 100 COMMENT '最大用户数',
  `max_roles` int(11) DEFAULT 20 COMMENT '最大角色数',
  `max_storage` bigint(20) DEFAULT 104857600 COMMENT '最大存储空间（字节）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_domain` (`domain`),
  KEY `idx_status` (`status`),
  KEY `idx_expired_at` (`expired_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户信息表';

-- ========================================
-- 2. 用户管理表
-- ========================================

-- 用户表
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `tenant_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '租户ID，0表示系统管理员',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码（bcrypt加密）',
  `nickname` varchar(100) DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像地址',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '用户状态：1-启用 2-禁用 3-锁定',
  `is_system` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统用户：0-否 1-是',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) DEFAULT NULL COMMENT '最后登录IP',
  `login_failures` int(11) NOT NULL DEFAULT 0 COMMENT '连续登录失败次数',
  `locked_until` datetime DEFAULT NULL COMMENT '锁定截止时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_username` (`tenant_id`, `username`),
  UNIQUE KEY `uk_tenant_email` (`tenant_id`, `email`),
  UNIQUE KEY `uk_tenant_phone` (`tenant_id`, `phone`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_system` (`is_system`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ========================================
-- 3. 角色权限管理表
-- ========================================

-- 角色表
CREATE TABLE `roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `tenant_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '租户ID，0表示系统角色',
  `name` varchar(100) NOT NULL COMMENT '角色名称',
  `code` varchar(50) NOT NULL COMMENT '角色编码',
  `description` varchar(255) DEFAULT NULL COMMENT '角色描述',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '角色状态：1-启用 2-禁用',
  `is_system` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统角色：0-否 1-是',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序号',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_code` (`tenant_id`, `code`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_system` (`is_system`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 权限表
CREATE TABLE `permissions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '权限ID',
  `name` varchar(100) NOT NULL COMMENT '权限名称',
  `code` varchar(100) NOT NULL COMMENT '权限编码',
  `type` varchar(20) NOT NULL DEFAULT 'api' COMMENT '权限类型：api-接口权限 menu-菜单权限',
  `resource` varchar(255) DEFAULT NULL COMMENT '资源路径',
  `action` varchar(50) DEFAULT NULL COMMENT '操作类型',
  `description` varchar(255) DEFAULT NULL COMMENT '权限描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 用户角色关联表
CREATE TABLE `user_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 角色权限关联表
CREATE TABLE `role_permissions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `permission_id` bigint(20) NOT NULL COMMENT '权限ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- ========================================
-- 4. 菜单管理表
-- ========================================

-- 菜单表
CREATE TABLE `menus` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '父菜单ID，0表示顶级菜单',
  `name` varchar(100) NOT NULL COMMENT '菜单名称',
  `path` varchar(255) DEFAULT NULL COMMENT '路由路径',
  `component` varchar(255) DEFAULT NULL COMMENT '组件路径',
  `icon` varchar(100) DEFAULT NULL COMMENT '菜单图标',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序号',
  `is_hidden` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否隐藏：0-显示 1-隐藏',
  `is_system` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统菜单：0-否 1-是',
  `permission_code` varchar(100) DEFAULT NULL COMMENT '关联权限编码',
  `meta` json DEFAULT NULL COMMENT '菜单元数据',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_is_system` (`is_system`),
  KEY `idx_permission_code` (`permission_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- 租户菜单可见性配置表
CREATE TABLE `tenant_menus` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `tenant_id` bigint(20) NOT NULL COMMENT '租户ID',
  `menu_id` bigint(20) NOT NULL COMMENT '菜单ID',
  `is_visible` tinyint(1) NOT NULL DEFAULT 1 COMMENT '是否可见：0-不可见 1-可见',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_menu` (`tenant_id`, `menu_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户菜单可见性配置表';

-- 角色菜单关联表
CREATE TABLE `role_menus` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `menu_id` bigint(20) NOT NULL COMMENT '菜单ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_menu` (`role_id`, `menu_id`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- ========================================
-- 5. 数据字典表
-- ========================================

-- 字典类型表
CREATE TABLE `dict_types` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字典类型ID',
  `name` varchar(100) NOT NULL COMMENT '字典名称',
  `type` varchar(100) NOT NULL COMMENT '字典类型',
  `description` varchar(255) DEFAULT NULL COMMENT '描述',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1-启用 2-禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_type` (`type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字典类型表';

-- 字典数据表
CREATE TABLE `dict_data` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字典数据ID',
  `dict_type` varchar(100) NOT NULL COMMENT '字典类型',
  `label` varchar(100) NOT NULL COMMENT '字典标签',
  `value` varchar(100) NOT NULL COMMENT '字典键值',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序号',
  `css_class` varchar(100) DEFAULT NULL COMMENT 'CSS类名',
  `list_class` varchar(100) DEFAULT NULL COMMENT '列表样式',
  `is_default` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否默认：0-否 1-是',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1-启用 2-禁用',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_type_value` (`dict_type`, `value`),
  KEY `idx_dict_type` (`dict_type`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字典数据表';

-- ========================================
-- 6. 文件管理表
-- ========================================

-- 文件表
CREATE TABLE `files` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `tenant_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '租户ID，0表示系统文件',
  `original_name` varchar(255) NOT NULL COMMENT '原始文件名',
  `file_name` varchar(255) NOT NULL COMMENT '存储文件名',
  `file_path` varchar(500) NOT NULL COMMENT '文件路径',
  `file_size` bigint(20) NOT NULL COMMENT '文件大小（字节）',
  `file_type` varchar(100) NOT NULL COMMENT '文件类型',
  `mime_type` varchar(100) NOT NULL COMMENT 'MIME类型',
  `md5` varchar(32) NOT NULL COMMENT '文件MD5值',
  `upload_user_id` bigint(20) NOT NULL COMMENT '上传用户ID',
  `usage_type` varchar(50) DEFAULT NULL COMMENT '使用类型：logo-租户Logo bg-背景图片 avatar-头像',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_upload_user_id` (`upload_user_id`),
  KEY `idx_usage_type` (`usage_type`),
  KEY `idx_md5` (`md5`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件表';

-- ========================================
-- 7. 日志管理表
-- ========================================

-- 操作日志表
CREATE TABLE `operation_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `tenant_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '租户ID，0表示系统操作',
  `user_id` bigint(20) NOT NULL COMMENT '操作用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `operation` varchar(50) NOT NULL COMMENT '操作类型',
  `method` varchar(10) NOT NULL COMMENT '请求方法',
  `url` varchar(500) NOT NULL COMMENT '请求URL',
  `params` json DEFAULT NULL COMMENT '请求参数',
  `result` text COMMENT '操作结果',
  `error_message` text COMMENT '错误信息',
  `ip` varchar(45) NOT NULL COMMENT 'IP地址',
  `user_agent` varchar(500) DEFAULT NULL COMMENT 'User-Agent',
  `duration` int(11) DEFAULT NULL COMMENT '执行时长（毫秒）',
  `status` tinyint(4) NOT NULL COMMENT '状态：1-成功 2-失败',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_operation` (`operation`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 登录日志表
CREATE TABLE `login_logs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `tenant_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '租户ID',
  `user_id` bigint(20) DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `ip` varchar(45) NOT NULL COMMENT 'IP地址',
  `user_agent` varchar(500) DEFAULT NULL COMMENT 'User-Agent',
  `location` varchar(100) DEFAULT NULL COMMENT '登录地点',
  `browser` varchar(100) DEFAULT NULL COMMENT '浏览器',
  `os` varchar(100) DEFAULT NULL COMMENT '操作系统',
  `status` tinyint(4) NOT NULL COMMENT '登录状态：1-成功 2-失败',
  `message` varchar(255) DEFAULT NULL COMMENT '提示信息',
  `login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  PRIMARY KEY (`id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_username` (`username`),
  KEY `idx_status` (`status`),
  KEY `idx_login_time` (`login_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- ========================================
-- 8. 初始化数据
-- ========================================

-- 插入默认租户（系统租户）
INSERT INTO `tenants` (`id`, `name`, `domain`, `status`, `config`, `created_at`, `updated_at`) VALUES
(0, '系统管理', 'system', 1, '{"title":"LightStack管理平台","logo":"","theme":"#1890ff"}', NOW(), NOW());

-- 插入系统管理员用户
INSERT INTO `users` (`id`, `tenant_id`, `username`, `password`, `nickname`, `email`, `status`, `is_system`, `created_at`, `updated_at`) VALUES
(1, 0, 'admin', '$2a$10$7JB720yubVSocvyT6Y62LeCULJk/3nWQ5oMF1.N3Q7GqOy7Qv8qAa', '超级管理员', 'admin@lightstack.com', 1, 1, NOW(), NOW());
-- 密码：admin123

-- 插入默认角色
INSERT INTO `roles` (`id`, `tenant_id`, `name`, `code`, `description`, `status`, `is_system`, `sort_order`, `created_at`, `updated_at`) VALUES
(1, 0, '超级管理员', 'super_admin', '拥有系统所有权限', 1, 1, 1, NOW(), NOW()),
(2, 0, '租户管理员', 'tenant_admin', '租户管理员，管理本租户下的用户和角色', 1, 1, 2, NOW(), NOW()),
(3, 0, '普通用户', 'user', '普通用户，只能查看和操作自己的信息', 1, 1, 3, NOW(), NOW());

-- 给超级管理员分配角色
INSERT INTO `user_roles` (`user_id`, `role_id`, `created_at`) VALUES (1, 1, NOW());

-- 插入基础权限
INSERT INTO `permissions` (`name`, `code`, `type`, `resource`, `action`, `description`, `created_at`, `updated_at`) VALUES
('用户管理-查看', 'user:list', 'api', '/api/users', 'GET', '查看用户列表', NOW(), NOW()),
('用户管理-创建', 'user:create', 'api', '/api/users', 'POST', '创建用户', NOW(), NOW()),
('用户管理-更新', 'user:update', 'api', '/api/users/*', 'PUT', '更新用户信息', NOW(), NOW()),
('用户管理-删除', 'user:delete', 'api', '/api/users/*', 'DELETE', '删除用户', NOW(), NOW()),
('角色管理-查看', 'role:list', 'api', '/api/roles', 'GET', '查看角色列表', NOW(), NOW()),
('角色管理-创建', 'role:create', 'api', '/api/roles', 'POST', '创建角色', NOW(), NOW()),
('角色管理-更新', 'role:update', 'api', '/api/roles/*', 'PUT', '更新角色信息', NOW(), NOW()),
('角色管理-删除', 'role:delete', 'api', '/api/roles/*', 'DELETE', '删除角色', NOW(), NOW()),
('菜单管理-查看', 'menu:list', 'api', '/api/menus', 'GET', '查看菜单列表', NOW(), NOW()),
('菜单管理-创建', 'menu:create', 'api', '/api/menus', 'POST', '创建菜单', NOW(), NOW()),
('菜单管理-更新', 'menu:update', 'api', '/api/menus/*', 'PUT', '更新菜单信息', NOW(), NOW()),
('菜单管理-删除', 'menu:delete', 'api', '/api/menus/*', 'DELETE', '删除菜单', NOW(), NOW()),
('租户管理-查看', 'tenant:list', 'api', '/api/tenants', 'GET', '查看租户列表', NOW(), NOW()),
('租户管理-创建', 'tenant:create', 'api', '/api/tenants', 'POST', '创建租户', NOW(), NOW()),
('租户管理-更新', 'tenant:update', 'api', '/api/tenants/*', 'PUT', '更新租户信息', NOW(), NOW()),
('租户管理-删除', 'tenant:delete', 'api', '/api/tenants/*', 'DELETE', '删除租户', NOW(), NOW());

-- 给超级管理员角色分配所有权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`, `created_at`)
SELECT 1, id, NOW() FROM `permissions`;

-- 插入数据字典类型
INSERT INTO `dict_types` (`name`, `type`, `description`, `status`, `created_at`, `updated_at`) VALUES
('用户状态', 'user_status', '用户状态字典', 1, NOW(), NOW()),
('租户状态', 'tenant_status', '租户状态字典', 1, NOW(), NOW()),
('角色状态', 'role_status', '角色状态字典', 1, NOW(), NOW());

-- 插入数据字典数据
INSERT INTO `dict_data` (`dict_type`, `label`, `value`, `sort_order`, `css_class`, `list_class`, `is_default`, `status`, `remark`, `created_at`, `updated_at`) VALUES
('user_status', '启用', '1', 1, 'success', 'success', 1, 1, '用户正常状态', NOW(), NOW()),
('user_status', '禁用', '2', 2, 'danger', 'danger', 0, 1, '用户禁用状态', NOW(), NOW()),
('user_status', '锁定', '3', 3, 'warning', 'warning', 0, 1, '用户锁定状态', NOW(), NOW()),
('tenant_status', '启用', '1', 1, 'success', 'success', 1, 1, '租户正常状态', NOW(), NOW()),
('tenant_status', '禁用', '2', 2, 'danger', 'danger', 0, 1, '租户禁用状态', NOW(), NOW()),
('tenant_status', '试用', '3', 3, 'warning', 'warning', 0, 1, '租户试用状态', NOW(), NOW()),
('tenant_status', '过期', '4', 4, 'info', 'info', 0, 1, '租户过期状态', NOW(), NOW()),
('role_status', '启用', '1', 1, 'success', 'success', 1, 1, '角色正常状态', NOW(), NOW()),
('role_status', '禁用', '2', 2, 'danger', 'danger', 0, 1, '角色禁用状态', NOW(), NOW());

-- 插入基础菜单
INSERT INTO `menus` (`id`, `parent_id`, `name`, `path`, `component`, `icon`, `sort_order`, `is_hidden`, `is_system`, `permission_code`, `meta`, `created_at`, `updated_at`) VALUES
(1, 0, '系统管理', '/system', 'Layout', 'system', 100, 0, 1, NULL, '{"title":"系统管理","icon":"system"}', NOW(), NOW()),
(2, 1, '用户管理', '/system/users', 'system/users/index', 'user', 1, 0, 1, 'user:list', '{"title":"用户管理","icon":"user"}', NOW(), NOW()),
(3, 1, '角色管理', '/system/roles', 'system/roles/index', 'role', 2, 0, 1, 'role:list', '{"title":"角色管理","icon":"role"}', NOW(), NOW()),
(4, 1, '菜单管理', '/system/menus', 'system/menus/index', 'menu', 3, 0, 1, 'menu:list', '{"title":"菜单管理","icon":"menu"}', NOW(), NOW()),
(5, 1, '租户管理', '/system/tenants', 'system/tenants/index', 'tenant', 4, 0, 1, 'tenant:list', '{"title":"租户管理","icon":"tenant"}', NOW(), NOW()),
(6, 1, '字典管理', '/system/dicts', 'system/dicts/index', 'dict', 5, 0, 1, NULL, '{"title":"字典管理","icon":"dict"}', NOW(), NOW()),
(7, 1, '操作日志', '/system/logs', 'system/logs/index', 'log', 6, 0, 1, NULL, '{"title":"操作日志","icon":"log"}', NOW(), NOW());

-- 给超级管理员角色分配菜单权限
INSERT INTO `role_menus` (`role_id`, `menu_id`, `created_at`)
SELECT 1, id, NOW() FROM `menus` WHERE `is_system` = 1;

SET FOREIGN_KEY_CHECKS = 1;