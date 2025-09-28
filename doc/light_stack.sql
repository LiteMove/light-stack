/*
 Navicat Premium Dump SQL

 Source Server         : light_stack
 Source Server Type    : MySQL
 Source Server Version : 50744 (5.7.44-log)
 Source Host           : matuto_db:3306
 Source Schema         : light_stack

 Target Server Type    : MySQL
 Target Server Version : 50744 (5.7.44-log)
 File Encoding         : 65001

 Date: 28/09/2025 16:20:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for dict_data
-- ----------------------------
DROP TABLE IF EXISTS `dict_data`;
CREATE TABLE `dict_data`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字典数据ID',
  `dict_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典类型',
  `label` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典标签',
  `value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典键值',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序号',
  `css_class` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'CSS类名',
  `list_class` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '列表样式',
  `is_default` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否默认：0-否 1-是',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1-启用 2-禁用',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_type_value`(`dict_type`, `value`) USING BTREE,
  INDEX `idx_dict_type`(`dict_type`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE,
  INDEX `idx_sort_order`(`sort_order`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '字典数据表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of dict_data
-- ----------------------------
INSERT INTO `dict_data` VALUES (1, 'user_status', '启用', '1', 1, 'success', 'success', 1, 1, '用户正常状态', '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_data` VALUES (2, 'user_status', '禁用', '2', 2, 'danger', 'danger', 0, 1, '用户禁用状态', '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_data` VALUES (3, 'user_status', '锁定', '3', 3, 'warning', 'warning', 0, 1, '用户锁定状态', '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_data` VALUES (4, 'tenant_status', '启用', '1', 1, 'success', 'success', 1, 1, '租户正常状态', '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_data` VALUES (5, 'tenant_status', '禁用', '2', 2, 'danger', 'danger', 0, 1, '租户禁用状态', '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_data` VALUES (6, 'tenant_status', '试用', '3', 3, 'warning', 'warning', 0, 1, '租户试用状态', '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_data` VALUES (7, 'tenant_status', '过期', '4', 4, 'info', 'info', 0, 1, '租户过期状态', '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_data` VALUES (8, 'role_status', '启用', '1', 1, 'success', 'success', 1, 1, '角色正常状态', '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_data` VALUES (9, 'role_status', '禁用', '2', 2, 'danger', 'danger', 0, 1, '角色禁用状态', '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);

-- ----------------------------
-- Table structure for dict_types
-- ----------------------------
DROP TABLE IF EXISTS `dict_types`;
CREATE TABLE `dict_types`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字典类型ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典名称',
  `type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典类型',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '描述',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态：1-启用 2-禁用',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_type`(`type`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '字典类型表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of dict_types
-- ----------------------------
INSERT INTO `dict_types` VALUES (1, '用户状态', 'user_status', '用户状态字典', 1, '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_types` VALUES (2, '租户状态', 'tenant_status', '租户状态字典', 1, '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `dict_types` VALUES (3, '角色状态', 'role_status', '角色状态字典', 1, '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);

-- ----------------------------
-- Table structure for files
-- ----------------------------
DROP TABLE IF EXISTS `files`;
CREATE TABLE `files`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  `tenant_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '租户ID，0表示系统文件',
  `original_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '原始文件名',
  `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '存储文件名',
  `file_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件路径',
  `file_size` bigint(20) NOT NULL COMMENT '文件大小（字节）',
  `file_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型',
  `mime_type` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'MIME类型',
  `md5` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件MD5值',
  `upload_user_id` bigint(20) NOT NULL COMMENT '上传用户ID',
  `usage_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '使用类型：logo-租户Logo bg-背景图片 avatar-头像',
  `storage_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '存储类型',
  `is_public` tinyint(1) NULL DEFAULT 0 COMMENT '是否公开访问，0：私有，1：公有',
  `access_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '完整访问URL',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tenant_id`(`tenant_id`) USING BTREE,
  INDEX `idx_upload_user_id`(`upload_user_id`) USING BTREE,
  INDEX `idx_usage_type`(`usage_type`) USING BTREE,
  INDEX `idx_md5`(`md5`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 29 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文件表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of files
-- ----------------------------
INSERT INTO `files` VALUES (24, 1, '上线我的 2.0.png', '1758706009945887000.png', 'private/tenant_1/2025/09/24/1758706009945887000.png', 1666407, 'png', 'image/png', '0abf3fe7172f56a49127ba0ea6e889df', 1, 'avatar', 'local', 0, 'http://127.0.0.1:8080/api/static/private/tenant_1/2025/09/24/1758706009945887000.png', '2025-09-24 17:26:50', '2025-09-24 17:51:20', '2025-09-24 17:51:21');
INSERT INTO `files` VALUES (25, 1, '个人公众号起名.png', '1758707254453896600.png', 'public/tenant_1/2025/09/24/1758707254453896600.png', 1351175, 'png', 'image/png', '65bd9ff02d9567752b7be06a2cb2b791', 1, 'avatar', 'local', 1, 'http://127.0.0.1:8080/api/static/public/tenant_1/2025/09/24/1758707254453896600.png', '2025-09-24 17:47:34', '2025-09-24 17:47:34', NULL);
INSERT INTO `files` VALUES (26, 1, '上线我的 2.0.png', '1758707486724848500.png', 'public/tenant_1/2025/09/24/1758707486724848500.png', 1666407, 'png', 'image/png', '0abf3fe7172f56a49127ba0ea6e889df', 1, 'avatar', 'local', 1, 'http://127.0.0.1:8080/api/static/public/tenant_1/2025/09/24/1758707486724848500.png', '2025-09-24 17:51:27', '2025-09-24 17:51:27', NULL);
INSERT INTO `files` VALUES (27, 1, 'logo.png', '1758708129776068200.png', 'public/tenant_1/2025/09/24/1758708129776068200.png', 9854, 'png', 'image/png', '28048fc01baf30a1d6365d20306a7b3e', 1, 'system-logo', 'local', 1, 'http://127.0.0.1:8080/api/static/public/tenant_1/2025/09/24/1758708129776068200.png', '2025-09-24 18:02:10', '2025-09-24 18:02:10', NULL);
INSERT INTO `files` VALUES (28, 1, '新建 文本文档.txt', '1759026609826753200.txt', 'private/tenant_1/2025/09/28/1759026609826753200.txt', 21, 'txt', 'text/plain', '5b67bd58f9aef918d42d13970ea88217', 1, 'document', 'local', 0, 'http://127.0.0.1:8080/api/static/private/tenant_1/2025/09/28/1759026609826753200.txt', '2025-09-28 10:30:10', '2025-09-28 10:30:10', NULL);

-- ----------------------------
-- Table structure for gen_histories
-- ----------------------------
DROP TABLE IF EXISTS `gen_histories`;
CREATE TABLE `gen_histories`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '历史ID',
  `table_config_id` bigint(20) NOT NULL COMMENT '表配置ID',
  `table_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '表名称',
  `business_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '业务名称',
  `generate_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '生成类型(all全部、backend后端、frontend前端)',
  `file_count` int(11) NULL DEFAULT 0 COMMENT '生成文件数量',
  `file_size` bigint(20) NULL DEFAULT 0 COMMENT '文件大小(字节)',
  `download_count` int(11) NULL DEFAULT 0 COMMENT '下载次数',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'success' COMMENT '生成状态(success成功、failed失败、processing处理中)',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '错误信息',
  `file_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '生成文件路径',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '备注',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `created_by` bigint(20) NULL DEFAULT NULL COMMENT '创建人',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_table_config_id`(`table_config_id`) USING BTREE,
  INDEX `idx_table_name`(`table_name`) USING BTREE,
  INDEX `idx_created_at`(`created_at`) USING BTREE,
  CONSTRAINT `fk_gen_history_config` FOREIGN KEY (`table_config_id`) REFERENCES `gen_table_configs` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '代码生成历史记录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of gen_histories
-- ----------------------------
INSERT INTO `gen_histories` VALUES (4, 2, 'login_logs', 'logs', 'all', 11, 63903, 0, 'success', '', 'generated\\logs_20250926233026_779887700.zip', '', '2025-09-26 23:30:27', NULL);

-- ----------------------------
-- Table structure for gen_table_columns
-- ----------------------------
DROP TABLE IF EXISTS `gen_table_columns`;
CREATE TABLE `gen_table_columns`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字段ID',
  `table_config_id` bigint(20) NOT NULL COMMENT '表配置ID',
  `column_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字段名称',
  `column_comment` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '字段描述',
  `column_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '字段类型',
  `go_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Go类型',
  `go_field` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'Go字段名',
  `is_pk` tinyint(1) NULL DEFAULT 0 COMMENT '是否主键',
  `is_increment` tinyint(1) NULL DEFAULT 0 COMMENT '是否自增',
  `is_required` tinyint(1) NULL DEFAULT 0 COMMENT '是否必填',
  `is_insert` tinyint(1) NULL DEFAULT 1 COMMENT '是否为插入字段',
  `is_edit` tinyint(1) NULL DEFAULT 1 COMMENT '是否为编辑字段',
  `is_list` tinyint(1) NULL DEFAULT 1 COMMENT '是否列表字段',
  `is_query` tinyint(1) NULL DEFAULT 0 COMMENT '是否查询字段',
  `query_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'EQ' COMMENT '查询方式(EQ等于、NE不等于、GT大于、LT小于、LIKE模糊、BETWEEN范围)',
  `html_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'input' COMMENT '显示类型(input、textarea、select、radio、checkbox、datetime)',
  `dict_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '字典类型',
  `sort` int(11) NULL DEFAULT 0 COMMENT '排序',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_table_config_id`(`table_config_id`) USING BTREE,
  INDEX `idx_column_name`(`column_name`) USING BTREE,
  CONSTRAINT `fk_gen_table_columns_config` FOREIGN KEY (`table_config_id`) REFERENCES `gen_table_configs` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 349 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '代码生成字段配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of gen_table_columns
-- ----------------------------
INSERT INTO `gen_table_columns` VALUES (300, 1, 'id', '字典数据ID', 'bigint(20)', 'int64', 'Id', 1, 1, 1, 1, 1, 1, 0, 'EQ', 'input', '', 1, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (301, 1, 'dict_type', '字典类型', 'varchar(100)', 'string', 'DictType', 0, 0, 1, 1, 1, 1, 1, 'LIKE', 'select', 'sys_type', 2, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (302, 1, 'label', '字典标签', 'varchar(100)', 'string', 'Label', 0, 0, 1, 1, 1, 1, 0, 'LIKE', 'input', '', 3, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (303, 1, 'value', '字典键值', 'varchar(100)', 'string', 'Value', 0, 0, 1, 1, 1, 1, 0, 'LIKE', 'input', '', 4, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (304, 1, 'sort_order', '排序号', 'int(11)', 'int', 'SortOrder', 0, 0, 0, 1, 1, 1, 0, 'EQ', 'input', '', 5, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (305, 1, 'css_class', 'CSS类名', 'varchar(100)', 'string', 'CssClass', 0, 0, 0, 1, 1, 1, 0, 'LIKE', 'input', '', 6, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (306, 1, 'list_class', '列表样式', 'varchar(100)', 'string', 'ListClass', 0, 0, 0, 1, 1, 1, 0, 'LIKE', 'input', '', 7, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (307, 1, 'is_default', '是否默认：0-否 1-是', 'tinyint(1)', 'bool', 'IsDefault', 0, 0, 0, 1, 1, 1, 0, 'EQ', 'radio', '', 8, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (308, 1, 'status', '状态：1-启用 2-禁用', 'tinyint(4)', 'int8', 'Status', 0, 0, 0, 1, 1, 1, 1, 'EQ', 'select', 'sys_status', 9, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (309, 1, 'remark', '备注', 'varchar(255)', 'string', 'Remark', 0, 0, 0, 1, 1, 1, 0, 'LIKE', 'textarea', '', 10, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (310, 1, 'created_at', '创建时间', 'datetime', 'time.Time', 'CreatedAt', 0, 0, 0, 1, 1, 1, 0, 'BETWEEN', 'datetime', '', 11, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (311, 1, 'updated_at', '更新时间', 'datetime', 'time.Time', 'UpdatedAt', 0, 0, 0, 1, 1, 1, 0, 'BETWEEN', 'datetime', '', 12, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (312, 1, 'deleted_at', '删除时间', 'datetime', 'time.Time', 'DeletedAt', 0, 0, 0, 1, 1, 1, 0, 'BETWEEN', 'datetime', '', 13, '2025-09-26 22:40:04', '2025-09-26 22:40:04');
INSERT INTO `gen_table_columns` VALUES (337, 2, 'id', '日志ID', 'bigint(20)', 'int64', 'Id', 1, 1, 1, 1, 1, 1, 0, 'EQ', 'input', '', 1, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (338, 2, 'tenant_id', '租户ID', 'bigint(20)', 'int64', 'TenantId', 0, 0, 0, 1, 1, 1, 0, 'EQ', 'input', '', 2, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (339, 2, 'user_id', '用户ID', 'bigint(20)', 'int64', 'UserId', 0, 0, 0, 1, 1, 1, 0, 'EQ', 'input', '', 3, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (340, 2, 'username', '用户名', 'varchar(50)', 'string', 'Username', 0, 0, 1, 1, 1, 1, 1, 'LIKE', 'input', '', 4, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (341, 2, 'ip', 'IP地址', 'varchar(45)', 'string', 'Ip', 0, 0, 1, 1, 1, 1, 0, 'LIKE', 'input', '', 5, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (342, 2, 'user_agent', 'User-Agent', 'varchar(500)', 'string', 'UserAgent', 0, 0, 0, 1, 1, 1, 0, 'LIKE', 'input', '', 6, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (343, 2, 'location', '登录地点', 'varchar(100)', 'string', 'Location', 0, 0, 0, 1, 1, 1, 0, 'LIKE', 'input', '', 7, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (344, 2, 'browser', '浏览器', 'varchar(100)', 'string', 'Browser', 0, 0, 0, 1, 1, 1, 0, 'LIKE', 'input', '', 8, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (345, 2, 'os', '操作系统', 'varchar(100)', 'string', 'Os', 0, 0, 0, 1, 1, 1, 0, 'LIKE', 'input', '', 9, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (346, 2, 'status', '登录状态：1-成功 2-失败', 'tinyint(4)', 'int8', 'Status', 0, 0, 1, 1, 1, 1, 1, 'EQ', 'select', 'sys_status', 10, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (347, 2, 'message', '提示信息', 'varchar(255)', 'string', 'Message', 0, 0, 0, 1, 1, 1, 0, 'LIKE', 'input', '', 11, '2025-09-26 23:30:24', '2025-09-26 23:30:24');
INSERT INTO `gen_table_columns` VALUES (348, 2, 'login_time', '登录时间', 'datetime', 'time.Time', 'LoginTime', 0, 0, 0, 1, 1, 1, 0, 'BETWEEN', 'datetime', '', 12, '2025-09-26 23:30:24', '2025-09-26 23:30:24');

-- ----------------------------
-- Table structure for gen_table_configs
-- ----------------------------
DROP TABLE IF EXISTS `gen_table_configs`;
CREATE TABLE `gen_table_configs`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `table_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '表名称',
  `table_comment` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '表描述',
  `business_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '业务名称',
  `module_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模块名称',
  `function_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '功能名称',
  `class_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类名',
  `package_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '包名',
  `author` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'system' COMMENT '作者',
  `parent_menu_id` bigint(20) NULL DEFAULT NULL COMMENT '父级菜单ID',
  `menu_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '菜单名称',
  `menu_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '菜单URL',
  `menu_icon` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '菜单图标',
  `permissions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '权限字符串(JSON数组)',
  `options` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '其他配置选项(JSON)',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '备注',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `created_by` bigint(20) NULL DEFAULT NULL COMMENT '创建人',
  `updated_by` bigint(20) NULL DEFAULT NULL COMMENT '更新人',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_table_name`(`table_name`) USING BTREE,
  INDEX `idx_business_name`(`business_name`) USING BTREE,
  INDEX `idx_module_name`(`module_name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '代码生成表配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of gen_table_configs
-- ----------------------------
INSERT INTO `gen_table_configs` VALUES (1, 'dict_data', '字典数据表', 'data', 'system', 'data管理', 'Data', 'system', 'system', NULL, 'data管理', '/data', 'Grid', '[\"system:data:list\",\"system:data:add\",\"system:data:edit\",\"system:data:delete\",\"system:data:view\"]', '{\"genPath\":\"\",\"genType\":\"\",\"tplType\":\"\",\"treeCode\":\"\",\"treeParent\":\"\",\"treeName\":\"\"}', '', '2025-09-26 11:41:37', '2025-09-26 22:40:04', NULL, NULL);
INSERT INTO `gen_table_configs` VALUES (2, 'login_logs', '登录日志表', 'logs', 'system', 'logs管理', 'Logs', 'system', 'system', 1, 'logs管理', '/logs', 'Grid', '[\"system:logs:list\",\"system:logs:add\",\"system:logs:edit\",\"system:logs:delete\",\"system:logs:view\"]', '{\"genPath\":\"\",\"genType\":\"\",\"tplType\":\"\",\"treeCode\":\"\",\"treeParent\":\"\",\"treeName\":\"\"}', '', '2025-09-26 23:24:26', '2025-09-26 23:30:24', NULL, NULL);

-- ----------------------------
-- Table structure for login_logs
-- ----------------------------
DROP TABLE IF EXISTS `login_logs`;
CREATE TABLE `login_logs`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `tenant_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '租户ID',
  `user_id` bigint(20) NULL DEFAULT NULL COMMENT '用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `user_agent` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'User-Agent',
  `location` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '登录地点',
  `browser` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '浏览器',
  `os` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作系统',
  `status` tinyint(4) NOT NULL COMMENT '登录状态：1-成功 2-失败',
  `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '提示信息',
  `login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tenant_id`(`tenant_id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_username`(`username`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE,
  INDEX `idx_login_time`(`login_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '登录日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of login_logs
-- ----------------------------

-- ----------------------------
-- Table structure for menus
-- ----------------------------
DROP TABLE IF EXISTS `menus`;
CREATE TABLE `menus`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '菜单/权限ID',
  `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '父菜单ID，0表示顶级菜单',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单/权限名称',
  `code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '权限编码（唯一标识）',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'menu' COMMENT '类型：directory-目录 menu-菜单 permission-权限',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '路由路径（菜单类型有效）',
  `status` tinyint(4) NULL DEFAULT 1 COMMENT '状态',
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '组件路径（菜单类型有效）',
  `icon` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '图标（目录/菜单类型有效）',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序号',
  `is_hidden` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否隐藏：0-显示 1-隐藏',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_code`(`code`) USING BTREE,
  INDEX `idx_parent_id`(`parent_id`) USING BTREE,
  INDEX `idx_type`(`type`) USING BTREE,
  INDEX `idx_sort_order`(`sort_order`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 30 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单权限表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of menus
-- ----------------------------
INSERT INTO `menus` VALUES (1, 0, '系统管理', 'system:management', 'directory', '/system', 1, 'Layout', 'system', 100, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:12', NULL);
INSERT INTO `menus` VALUES (2, 1, '用户管理', 'system:user:menu', 'menu', '/system/users', 1, 'system/users/index', 'user', 0, 0, '0000-00-00 00:00:00', '2025-09-24 22:55:35', NULL);
INSERT INTO `menus` VALUES (3, 27, '角色管理', 'system:role:menu', 'menu', '/system/roles', 1, 'system/roles/index', 'role', 0, 0, '0000-00-00 00:00:00', '2025-09-24 22:55:00', NULL);
INSERT INTO `menus` VALUES (4, 27, '菜单权限', 'system:menu:menu', 'menu', '/system/menus', 1, 'system/menus/index', 'menu', 0, 0, '0000-00-00 00:00:00', '2025-09-24 22:55:11', NULL);
INSERT INTO `menus` VALUES (5, 27, '租户管理', 'system:tenant:menu', 'menu', '/system/tenants', 1, 'system/tenants/index', 'avatar', 0, 0, '0000-00-00 00:00:00', '2025-09-24 22:55:17', NULL);
INSERT INTO `menus` VALUES (6, 1, '字典管理', 'system:dict:menu', 'menu', '/system/dicts', 1, 'system/dicts/index', 'Collection', 0, 0, '0000-00-00 00:00:00', '2025-09-25 15:23:37', NULL);
INSERT INTO `menus` VALUES (7, 1, '操作日志', 'system:log:menu', 'menu', '/system/logs', 1, 'system/logs/index', 'log', 6, 1, '2025-09-18 20:21:12', '2025-09-20 18:11:03', NULL);
INSERT INTO `menus` VALUES (8, 2, '用户管理-列表', 'system:user:list', 'permission', '', 1, '', '', 0, 0, '0000-00-00 00:00:00', '2025-09-24 21:31:56', NULL);
INSERT INTO `menus` VALUES (9, 2, '用户管理-创建', 'system:user:create', 'permission', NULL, 1, NULL, NULL, 2, 0, '2025-09-18 20:21:12', '2025-09-24 21:28:54', NULL);
INSERT INTO `menus` VALUES (10, 2, '用户管理-更新', 'system:user:update', 'permission', NULL, 1, NULL, NULL, 3, 0, '2025-09-18 20:21:12', '2025-09-24 21:28:55', NULL);
INSERT INTO `menus` VALUES (11, 2, '用户管理-删除', 'system:user:delete', 'permission', NULL, 1, NULL, NULL, 4, 0, '2025-09-18 20:21:12', '2025-09-24 21:28:57', NULL);
INSERT INTO `menus` VALUES (12, 3, '角色管理-查看', 'role:list', 'permission', NULL, 1, NULL, NULL, 1, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:12', NULL);
INSERT INTO `menus` VALUES (13, 3, '角色管理-创建', 'role:create', 'permission', NULL, 1, NULL, NULL, 2, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:12', NULL);
INSERT INTO `menus` VALUES (14, 3, '角色管理-更新', 'role:update', 'permission', NULL, 1, NULL, NULL, 3, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:12', NULL);
INSERT INTO `menus` VALUES (15, 3, '角色管理-删除', 'role:delete', 'permission', NULL, 1, NULL, NULL, 4, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:12', NULL);
INSERT INTO `menus` VALUES (16, 4, '菜单权限-查看', 'menu:list', 'permission', NULL, 1, NULL, NULL, 1, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:12', NULL);
INSERT INTO `menus` VALUES (17, 4, '菜单权限-创建', 'menu:create', 'permission', NULL, 1, NULL, NULL, 2, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:12', NULL);
INSERT INTO `menus` VALUES (18, 4, '菜单权限-更新', 'menu:update', 'permission', NULL, 1, NULL, NULL, 3, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:12', NULL);
INSERT INTO `menus` VALUES (19, 4, '菜单权限-删除', 'menu:delete', 'permission', NULL, 1, NULL, NULL, 4, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:13', NULL);
INSERT INTO `menus` VALUES (20, 5, '租户管理-查看', 'tenant:list', 'permission', NULL, 1, NULL, NULL, 1, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:13', NULL);
INSERT INTO `menus` VALUES (21, 5, '租户管理-创建', 'tenant:create', 'permission', NULL, 1, NULL, NULL, 2, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:13', NULL);
INSERT INTO `menus` VALUES (22, 5, '租户管理-更新', 'tenant:update', 'permission', NULL, 1, NULL, NULL, 3, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:13', NULL);
INSERT INTO `menus` VALUES (23, 5, '租户管理-删除', 'tenant:delete', 'permission', NULL, 1, NULL, NULL, 4, 0, '2025-09-18 20:21:12', '2025-09-18 22:16:13', NULL);
INSERT INTO `menus` VALUES (24, 2, '用户角色-分配', 'system:user:role:assign', 'permission', '', 1, '', '', 0, 0, '0000-00-00 00:00:00', '2025-09-24 21:32:59', NULL);
INSERT INTO `menus` VALUES (25, 2, '用户角色-查看', 'system:user:role:list', 'permission', '', 1, '', '', 0, 0, '0000-00-00 00:00:00', '2025-09-24 21:33:05', NULL);
INSERT INTO `menus` VALUES (27, 0, '权限管理', '__', 'directory', '/perm', 1, '', 'UserFilled', 0, 0, '0000-00-00 00:00:00', '2025-09-24 23:15:34', NULL);
INSERT INTO `menus` VALUES (28, 1, '文件管理', 'file_management', 'menu', '/files', 1, 'system/files/index', 'FolderOpened', 0, 0, '0000-00-00 00:00:00', '2025-09-24 22:55:29', NULL);
INSERT INTO `menus` VALUES (29, 2, '用户管理-重置密码', 'system:user:reset', 'permission', '', 1, '', '', 0, 0, '0000-00-00 00:00:00', '2025-09-25 11:39:23', NULL);

-- ----------------------------
-- Table structure for operation_logs
-- ----------------------------
DROP TABLE IF EXISTS `operation_logs`;
CREATE TABLE `operation_logs`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `tenant_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '租户ID，0表示系统操作',
  `user_id` bigint(20) NOT NULL COMMENT '操作用户ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `operation` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作类型',
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求方法',
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '请求URL',
  `params` json NULL COMMENT '请求参数',
  `result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '操作结果',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '错误信息',
  `ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'IP地址',
  `user_agent` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT 'User-Agent',
  `duration` int(11) NULL DEFAULT NULL COMMENT '执行时长（毫秒）',
  `status` tinyint(4) NOT NULL COMMENT '状态：1-成功 2-失败',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_tenant_id`(`tenant_id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_operation`(`operation`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE,
  INDEX `idx_created_at`(`created_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '操作日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of operation_logs
-- ----------------------------

-- ----------------------------
-- Table structure for role_menus
-- ----------------------------
DROP TABLE IF EXISTS `role_menus`;
CREATE TABLE `role_menus`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `menu_id` bigint(20) NOT NULL COMMENT '菜单/权限ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_role_menu`(`role_id`, `menu_id`) USING BTREE,
  INDEX `idx_role_id`(`role_id`) USING BTREE,
  INDEX `idx_menu_id`(`menu_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 162 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色菜单权限关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of role_menus
-- ----------------------------
INSERT INTO `role_menus` VALUES (77, 2, 2, '2025-09-22 14:19:21');
INSERT INTO `role_menus` VALUES (78, 2, 8, '2025-09-22 14:19:21');
INSERT INTO `role_menus` VALUES (79, 2, 9, '2025-09-22 14:19:21');
INSERT INTO `role_menus` VALUES (80, 2, 10, '2025-09-22 14:19:21');
INSERT INTO `role_menus` VALUES (81, 2, 11, '2025-09-22 14:19:21');
INSERT INTO `role_menus` VALUES (82, 2, 24, '2025-09-22 14:19:21');
INSERT INTO `role_menus` VALUES (83, 2, 25, '2025-09-22 14:19:21');
INSERT INTO `role_menus` VALUES (130, 3, 8, '2025-09-24 22:19:25');
INSERT INTO `role_menus` VALUES (131, 3, 2, '2025-09-24 22:19:25');
INSERT INTO `role_menus` VALUES (132, 3, 12, '2025-09-24 22:19:25');
INSERT INTO `role_menus` VALUES (133, 3, 3, '2025-09-24 22:19:25');
INSERT INTO `role_menus` VALUES (134, 3, 1, '2025-09-24 22:19:25');
INSERT INTO `role_menus` VALUES (135, 3, 13, '2025-09-24 22:19:25');
INSERT INTO `role_menus` VALUES (136, 1, 8, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (137, 1, 2, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (138, 1, 24, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (139, 1, 25, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (140, 1, 28, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (141, 1, 12, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (142, 1, 3, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (143, 1, 1, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (144, 1, 16, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (145, 1, 4, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (146, 1, 20, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (147, 1, 5, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (148, 1, 9, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (149, 1, 13, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (150, 1, 17, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (151, 1, 21, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (152, 1, 10, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (153, 1, 14, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (154, 1, 18, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (155, 1, 22, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (156, 1, 11, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (157, 1, 15, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (158, 1, 19, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (159, 1, 23, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (160, 1, 6, '2025-09-24 22:41:19');
INSERT INTO `role_menus` VALUES (161, 1, 7, '2025-09-24 22:41:19');

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色编码',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '角色描述',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '角色状态：1-启用 2-禁用',
  `is_system` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统角色：0-否 1-是',
  `sort_order` int(11) NOT NULL DEFAULT 0 COMMENT '排序号',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_code`(`code`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE,
  INDEX `idx_is_system`(`is_system`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES (1, '超级管理员', 'super_admin', '拥有系统所有权限，可管理所有租户、用户、角色和菜单权限', 1, 1, 1, '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `roles` VALUES (2, '租户管理员', 'tenant_admin', '租户管理员，可管理本租户下的用户（创建、修改、删除），可以给用户分配非系统角色', 1, 1, 2, '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);
INSERT INTO `roles` VALUES (3, '普通用户', 'user', '普通用户，只能查看和操作自己的信息', 1, 1, 3, '2025-09-18 20:21:12', '2025-09-18 20:21:12', NULL);

-- ----------------------------
-- Table structure for tenants
-- ----------------------------
DROP TABLE IF EXISTS `tenants`;
CREATE TABLE `tenants`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '租户ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租户名称',
  `domain` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '租户域名',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '租户状态：1-启用 2-禁用 3-试用 4-过期',
  `expired_at` datetime NULL DEFAULT NULL COMMENT '过期时间',
  `config` json NULL COMMENT '租户配置信息（Logo、主题色等）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_domain`(`domain`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE,
  INDEX `idx_expired_at`(`expired_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '租户信息表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tenants
-- ----------------------------
INSERT INTO `tenants` VALUES (1, 'LightStack', 'system', 1, NULL, '{\"logo\": \"http://127.0.0.1:8080/api/static/public/tenant_1/2025/09/24/1758708129776068200.png\", \"copyright\": \"\", \"systemName\": \"轻栈管理平台\", \"description\": \"\", \"fileStorage\": {\"type\": \"local\", \"maxFileSize\": 52428800, \"ossProvider\": \"aliyun\", \"allowedTypes\": [\".jpg\", \".jpeg\", \".gif\", \".pdf\", \".doc\", \".docx\", \".xlsx\", \".txt\", \".png\", \".xls\"], \"defaultPublic\": false, \"localAccessDomain\": \"http://127.0.0.1:8080\"}}', '2025-09-18 20:21:12', '2025-09-25 12:31:42', NULL);
INSERT INTO `tenants` VALUES (2, 'Test', 'test.light-stack.com', 1, '2025-09-24 10:00:00', '{\"logo\": \"\", \"copyright\": \"\", \"systemName\": \"\", \"description\": \"\", \"fileStorage\": {\"type\": \"local\", \"maxFileSize\": 52428800, \"ossProvider\": \"aliyun\", \"allowedTypes\": [\".jpg\", \".pdf\", \".doc\", \".docx\", \".xlsx\", \".txt\", \".png\", \".xls\", \".jpeg\", \".gif\"], \"defaultPublic\": true, \"localAccessDomain\": \"http://127.0.0.1:8080\"}}', '2025-09-19 17:46:27', '2025-09-24 15:10:43', NULL);
INSERT INTO `tenants` VALUES (3, 'Matuto', 'matuto.com', 1, '2025-10-03 15:59:59', '{\"logo\": \"\", \"copyright\": \"\", \"systemName\": \"\", \"description\": \"\", \"fileStorage\": {\"type\": \"local\", \"maxFileSize\": 52428800, \"ossProvider\": \"aliyun\", \"allowedTypes\": [\".jpg\", \".gif\", \".pdf\", \".doc\", \".xlsx\", \".txt\", \".docx\", \".jpeg\", \".png\", \".xls\"], \"defaultPublic\": false, \"localAccessDomain\": \"http://127.0.0.1:8080\"}}', '2025-09-21 08:43:27', '2025-09-24 15:21:21', NULL);

-- ----------------------------
-- Table structure for user_roles
-- ----------------------------
DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` bigint(20) NOT NULL COMMENT '用户ID',
  `role_id` bigint(20) NOT NULL COMMENT '角色ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_user_role`(`user_id`, `role_id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_role_id`(`role_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 27 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户角色关联表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_roles
-- ----------------------------
INSERT INTO `user_roles` VALUES (7, 10, 2, '2025-09-20 10:49:38');
INSERT INTO `user_roles` VALUES (21, 1, 1, '2025-09-22 09:55:23');
INSERT INTO `user_roles` VALUES (22, 1, 2, '2025-09-22 11:36:00');
INSERT INTO `user_roles` VALUES (26, 12, 3, '2025-09-28 10:29:52');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `tenant_id` bigint(20) NOT NULL DEFAULT 1 COMMENT '租户ID，0表示系统管理员',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码（bcrypt加密）',
  `nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '邮箱',
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT '' COMMENT '手机号',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '头像地址',
  `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '用户状态：1-启用 2-禁用 3-锁定',
  `is_system` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否系统用户：0-否 1-是',
  `last_login_at` datetime NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '最后登录IP',
  `login_failures` int(11) NOT NULL DEFAULT 0 COMMENT '连续登录失败次数',
  `locked_until` datetime NULL DEFAULT NULL COMMENT '锁定截止时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_tenant_username`(`tenant_id`, `username`) USING BTREE,
  UNIQUE INDEX `uk_tenant_email`(`tenant_id`, `email`) USING BTREE,
  UNIQUE INDEX `uk_tenant_phone`(`tenant_id`, `phone`) USING BTREE,
  INDEX `idx_tenant_id`(`tenant_id`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE,
  INDEX `idx_is_system`(`is_system`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 1, 'admin', '$2a$10$Ck5B5o1Md2O7K.tER/Hug.5phi4hRazcY04WdF46ykrmV3KdPo.3G', '超级管理员', 'admin@lightstack.com', '15688888888', 'http://127.0.0.1:8080/api/static/public/tenant_1/2025/09/24/1758707254453896600.png', 1, 1, '2025-09-28 15:53:08', '', 0, NULL, '2025-09-18 20:21:12', '2025-09-28 15:53:08', NULL);
INSERT INTO `users` VALUES (10, 2, 'test', '$2a$10$Ck5B5o1Md2O7K.tER/Hug.5phi4hRazcY04WdF46ykrmV3KdPo.3G', 'test', NULL, NULL, '', 1, 0, '2025-09-20 11:52:40', '', 0, NULL, '2025-09-20 10:49:05', '2025-09-23 18:16:28', NULL);
INSERT INTO `users` VALUES (11, 2, 'test01', '$2a$10$0EcaMxX5aGfGQzO2wG85ye0OOnhY40TH0rUWJYthVIPHimaVRgMq2', 'Test01', NULL, NULL, '', 1, 0, NULL, '', 0, NULL, '2025-09-20 10:49:20', '2025-09-20 10:49:20', NULL);
INSERT INTO `users` VALUES (12, 1, 'test', '$2a$10$4daTLk90ZT.qEUFlVniYjO/4DJ/s1/d1BuwMhGKaou45.BjghDGka', '测试用户', 'test@qq.com', NULL, 'http://127.0.0.1:8080/api/static/public/tenant_1/2025/09/24/1758707486724848500.png', 1, 0, '2025-09-24 22:31:15', '', 0, NULL, '2025-09-21 08:50:31', '2025-09-28 15:53:19', NULL);

SET FOREIGN_KEY_CHECKS = 1;
