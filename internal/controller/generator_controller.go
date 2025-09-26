package controller

import (
	"fmt"
	"github.com/LiteMove/light-stack/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/LiteMove/light-stack/internal/generator"
	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/gin-gonic/gin"
)

// GeneratorController 代码生成器控制器
type GeneratorController struct {
	dbService     *service.DBAnalyzerService
	configService *service.GenConfigService
	codeGenerator *generator.CodeGenerator
	filePackager  *generator.FilePackager
	menuService   service.MenuService
}

// NewGeneratorController 创建代码生成器控制器
func NewGeneratorController(
	dbService *service.DBAnalyzerService,
	configService *service.GenConfigService,
	codeGenerator *generator.CodeGenerator,
	filePackager *generator.FilePackager,
	menuService service.MenuService,
) *GeneratorController {
	return &GeneratorController{
		dbService:     dbService,
		configService: configService,
		codeGenerator: codeGenerator,
		filePackager:  filePackager,
		menuService:   menuService,
	}
}

// GetTableList 获取数据库表列表
func (c *GeneratorController) GetTableList(ctx *gin.Context) {
	tables, err := c.dbService.GetTableList()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(ctx, tables)
}

// GetTableColumns 获取表字段信息
func (c *GeneratorController) GetTableColumns(ctx *gin.Context) {
	tableName := ctx.Param("tableName")
	if tableName == "" {
		response.BadRequest(ctx, "表名不能为空")
		return
	}

	columns, err := c.dbService.GetTableColumns(tableName)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	response.Success(ctx, columns)
}

// GetTableInfo 获取表详细信息
func (c *GeneratorController) GetTableInfo(ctx *gin.Context) {
	tableName := ctx.Param("tableName")
	if tableName == "" {
		response.BadRequest(ctx, "表名不能为空")
		return
	}

	tableInfo, err := c.dbService.GetTableInfo(tableName)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, tableInfo)
}

// GenerateCode 生成代码
func (c *GeneratorController) GenerateCode(ctx *gin.Context) {
	var req GenerateCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	// 获取配置
	config, err := c.configService.GetConfig(req.ConfigID)
	if err != nil {
		response.BadRequest(ctx, "获取配置失败: "+err.Error())
		return
	}

	// 验证配置
	if err := c.codeGenerator.ValidateConfig(config); err != nil {
		response.BadRequest(ctx, "配置验证失败: "+err.Error())
		return
	}

	// 生成代码
	options := &generator.GenerateOptions{
		GenerateBackend:  req.GenerateBackend,
		GenerateFrontend: req.GenerateFrontend,
		GenerateSQL:      req.GenerateSQL,
		OutputFormat:     req.OutputFormat,
	}

	result, err := c.codeGenerator.GenerateCode(config, options)
	if err != nil {
		response.BadRequest(ctx, "代码生成失败: "+err.Error())
		return
	}

	// 创建历史记录
	history := &model.GenHistory{
		TableConfigID: config.ID,
		TableName:     config.TableName,
		BusinessName:  config.BusinessName,
		GenerateType:  c.getGenerateType(options),
		FileCount:     result.FileCount,
		FileSize:      result.TotalSize,
		Status:        model.StatusSuccess,
	}

	// 获取当前用户ID
	if userID, exists := ctx.Get("userID"); exists {
		if uid, ok := userID.(int64); ok {
			history.CreatedBy = &uid
		}
	}

	// 如果需要打包文件
	if req.OutputFormat == "zip" {
		zipFileName := fmt.Sprintf("%s_%s_%s.zip",
			config.BusinessName,
			time.Now().Format("20060102150405"),
			fmt.Sprintf("%d", time.Now().UnixNano())[10:],
		)

		zipPath, err := c.filePackager.PackageToZip(result, zipFileName)
		if err != nil {
			response.BadRequest(ctx, "打包文件失败: "+err.Error())
			return
		}

		history.FilePath = zipPath
	}

	// 保存历史记录
	if _, err := c.configService.CreateHistory(history); err != nil {
		// 历史记录保存失败不影响主流程，只记录日志
		fmt.Printf("保存生成历史失败: %v\n", err)
	}

	response.Success(ctx, gin.H{
		"taskId":    result.TableName + "_" + strconv.FormatInt(time.Now().Unix(), 10),
		"fileCount": result.FileCount,
		"totalSize": result.TotalSize,
		"duration":  result.Duration.String(),
		"files":     result.Files,
	})
}

