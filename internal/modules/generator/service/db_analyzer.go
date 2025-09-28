package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/LiteMove/light-stack/internal/modules/generator/model"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/shared/utils"
	"gorm.io/gorm"
)

// DBAnalyzerService 数据库分析服务
type DBAnalyzerService struct {
	repo *repository.DBAnalyzerRepository
	db   *gorm.DB
}

// NewDBAnalyzerService 创建数据库分析服务
func NewDBAnalyzerService(repo *repository.DBAnalyzerRepository, db *gorm.DB) *DBAnalyzerService {
	return &DBAnalyzerService{
		repo: repo,
		db:   db,
	}
}

// GetTableList 获取数据库表列表
func (s *DBAnalyzerService) GetTableList() ([]model.TableInfo, error) {
	tables, err := s.repo.GetTableList()
	if err != nil {
		return nil, fmt.Errorf("获取表列表失败: %v", err)
	}

	var tableInfos []model.TableInfo
	for _, table := range tables {
		tableInfos = append(tableInfos, model.TableInfo{
			TableName:    table.TableName,
			TableComment: table.TableComment,
			CreateTime:   table.CreateTime,
			UpdateTime:   table.UpdateTime,
		})
	}

	return tableInfos, nil
}

// GetTableInfo 获取表详细信息
func (s *DBAnalyzerService) GetTableInfo(tableName string) (*model.TableInfo, error) {
	// 获取表基本信息
	tableInfo, err := s.repo.GetTableInfo(tableName)
	if err != nil {
		return nil, fmt.Errorf("获取表信息失败: %v", err)
	}

	// 获取字段信息
	columns, err := s.repo.GetTableColumns(tableName)
	if err != nil {
		return nil, fmt.Errorf("获取表字段信息失败: %v", err)
	}

	// 转换字段信息
	var columnInfos []model.ColumnInfo
	for _, col := range columns {
		columnInfo := s.convertColumnInfo(col)
		columnInfos = append(columnInfos, columnInfo)
	}

	return &model.TableInfo{
		TableName:    tableInfo.TableName,
		TableComment: tableInfo.TableComment,
		CreateTime:   tableInfo.CreateTime,
		UpdateTime:   tableInfo.UpdateTime,
		Columns:      columnInfos,
	}, nil
}

// GetTableColumns 获取表字段信息
func (s *DBAnalyzerService) GetTableColumns(tableName string) ([]model.ColumnInfo, error) {
	columns, err := s.repo.GetTableColumns(tableName)
	if err != nil {
		return nil, fmt.Errorf("获取表字段信息失败: %v", err)
	}

	var columnInfos []model.ColumnInfo
	for _, col := range columns {
		columnInfo := s.convertColumnInfo(col)
		columnInfos = append(columnInfos, columnInfo)
	}

	return columnInfos, nil
}

// convertColumnInfo 转换字段信息
func (s *DBAnalyzerService) convertColumnInfo(col repository.TableColumn) model.ColumnInfo {
	columnInfo := model.ColumnInfo{
		ColumnName:      col.ColumnName,
		ColumnType:      col.ColumnType,
		ColumnComment:   col.ColumnComment,
		IsNullable:      col.IsNullable,
		ColumnDefault:   col.ColumnDefault,
		ColumnKey:       col.ColumnKey,
		Extra:           col.Extra,
		OrdinalPosition: col.OrdinalPosition,
	}

	// 转换Go类型
	columnInfo.GoType = s.convertGoType(col.ColumnType)
	columnInfo.GoField = s.convertGoField(col.ColumnName)

	// 判断字段属性
	columnInfo.IsPk = col.ColumnKey == "PRI"
	columnInfo.IsIncrement = strings.Contains(col.Extra, "auto_increment")
	columnInfo.IsRequired = col.IsNullable == "NO" && col.ColumnDefault == ""

	// 设置默认配置
	columnInfo.IsInsert = !columnInfo.IsPk || !columnInfo.IsIncrement
	columnInfo.IsEdit = !columnInfo.IsPk
	columnInfo.IsList = true
	columnInfo.IsQuery = s.isQueryField(col.ColumnName)

	// 设置查询类型
	columnInfo.QueryType = s.getQueryType(columnInfo.GoType)

	// 设置HTML类型
	columnInfo.HtmlType = s.getHtmlType(col.ColumnName, columnInfo.GoType, col.ColumnType)

	// 设置字典类型（如果需要）
	columnInfo.DictType = s.getDictType(col.ColumnName)

	return columnInfo
}

