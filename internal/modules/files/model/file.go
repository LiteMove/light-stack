package model

import (
	"time"

	systemModel "github.com/LiteMove/light-stack/internal/modules/system/model"
	"github.com/LiteMove/light-stack/internal/shared/model"
)

// File 文件模型
type File struct {
	model.TenantBaseModel
	OriginalName string `json:"originalName" gorm:"not null;size:255" validate:"required,max=255"`
	FileName     string `json:"fileName" gorm:"not null;size:255" validate:"required,max=255"`
	FilePath     string `json:"filePath" gorm:"not null;size:500" validate:"required,max=500"`
	FileSize     int64  `json:"fileSize" gorm:"not null" validate:"required,min=0"`
	FileType     string `json:"fileType" gorm:"not null;size:100" validate:"required,max=100"`
	MimeType     string `json:"mimeType" gorm:"not null;size:100" validate:"required,max=100"`
	MD5          string `json:"md5" gorm:"not null;size:32;index:idx_md5" validate:"required,len=32"`
	UploadUserID uint64 `json:"uploadUserId" gorm:"not null;index:idx_upload_user_id" validate:"required"`
	UsageType    string `json:"usageType" gorm:"size:50;index:idx_usage_type" validate:"max=50"`
	StorageType  string `json:"storageType" gorm:"not null;default:'local';size:20" validate:"required,oneof=local oss"`
	IsPublic     bool   `json:"isPublic" gorm:"not null;default:false"`
	AccessURL    string `json:"accessUrl" gorm:"size:1000"`

	// 关联关系
	UploadUser *systemModel.User `json:"upload_user,omitempty" gorm:"foreignKey:UploadUserID"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// FileProfile 文件资料（简化版本）
type FileProfile struct {
	ID           uint64    `json:"id"`
	TenantID     uint64    `json:"tenantId"`
	OriginalName string    `json:"originalName"`
	FileName     string    `json:"fileName"`
	FilePath     string    `json:"filePath"`
	FileSize     int64     `json:"fileSize"`
	FileType     string    `json:"fileType"`
	MimeType     string    `json:"mimeType"`
	MD5          string    `json:"md5"`
	UploadUserID uint64    `json:"uploadUserId"`
	UsageType    string    `json:"usageType"`
	StorageType  string    `json:"storageType"`
	IsPublic     bool      `json:"isPublic"`
	AccessURL    string    `json:"accessUrl"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// ToProfile 转换为文件资料
func (f *File) ToProfile() FileProfile {
	return FileProfile{
		ID:           f.ID,
		TenantID:     f.TenantID,
		OriginalName: f.OriginalName,
		FileName:     f.FileName,
		FilePath:     f.FilePath,
		FileSize:     f.FileSize,
		FileType:     f.FileType,
		MimeType:     f.MimeType,
		MD5:          f.MD5,
		UploadUserID: f.UploadUserID,
		UsageType:    f.UsageType,
		StorageType:  f.StorageType,
		IsPublic:     f.IsPublic,
		AccessURL:    f.AccessURL,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}
}

// GetSizeInKB 获取文件大小（KB）
func (f *File) GetSizeInKB() float64 {
	return float64(f.FileSize) / 1024
}

// GetSizeInMB 获取文件大小（MB）
func (f *File) GetSizeInMB() float64 {
	return float64(f.FileSize) / (1024 * 1024)
}

// IsImage 检查是否为图片文件
func (f *File) IsImage() bool {
	imageTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml"}
	for _, imgType := range imageTypes {
		if f.MimeType == imgType {
			return true
		}
	}
	return false
}

// IsDocument 检查是否为文档文件
func (f *File) IsDocument() bool {
	docTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"text/plain",
	}
	for _, docType := range docTypes {
		if f.MimeType == docType {
			return true
		}
	}
	return false
}
