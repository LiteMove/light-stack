-- ============================================
-- 代码生成器数据库表结构
-- ============================================

-- 数据表配置表
CREATE TABLE `gen_table_configs` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
    `table_name` varchar(64) NOT NULL COMMENT '表名称',
    `table_comment` varchar(255) DEFAULT '' COMMENT '表描述',
    `business_name` varchar(64) NOT NULL COMMENT '业务名称',
    `module_name` varchar(64) NOT NULL COMMENT '模块名称',
    `function_name` varchar(64) NOT NULL COMMENT '功能名称',
    `class_name` varchar(64) NOT NULL COMMENT '类名',
    `package_name` varchar(64) NOT NULL COMMENT '包名',
    `author` varchar(64) DEFAULT 'system' COMMENT '作者',
    `parent_menu_id` bigint(20) DEFAULT NULL COMMENT '父级菜单ID',
    `menu_name` varchar(64) DEFAULT '' COMMENT '菜单名称',
    `menu_url` varchar(255) DEFAULT '' COMMENT '菜单URL',
    `menu_icon` varchar(64) DEFAULT '' COMMENT '菜单图标',
    `permissions` text COMMENT '权限字符串(JSON数组)',
    `options` text COMMENT '其他配置选项(JSON)',
    `remark` varchar(500) DEFAULT '' COMMENT '备注',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `created_by` bigint(20) DEFAULT NULL COMMENT '创建人',
    `updated_by` bigint(20) DEFAULT NULL COMMENT '更新人',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_table_name` (`table_name`),
    INDEX `idx_business_name` (`business_name`),
    INDEX `idx_module_name` (`module_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代码生成表配置';

-- 字段配置表
CREATE TABLE `gen_table_columns` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '字段ID',
    `table_config_id` bigint(20) NOT NULL COMMENT '表配置ID',
    `column_name` varchar(64) NOT NULL COMMENT '字段名称',
    `column_comment` varchar(255) DEFAULT '' COMMENT '字段描述',
    `column_type` varchar(32) NOT NULL COMMENT '字段类型',
    `go_type` varchar(32) NOT NULL COMMENT 'Go类型',
    `go_field` varchar(64) NOT NULL COMMENT 'Go字段名',
    `is_pk` tinyint(1) DEFAULT 0 COMMENT '是否主键',
    `is_increment` tinyint(1) DEFAULT 0 COMMENT '是否自增',
    `is_required` tinyint(1) DEFAULT 0 COMMENT '是否必填',
    `is_insert` tinyint(1) DEFAULT 1 COMMENT '是否为插入字段',
    `is_edit` tinyint(1) DEFAULT 1 COMMENT '是否为编辑字段',
    `is_list` tinyint(1) DEFAULT 1 COMMENT '是否列表字段',
    `is_query` tinyint(1) DEFAULT 0 COMMENT '是否查询字段',
    `query_type` varchar(32) DEFAULT 'EQ' COMMENT '查询方式(EQ等于、NE不等于、GT大于、LT小于、LIKE模糊、BETWEEN范围)',
    `html_type` varchar(32) DEFAULT 'input' COMMENT '显示类型(input、textarea、select、radio、checkbox、datetime)',
    `dict_type` varchar(64) DEFAULT '' COMMENT '字典类型',
    `sort` int(11) DEFAULT 0 COMMENT '排序',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_table_config_id` (`table_config_id`),
    INDEX `idx_column_name` (`column_name`),
    CONSTRAINT `fk_gen_table_columns_config` FOREIGN KEY (`table_config_id`) REFERENCES `gen_table_configs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代码生成字段配置';

-- 生成历史记录表
CREATE TABLE `gen_history` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '历史ID',
    `table_config_id` bigint(20) NOT NULL COMMENT '表配置ID',
    `table_name` varchar(64) NOT NULL COMMENT '表名称',
    `business_name` varchar(64) NOT NULL COMMENT '业务名称',
    `generate_type` varchar(32) NOT NULL COMMENT '生成类型(all全部、backend后端、frontend前端)',
    `file_count` int(11) DEFAULT 0 COMMENT '生成文件数量',
    `file_size` bigint(20) DEFAULT 0 COMMENT '文件大小(字节)',
    `download_count` int(11) DEFAULT 0 COMMENT '下载次数',
    `status` varchar(32) NOT NULL DEFAULT 'success' COMMENT '生成状态(success成功、failed失败、processing处理中)',
    `error_message` text COMMENT '错误信息',
    `file_path` varchar(500) DEFAULT '' COMMENT '生成文件路径',
    `remark` varchar(500) DEFAULT '' COMMENT '备注',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_by` bigint(20) DEFAULT NULL COMMENT '创建人',
    PRIMARY KEY (`id`),
    INDEX `idx_table_config_id` (`table_config_id`),
    INDEX `idx_table_name` (`table_name`),
    INDEX `idx_created_at` (`created_at`),
    CONSTRAINT `fk_gen_history_config` FOREIGN KEY (`table_config_id`) REFERENCES `gen_table_configs` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='代码生成历史记录';