package controller

import (
	"fmt"
	"github.com/LiteMove/light-stack/pkg/response"
	"net/http"
	"strconv"
	"time"

	"github.com/LiteMove/light-stack/internal/modules/generator/engine"
	"github.com/LiteMove/light-stack/internal/modules/generator/model"
	"github.com/LiteMove/light-stack/internal/modules/generator/service"
	systemModel "github.com/LiteMove/light-stack/internal/modules/system/model"
	"github.com/gin-gonic/gin"
)

// MenuService 菜单服务接口（跨模块依赖）
type MenuService interface {
	CreateMenu(menu *systemModel.Menu) error
	GetMenuTree() ([]systemModel.MenuTreeNode, error)
}

// GeneratorController 代码生成器控制器
type GeneratorController struct {
	dbService     *service.DBAnalyzerService
	configService *service.GenConfigService
	codeGenerator *generator.CodeGenerator
	filePackager  *generator.FilePackager
	menuService   MenuService
}

// NewGeneratorController 创建代码生成器控制器
func NewGeneratorController(
	dbService *service.DBAnalyzerService,
	configService *service.GenConfigService,
	codeGenerator *generator.CodeGenerator,
	filePackager *generator.FilePackager,
	menuService MenuService,
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

	// 总是创建ZIP文件用于下载（不管前端是否指定outputFormat）
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

	// 获取配置
	config, err := c.configService.GetConfig(configID)
	if err != nil {
		response.BadRequest(ctx, "获取配置失败: "+err.Error())
		return
	}

	// 预览所有模板的代码
	files, err := c.codeGenerator.PreviewCode(config)
	if err != nil {
		response.BadRequest(ctx, "预览代码失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"files": files,
	})
}

// DownloadCode 下载代码
func (c *GeneratorController) DownloadCode(ctx *gin.Context) {
	taskID := ctx.Param("taskId")
	if taskID == "" {
		response.BadRequest(ctx, "任务ID不能为空")
		return
	}

	// 尝试将taskId解析为历史记录ID
	historyID, err := strconv.ParseInt(taskID, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的任务ID")
		return
	}

	// 根据历史记录ID查找文件路径
	history, err := c.configService.GetHistoryByID(historyID)
	if err != nil {
		response.BadRequest(ctx, "找不到对应的生成记录: "+err.Error())
		return
	}

	if history.FilePath == "" {
		response.BadRequest(ctx, "该记录没有生成的文件包")
		return
	}

	// 检查文件是否存在
	if !c.filePackager.FileExists(history.FilePath) {
		response.BadRequest(ctx, "文件不存在或已被删除")
		return
	}

	// 设置下载头
	fileName := fmt.Sprintf("%s_code.zip", history.BusinessName)
	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	// 发送文件
	ctx.File(history.FilePath)
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

func (c *GeneratorController) DeleteHistory(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(context, "无效的历史记录ID")
		return
	}
	if err := c.configService.DeleteHistory(id); err != nil {
		response.BadRequest(context, "删除历史记录失败: "+err.Error())
		return
	}
	response.Success(context, nil)

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
