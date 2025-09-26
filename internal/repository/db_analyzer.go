package repository

import (
	"fmt"

	"gorm.io/gorm"
)

// DBAnalyzerRepository 数据库分析仓储
type DBAnalyzerRepository struct {
	db *gorm.DB
}

// NewDBAnalyzerRepository 创建数据库分析仓储
func NewDBAnalyzerRepository(db *gorm.DB) *DBAnalyzerRepository {
	return &DBAnalyzerRepository{
		db: db,
	}
}

// TableBasicInfo 数据库表基本信息
type TableBasicInfo struct {
	TableName    string `json:"tableName"`    // 表名
	TableComment string `json:"tableComment"` // 表注释
	CreateTime   string `json:"createTime"`   // 创建时间
	UpdateTime   string `json:"updateTime"`   // 更新时间
}

// TableColumn 数据库字段信息
type TableColumn struct {
	ColumnName      string `json:"columnName"`      // 字段名称
	ColumnType      string `json:"columnType"`      // 字段类型
	ColumnComment   string `json:"columnComment"`   // 字段注释
	IsNullable      string `json:"isNullable"`      // 是否可空
	ColumnDefault   string `json:"columnDefault"`   // 默认值
	ColumnKey       string `json:"columnKey"`       // 键类型
	Extra           string `json:"extra"`           // 扩展信息
	OrdinalPosition int    `json:"ordinalPosition"` // 字段位置
}

// GetTableList 获取数据库表列表
func (r *DBAnalyzerRepository) GetTableList() ([]TableBasicInfo, error) {
	var tables []TableBasicInfo

	query := `
		SELECT
			TABLE_NAME as table_name,
			TABLE_COMMENT as table_comment,
			CREATE_TIME as create_time,
			UPDATE_TIME as update_time
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_TYPE = 'BASE TABLE'
		ORDER BY TABLE_NAME
	`

	err := r.db.Raw(query).Scan(&tables).Error
	if err != nil {
		return nil, fmt.Errorf("查询表列表失败: %v", err)
	}

	return tables, nil
}

// GetTableInfo 获取表基本信息
func (r *DBAnalyzerRepository) GetTableInfo(tableName string) (*TableBasicInfo, error) {
	var table TableBasicInfo

	query := `
		SELECT
			TABLE_NAME as table_name,
			TABLE_COMMENT as table_comment,
			CREATE_TIME as create_time,
			UPDATE_TIME as update_time
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = ?
	`

	err := r.db.Raw(query, tableName).Scan(&table).Error
	if err != nil {
		return nil, fmt.Errorf("查询表信息失败: %v", err)
	}

	if table.TableName == "" {
		return nil, fmt.Errorf("表 '%s' 不存在", tableName)
	}

	return &table, nil
}

// GetTableColumns 获取表字段信息
func (r *DBAnalyzerRepository) GetTableColumns(tableName string) ([]TableColumn, error) {
	var columns []TableColumn

	query := `
		SELECT
			COLUMN_NAME as column_name,
			COLUMN_TYPE as column_type,
			IFNULL(COLUMN_COMMENT, '') as column_comment,
			IS_NULLABLE as is_nullable,
			IFNULL(COLUMN_DEFAULT, '') as column_default,
			COLUMN_KEY as column_key,
			EXTRA as extra,
			ORDINAL_POSITION as ordinal_position
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION
	`

	err := r.db.Raw(query, tableName).Scan(&columns).Error
	if err != nil {
		return nil, fmt.Errorf("查询表字段信息失败: %v", err)
	}

	return columns, nil
}

// TableExists 检查表是否存在
func (r *DBAnalyzerRepository) TableExists(tableName string) (bool, error) {
	var count int64

	query := `
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = ?
	`

	err := r.db.Raw(query, tableName).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查表是否存在失败: %v", err)
	}

	return count > 0, nil
}

