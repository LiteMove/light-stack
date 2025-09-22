package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/repository"
)

// FileService 文件服务
type FileService struct {
	fileRepo  *repository.FileRepository
	uploadDir string
}

// NewFileService 创建文件服务实例
func NewFileService(fileRepo *repository.FileRepository) *FileService {
	uploadDir := "uploads" // 默认上传目录
	if dir := os.Getenv("UPLOAD_DIR"); dir != "" {
		uploadDir = dir
	}

	// 确保上传目录存在
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	return &FileService{
		fileRepo:  fileRepo,
		uploadDir: uploadDir,
	}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(file *multipart.FileHeader, userID, tenantID uint64, usageType string) (*model.File, error) {
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// 计算文件MD5
	md5Hash, err := s.calculateMD5(src)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate MD5: %v", err)
	}

	// 检查文件是否已存在（基于MD5和租户）
	existingFile, err := s.fileRepo.GetByMD5AndTenant(md5Hash, tenantID)
	if err == nil && existingFile != nil {
		// 文件已存在，直接返回现有文件信息
		return existingFile, nil
	}

	// 重置文件指针到开始位置
	src.Seek(0, 0)

	// 生成唯一文件名
	fileName := s.generateFileName(file.Filename)

	// 创建目标文件路径
	dateDir := time.Now().Format("2006/01/02")
	destDir := filepath.Join(s.uploadDir, dateDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create destination directory: %v", err)
	}

	destPath := filepath.Join(destDir, fileName)

	// 创建目标文件
	dst, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// 获取文件信息
	fileInfo, err := dst.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %v", err)
	}

	// 创建文件记录
	fileModel := &model.File{
		TenantBaseModel: model.TenantBaseModel{
			TenantID: tenantID,
		},
		OriginalName: file.Filename,
		FileName:     fileName,
		FilePath:     destPath,
		FileSize:     fileInfo.Size(),
		FileType:     s.getFileExtension(file.Filename),
		MimeType:     s.getMimeType(file.Header.Get("Content-Type"), file.Filename),
		MD5:          md5Hash,
		UploadUserID: userID,
		UsageType:    usageType,
	}

	// 保存到数据库
	if err := s.fileRepo.Create(fileModel); err != nil {
		// 删除已上传的文件
		os.Remove(destPath)
		return nil, fmt.Errorf("failed to save file record: %v", err)
	}

	return fileModel, nil
}

// GetFileByID 根据ID获取文件
func (s *FileService) GetFileByID(id uint64) (*model.File, error) {
	return s.fileRepo.GetByID(id)
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(id uint64) error {
	// 获取文件信息
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("file not found: %v", err)
	}

	// 删除物理文件
	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete physical file: %v", err)
	}

	// 删除数据库记录
	if err := s.fileRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete file record: %v", err)
	}

	return nil
}

// GetFilesByUser 获取用户上传的文件列表
func (s *FileService) GetFilesByUser(userID, tenantID uint64, page, pageSize int) ([]*model.File, int64, error) {
	offset := (page - 1) * pageSize
	return s.fileRepo.GetFilesByUser(userID, tenantID, offset, pageSize)
}

// GetAllFiles 获取所有文件列表（管理员功能）
func (s *FileService) GetAllFiles(tenantID uint64, page, pageSize int, filters map[string]interface{}) ([]*model.File, int64, error) {
	offset := (page - 1) * pageSize
	return s.fileRepo.GetAllFiles(tenantID, offset, pageSize, filters)
}

// calculateMD5 计算文件MD5
func (s *FileService) calculateMD5(file io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// generateFileName 生成唯一文件名
func (s *FileService) generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d%s", timestamp, ext)
}

// getFileExtension 获取文件扩展名
func (s *FileService) getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 0 {
		return ext[1:] // 去掉点号
	}
	return ""
}

// getMimeType 获取MIME类型
func (s *FileService) getMimeType(contentType, filename string) string {
	if contentType != "" {
		return contentType
	}

	// 根据文件扩展名推断MIME类型
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}
