package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LiteMove/light-stack/internal/modules/generator/model"
)

// CodeGenerator 代码生成器
type CodeGenerator struct {
	templateEngine *TemplateEngine
	outputDir      string
}

// NewCodeGenerator 创建代码生成器
func NewCodeGenerator(templateEngine *TemplateEngine) *CodeGenerator {
	return &CodeGenerator{
		templateEngine: templateEngine,
		outputDir:      "generated",
	}
}

// GenerateCode 生成代码 - 生成所有可用模板的代码
func (g *CodeGenerator) GenerateCode(config *model.GenTableConfig, options *GenerateOptions) (*GenerateResult, error) {
	// 准备模板数据
	templateData := g.templateEngine.PrepareTemplateData(config)

	result := &GenerateResult{
		TableName:    config.TableName,
		BusinessName: config.BusinessName,
		Files:        make(map[string]string),
		StartTime:    time.Now(),
	}

	// 生成所有可用模板的代码，不再根据选项分类
	if err := g.generateAllTemplates(templateData, result); err != nil {
		result.Success = false
		result.ErrorMessage = err.Error()
		return result, err
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Success = true
	result.FileCount = len(result.Files)

	// 计算总文件大小
	for _, content := range result.Files {
		result.TotalSize += int64(len(content))
	}

	return result, nil
}

// generateBackendCode 生成后端代码
func (g *CodeGenerator) generateBackendCode(data *model.TemplateData, result *GenerateResult) error {
	backendTemplates := map[string]string{
		"model":      fmt.Sprintf("internal/model/%s.go", strings.ToLower(data.BusinessName)),
		"service":    fmt.Sprintf("internal/service/%s_service.go", strings.ToLower(data.BusinessName)),
		"controller": fmt.Sprintf("internal/controller/%s_controller.go", strings.ToLower(data.BusinessName)),
		"repository": fmt.Sprintf("internal/repository/%s_repository.go", strings.ToLower(data.BusinessName)),
		"request":    fmt.Sprintf("internal/model/%s_request.go", strings.ToLower(data.BusinessName)),
	}

	for templateName, fileName := range backendTemplates {
		// 检查模板是否存在
		if !g.templateEngine.HasTemplate(templateName) {
			fmt.Printf("警告: 后端模板 %s 不存在，跳过生成文件: %s\n", templateName, fileName)
			continue
		}

		content, err := g.templateEngine.RenderTemplate(templateName, data)
		if err != nil {
			return fmt.Errorf("生成后端文件 %s 失败: %v", fileName, err)
		}
		result.Files[fileName] = content
	}

	return nil
}

// generateFrontendCode 生成前端代码
func (g *CodeGenerator) generateFrontendCode(data *model.TemplateData, result *GenerateResult) error {
	frontendTemplates := map[string]string{
		"list_vue":   fmt.Sprintf("web/src/views/%s/%sList.vue", strings.ToLower(data.ModuleName), data.ClassName),
		"form_vue":   fmt.Sprintf("web/src/views/%s/components/%sForm.vue", strings.ToLower(data.ModuleName), data.ClassName),
		"detail_vue": fmt.Sprintf("web/src/views/%s/components/%sDetail.vue", strings.ToLower(data.ModuleName), data.ClassName),
		"api_ts":     fmt.Sprintf("web/src/api/%s.ts", strings.ToLower(data.BusinessName)),
		"types_ts":   fmt.Sprintf("web/src/types/%s.ts", strings.ToLower(data.BusinessName)),
	}

	for templateName, fileName := range frontendTemplates {
		// 检查模板是否存在
		if !g.templateEngine.HasTemplate(templateName) {
			fmt.Printf("警告: 前端模板 %s 不存在，跳过生成文件: %s\n", templateName, fileName)
			continue
		}

		content, err := g.templateEngine.RenderTemplate(templateName, data)
		if err != nil {
			return fmt.Errorf("生成前端文件 %s 失败: %v", fileName, err)
		}
		result.Files[fileName] = content
	}

	return nil
}

// generateSQLCode 生成SQL代码
func (g *CodeGenerator) generateSQLCode(data *model.TemplateData, result *GenerateResult) error {
	sqlTemplates := map[string]string{
		"menu_sql": fmt.Sprintf("sql/%s_menu.sql", strings.ToLower(data.BusinessName)),
	}

	for templateName, fileName := range sqlTemplates {
		// 检查模板是否存在
		if !g.templateEngine.HasTemplate(templateName) {
			fmt.Printf("警告: SQL模板 %s 不存在，跳过生成文件: %s\n", templateName, fileName)
			continue
		}

		content, err := g.templateEngine.RenderTemplate(templateName, data)
		if err != nil {
			return fmt.Errorf("生成SQL文件 %s 失败: %v", fileName, err)
		}
		result.Files[fileName] = content
	}

	return nil
}

// SaveFiles 保存生成的文件到磁盘
func (g *CodeGenerator) SaveFiles(result *GenerateResult, baseDir string) error {
	for fileName, content := range result.Files {
		fullPath := filepath.Join(baseDir, fileName)

		// 确保目录存在
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录 %s 失败: %v", dir, err)
		}

		// 写入文件
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("写入文件 %s 失败: %v", fullPath, err)
		}
	}

	return nil
}

