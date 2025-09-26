package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/utils"
)

// GenConfigService 代码生成配置服务
type GenConfigService struct {
	repo      *repository.GenConfigRepository
	dbService *DBAnalyzerService
}

// NewGenConfigService 创建代码生成配置服务
func NewGenConfigService(repo *repository.GenConfigRepository, dbService *DBAnalyzerService) *GenConfigService {
	return &GenConfigService{
		repo:      repo,
		dbService: dbService,
	}
}

// CreateConfig 创建配置
func (s *GenConfigService) CreateConfig(req *CreateConfigRequest) (*model.GenTableConfig, error) {
	// 验证表名
	if err := s.dbService.ValidateTableName(req.TableName); err != nil {
		return nil, err
	}

	// 检查是否已存在配置，如果存在则更新
	existing, _ := s.repo.GetByTableName(req.TableName)
	if existing != nil {
		// 构建更新请求
		updateReq := &UpdateConfigRequest{
			BusinessName: req.BusinessName,
			ModuleName:   req.ModuleName,
			FunctionName: req.FunctionName,
			Author:       req.Author,
			ParentMenuID: req.ParentMenuID,
			MenuName:     req.MenuName,
			MenuURL:      req.MenuURL,
			MenuIcon:     req.MenuIcon,
			Options:      req.Options,
			Remark:       req.Remark,
			UpdatedBy:    req.CreatedBy, // 使用创建人作为更新人
		}

		// 获取表信息以更新字段配置
		tableInfo, err := s.dbService.GetTableInfo(req.TableName)
		if err != nil {
			return nil, fmt.Errorf("获取表信息失败: %v", err)
		}

		// 构建字段配置
		var columns []UpdateColumnConfigRequest
		for _, col := range tableInfo.Columns {
			column := UpdateColumnConfigRequest{
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
			columns = append(columns, column)
		}
		updateReq.Columns = columns

		// 调用更新方法
		return s.UpdateConfig(existing.ID, updateReq)
	}

	// 获取表信息
	tableInfo, err := s.dbService.GetTableInfo(req.TableName)
	if err != nil {
		return nil, fmt.Errorf("获取表信息失败: %v", err)
	}

	// 构建配置
	// 序列化权限
	permissionsJSON, err := json.Marshal(utils.GeneratePermissions(req.ModuleName, req.BusinessName))
	if err != nil {
		return nil, fmt.Errorf("序列化权限失败: %v", err)
	}

	// 序列化选项
	optionsJSON, err := json.Marshal(req.Options)
	if err != nil {
		return nil, fmt.Errorf("序列化选项失败: %v", err)
	}

	config := &model.GenTableConfig{
		TableName:    req.TableName,
		TableComment: tableInfo.TableComment,
		BusinessName: req.BusinessName,
		ModuleName:   req.ModuleName,
		FunctionName: req.FunctionName,
		ClassName:    utils.ToPascalCase(req.BusinessName),
		PackageName:  strings.ToLower(req.ModuleName),
		Author:       utils.DefaultString(req.Author, "system"),
		ParentMenuID: req.ParentMenuID,
		MenuName:     req.MenuName,
		MenuURL:      req.MenuURL,
		MenuIcon:     req.MenuIcon,
		Permissions:  string(permissionsJSON),
		Options:      string(optionsJSON),
		Remark:       req.Remark,
		CreatedBy:    req.CreatedBy,
	}

	// 构建字段配置
	var columns []model.GenTableColumn
	for i, col := range tableInfo.Columns {
		column := model.GenTableColumn{
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
			Sort:          i + 1,
		}
		columns = append(columns, column)
	}

	config.Columns = columns

	// 保存配置
	result, err := s.repo.Create(config)
	if err != nil {
		return nil, fmt.Errorf("保存配置失败: %v", err)
	}

	return result, nil
}

// UpdateConfig 更新配置
func (s *GenConfigService) UpdateConfig(id int64, req *UpdateConfigRequest) (*model.GenTableConfig, error) {
	// 获取现有配置
	config, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("配置不存在: %v", err)
	}

	// 更新基本信息
	if req.BusinessName != "" {
		config.BusinessName = req.BusinessName
		config.ClassName = utils.ToPascalCase(req.BusinessName)
		permissionsJSON, _ := json.Marshal(utils.GeneratePermissions(config.ModuleName, req.BusinessName))
		config.Permissions = string(permissionsJSON)
	}
	if req.ModuleName != "" {
		config.ModuleName = req.ModuleName
		config.PackageName = strings.ToLower(req.ModuleName)
		permissionsJSON, _ := json.Marshal(utils.GeneratePermissions(req.ModuleName, config.BusinessName))
		config.Permissions = string(permissionsJSON)
	}
	if req.FunctionName != "" {
		config.FunctionName = req.FunctionName
	}
	if req.Author != "" {
		config.Author = req.Author
	}
	if req.ParentMenuID != nil {
		config.ParentMenuID = req.ParentMenuID
	}
	if req.MenuName != "" {
		config.MenuName = req.MenuName
	}
	if req.MenuURL != "" {
		config.MenuURL = req.MenuURL
	}
	if req.MenuIcon != "" {
		config.MenuIcon = req.MenuIcon
	}
	if req.Remark != "" {
		config.Remark = req.Remark
	}
	config.UpdatedBy = req.UpdatedBy

	// 更新字段配置
	if req.Columns != nil && len(req.Columns) > 0 {
		// 删除现有字段配置
		if err := s.repo.DeleteColumns(id); err != nil {
			return nil, fmt.Errorf("删除字段配置失败: %v", err)
		}

		// 添加新的字段配置
		var columns []model.GenTableColumn
		for i, colReq := range req.Columns {
			column := model.GenTableColumn{
				TableConfigID: id,
				ColumnName:    colReq.ColumnName,
				ColumnComment: colReq.ColumnComment,
				ColumnType:    colReq.ColumnType,
				GoType:        colReq.GoType,
				GoField:       colReq.GoField,
				IsPk:          colReq.IsPk,
				IsIncrement:   colReq.IsIncrement,
				IsRequired:    colReq.IsRequired,
				IsInsert:      colReq.IsInsert,
				IsEdit:        colReq.IsEdit,
				IsList:        colReq.IsList,
				IsQuery:       colReq.IsQuery,
				QueryType:     colReq.QueryType,
				HtmlType:      colReq.HtmlType,
				DictType:      colReq.DictType,
				Sort:          i + 1,
			}
			columns = append(columns, column)
		}
		config.Columns = columns
	}

	// 保存更新
	result, err := s.repo.Update(config)
	if err != nil {
		return nil, fmt.Errorf("更新配置失败: %v", err)
	}

	return result, nil
}

