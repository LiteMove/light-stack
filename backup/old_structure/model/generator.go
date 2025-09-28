package model

import (
	"encoding/json"
	"time"
)

// GenTableConfig 代码生成表配置
type GenTableConfig struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:配置ID"`
	TableName    string    `json:"tableName" gorm:"size:64;not null;uniqueIndex:uk_table_name;comment:表名称"`
	TableComment string    `json:"tableComment" gorm:"size:255;default:'';comment:表描述"`
	BusinessName string    `json:"businessName" gorm:"size:64;not null;index:idx_business_name;comment:业务名称"`
	ModuleName   string    `json:"moduleName" gorm:"size:64;not null;index:idx_module_name;comment:模块名称"`
	FunctionName string    `json:"functionName" gorm:"size:64;not null;comment:功能名称"`
	ClassName    string    `json:"className" gorm:"size:64;not null;comment:类名"`
	PackageName  string    `json:"packageName" gorm:"size:64;not null;comment:包名"`
	Author       string    `json:"author" gorm:"size:64;default:system;comment:作者"`
	ParentMenuID *int64    `json:"parentMenuId" gorm:"comment:父级菜单ID"`
	MenuName     string    `json:"menuName" gorm:"size:64;default:'';comment:菜单名称"`
	MenuURL      string    `json:"menuUrl" gorm:"size:255;default:'';comment:菜单URL"`
	MenuIcon     string    `json:"menuIcon" gorm:"size:64;default:'';comment:菜单图标"`
	Permissions  string    `json:"permissions" gorm:"type:text;comment:权限字符串(JSON数组)"`
	Options      string    `json:"options" gorm:"type:text;comment:其他配置选项(JSON)"`
	Remark       string    `json:"remark" gorm:"size:500;default:'';comment:备注"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime;comment:更新时间"`
	CreatedBy    *int64    `json:"createdBy" gorm:"comment:创建人"`
	UpdatedBy    *int64    `json:"updatedBy" gorm:"comment:更新人"`

	// 关联字段配置
	Columns []GenTableColumn `json:"columns" gorm:"foreignKey:TableConfigID;constraint:OnDelete:CASCADE"`
}

// GetPermissions 获取权限列表
func (g *GenTableConfig) GetPermissions() []string {
	var permissions []string
	if g.Permissions != "" {
		json.Unmarshal([]byte(g.Permissions), &permissions)
	}
	return permissions
}

// SetPermissions 设置权限列表
func (g *GenTableConfig) SetPermissions(permissions []string) error {
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}
	g.Permissions = string(data)
	return nil
}

// GetOptions 获取选项配置
func (g *GenTableConfig) GetOptions() OptionConfig {
	var options OptionConfig
	if g.Options != "" {
		json.Unmarshal([]byte(g.Options), &options)
	}
	return options
}

// SetOptions 设置选项配置
func (g *GenTableConfig) SetOptions(options OptionConfig) error {
	data, err := json.Marshal(options)
	if err != nil {
		return err
	}
	g.Options = string(data)
	return nil
}

// GenTableColumn 代码生成字段配置
type GenTableColumn struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:字段ID"`
	TableConfigID int64     `json:"tableConfigId" gorm:"not null;index:idx_table_config_id;comment:表配置ID"`
	ColumnName    string    `json:"columnName" gorm:"size:64;not null;index:idx_column_name;comment:字段名称"`
	ColumnComment string    `json:"columnComment" gorm:"size:255;default:'';comment:字段描述"`
	ColumnType    string    `json:"columnType" gorm:"size:32;not null;comment:字段类型"`
	GoType        string    `json:"goType" gorm:"size:32;not null;comment:Go类型"`
	GoField       string    `json:"goField" gorm:"size:64;not null;comment:Go字段名"`
	IsPk          bool      `json:"isPk" gorm:"default:false;comment:是否主键"`
	IsIncrement   bool      `json:"isIncrement" gorm:"default:false;comment:是否自增"`
	IsRequired    bool      `json:"isRequired" gorm:"default:false;comment:是否必填"`
	IsInsert      bool      `json:"isInsert" gorm:"default:true;comment:是否为插入字段"`
	IsEdit        bool      `json:"isEdit" gorm:"default:true;comment:是否为编辑字段"`
	IsList        bool      `json:"isList" gorm:"default:true;comment:是否列表字段"`
	IsQuery       bool      `json:"isQuery" gorm:"default:false;comment:是否查询字段"`
	QueryType     string    `json:"queryType" gorm:"size:32;default:EQ;comment:查询方式"`
	HtmlType      string    `json:"htmlType" gorm:"size:32;default:input;comment:显示类型"`
	DictType      string    `json:"dictType" gorm:"size:64;default:'';comment:字典类型"`
	Sort          int       `json:"sort" gorm:"default:0;comment:排序"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime;comment:更新时间"`
}

// GenHistory 代码生成历史记录
type GenHistory struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:历史ID"`
	TableConfigID int64     `json:"tableConfigId" gorm:"not null;index:idx_table_config_id;comment:表配置ID"`
	TableName     string    `json:"tableName" gorm:"size:64;not null;index:idx_table_name;comment:表名称"`
	BusinessName  string    `json:"businessName" gorm:"size:64;not null;comment:业务名称"`
	GenerateType  string    `json:"generateType" gorm:"size:32;not null;comment:生成类型"`
	FileCount     int       `json:"fileCount" gorm:"default:0;comment:生成文件数量"`
	FileSize      int64     `json:"fileSize" gorm:"default:0;comment:文件大小(字节)"`
	DownloadCount int       `json:"downloadCount" gorm:"default:0;comment:下载次数"`
	Status        string    `json:"status" gorm:"size:32;not null;default:success;comment:生成状态"`
	ErrorMessage  string    `json:"errorMessage" gorm:"type:text;comment:错误信息"`
	FilePath      string    `json:"filePath" gorm:"size:500;default:'';comment:生成文件路径"`
	Remark        string    `json:"remark" gorm:"size:500;default:'';comment:备注"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime;index:idx_created_at;comment:创建时间"`
	CreatedBy     *int64    `json:"createdBy" gorm:"comment:创建人"`
}

