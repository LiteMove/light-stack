package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"
)

// {{.ClassName}}Controller {{.FunctionName}}控制器
type {{.ClassName}}Controller struct {
	service *service.{{.ClassName}}Service
}

// New{{.ClassName}}Controller 创建{{.FunctionName}}控制器
func New{{.ClassName}}Controller(service *service.{{.ClassName}}Service) *{{.ClassName}}Controller {
	return &{{.ClassName}}Controller{
		service: service,
	}
}

// Create 创建{{.FunctionName}}
// @Summary 创建{{.FunctionName}}
// @Description 创建新的{{.FunctionName}}
// @Tags {{.FunctionName}}管理
// @Accept json
// @Produce json
// @Param {{uncapitalize .BusinessName}} body model.{{.ClassName}} true "{{.FunctionName}}信息"
// @Success 200 {object} response.Response
// @Router /api/{{toKebabCase .BusinessName}} [post]
func (c *{{.ClassName}}Controller) Create(ctx *gin.Context) {
	var {{uncapitalize .BusinessName}} model.{{.ClassName}}

	if err := ctx.ShouldBindJSON(&{{uncapitalize .BusinessName}}); err != nil {
		response.Error(ctx, response.CodeInvalidParam, "参数错误: "+err.Error())
		return
	}

	if err := c.service.Create(&{{uncapitalize .BusinessName}}); err != nil {
		response.Error(ctx, response.CodeServerError, "创建{{.FunctionName}}失败: "+err.Error())
		return
	}

	response.Success(ctx, {{uncapitalize .BusinessName}})
}

// GetByID 根据ID获取{{.FunctionName}}
// @Summary 根据ID获取{{.FunctionName}}
// @Description 根据ID获取{{.FunctionName}}详情
// @Tags {{.FunctionName}}管理
// @Accept json
// @Produce json
// @Param id path int true "{{.FunctionName}}ID"
// @Success 200 {object} response.Response{data=model.{{.ClassName}}}
// @Router /api/{{toKebabCase .BusinessName}}/{id} [get]
func (c *{{.ClassName}}Controller) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, response.CodeInvalidParam, "ID格式错误")
		return
	}

	{{uncapitalize .BusinessName}}, err := c.service.GetByID(id)
	if err != nil {
		response.Error(ctx, response.CodeNotFound, err.Error())
		return
	}

	response.Success(ctx, {{uncapitalize .BusinessName}})
}

// Update 更新{{.FunctionName}}
// @Summary 更新{{.FunctionName}}
// @Description 更新{{.FunctionName}}信息
// @Tags {{.FunctionName}}管理
// @Accept json
// @Produce json
// @Param id path int true "{{.FunctionName}}ID"
// @Param {{uncapitalize .BusinessName}} body model.{{.ClassName}} true "{{.FunctionName}}信息"
// @Success 200 {object} response.Response
// @Router /api/{{toKebabCase .BusinessName}}/{id} [put]
func (c *{{.ClassName}}Controller) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, response.CodeInvalidParam, "ID格式错误")
		return
	}

	var {{uncapitalize .BusinessName}} model.{{.ClassName}}
	if err := ctx.ShouldBindJSON(&{{uncapitalize .BusinessName}}); err != nil {
		response.Error(ctx, response.CodeInvalidParam, "参数错误: "+err.Error())
		return
	}

	{{uncapitalize .BusinessName}}.ID = id
	if err := c.service.Update(&{{uncapitalize .BusinessName}}); err != nil {
		response.Error(ctx, response.CodeServerError, "更新{{.FunctionName}}失败: "+err.Error())
		return
	}

	response.Success(ctx, {{uncapitalize .BusinessName}})
}

// Delete 删除{{.FunctionName}}
// @Summary 删除{{.FunctionName}}
// @Description 删除{{.FunctionName}}
// @Tags {{.FunctionName}}管理
// @Accept json
// @Produce json
// @Param id path int true "{{.FunctionName}}ID"
// @Success 200 {object} response.Response
// @Router /api/{{toKebabCase .BusinessName}}/{id} [delete]
func (c *{{.ClassName}}Controller) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, response.CodeInvalidParam, "ID格式错误")
		return
	}

	if err := c.service.Delete(id); err != nil {
		response.Error(ctx, response.CodeServerError, "删除{{.FunctionName}}失败: "+err.Error())
		return
	}

	response.Success(ctx, nil)
}

// GetList 获取{{.FunctionName}}列表
// @Summary 获取{{.FunctionName}}列表
// @Description 获取{{.FunctionName}}列表，支持分页和查询
// @Tags {{.FunctionName}}管理
// @Accept json
// @Produce json
{{- if .HasQuery }}
// @Param query query model.{{.ClassName}}Query false "查询参数"
{{- else }}
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
{{- end }}
// @Success 200 {object} response.PageResponse{data=[]model.{{.ClassName}}}
// @Router /api/{{toKebabCase .BusinessName}} [get]
func (c *{{.ClassName}}Controller) GetList(ctx *gin.Context) {
{{- if .HasQuery }}
	var query model.{{.ClassName}}Query

	if err := ctx.ShouldBindQuery(&query); err != nil {
		response.Error(ctx, response.CodeInvalidParam, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页参数
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	{{pluralize (uncapitalize .BusinessName)}}, total, err := c.service.GetList(&query)
{{- else }}
	page := 1
	pageSize := 10

	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if sizeStr := ctx.Query("pageSize"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			pageSize = s
		}
	}

	{{pluralize (uncapitalize .BusinessName)}}, total, err := c.service.GetList(page, pageSize)
{{- end }}
	if err != nil {
		response.Error(ctx, response.CodeServerError, "查询{{.FunctionName}}列表失败: "+err.Error())
		return
	}

	response.PageSuccess(ctx, {{pluralize (uncapitalize .BusinessName)}}, total)
}