// GetConfig 获取配置详情
func (s *GenConfigService) GetConfig(id int64) (*model.GenTableConfig, error) {
	config, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %v", err)
	}
	return config, nil
}

// GetConfigByTableName 根据表名获取配置
func (s *GenConfigService) GetConfigByTableName(tableName string) (*model.GenTableConfig, error) {
	config, err := s.repo.GetByTableName(tableName)
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %v", err)
	}
	return config, nil
}

// GetConfigList 获取配置列表
func (s *GenConfigService) GetConfigList(req *GetConfigListRequest) ([]*model.GenTableConfig, int64, error) {
	return s.repo.GetList(req.Page, req.Size, req.TableName, req.BusinessName)
}

// DeleteConfig 删除配置
func (s *GenConfigService) DeleteConfig(id int64) error {
	// 检查配置是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("配置不存在: %v", err)
	}

	// 删除配置（级联删除字段配置）
	return s.repo.Delete(id)
}

// ImportTableConfig 导入表配置
func (s *GenConfigService) ImportTableConfig(tableName string, req *ImportTableConfigRequest) (*model.GenTableConfig, error) {
	// 验证表名
	if err := s.dbService.ValidateTableName(tableName); err != nil {
		return nil, err
	}

	// 检查是否已存在配置
	existing, _ := s.repo.GetByTableName(tableName)
	if existing != nil {
		return nil, fmt.Errorf("表 '%s' 的配置已存在", tableName)
	}

	// 获取表信息
	tableInfo, err := s.dbService.GetTableInfo(tableName)
	if err != nil {
		return nil, fmt.Errorf("获取表信息失败: %v", err)
	}

	// 生成默认配置
	businessName := s.generateBusinessName(tableName)
	moduleName := utils.DefaultString(req.ModuleName, "system")

	createReq := &CreateConfigRequest{
		TableName:    tableName,
		BusinessName: businessName,
		ModuleName:   moduleName,
		FunctionName: tableInfo.TableComment,
		Author:       req.Author,
		ParentMenuID: req.ParentMenuID,
		MenuName:     tableInfo.TableComment,
		MenuURL:      "/" + strings.ToLower(moduleName) + "/" + strings.ToLower(businessName),
		MenuIcon:     "table",
		Options:      model.OptionConfig{}, // 空的结构体，在CreateConfig中会被序列化为JSON
		Remark:       fmt.Sprintf("从表 %s 导入的配置", tableName),
		CreatedBy:    req.CreatedBy,
	}

	return s.CreateConfig(createReq)
}

// generateBusinessName 生成业务名称
func (s *GenConfigService) generateBusinessName(tableName string) string {
	// 移除常见前缀
	prefixes := []string{"sys_", "t_", "tb_", "tbl_"}
	name := tableName
	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(name), prefix) {
			name = name[len(prefix):]
			break
		}
	}

	// 转换为驼峰命名
	return utils.Uncapitalize(utils.ToCamelCase(name))
}

