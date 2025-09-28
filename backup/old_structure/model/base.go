package model

import (
	"gorm.io/gorm"
	"time"
)

// BaseModel 基础模型（包含ID和基本时间字段）
type BaseModel struct {
	ID        uint64         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TenantBaseModel 租户基础模型（包含租户ID）
type TenantBaseModel struct {
	BaseModel
	TenantID uint64 `json:"tenant_id" gorm:"not null;default:0;index"`
}

// BeforeCreate 创建前的钩子（设置默认值）
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate 更新前的钩子
func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	b.UpdatedAt = time.Now()
	return nil
}
