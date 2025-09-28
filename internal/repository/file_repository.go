package repository

import (
	"strings"

	"github.com/LiteMove/light-stack/internal/modules/files/model"
	"gorm.io/gorm"
)

// FileRepository 文件数据访问层
type FileRepository struct {
	db *gorm.DB
}

// NewFileRepository 创建文件数据访问层实例
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{db: db}
}

// Create 创建文件记录
func (r *FileRepository) Create(file *model.File) error {
	return r.db.Create(file).Error
}

// GetByID 根据ID获取文件
func (r *FileRepository) GetByID(id uint64) (*model.File, error) {
	var file model.File
	err := r.db.Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// GetByMD5AndTenant 根据MD5和租户ID查找文件
func (r *FileRepository) GetByMD5AndTenant(md5 string, tenantID uint64) (*model.File, error) {
	var file model.File
	err := r.db.Where("md5 = ? AND tenant_id = ?", md5, tenantID).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// Delete 删除文件记录
func (r *FileRepository) Delete(id uint64) error {
	return r.db.Delete(&model.File{}, id).Error
}

// GetFilesByUser 获取用户上传的文件列表
func (r *FileRepository) GetFilesByUser(userID uint64, tenantID uint64, offset, limit int) ([]*model.File, int64, error) {
	var files []*model.File
	var total int64

	db := r.db.Where("upload_user_id = ? AND tenant_id = ?", userID, tenantID)

	// 获取总数
	if err := db.Model(&model.File{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// GetAllFiles 获取所有文件列表（按租户）
func (r *FileRepository) GetAllFiles(tenantID uint64, offset, limit int, filters map[string]interface{}) ([]*model.File, int64, error) {
	var files []*model.File
	var total int64

	db := r.db.Where("tenant_id = ?", tenantID)

	// 应用过滤条件
	if filename, ok := filters["filename"]; ok && filename != "" {
		db = db.Where("original_name LIKE ?", "%"+filename.(string)+"%")
	}
	if fileType, ok := filters["file_type"]; ok && fileType != "" {
		// 处理多个文件类型，用逗号分隔
		fileTypeStr := fileType.(string)
		if fileTypeStr != "" {
			types := strings.Split(fileTypeStr, ",")
			if len(types) > 1 {
				db = db.Where("file_type IN ?", types)
			} else {
				db = db.Where("file_type = ?", fileTypeStr)
			}
		}
	}
	if usageType, ok := filters["usage_type"]; ok && usageType != "" {
		db = db.Where("usage_type = ?", usageType)
	}
	if uploadUserID, ok := filters["upload_user_id"]; ok && uploadUserID != "" {
		db = db.Where("upload_user_id = ?", uploadUserID)
	}

	// 获取总数
	if err := db.Model(&model.File{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，包含上传用户信息
	if err := db.Preload("UploadUser").Offset(offset).Limit(limit).Order("created_at DESC").Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// UpdateFile 更新文件信息
func (r *FileRepository) UpdateFile(file *model.File) error {
	return r.db.Save(file).Error
}

// GetTotalCount 获取文件总数（超管使用）
func (r *FileRepository) GetTotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.File{}).Count(&count).Error
	return count, err
}

// GetCountByTenantID 根据租户ID获取文件数量
func (r *FileRepository) GetCountByTenantID(tenantID uint64) (int64, error) {
	var count int64
	err := r.db.Model(&model.File{}).Where("tenant_id = ?", tenantID).Count(&count).Error
	return count, err
}
