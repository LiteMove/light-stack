package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

// Tenant 租户模型
type Tenant struct {
	BaseModel
	Name      string     `json:"name" gorm:"not null;size:100" validate:"required,min=1,max=100"`
	Domain    string     `json:"domain" gorm:"size:100;uniqueIndex:uk_domain" validate:"max=100"`
	Status    int        `json:"status" gorm:"not null;default:1;index" validate:"required,oneof=1 2 3 4"`
	ExpiredAt *time.Time `json:"expiredAt" gorm:"index"`
	Config    string     `json:"config" gorm:"type:json"`

	// 关联关系
	Users []User `json:"users,omitempty" gorm:"foreignKey:TenantID"`
}

// TableName 指定表名
func (Tenant) TableName() string {
	return "tenants"
}

// TenantProfile 租户资料（简化版本）
type TenantProfile struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Domain    string     `json:"domain"`
	Status    int        `json:"status"`
	ExpiredAt *time.Time `json:"expiredAt"`
	Config    string     `json:"config"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// ToProfile 转换为租户资料
func (t *Tenant) ToProfile() TenantProfile {
	return TenantProfile{
		ID:        t.ID,
		Name:      t.Name,
		Domain:    t.Domain,
		Status:    t.Status,
		ExpiredAt: t.ExpiredAt,
		Config:    t.Config,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// IsActive 检查租户是否为活跃状态
func (t *Tenant) IsActive() bool {
	return t.Status == 1
}

// IsExpired 检查租户是否已过期
func (t *Tenant) IsExpired() bool {
	if t.ExpiredAt == nil {
		return false
	}
	return t.ExpiredAt.Before(time.Now())
}

// BeforeCreate 创建前的钩子
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	if t.Status == 0 {
		t.Status = 1 // 默认启用
	}
	return nil
}

// SystemTenantId 系统租户ID
var SystemTenantId = uint64(1)

// TenantStatus 租户状态 1-启用 2-禁用 3-试用 4-过期
const (
	TenantStatusActive   = 1
	TenantStatusDisabled = 2
	TenantStatusTrial    = 3
	TenantStatusExpired  = 4
)

// FileStorageConfig 文件存储配置
type FileStorageConfig struct {
	Type          string   `json:"type"` // local/oss
	BaseURL       string   `json:"baseUrl"`
	DefaultPublic bool     `json:"defaultPublic"` // 默认是否公开
	MaxFileSize   int64    `json:"maxFileSize"`   // 最大文件大小(字节)
	AllowedTypes  []string `json:"allowedTypes"`  // 允许的文件类型

	// 本地存储配置
	LocalPath string `json:"localPath,omitempty"`

	// OSS配置 - 支持多平台
	OSSProvider     string `json:"ossProvider,omitempty"` // aliyun/tencent/aws/qiniu/upyun
	OSSEndpoint     string `json:"ossEndpoint,omitempty"`
	OSSRegion       string `json:"ossRegion,omitempty"` // AWS S3等需要
	OSSBucket       string `json:"ossBucket,omitempty"`
	OSSAccessKey    string `json:"ossAccessKey,omitempty"`
	OSSSecretKey    string `json:"ossSecretKey,omitempty"`
	OSSCustomDomain string `json:"ossCustomDomain,omitempty"` // 自定义域名
}

// TenantConfig 租户配置结构
type TenantConfig struct {
	FileStorage FileStorageConfig `json:"fileStorage"`
	// 系统基本信息
	SystemName  string `json:"systemName"`  // 系统名称
	Logo        string `json:"logo"`        // 系统Logo URL
	Description string `json:"description"` // 系统描述
	Copyright   string `json:"copyright"`   // 版权信息
}

// GetConfig 获取租户配置
func (t *Tenant) GetConfig() (*TenantConfig, error) {
	if t.Config == "" {
		return &TenantConfig{}, nil
	}

	var config TenantConfig
	if err := json.Unmarshal([]byte(t.Config), &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SetConfig 设置租户配置
func (t *Tenant) SetConfig(config *TenantConfig) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	t.Config = string(configBytes)
	return nil
}

// GetFileStorageConfig 获取文件存储配置（包含默认值）
func (t *Tenant) GetFileStorageConfig() (*FileStorageConfig, error) {
	config, err := t.GetConfig()
	if err != nil {
		return nil, err
	}

	// 应用默认值
	if config.FileStorage.Type == "" {
		config.FileStorage.Type = "local"
	}
	if config.FileStorage.MaxFileSize == 0 {
		config.FileStorage.MaxFileSize = 50 << 20 // 50MB
	}
	if len(config.FileStorage.AllowedTypes) == 0 {
		config.FileStorage.AllowedTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt"}
	}
	if config.FileStorage.BaseURL == "" && config.FileStorage.Type == "local" {
		config.FileStorage.BaseURL = "/static"
	}
	if config.FileStorage.LocalPath == "" {
		config.FileStorage.LocalPath = "uploads"
	}

	return &config.FileStorage, nil
}