// PreviewCode 预览代码
func (c *GeneratorController) PreviewCode(ctx *gin.Context) {
	configIDStr := ctx.Param("configId")
	configID, err := strconv.ParseInt(configIDStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的配置ID")
		return
	}

	templateName := ctx.Query("template")
	if templateName == "" {
		templateName = "model" // 默认预览model模板
	}

	// 获取配置
	config, err := c.configService.GetConfig(configID)
	if err != nil {
		response.BadRequest(ctx, "获取配置失败: "+err.Error())
		return
	}

	// 预览代码
	content, err := c.codeGenerator.PreviewCode(config, templateName)
	if err != nil {
		response.BadRequest(ctx, "预览代码失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"template": templateName,
		"content":  content,
	})
}

// DownloadCode 下载代码
func (c *GeneratorController) DownloadCode(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	if taskID == "" {
		response.BadRequest(ctx, "任务ID不能为空")
		return
	}

	// 这里应该根据taskID查找对应的ZIP文件
	// 简化实现，直接从query参数获取文件路径
	filePath := ctx.Query("filePath")
	if filePath == "" {
		response.BadRequest(ctx, "文件路径不能为空")
		return
	}

	// 设置下载头
	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", taskID+".zip"))

	// 发送文件
	ctx.File(filePath)
}

// GetSystemMenus 获取系统现有菜单树
func (c *GeneratorController) GetSystemMenus(ctx *gin.Context) {
	// 复用菜单管理的代码获取菜单树
	tree, err := c.menuService.GetMenuTree()
	if err != nil {
		response.BadRequest(ctx, "获取菜单树失败: "+err.Error())
		return
	}

	response.Success(ctx, tree)
}

// GetHistory 获取生成历史记录
func (c *GeneratorController) GetHistory(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")
	tableName := ctx.Query("tableName")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	histories, total, err := c.configService.GetHistoryList(page, size, tableName)
	if err != nil {
		response.BadRequest(ctx, "获取生成历史失败: "+err.Error())
		return
	}

	response.SuccessWithPage(ctx, histories, total, page, size)
}

// GetAvailableTemplates 获取可用模板列表
func (c *GeneratorController) GetAvailableTemplates(ctx *gin.Context) {
	templates := c.codeGenerator.GetAvailableTemplates()

	response.Success(ctx, templates)
}

// getGenerateType 获取生成类型
func (c *GeneratorController) getGenerateType(options *generator.GenerateOptions) string {
	if options.GenerateBackend && options.GenerateFrontend && options.GenerateSQL {
		return model.GenerateTypeAll
	} else if options.GenerateBackend {
		return model.GenerateTypeBackend
	} else if options.GenerateFrontend {
		return model.GenerateTypeFrontend
	}
	return model.GenerateTypeAll
}

// 请求结构体定义

// GenerateCodeRequest 生成代码请求
type GenerateCodeRequest struct {
	ConfigID         int64  `json:"configId" binding:"required"` // 配置ID
	GenerateBackend  bool   `json:"generateBackend"`             // 生成后端代码
	GenerateFrontend bool   `json:"generateFrontend"`            // 生成前端代码
	GenerateSQL      bool   `json:"generateSQL"`                 // 生成SQL代码
	OutputFormat     string `json:"outputFormat"`                // 输出格式
}
