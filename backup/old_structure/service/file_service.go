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
	"github.com/LiteMove/light-stack/internal/storage"
)

// FileService 文件服务
type FileService struct {
	fileRepo      *repository.FileRepository
	tenantService TenantService
}

// NewFileService 创建文件服务实例
func NewFileService(fileRepo *repository.FileRepository, tenantService TenantService) *FileService {
	return &FileService{
		fileRepo:      fileRepo,
		tenantService: tenantService,
	}
}

// UploadFile 上传文件（支持新的存储架构）
func (s *FileService) UploadFile(file *multipart.FileHeader, userID, tenantID uint64, usageType string, isPublic bool) (*model.File, error) {
	// 获取租户的存储配置
	tenant, err := s.tenantService.GetTenant(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	storageConfig, err := tenant.GetFileStorageConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage config: %w", err)
	}

	// 如果没有明确指定isPublic，则使用租户配置的默认值
	// 特殊处理：某些类型的文件应该强制为公开访问
	if usageType == "system-logo" || usageType == "avatar" {
		isPublic = true // 系统logo和头像强制公开访问
	} else if !isPublic {
		isPublic = storageConfig.DefaultPublic
	}

	// 验证文件大小限制
	if file.Size > storageConfig.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds limit: %d > %d", file.Size, storageConfig.MaxFileSize)
	}

	// 验证文件类型
	fileExt := s.getFileExtension(file.Filename)
	if !s.isAllowedFileType(fileExt, storageConfig.AllowedTypes) {
		return nil, fmt.Errorf("file type not allowed: %s", fileExt)
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// 计算文件MD5
	md5Hash, err := s.calculateMD5(src)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate MD5: %w", err)
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

	// 生成存储路径
	dateDir := time.Now().Format("2006/01/02")
	storagePath := storage.GenerateStoragePath(tenantID, dateDir, fileName, isPublic)

	// 创建存储管理器
	storageManager, err := storage.NewManager(storageConfig)
	if err != nil {
		// 提供更友好的错误信息
		if strings.Contains(err.Error(), "租户本地访问域名配置不能为空") {
			return nil, fmt.Errorf("租户配置错误：请在租户配置中设置本地访问域名(LocalAccessDomain)，例如：http://127.0.0.1:8080")
		}
		return nil, fmt.Errorf("存储管理器初始化失败: %w", err)
	}

	// 上传文件到存储系统
	accessURL, err := storageManager.Upload(src, storagePath, isPublic)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// 创建文件记录
	fileModel := &model.File{
		TenantBaseModel: model.TenantBaseModel{
			TenantID: tenantID,
		},
		OriginalName: file.Filename,
		FileName:     fileName,
		FilePath:     storagePath,
		FileSize:     file.Size,
		FileType:     fileExt,
		MimeType:     s.getMimeType(file.Header.Get("Content-Type"), file.Filename),
		MD5:          md5Hash,
		UploadUserID: userID,
		UsageType:    usageType,
		StorageType:  storageConfig.Type,
		IsPublic:     isPublic,
		AccessURL:    accessURL,
	}

	// 保存到数据库
	if err := s.fileRepo.Create(fileModel); err != nil {
		// 删除已上传的文件
		storageManager.Delete(storagePath)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	return fileModel, nil
}

// GetFileByID 根据ID获取文件
func (s *FileService) GetFileByID(id uint64) (*model.File, error) {
	return s.fileRepo.GetByID(id)
}

// DeleteFile 删除文件（支持新的存储架构）
func (s *FileService) DeleteFile(id uint64) error {
	// 获取文件信息
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// 获取租户的存储配置
	tenant, err := s.tenantService.GetTenant(file.TenantID)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	storageConfig, err := tenant.GetFileStorageConfig()
	if err != nil {
		return fmt.Errorf("failed to get storage config: %w", err)
	}

	// 创建存储管理器
	storageManager, err := storage.NewManager(storageConfig)
	if err != nil {
		return fmt.Errorf("failed to create storage manager: %w", err)
	}

	// 删除物理文件
	if err := storageManager.Delete(file.FilePath); err != nil {
		// 记录错误但不阻止数据库删除
		fmt.Printf("Warning: failed to delete physical file %s: %v\n", file.FilePath, err)
	}

	// 删除数据库记录
	if err := s.fileRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
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

// GetPrivateFileContent 获取私有文件内容（带权限验证）
func (s *FileService) GetPrivateFileContent(fileID, userID, tenantID uint64) (*model.File, []byte, error) {
	// 获取文件信息
	file, err := s.GetFileByID(fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("file not found")
	}

	// 验证文件是否属于指定租户
	if file.TenantID != tenantID {
		return nil, nil, fmt.Errorf("access denied")
	}

	// 验证文件是否为私有文件
	if file.IsPublic {
		return nil, nil, fmt.Errorf("public file should be accessed directly")
	}

	// 验证用户权限（基本权限验证，用户必须属于同一租户）
	// 这里可以添加更复杂的权限逻辑，比如检查用户角色等
	// 目前简化为：只要是同一租户的用户就可以访问私有文件

	// 读取文件内容
	fileContent, err := s.readFileContent(file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return file, fileContent, nil
}

// readFileContent 从存储中读取文件内容
func (s *FileService) readFileContent(file *model.File) ([]byte, error) {
	// 获取租户的存储配置
	tenant, err := s.tenantService.GetTenant(file.TenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	storageConfig, err := tenant.GetFileStorageConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get storage config: %w", err)
	}

	// 创建存储管理器
	storageManager, err := storage.NewManager(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage manager: %w", err)
	}

	// 获取文件完整路径
	fullPath := storageManager.GetFullPath(file.FilePath, file.IsPublic)

	// 读取文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return content, nil
}

// isAllowedFileType 检查文件类型是否允许
func (s *FileService) isAllowedFileType(fileExt string, allowedTypes []string) bool {
	if len(allowedTypes) == 0 {
		return true // 如果没有限制，则允许所有类型
	}

	ext := "." + strings.ToLower(fileExt)
	for _, allowedType := range allowedTypes {
		if strings.ToLower(allowedType) == ext {
			return true
		}
	}
	return false
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
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
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
	case ".zip":
		return "application/zip"
	case ".rar":
		return "application/x-rar-compressed"
	default:
		return "application/octet-stream"
	}
}
