package controller

import (
	"strconv"

	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// DictController 字典控制器
type DictController struct {
	dictService service.DictService
	validator   *validator.Validate
}

// NewDictController 创建字典控制器
func NewDictController(dictService service.DictService) *DictController {
	return &DictController{
		dictService: dictService,
		validator:   validator.New(),
	}
}

// === 请求结构体 ===

// CreateDictTypeRequest 创建字典类型请求
type CreateDictTypeRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Type        string `json:"type" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
	Status      int    `json:"status" validate:"required,oneof=1 2"`
}

// UpdateDictTypeRequest 更新字典类型请求
type UpdateDictTypeRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Type        string `json:"type" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
	Status      int    `json:"status" validate:"required,oneof=1 2"`
}

// DictTypeListRequest 字典类型列表请求
type DictTypeListRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"pageSize" validate:"min=1,max=100"`
	Status   int    `form:"status" validate:"oneof=0 1 2"`
	Name     string `form:"name"`
}

// CreateDictDataRequest 创建字典数据请求
type CreateDictDataRequest struct {
	DictType  string `json:"dictType" validate:"required,max=100"`
	Label     string `json:"label" validate:"required,max=100"`
	Value     string `json:"value" validate:"required,max=100"`
	SortOrder int    `json:"sortOrder" validate:"min=0"`
	CssClass  string `json:"cssClass" validate:"max=100"`
	ListClass string `json:"listClass" validate:"max=100"`
	IsDefault bool   `json:"isDefault"`
	Status    int    `json:"status" validate:"required,oneof=1 2"`
	Remark    string `json:"remark" validate:"max=255"`
}

// UpdateDictDataRequest 更新字典数据请求
type UpdateDictDataRequest struct {
	DictType  string `json:"dictType" validate:"required,max=100"`
	Label     string `json:"label" validate:"required,max=100"`
	Value     string `json:"value" validate:"required,max=100"`
	SortOrder int    `json:"sortOrder" validate:"min=0"`
	CssClass  string `json:"cssClass" validate:"max=100"`
	ListClass string `json:"listClass" validate:"max=100"`
	IsDefault bool   `json:"isDefault"`
	Status    int    `json:"status" validate:"required,oneof=1 2"`
	Remark    string `json:"remark" validate:"max=255"`
}

// DictDataListRequest 字典数据列表请求
type DictDataListRequest struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"pageSize" validate:"min=1,max=100"`
	Status   int    `form:"status" validate:"oneof=0 1 2"`
	Label    string `form:"label"`
}

// BatchUpdateStatusRequest 批量更新状态请求
type BatchUpdateStatusRequest struct {
	IDs    []uint64 `json:"ids" validate:"required,min=1"`
	Status int      `json:"status" validate:"required,oneof=1 2"`
}

// === 字典类型相关接口 ===

// CreateType 创建字典类型
func (c *DictController) CreateType(ctx *gin.Context) {
	var req CreateDictTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 创建字典类型对象
	dictType := &model.DictType{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Status:      req.Status,
	}

	// 调用服务创建字典类型
	if err := c.dictService.CreateType(dictType); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, dictType.ToProfile())
}

// UpdateType 更新字典类型
func (c *DictController) UpdateType(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的字典类型ID")
		return
	}

	var req UpdateDictTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 获取现有记录
	dictType, err := c.dictService.GetType(id)
	if err != nil {
		response.NotFound(ctx, "字典类型不存在")
		return
	}

	// 更新字段
	dictType.Name = req.Name
	dictType.Type = req.Type
	dictType.Description = req.Description
	dictType.Status = req.Status

	// 调用服务更新字典类型
	if err := c.dictService.UpdateType(dictType); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, dictType.ToProfile())
}

// DeleteType 删除字典类型
func (c *DictController) DeleteType(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的字典类型ID")
		return
	}

	// 调用服务删除字典类型
	if err := c.dictService.DeleteType(id); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// GetType 获取字典类型详情
func (c *DictController) GetType(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的字典类型ID")
		return
	}

	// 调用服务获取字典类型
	dictType, err := c.dictService.GetType(id)
	if err != nil {
		response.NotFound(ctx, "字典类型不存在")
		return
	}

	response.Success(ctx, dictType.ToProfile())
}