// generateAllTemplates 生成所有可用模板的代码
func (g *CodeGenerator) generateAllTemplates(data *model.TemplateData, result *GenerateResult) error {
	// 定义所有模板及其对应的文件名
	allTemplates := map[string]string{
		// 后端模板
		"model":      fmt.Sprintf("internal/model/%s.go", strings.ToLower(data.BusinessName)),
		"service":    fmt.Sprintf("internal/service/%s_service.go", strings.ToLower(data.BusinessName)),
		"controller": fmt.Sprintf("internal/controller/%s_controller.go", strings.ToLower(data.BusinessName)),
		"repository": fmt.Sprintf("internal/repository/%s_repository.go", strings.ToLower(data.BusinessName)),
		"request":    fmt.Sprintf("internal/model/%s_request.go", strings.ToLower(data.BusinessName)),

		// 前端模板
		"list_vue":   fmt.Sprintf("web/src/views/%s/%sList.vue", strings.ToLower(data.ModuleName), data.ClassName),
		"form_vue":   fmt.Sprintf("web/src/views/%s/components/%sForm.vue", strings.ToLower(data.ModuleName), data.ClassName),
		"detail_vue": fmt.Sprintf("web/src/views/%s/components/%sDetail.vue", strings.ToLower(data.ModuleName), data.ClassName),
		"api_ts":     fmt.Sprintf("web/src/api/%s.ts", strings.ToLower(data.BusinessName)),
		"types_ts":   fmt.Sprintf("web/src/types/%s.ts", strings.ToLower(data.BusinessName)),

		// SQL模板
		"menu_sql": fmt.Sprintf("sql/%s_menu.sql", strings.ToLower(data.BusinessName)),
	}

	// 遍历所有模板并生成代码
	for templateName, fileName := range allTemplates {
		// 检查模板是否存在
		if !g.templateEngine.HasTemplate(templateName) {
			fmt.Printf("警告: 模板 %s 不存在，跳过生成文件: %s\n", templateName, fileName)
			continue
		}

		content, err := g.templateEngine.RenderTemplate(templateName, data)
		if err != nil {
			return fmt.Errorf("生成文件 %s 失败: %v", fileName, err)
		}
		result.Files[fileName] = content
	}

	return nil
}

// PreviewCode 预览所有模板的代码
func (g *CodeGenerator) PreviewCode(config *model.GenTableConfig) (map[string]string, error) {
	templateData := g.templateEngine.PrepareTemplateData(config)

	// 渲染所有可用模板
	return g.templateEngine.RenderAllTemplates(templateData)
}

// GetAvailableTemplates 获取可用的模板列表
func (g *CodeGenerator) GetAvailableTemplates() []string {
	var templates []string
	for name := range g.templateEngine.templates {
		templates = append(templates, name)
	}
	return templates
}

// ValidateConfig 验证配置
func (g *CodeGenerator) ValidateConfig(config *model.GenTableConfig) error {
	if config.TableName == "" {
		return fmt.Errorf("表名不能为空")
	}
	if config.BusinessName == "" {
		return fmt.Errorf("业务名称不能为空")
	}
	if config.ModuleName == "" {
		return fmt.Errorf("模块名称不能为空")
	}
	if config.ClassName == "" {
		return fmt.Errorf("类名不能为空")
	}
	if len(config.Columns) == 0 {
		return fmt.Errorf("字段配置不能为空")
	}

	// 验证字段配置
	hasPk := false
	for _, col := range config.Columns {
		if col.ColumnName == "" {
			return fmt.Errorf("字段名不能为空")
		}
		if col.GoType == "" {
			return fmt.Errorf("Go类型不能为空")
		}
		if col.IsPk {
			hasPk = true
		}
	}

	if !hasPk {
		return fmt.Errorf("必须指定主键字段")
	}

	return nil
}

// GenerateOptions 生成选项
type GenerateOptions struct {
	GenerateBackend  bool   `json:"generateBackend"`  // 生成后端代码
	GenerateFrontend bool   `json:"generateFrontend"` // 生成前端代码
	GenerateSQL      bool   `json:"generateSQL"`      // 生成SQL代码
	OutputFormat     string `json:"outputFormat"`     // 输出格式: files, zip
	CustomOutputDir  string `json:"customOutputDir"`  // 自定义输出目录
}

// GenerateResult 生成结果
type GenerateResult struct {
	TableName    string            `json:"tableName"`    // 表名
	BusinessName string            `json:"businessName"` // 业务名称
	Files        map[string]string `json:"files"`        // 生成的文件内容
	FileCount    int               `json:"fileCount"`    // 文件数量
	TotalSize    int64             `json:"totalSize"`    // 总大小（字节）
	Success      bool              `json:"success"`      // 是否成功
	ErrorMessage string            `json:"errorMessage"` // 错误信息
	StartTime    time.Time         `json:"startTime"`    // 开始时间
	EndTime      time.Time         `json:"endTime"`      // 结束时间
	Duration     time.Duration     `json:"duration"`     // 耗时
}

// DefaultGenerateOptions 默认生成选项
func DefaultGenerateOptions() *GenerateOptions {
	return &GenerateOptions{
		GenerateBackend:  true,
		GenerateFrontend: true,
		GenerateSQL:      true,
		OutputFormat:     "zip",
	}
}