// OptionConfig 其他配置选项
type OptionConfig struct {
	GenPath    string `json:"genPath"`    // 生成路径
	GenType    string `json:"genType"`    // 生成类型
	TplType    string `json:"tplType"`    // 模板类型
	TreeCode   string `json:"treeCode"`   // 树表编码
	TreeParent string `json:"treeParent"` // 树表父级
	TreeName   string `json:"treeName"`   // 树表名称
}

// TableInfo 数据库表信息
type TableInfo struct {
	TableName    string       `json:"tableName"`    // 表名
	TableComment string       `json:"tableComment"` // 表注释
	CreateTime   string       `json:"createTime"`   // 创建时间
	UpdateTime   string       `json:"updateTime"`   // 更新时间
	Columns      []ColumnInfo `json:"columns"`      // 字段信息
}

// ColumnInfo 数据库字段信息
type ColumnInfo struct {
	ColumnName      string `json:"columnName"`      // 字段名称
	ColumnType      string `json:"columnType"`      // 字段类型
	ColumnComment   string `json:"columnComment"`   // 字段注释
	IsNullable      string `json:"isNullable"`      // 是否可空
	ColumnDefault   string `json:"columnDefault"`   // 默认值
	ColumnKey       string `json:"columnKey"`       // 键类型
	Extra           string `json:"extra"`           // 扩展信息
	OrdinalPosition int    `json:"ordinalPosition"` // 字段位置

	// 转换后的信息
	GoType      string `json:"goType"`      // Go类型
	GoField     string `json:"goField"`     // Go字段名
	IsPk        bool   `json:"isPk"`        // 是否主键
	IsIncrement bool   `json:"isIncrement"` // 是否自增
	IsRequired  bool   `json:"isRequired"`  // 是否必填
	IsInsert    bool   `json:"isInsert"`    // 是否插入字段
	IsEdit      bool   `json:"isEdit"`      // 是否编辑字段
	IsList      bool   `json:"isList"`      // 是否列表字段
	IsQuery     bool   `json:"isQuery"`     // 是否查询字段
	QueryType   string `json:"queryType"`   // 查询类型
	HtmlType    string `json:"htmlType"`    // HTML类型
	DictType    string `json:"dictType"`    // 字典类型
}

// TemplateData 模板数据结构
type TemplateData struct {
	PackageName  string       `json:"packageName"`  // 包名
	ClassName    string       `json:"className"`    // 类名
	TableName    string       `json:"tableName"`    // 表名
	BusinessName string       `json:"businessName"` // 业务名
	ModuleName   string       `json:"moduleName"`   // 模块名
	FunctionName string       `json:"functionName"` // 功能名
	Author       string       `json:"author"`       // 作者
	Date         string       `json:"date"`         // 日期
	ParentMenuID int64        `json:"parentMenuId"` // 父级菜单ID
	MenuName     string       `json:"menuName"`     // 菜单名称
	MenuURL      string       `json:"menuUrl"`      // 菜单URL
	MenuIcon     string       `json:"menuIcon"`     // 菜单图标
	Permissions  []string     `json:"permissions"`  // 权限字符串列表
	Fields       []ColumnInfo `json:"fields"`       // 字段信息
	HasQuery     bool         `json:"hasQuery"`     // 是否有查询字段
	QueryFields  []ColumnInfo `json:"queryFields"`  // 查询字段
	ListFields   []ColumnInfo `json:"listFields"`   // 列表字段
	FormFields   []ColumnInfo `json:"formFields"`   // 表单字段
	Options      OptionConfig `json:"options"`      // 配置选项
}

// 常量定义
const (
	// 查询类型
	QueryTypeEQ      = "EQ"      // 等于
	QueryTypeNE      = "NE"      // 不等于
	QueryTypeGT      = "GT"      // 大于
	QueryTypeGTE     = "GTE"     // 大于等于
	QueryTypeLT      = "LT"      // 小于
	QueryTypeLTE     = "LTE"     // 小于等于
	QueryTypeLike    = "LIKE"    // 模糊查询
	QueryTypeBetween = "BETWEEN" // 范围查询

	// HTML类型
	HtmlTypeInput    = "input"    // 输入框
	HtmlTypeTextarea = "textarea" // 文本域
	HtmlTypeSelect   = "select"   // 下拉框
	HtmlTypeRadio    = "radio"    // 单选框
	HtmlTypeCheckbox = "checkbox" // 复选框
	HtmlTypeDatetime = "datetime" // 日期时间
	HtmlTypeUpload   = "upload"   // 文件上传

	// 生成类型
	GenerateTypeAll      = "all"      // 全部
	GenerateTypeBackend  = "backend"  // 后端
	GenerateTypeFrontend = "frontend" // 前端

	// 生成状态
	StatusSuccess    = "success"    // 成功
	StatusFailed     = "failed"     // 失败
	StatusProcessing = "processing" // 处理中
)