// ValidateConfig 验证配置
func (s *GenConfigService) ValidateConfig(config *model.GenTableConfig) error {
	if config.TableName == "" {
		return fmt.Errorf("表名不能为空")
	}
	if config.BusinessName == "" {
		return fmt.Errorf("业务名称不能为空")
	}
	if config.ModuleName == "" {
		return fmt.Errorf("模块名称不能为空")
	}
	if config.FunctionName == "" {
		return fmt.Errorf("功能名称不能为空")
	}
	if config.ClassName == "" {
		return fmt.Errorf("类名不能为空")
	}
	if config.PackageName == "" {
		return fmt.Errorf("包名不能为空")
	}

	return nil
}

// CreateHistory 创建历史记录
func (s *GenConfigService) CreateHistory(history *model.GenHistory) (*model.GenHistory, error) {
	return s.repo.CreateHistory(history)
}

// GetHistoryList 获取历史记录列表
func (s *GenConfigService) GetHistoryList(page, size int, tableName string) ([]*model.GenHistory, int64, error) {
	return s.repo.GetHistoryList(page, size, tableName)
}

// 请求结构体定义

// CreateConfigRequest 创建配置请求
type CreateConfigRequest struct {
	TableName    string             `json:"tableName" binding:"required"`    // 表名
	BusinessName string             `json:"businessName" binding:"required"` // 业务名称
	ModuleName   string             `json:"moduleName" binding:"required"`   // 模块名称
	FunctionName string             `json:"functionName" binding:"required"` // 功能名称
	Author       string             `json:"author"`                          // 作者
	ParentMenuID *int64             `json:"parentMenuId"`                    // 父级菜单ID
	MenuName     string             `json:"menuName"`                        // 菜单名称
	MenuURL      string             `json:"menuUrl"`                         // 菜单URL
	MenuIcon     string             `json:"menuIcon"`                        // 菜单图标
	Options      model.OptionConfig `json:"options"`                         // 其他选项
	Remark       string             `json:"remark"`                          // 备注
	CreatedBy    *int64             `json:"createdBy"`                       // 创建人
}

// UpdateConfigRequest 更新配置请求
type UpdateConfigRequest struct {
	BusinessName string                      `json:"businessName"` // 业务名称
	ModuleName   string                      `json:"moduleName"`   // 模块名称
	FunctionName string                      `json:"functionName"` // 功能名称
	Author       string                      `json:"author"`       // 作者
	ParentMenuID *int64                      `json:"parentMenuId"` // 父级菜单ID
	MenuName     string                      `json:"menuName"`     // 菜单名称
	MenuURL      string                      `json:"menuUrl"`      // 菜单URL
	MenuIcon     string                      `json:"menuIcon"`     // 菜单图标
	Options      model.OptionConfig          `json:"options"`      // 其他选项
	Remark       string                      `json:"remark"`       // 备注
	UpdatedBy    *int64                      `json:"updatedBy"`    // 更新人
	Columns      []UpdateColumnConfigRequest `json:"columns"`      // 字段配置
}

// UpdateColumnConfigRequest 更新字段配置请求
type UpdateColumnConfigRequest struct {
	ColumnName    string `json:"columnName"`    // 字段名称
	ColumnComment string `json:"columnComment"` // 字段注释
	ColumnType    string `json:"columnType"`    // 字段类型
	GoType        string `json:"goType"`        // Go类型
	GoField       string `json:"goField"`       // Go字段名
	IsPk          bool   `json:"isPk"`          // 是否主键
	IsIncrement   bool   `json:"isIncrement"`   // 是否自增
	IsRequired    bool   `json:"isRequired"`    // 是否必填
	IsInsert      bool   `json:"isInsert"`      // 是否插入字段
	IsEdit        bool   `json:"isEdit"`        // 是否编辑字段
	IsList        bool   `json:"isList"`        // 是否列表字段
	IsQuery       bool   `json:"isQuery"`       // 是否查询字段
	QueryType     string `json:"queryType"`     // 查询类型
	HtmlType      string `json:"htmlType"`      // HTML类型
	DictType      string `json:"dictType"`      // 字典类型
}

// GetConfigListRequest 获取配置列表请求
type GetConfigListRequest struct {
	Page         int    `json:"page" form:"page"`                 // 页码
	Size         int    `json:"size" form:"size"`                 // 每页数量
	TableName    string `json:"tableName" form:"tableName"`       // 表名
	BusinessName string `json:"businessName" form:"businessName"` // 业务名称
}

// ImportTableConfigRequest 导入表配置请求
type ImportTableConfigRequest struct {
	ModuleName   string `json:"moduleName"`   // 模块名称
	Author       string `json:"author"`       // 作者
	ParentMenuID *int64 `json:"parentMenuId"` // 父级菜单ID
	CreatedBy    *int64 `json:"createdBy"`    // 创建人
}