// convertGoType 转换MySQL类型为Go类型
func (s *DBAnalyzerService) convertGoType(columnType string) string {
	// 转换为小写并移除长度限制
	lowerType := strings.ToLower(columnType)

	// 移除括号内容
	re := regexp.MustCompile(`\([^)]*\)`)
	baseType := re.ReplaceAllString(lowerType, "")

	// 移除unsigned等修饰符
	baseType = strings.TrimSpace(strings.ReplaceAll(baseType, "unsigned", ""))

	switch baseType {
	case "tinyint":
		if strings.Contains(lowerType, "tinyint(1)") {
			return "bool"
		}
		return "int8"
	case "smallint":
		return "int16"
	case "mediumint", "int":
		return "int"
	case "bigint":
		return "int64"
	case "float":
		return "float32"
	case "double", "decimal":
		return "float64"
	case "char", "varchar", "text", "longtext", "mediumtext", "tinytext":
		return "string"
	case "date", "datetime", "timestamp", "time":
		return "time.Time"
	case "json":
		return "string"
	default:
		return "string"
	}
}

// convertGoField 转换字段名为Go字段名（驼峰命名）
func (s *DBAnalyzerService) convertGoField(columnName string) string {
	return utils.ToCamelCase(columnName)
}

// isQueryField 判断是否为查询字段
func (s *DBAnalyzerService) isQueryField(columnName string) bool {
	queryFields := []string{"name", "title", "status", "type", "code", "keyword"}
	columnLower := strings.ToLower(columnName)

	for _, field := range queryFields {
		if strings.Contains(columnLower, field) {
			return true
		}
	}

	return false
}

// getQueryType 获取查询类型
func (s *DBAnalyzerService) getQueryType(goType string) string {
	switch goType {
	case "string":
		return model.QueryTypeLike
	case "int", "int8", "int16", "int32", "int64", "float32", "float64":
		return model.QueryTypeEQ
	case "bool":
		return model.QueryTypeEQ
	case "time.Time":
		return model.QueryTypeBetween
	default:
		return model.QueryTypeEQ
	}
}

// getHtmlType 获取HTML类型
func (s *DBAnalyzerService) getHtmlType(columnName, goType, columnType string) string {
	columnLower := strings.ToLower(columnName)

	// 根据字段名判断
	if strings.Contains(columnLower, "password") {
		return "password"
	}
	if strings.Contains(columnLower, "email") {
		return "email"
	}
	if strings.Contains(columnLower, "url") {
		return "url"
	}
	if strings.Contains(columnLower, "image") || strings.Contains(columnLower, "avatar") {
		return model.HtmlTypeUpload
	}
	if strings.Contains(columnLower, "content") || strings.Contains(columnLower, "description") ||
		strings.Contains(columnLower, "remark") || strings.Contains(columnLower, "detail") {
		return model.HtmlTypeTextarea
	}
	if strings.Contains(columnLower, "status") || strings.Contains(columnLower, "type") {
		return model.HtmlTypeSelect
	}
	if strings.Contains(columnLower, "date") || strings.Contains(columnLower, "time") {
		return model.HtmlTypeDatetime
	}

	// 根据Go类型判断
	switch goType {
	case "bool":
		return model.HtmlTypeRadio
	case "time.Time":
		return model.HtmlTypeDatetime
	default:
		// 根据MySQL类型判断
		if strings.Contains(strings.ToLower(columnType), "text") {
			return model.HtmlTypeTextarea
		}
		return model.HtmlTypeInput
	}
}

// getDictType 获取字典类型
func (s *DBAnalyzerService) getDictType(columnName string) string {
	columnLower := strings.ToLower(columnName)

	dictMappings := map[string]string{
		"status": "sys_status",
		"type":   "sys_type",
		"sex":    "sys_user_sex",
		"gender": "sys_user_sex",
	}

	for key, value := range dictMappings {
		if strings.Contains(columnLower, key) {
			return value
		}
	}

	return ""
}

// ValidateTableName 验证表名
func (s *DBAnalyzerService) ValidateTableName(tableName string) error {
	if tableName == "" {
		return fmt.Errorf("表名不能为空")
	}

	// 检查表名格式
	matched, err := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9_]*$", tableName)
	if err != nil {
		return fmt.Errorf("正则表达式错误: %v", err)
	}
	if !matched {
		return fmt.Errorf("表名格式不正确，只能包含字母、数字和下划线，且以字母开头")
	}

	// 检查表是否存在
	exists, err := s.repo.TableExists(tableName)
	if err != nil {
		return fmt.Errorf("检查表是否存在失败: %v", err)
	}
	if !exists {
		return fmt.Errorf("表 '%s' 不存在", tableName)
	}

	return nil
}

// GetDatabaseName 获取当前数据库名
func (s *DBAnalyzerService) GetDatabaseName() (string, error) {
	var dbName string
	err := s.db.Raw("SELECT DATABASE()").Scan(&dbName).Error
	if err != nil {
		return "", fmt.Errorf("获取数据库名失败: %v", err)
	}
	return dbName, nil
}