// GetTypeList 获取字典类型列表
func (c *DictController) GetTypeList(ctx *gin.Context) {
	var req DictTypeListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用服务获取字典类型列表
	dictTypes, total, err := c.dictService.GetTypeList(req.Page, req.PageSize, req.Status, req.Name)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	// 转换为Profile格式
	profiles := make([]model.DictTypeProfile, len(dictTypes))
	for i, dt := range dictTypes {
		profiles[i] = dt.ToProfile()
	}

	response.SuccessWithPage(ctx, profiles, total, req.Page, req.PageSize)
}

// === 字典数据相关接口 ===

// CreateData 创建字典数据
func (c *DictController) CreateData(ctx *gin.Context) {
	var req CreateDictDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 创建字典数据对象
	dictData := &model.DictData{
		DictType:  req.DictType,
		Label:     req.Label,
		Value:     req.Value,
		SortOrder: req.SortOrder,
		CssClass:  req.CssClass,
		ListClass: req.ListClass,
		IsDefault: req.IsDefault,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	// 调用服务创建字典数据
	if err := c.dictService.CreateData(dictData); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, dictData.ToProfile())
}

// UpdateData 更新字典数据
func (c *DictController) UpdateData(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的字典数据ID")
		return
	}

	var req UpdateDictDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 获取现有记录
	dictData, err := c.dictService.GetData(id)
	if err != nil {
		response.NotFound(ctx, "字典数据不存在")
		return
	}

	// 更新字段
	dictData.DictType = req.DictType
	dictData.Label = req.Label
	dictData.Value = req.Value
	dictData.SortOrder = req.SortOrder
	dictData.CssClass = req.CssClass
	dictData.ListClass = req.ListClass
	dictData.IsDefault = req.IsDefault
	dictData.Status = req.Status
	dictData.Remark = req.Remark

	// 调用服务更新字典数据
	if err := c.dictService.UpdateData(dictData); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, dictData.ToProfile())
}

// DeleteData 删除字典数据
func (c *DictController) DeleteData(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的字典数据ID")
		return
	}

	// 调用服务删除字典数据
	if err := c.dictService.DeleteData(id); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// GetData 获取字典数据详情
func (c *DictController) GetData(ctx *gin.Context) {
	// 获取ID参数
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的字典数据ID")
		return
	}

	// 调用服务获取字典数据
	dictData, err := c.dictService.GetData(id)
	if err != nil {
		response.NotFound(ctx, "字典数据不存在")
		return
	}

	response.Success(ctx, dictData.ToProfile())
}

// GetDataList 获取字典数据列表
func (c *DictController) GetDataList(ctx *gin.Context) {
	// 获取字典类型参数
	dictType := ctx.Param("type")
	if dictType == "" {
		response.BadRequest(ctx, "字典类型参数不能为空")
		return
	}

	var req DictDataListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 调用服务获取字典数据列表
	dictData, total, err := c.dictService.GetDataList(dictType, req.Page, req.PageSize, req.Status, req.Label)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	// 转换为Profile格式
	profiles := make([]model.DictDataProfile, len(dictData))
	for i, dd := range dictData {
		profiles[i] = dd.ToProfile()
	}

	response.SuccessWithPage(ctx, profiles, total, req.Page, req.PageSize)
}

// === 批量操作接口 ===

// BatchUpdateDataStatus 批量更新字典数据状态
func (c *DictController) BatchUpdateDataStatus(ctx *gin.Context) {
	var req BatchUpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 调用服务批量更新状态
	if err := c.dictService.BatchUpdateDataStatus(req.IDs, req.Status); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// BatchDeleteData 批量删除字典数据
func (c *DictController) BatchDeleteData(ctx *gin.Context) {
	var req struct {
		IDs []uint64 `json:"ids" validate:"required,min=1"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数格式错误: "+err.Error())
		return
	}

	// 参数验证
	if err := c.validator.Struct(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	// 调用服务批量删除
	if err := c.dictService.BatchDeleteData(req.IDs); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, nil)
}

// === 前端下拉框接口 ===

// GetDictOptions 获取字典选项（用于前端下拉框）
func (c *DictController) GetDictOptions(ctx *gin.Context) {
	// 获取字典类型参数
	dictType := ctx.Param("type")
	if dictType == "" {
		response.BadRequest(ctx, "字典类型参数不能为空")
		return
	}

	// 调用服务获取字典选项
	options, err := c.dictService.GetDictOptions(dictType)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, options)
}
