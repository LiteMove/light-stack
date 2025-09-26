package controller

import (
	"net/http"
	"strconv"

	"github.com/LiteMove/light-stack/internal/service"
	"github.com/gin-gonic/gin"
)

// GenConfigController 代码生成配置控制器
type GenConfigController struct {
	service *service.GenConfigService
}

// NewGenConfigController 创建代码生成配置控制器
func NewGenConfigController(service *service.GenConfigService) *GenConfigController {
	return &GenConfigController{
		service: service,
	}
}

// CreateConfig 创建配置
func (c *GenConfigController) CreateConfig(ctx *gin.Context) {
	var req service.CreateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 获取当前用户ID（从JWT中获取）
	if userID, exists := ctx.Get("userID"); exists {
		if uid, ok := userID.(int64); ok {
			req.CreatedBy = &uid
		}
	}

	config, err := c.service.CreateConfig(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "创建成功",
		"data": config,
	})
}

// UpdateConfig 更新配置
func (c *GenConfigController) UpdateConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}

	var req service.UpdateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 获取当前用户ID
	if userID, exists := ctx.Get("userID"); exists {
		if uid, ok := userID.(int64); ok {
			req.UpdatedBy = &uid
		}
	}

	config, err := c.service.UpdateConfig(id, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新成功",
		"data": config,
	})
}

// GetConfig 获取配置详情
func (c *GenConfigController) GetConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}

	config, err := c.service.GetConfig(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": config,
	})
}

// GetConfigList 获取配置列表
func (c *GenConfigController) GetConfigList(ctx *gin.Context) {
	var req service.GetConfigListRequest

	// 解析查询参数
	pageStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 {
		size = 10
	}

	req.Page = page
	req.Size = size
	req.TableName = ctx.Query("tableName")
	req.BusinessName = ctx.Query("businessName")

	configs, total, err := c.service.GetConfigList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": gin.H{
			"records": configs,
			"total":   total,
			"page":    page,
			"size":    size,
		},
	})
}

// DeleteConfig 删除配置
func (c *GenConfigController) DeleteConfig(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}

	if err := c.service.DeleteConfig(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "删除成功",
	})
}

// ImportTableConfig 导入表配置
func (c *GenConfigController) ImportTableConfig(ctx *gin.Context) {
	tableName := ctx.Param("tableName")
	if tableName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "表名不能为空"})
		return
	}

	var req service.ImportTableConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 获取当前用户ID
	if userID, exists := ctx.Get("userID"); exists {
		if uid, ok := userID.(int64); ok {
			req.CreatedBy = &uid
		}
	}

	config, err := c.service.ImportTableConfig(tableName, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "导入成功",
		"data": config,
	})
}

// GetConfigByTableName 根据表名获取配置
func (c *GenConfigController) GetConfigByTableName(ctx *gin.Context) {
	tableName := ctx.Param("tableName")
	if tableName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "表名不能为空"})
		return
	}

	config, err := c.service.GetConfigByTableName(tableName)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "获取成功",
		"data": config,
	})
}
