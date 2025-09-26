-- {{.FunctionName}}菜单配置
INSERT INTO `menus` (`parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `sort_order`, `status`, `tenant_id`, `created_at`, `updated_at`)
VALUES ({{.ParentMenuID}}, '{{.MenuName}}', '{{toSnakeCase .BusinessName}}', 'menu', '{{.MenuURL}}', '{{.ModuleName}}/{{toKebabCase .BusinessName}}', '{{.MenuIcon}}', 100, 1, 0, NOW(), NOW());

-- 获取插入的菜单ID
SET @menu_id = LAST_INSERT_ID();

{{- range .Permissions }}
-- 插入权限: {{.}}
INSERT INTO `permissions` (`name`, `code`, `type`, `resource`, `action`, `menu_id`, `status`, `tenant_id`, `created_at`, `updated_at`)
VALUES ('{{$.FunctionName}}{{if hasSuffix . ":list"}}列表{{else if hasSuffix . ":add"}}新增{{else if hasSuffix . ":edit"}}编辑{{else if hasSuffix . ":delete"}}删除{{else if hasSuffix . ":view"}}查看{{end}}', '{{.}}', 'operation', '{{$.ModuleName}}:{{$.BusinessName}}', '{{if hasSuffix . ":list"}}list{{else if hasSuffix . ":add"}}add{{else if hasSuffix . ":edit"}}edit{{else if hasSuffix . ":delete"}}delete{{else if hasSuffix . ":view"}}view{{end}}', @menu_id, 1, 0, NOW(), NOW());
{{- end }}

-- 为超级管理员角色分配菜单权限
INSERT INTO `role_menus` (`role_id`, `menu_id`)
SELECT r.id, @menu_id FROM `roles` r WHERE r.code = 'super_admin' AND r.tenant_id = 0;