package model

import "github.com/LiteMove/light-stack/internal/shared/model"

import (
	"gorm.io/gorm"
	"time"
)

// DictType 字典类型模型
type DictType struct {
	model.BaseModel
	Name        string `json:"name" gorm:"not null;size:100" validate:"required,min=1,max=100"`
	Type        string `json:"type" gorm:"not null;size:100;uniqueIndex:uk_type" validate:"required,min=1,max=100"`
	Description string `json:"description" gorm:"size:255" validate:"max=255"`
	Status      int    `json:"status" gorm:"not null;default:1;index" validate:"required,oneof=1 2"`

	// 关联关系
	DictData []DictData `json:"dict_data,omitempty" gorm:"foreignKey:DictType;references:Type"`
}

// TableName 指定表名
func (DictType) TableName() string {
	return "dict_types"
}

// DictData 字典数据模型
type DictData struct {
	model.BaseModel
	DictType  string `json:"dictType" gorm:"not null;size:100;uniqueIndex:uk_type_value;index:idx_dict_type" validate:"required,max=100"`
	Label     string `json:"label" gorm:"not null;size:100" validate:"required,max=100"`
	Value     string `json:"value" gorm:"not null;size:100;uniqueIndex:uk_type_value" validate:"required,max=100"`
	SortOrder int    `json:"sortOrder" gorm:"not null;default:0;index:idx_sort_order"`
	CssClass  string `json:"cssClass" gorm:"size:100" validate:"max=100"`
	ListClass string `json:"listClass" gorm:"size:100" validate:"max=100"`
	IsDefault bool   `json:"isDefault" gorm:"not null;default:false"`
	Status    int    `json:"status" gorm:"not null;default:1;index" validate:"required,oneof=1 2"`
	Remark    string `json:"remark" gorm:"size:255" validate:"max=255"`
}

// TableName 指定表名
func (DictData) TableName() string {
	return "dict_data"
}

// DictTypeProfile 字典类型资料（简化版本）
type DictTypeProfile struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ToProfile 转换为字典类型资料
func (dt *DictType) ToProfile() DictTypeProfile {
	return DictTypeProfile{
		ID:          dt.ID,
		Name:        dt.Name,
		Type:        dt.Type,
		Description: dt.Description,
		Status:      dt.Status,
		CreatedAt:   dt.CreatedAt,
		UpdatedAt:   dt.UpdatedAt,
	}
}

// DictDataProfile 字典数据资料（简化版本）
type DictDataProfile struct {
	ID        uint64    `json:"id"`
	DictType  string    `json:"dictType"`
	Label     string    `json:"label"`
	Value     string    `json:"value"`
	SortOrder int       `json:"sortOrder"`
	CssClass  string    `json:"cssClass"`
	ListClass string    `json:"listClass"`
	IsDefault bool      `json:"isDefault"`
	Status    int       `json:"status"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToProfile 转换为字典数据资料
func (dd *DictData) ToProfile() DictDataProfile {
	return DictDataProfile{
		ID:        dd.ID,
		DictType:  dd.DictType,
		Label:     dd.Label,
		Value:     dd.Value,
		SortOrder: dd.SortOrder,
		CssClass:  dd.CssClass,
		ListClass: dd.ListClass,
		IsDefault: dd.IsDefault,
		Status:    dd.Status,
		Remark:    dd.Remark,
		CreatedAt: dd.CreatedAt,
		UpdatedAt: dd.UpdatedAt,
	}
}

// BeforeCreate 创建前的钩子
func (dt *DictType) BeforeCreate(tx *gorm.DB) error {
	if dt.Status == 0 {
		dt.Status = 1 // 默认启用
	}
	return nil
}

// BeforeCreate 创建前的钩子
func (dd *DictData) BeforeCreate(tx *gorm.DB) error {
	if dd.Status == 0 {
		dd.Status = 1 // 默认启用
	}
	return nil
}
