-- 菜单 SQL
-- 插入父级菜单
INSERT INTO `sys_menu` (`menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `menu_type`, `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
('{{.FunctionName}}', {{.ParentMenuID}}, 1, '/{{toLower .ModuleName}}/{{toLower .BusinessName}}', 'Layout', 0, 'M', 1, 1, NULL, '{{.MenuIcon}}', 'admin', NOW(), NULL, NULL, '{{.FunctionName}}菜单');

-- 获取刚插入的菜单ID (假设为 @parent_menu_id)
SET @parent_menu_id = LAST_INSERT_ID();

-- 插入列表页面菜单
INSERT INTO `sys_menu` (`menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `menu_type`, `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
('{{.FunctionName}}列表', @parent_menu_id, 1, '{{.MenuURL}}', '{{.ModuleName}}/{{toLower .BusinessName}}/index', 0, 'C', 1, 1, '{{generatePermission .ModuleName .BusinessName "list"}}', 'List', 'admin', NOW(), NULL, NULL, '{{.FunctionName}}列表页面');

-- 获取列表菜单ID
SET @list_menu_id = LAST_INSERT_ID();

-- 插入查询按钮权限
INSERT INTO `sys_menu` (`menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `menu_type`, `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
('{{.FunctionName}}查询', @list_menu_id, 1, '', NULL, 0, 'F', 1, 1, '{{generatePermission .ModuleName .BusinessName "query"}}', '', 'admin', NOW(), NULL, NULL, '{{.FunctionName}}查询按钮');

-- 插入新增按钮权限
INSERT INTO `sys_menu` (`menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `menu_type`, `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
('{{.FunctionName}}新增', @list_menu_id, 2, '', NULL, 0, 'F', 1, 1, '{{generatePermission .ModuleName .BusinessName "add"}}', '', 'admin', NOW(), NULL, NULL, '{{.FunctionName}}新增按钮');

-- 插入编辑按钮权限
INSERT INTO `sys_menu` (`menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `menu_type`, `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
('{{.FunctionName}}编辑', @list_menu_id, 3, '', NULL, 0, 'F', 1, 1, '{{generatePermission .ModuleName .BusinessName "edit"}}', '', 'admin', NOW(), NULL, NULL, '{{.FunctionName}}编辑按钮');

-- 插入删除按钮权限
INSERT INTO `sys_menu` (`menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `menu_type`, `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
('{{.FunctionName}}删除', @list_menu_id, 4, '', NULL, 0, 'F', 1, 1, '{{generatePermission .ModuleName .BusinessName "remove"}}', '', 'admin', NOW(), NULL, NULL, '{{.FunctionName}}删除按钮');

-- 插入导出按钮权限
INSERT INTO `sys_menu` (`menu_name`, `parent_id`, `order_num`, `path`, `component`, `is_frame`, `menu_type`, `visible`, `status`, `perms`, `icon`, `create_by`, `create_time`, `update_by`, `update_time`, `remark`)
VALUES
('{{.FunctionName}}导出', @list_menu_id, 5, '', NULL, 0, 'F', 1, 1, '{{generatePermission .ModuleName .BusinessName "export"}}', '', 'admin', NOW(), NULL, NULL, '{{.FunctionName}}导出按钮');

-- 为管理员角色分配菜单权限 (假设管理员角色ID为1)
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT 1, id FROM `sys_menu`
WHERE `perms` LIKE '{{toLower .ModuleName}}:{{toLower .BusinessName}}:%'
   OR `path` = '/{{toLower .ModuleName}}/{{toLower .BusinessName}}'
   OR `parent_id` IN (
       SELECT id FROM (SELECT id FROM `sys_menu` WHERE `path` = '/{{toLower .ModuleName}}/{{toLower .BusinessName}}') AS tmp
   );

-- 输出提示信息
SELECT '{{.FunctionName}}菜单和权限已成功创建' AS message;