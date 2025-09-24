package storage

import (
	"fmt"
	sysConfig "github.com/LiteMove/light-stack/internal/config"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/LiteMove/light-stack/internal/model"
)

// LocalProvider 本地存储提供者
type LocalProvider struct {
	storageConfig *model.FileStorageConfig
}

// NewLocalProvider 创建本地存储提供者
func NewLocalProvider(config *model.FileStorageConfig) (*LocalProvider, error) {
	// 验证租户域名配置（必填）
	if config.LocalAccessDomain == "" {
		return nil, fmt.Errorf("租户本地访问域名配置不能为空，请在租户配置中设置 LocalAccessDomain")
	}

	// 确保基础目录存在
	sysConfig := sysConfig.Get()
	basePath := sysConfig.File.LocalPath
	if basePath == "" {
		basePath = "uploads"
	}

	// 创建public和private目录
	publicDir := filepath.Join(basePath, "public")
	privateDir := filepath.Join(basePath, "private")

	if err := os.MkdirAll(publicDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create public directory: %w", err)
	}

	if err := os.MkdirAll(privateDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create private directory: %w", err)
	}

	return &LocalProvider{
		storageConfig: config,
	}, nil
}

// Upload 上传文件到本地存储
func (p *LocalProvider) Upload(file io.Reader, path string, isPublic bool) (string, error) {
	fullPath := p.GetFullPath(path, isPublic)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", fullPath, err)
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %w", err)
	}

	// 生成访问URL
	accessURL := p.GetURL(path, isPublic)
	return accessURL, nil
}

// GetURL 获取文件访问URL
func (p *LocalProvider) GetURL(path string, isPublic bool) string {
	// 获取系统配置的静态路由路径
	sysConfig := sysConfig.Get()
	staticPath := sysConfig.File.BaseURL
	if staticPath == "" {
		staticPath = "/api/static" // 默认静态路由路径
	}

	// 租户域名 + 系统静态路径 (租户域名已在创建时验证为必填)
	baseURL := p.storageConfig.LocalAccessDomain + staticPath

	// 将Windows路径分隔符转换为URL路径分隔符
	urlPath := strings.ReplaceAll(path, "\\", "/")

	// 确保baseURL和urlPath之间有正确的分隔符
	if strings.HasSuffix(baseURL, "/") {
		return fmt.Sprintf("%s%s", baseURL, urlPath)
	} else {
		return fmt.Sprintf("%s/%s", baseURL, urlPath)
	}
}

// Delete 删除文件
func (p *LocalProvider) Delete(path string) error {
	// path已经包含完整的路径信息，直接删除
	fullPath := p.GetFullPath(path, false) // isPublic参数已不影响结果

	// 检查文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", fullPath)
	}

	// 删除文件
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", fullPath, err)
	}

	return nil
}

// GetFullPath 获取文件的完整本地路径
func (p *LocalProvider) GetFullPath(path string, isPublic bool) string {
	// 确保基础目录存在
	sysConfig := sysConfig.Get()
	basePath := sysConfig.File.LocalPath
	if basePath == "" {
		basePath = "uploads"
	}

	// path已经包含了public/private信息，直接拼接基础路径
	return filepath.Join(basePath, path)
}
