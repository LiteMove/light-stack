package model

import (
	"time"
)

// File 文件模型
type File struct {
	TenantBaseModel
	OriginalName string `json:"original_name" gorm:"not null;size:255" validate:"required,max=255"`
	FileName     string `json:"file_name" gorm:"not null;size:255" validate:"required,max=255"`
	FilePath     string `json:"file_path" gorm:"not null;size:500" validate:"required,max=500"`
	FileSize     int64  `json:"file_size" gorm:"not null" validate:"required,min=0"`
	FileType     string `json:"file_type" gorm:"not null;size:100" validate:"required,max=100"`
	MimeType     string `json:"mime_type" gorm:"not null;size:100" validate:"required,max=100"`
	MD5          string `json:"md5" gorm:"not null;size:32;index:idx_md5" validate:"required,len=32"`
	UploadUserID uint64 `json:"upload_user_id" gorm:"not null;index:idx_upload_user_id" validate:"required"`
	UsageType    string `json:"usage_type" gorm:"size:50;index:idx_usage_type" validate:"max=50"`

	// 关联关系
	UploadUser *User `json:"upload_user,omitempty" gorm:"foreignKey:UploadUserID"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// FileProfile 文件资料（简化版本）
type FileProfile struct {
	ID           uint64    `json:"id"`
	TenantID     uint64    `json:"tenant_id"`
	OriginalName string    `json:"original_name"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	MimeType     string    `json:"mime_type"`
	MD5          string    `json:"md5"`
	UploadUserID uint64    `json:"upload_user_id"`
	UsageType    string    `json:"usage_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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
