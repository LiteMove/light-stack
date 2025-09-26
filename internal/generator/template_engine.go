package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/utils"
)

// TemplateEngine 模板引擎
type TemplateEngine struct {
	templates map[string]*template.Template
}

// NewTemplateEngine 创建模板引擎
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		templates: make(map[string]*template.Template),
	}
}

// LoadTemplates 加载模板文件
func (e *TemplateEngine) LoadTemplates(templateDir string) error {
	// 定义模板文件路径
	templateFiles := map[string]string{
		// 后端模板
		"model":      filepath.Join(templateDir, "backend", "model.go.tpl"),
		"service":    filepath.Join(templateDir, "backend", "service.go.tpl"),
		"controller": filepath.Join(templateDir, "backend", "controller.go.tpl"),
		"repository": filepath.Join(templateDir, "backend", "repository.go.tpl"),
		"request":    filepath.Join(templateDir, "backend", "request.go.tpl"),

		// SQL模板
		"menu_sql": filepath.Join(templateDir, "sql", "menu.sql.tpl"),
	}

	// 加载每个模板
	for name, path := range templateFiles {
		// 检查文件是否存在
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// 如果模板文件不存在，跳过加载，但记录警告
			fmt.Printf("警告: 模板文件不存在，跳过加载: %s\n", path)
			continue
		}

		tmpl, err := e.loadTemplate(name, path)
		if err != nil {
			return fmt.Errorf("加载模板 %s 失败: %v", name, err)
		}
		e.templates[name] = tmpl
	}

	return nil
}

// loadTemplate 加载单个模板文件
func (e *TemplateEngine) loadTemplate(name, path string) (*template.Template, error) {
	// 创建模板函数映射
	funcMap := template.FuncMap{
		"toCamelCase":        utils.ToCamelCase,
		"toPascalCase":       utils.ToPascalCase,
		"toSnakeCase":        utils.ToSnakeCase,
		"toKebabCase":        utils.ToKebabCase,
		"toLower":            strings.ToLower,
		"toUpper":            strings.ToUpper,
		"pluralize":          utils.Pluralize,
		"uncapitalize":       utils.Uncapitalize,
		"now":                time.Now,
		"formatDate":         formatDate,
		"contains":           strings.Contains,
		"hasPrefix":          strings.HasPrefix,
		"hasSuffix":          strings.HasSuffix,
		"join":               strings.Join,
		"split":              strings.Split,
		"replace":            strings.ReplaceAll,
		"trim":               strings.TrimSpace,
		"isEmptyString":      utils.IsEmpty,
		"defaultString":      utils.DefaultString,
		"add":                add,
		"sub":                sub,
		"mul":                mul,
		"div":                div,
		"mod":                mod,
		"eq":                 eq,
		"ne":                 ne,
		"gt":                 gt,
		"ge":                 ge,
		"lt":                 lt,
		"le":                 le,
		"and":                and,
		"or":                 or,
		"not":                not,
		"generateGoField":    generateGoField,
		"generateJSField":    generateJSField,
		"generatePermission": generatePermission,
		"isQueryField":       isQueryField,
		"isListField":        isListField,
		"isFormField":        isFormField,
		"isRequiredField":    isRequiredField,
		"getHtmlInputType":   getHtmlInputType,
		"getValidationRules": getValidationRules,
		"getDefaultValue":    getDefaultValue,
	}

	// 解析模板文件
	tmpl, err := template.New(name).Funcs(funcMap).ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("解析模板文件失败: %v", err)
	}

	return tmpl, nil
}

