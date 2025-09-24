package storage

import (
	"fmt"
	sysConfig "github.com/LiteMove/light-stack/internal/config"
	"io"
	"os"
	"path/filepath"

	"github.com/LiteMove/light-stack/internal/model"
)

// LocalProvider 本地存储提供者
type LocalProvider struct {
	storageConfig *model.FileStorageConfig
}

// NewLocalProvider 创建本地存储提供者
func NewLocalProvider(config *model.FileStorageConfig) (*LocalProvider, error) {
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
	sysConfig := sysConfig.Get()
	baseURL := sysConfig.File.BaseURL
	if baseURL == "" {
		baseURL = "/static"
	}

	accessType := "private"
	if isPublic {
		accessType = "public"
	}

	return fmt.Sprintf("%s/%s/%s", baseURL, accessType, path)
}

// Delete 删除文件
func (p *LocalProvider) Delete(path string) error {
	// 尝试删除public和private两个位置的文件
	publicPath := p.GetFullPath(path, true)
	privatePath := p.GetFullPath(path, false)

	// 删除public文件（如果存在）
	if _, err := os.Stat(publicPath); err == nil {
		if err := os.Remove(publicPath); err != nil {
			return fmt.Errorf("failed to delete public file %s: %w", publicPath, err)
		}
	}

	// 删除private文件（如果存在）
	if _, err := os.Stat(privatePath); err == nil {
		if err := os.Remove(privatePath); err != nil {
			return fmt.Errorf("failed to delete private file %s: %w", privatePath, err)
		}
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

	accessType := "private"
	if isPublic {
		accessType = "public"
	}

	return filepath.Join(basePath, accessType, path)
}