// GetTableIndexes 获取表索引信息
func (r *DBAnalyzerRepository) GetTableIndexes(tableName string) ([]TableIndex, error) {
	var indexes []TableIndex

	query := `
		SELECT
			INDEX_NAME as index_name,
			COLUMN_NAME as column_name,
			NON_UNIQUE as non_unique,
			SEQ_IN_INDEX as seq_in_index,
			COLLATION as collation,
			CARDINALITY as cardinality,
			SUB_PART as sub_part,
			PACKED as packed,
			NULLABLE as nullable,
			INDEX_TYPE as index_type,
			INDEX_COMMENT as index_comment
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY INDEX_NAME, SEQ_IN_INDEX
	`

	err := r.db.Raw(query, tableName).Scan(&indexes).Error
	if err != nil {
		return nil, fmt.Errorf("查询表索引信息失败: %v", err)
	}

	return indexes, nil
}

// TableIndex 表索引信息
type TableIndex struct {
	IndexName    string `json:"indexName"`    // 索引名称
	ColumnName   string `json:"columnName"`   // 字段名称
	NonUnique    int    `json:"nonUnique"`    // 是否唯一索引
	SeqInIndex   int    `json:"seqInIndex"`   // 索引中的位置
	Collation    string `json:"collation"`    // 排序方式
	Cardinality  int64  `json:"cardinality"`  // 基数
	SubPart      *int   `json:"subPart"`      // 子部分
	Packed       string `json:"packed"`       // 打包方式
	Nullable     string `json:"nullable"`     // 是否可空
	IndexType    string `json:"indexType"`    // 索引类型
	IndexComment string `json:"indexComment"` // 索引注释
}

// GetDatabaseTables 获取指定数据库的表列表
func (r *DBAnalyzerRepository) GetDatabaseTables(databaseName string) ([]TableBasicInfo, error) {
	var tables []TableBasicInfo

	query := `
		SELECT
			TABLE_NAME as table_name,
			TABLE_COMMENT as table_comment,
			CREATE_TIME as create_time,
			UPDATE_TIME as update_time
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = ?
		AND TABLE_TYPE = 'BASE TABLE'
		ORDER BY TABLE_NAME
	`

	err := r.db.Raw(query, databaseName).Scan(&tables).Error
	if err != nil {
		return nil, fmt.Errorf("查询数据库表列表失败: %v", err)
	}

	return tables, nil
}

// GetTableSize 获取表大小信息
func (r *DBAnalyzerRepository) GetTableSize(tableName string) (*TableSize, error) {
	var size TableSize

	query := `
		SELECT
			TABLE_NAME as table_name,
			ENGINE as engine,
			TABLE_ROWS as table_rows,
			DATA_LENGTH as data_length,
			INDEX_LENGTH as index_length,
			DATA_LENGTH + INDEX_LENGTH as total_length
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
	`

	err := r.db.Raw(query, tableName).Scan(&size).Error
	if err != nil {
		return nil, fmt.Errorf("查询表大小信息失败: %v", err)
	}

	return &size, nil
}

// TableSize 表大小信息
type TableSize struct {
	TableName   string `json:"tableName"`   // 表名
	Engine      string `json:"engine"`      // 存储引擎
	TableRows   int64  `json:"tableRows"`   // 行数
	DataLength  int64  `json:"dataLength"`  // 数据大小
	IndexLength int64  `json:"indexLength"` // 索引大小
	TotalLength int64  `json:"totalLength"` // 总大小
}

// SearchTables 搜索表名
func (r *DBAnalyzerRepository) SearchTables(keyword string) ([]TableBasicInfo, error) {
	var tables []TableBasicInfo

	query := `
		SELECT
			TABLE_NAME as table_name,
			TABLE_COMMENT as table_comment,
			CREATE_TIME as create_time,
			UPDATE_TIME as update_time
		FROM INFORMATION_SCHEMA.TABLES
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_TYPE = 'BASE TABLE'
		AND (TABLE_NAME LIKE ? OR TABLE_COMMENT LIKE ?)
		ORDER BY TABLE_NAME
	`

	searchPattern := "%" + keyword + "%"
	err := r.db.Raw(query, searchPattern, searchPattern).Scan(&tables).Error
	if err != nil {
		return nil, fmt.Errorf("搜索表失败: %v", err)
	}

	return tables, nil
}