// RenderTemplate 渲染模板
func (e *TemplateEngine) RenderTemplate(templateName string, data *model.TemplateData) (string, error) {
	tmpl, exists := e.templates[templateName]
	if !exists {
		return "", fmt.Errorf("模板 '%s' 不存在", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("渲染模板失败: %v", err)
	}

	return buf.String(), nil
}

// HasTemplate 检查模板是否存在
func (e *TemplateEngine) HasTemplate(templateName string) bool {
	_, exists := e.templates[templateName]
	return exists
}

// RenderAllTemplates 渲染所有模板
func (e *TemplateEngine) RenderAllTemplates(data *model.TemplateData) (map[string]string, error) {
	results := make(map[string]string)

	for name := range e.templates {
		content, err := e.RenderTemplate(name, data)
		if err != nil {
			return nil, fmt.Errorf("渲染模板 %s 失败: %v", name, err)
		}
		results[name] = content
	}

	return results, nil
}

// PrepareTemplateData 准备模板数据
func (e *TemplateEngine) PrepareTemplateData(config *model.GenTableConfig) *model.TemplateData {
	data := &model.TemplateData{
		PackageName:  config.PackageName,
		ClassName:    config.ClassName,
		TableName:    config.TableName,
		BusinessName: config.BusinessName,
		ModuleName:   config.ModuleName,
		FunctionName: config.FunctionName,
		Author:       config.Author,
		Date:         formatDate(time.Now()),
		ParentMenuID: *config.ParentMenuID,
		MenuName:     config.MenuName,
		MenuURL:      config.MenuURL,
		MenuIcon:     config.MenuIcon,
		Permissions:  config.GetPermissions(),
		Options:      config.GetOptions(),
	}

	// 转换字段信息
	var fields []model.ColumnInfo
	var queryFields []model.ColumnInfo
	var listFields []model.ColumnInfo
	var formFields []model.ColumnInfo
	hasQuery := false

	for _, col := range config.Columns {
		field := model.ColumnInfo{
			ColumnName:    col.ColumnName,
			ColumnComment: col.ColumnComment,
			ColumnType:    col.ColumnType,
			GoType:        col.GoType,
			GoField:       col.GoField,
			IsPk:          col.IsPk,
			IsIncrement:   col.IsIncrement,
			IsRequired:    col.IsRequired,
			IsInsert:      col.IsInsert,
			IsEdit:        col.IsEdit,
			IsList:        col.IsList,
			IsQuery:       col.IsQuery,
			QueryType:     col.QueryType,
			HtmlType:      col.HtmlType,
			DictType:      col.DictType,
		}

		fields = append(fields, field)

		if field.IsQuery {
			queryFields = append(queryFields, field)
			hasQuery = true
		}

		if field.IsList {
			listFields = append(listFields, field)
		}

		if field.IsInsert || field.IsEdit {
			formFields = append(formFields, field)
		}
	}

	data.Fields = fields
	data.QueryFields = queryFields
	data.ListFields = listFields
	data.FormFields = formFields
	data.HasQuery = hasQuery

	return data
}

// 模板函数定义

// formatDate 格式化日期
func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// 数学运算函数
func add(a, b int) int { return a + b }
func sub(a, b int) int { return a - b }
func mul(a, b int) int { return a * b }
func div(a, b int) int { return a / b }
func mod(a, b int) int { return a % b }

// 比较函数
func eq(a, b interface{}) bool { return a == b }
func ne(a, b interface{}) bool { return a != b }
func gt(a, b int) bool         { return a > b }
func ge(a, b int) bool         { return a >= b }
func lt(a, b int) bool         { return a < b }
func le(a, b int) bool         { return a <= b }

// 逻辑函数
func and(a, b bool) bool { return a && b }
func or(a, b bool) bool  { return a || b }
func not(a bool) bool    { return !a }

// generateGoField 生成Go字段
func generateGoField(columnName string) string {
	return utils.ToCamelCase(columnName)
}

// generateJSField 生成JavaScript字段
func generateJSField(columnName string) string {
	return utils.Uncapitalize(utils.ToCamelCase(columnName))
}

// generatePermission 生成权限字符串
func generatePermission(moduleName, businessName, operation string) string {
	return fmt.Sprintf("%s:%s:%s", strings.ToLower(moduleName), strings.ToLower(businessName), operation)
}

// isQueryField 判断是否为查询字段
func isQueryField(field model.ColumnInfo) bool {
	return field.IsQuery
}

// isListField 判断是否为列表字段
func isListField(field model.ColumnInfo) bool {
	return field.IsList
}

// isFormField 判断是否为表单字段
func isFormField(field model.ColumnInfo) bool {
	return field.IsInsert || field.IsEdit
}

// isRequiredField 判断是否为必填字段
func isRequiredField(field model.ColumnInfo) bool {
	return field.IsRequired
}

// getHtmlInputType 获取HTML输入类型
func getHtmlInputType(htmlType string) string {
	switch htmlType {
	case "password":
		return "password"
	case "email":
		return "email"
	case "url":
		return "url"
	case "number":
		return "number"
	case "datetime":
		return "datetime-local"
	case "date":
		return "date"
	case "time":
		return "time"
	default:
		return "text"
	}
}

// getValidationRules 获取验证规则
func getValidationRules(field model.ColumnInfo) []string {
	var rules []string

	if field.IsRequired {
		rules = append(rules, "required")
	}

	switch field.GoType {
	case "int", "int8", "int16", "int32", "int64":
		rules = append(rules, "numeric")
	case "float32", "float64":
		rules = append(rules, "numeric")
	case "string":
		if strings.Contains(field.ColumnName, "email") {
			rules = append(rules, "email")
		}
		if strings.Contains(field.ColumnName, "url") {
			rules = append(rules, "url")
		}
	}

	return rules
}

// getDefaultValue 获取默认值
func getDefaultValue(field model.ColumnInfo) string {
	switch field.GoType {
	case "bool":
		return "false"
	case "int", "int8", "int16", "int32", "int64":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "string":
		return "\"\""
	case "time.Time":
		return "time.Time{}"
	default:
		return "nil"
	}
}
