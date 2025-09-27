package storage

import (
	"fmt"
	"io"

	"github.com/LiteMove/light-stack/internal/model"
)

// Provider 存储提供者接口
type Provider interface {
	Upload(file io.Reader, path string, isPublic bool) (string, error)
	GetURL(path string, isPublic bool) string
	Delete(path string) error
	GetFullPath(path string, isPublic bool) string
}

// Manager 存储管理器
type Manager struct {
	provider Provider
	config   *model.FileStorageConfig
}

// NewManager 创建新的存储管理器
func NewManager(config *model.FileStorageConfig) (*Manager, error) {
	var provider Provider
	var err error

	switch config.Type {
	case "local":
		provider, err = NewLocalProvider(config)
	case "oss":
		// 根据OSS提供商创建对应的Provider
		switch config.OSSProvider {
		case "aliyun":
			provider, err = NewAliyunOSSProvider(config)
		case "tencent":
			// TODO: 实现腾讯云COS Provider
			return nil, fmt.Errorf("tencent COS provider not implemented yet")
		case "aws":
			// TODO: 实现AWS S3 Provider
			return nil, fmt.Errorf("AWS S3 provider not implemented yet")
		case "qiniu":
			// TODO: 实现七牛云Kodo Provider
			return nil, fmt.Errorf("qiniu Kodo provider not implemented yet")
		case "upyun":
			// TODO: 实现又拍云USS Provider
			return nil, fmt.Errorf("upyun USS provider not implemented yet")
		default:
			return nil, fmt.Errorf("unsupported OSS provider: %s", config.OSSProvider)
		}
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config.Type)
	}

	if err != nil {
		return nil, err
	}

	return &Manager{
		provider: provider,
		config:   config,
	}, nil
}

// Upload 上传文件
func (m *Manager) Upload(file io.Reader, path string, isPublic bool) (string, error) {
	return m.provider.Upload(file, path, isPublic)
}

// GetURL 获取访问URL
func (m *Manager) GetURL(path string, isPublic bool) string {
	return m.provider.GetURL(path, isPublic)
}

// Delete 删除文件
func (m *Manager) Delete(path string) error {
	return m.provider.Delete(path)
}

// GetFullPath 获取完整路径
func (m *Manager) GetFullPath(path string, isPublic bool) string {
	return m.provider.GetFullPath(path, isPublic)
}

// GenerateStoragePath 生成存储路径
func GenerateStoragePath(tenantID uint64, dateDir, filename string, isPublic bool) string {
	accessType := "private"
	if isPublic {
		accessType = "public"
	}
	// 使用斜杠作为路径分隔符，确保URL路径正确
	return fmt.Sprintf("%s/tenant_%d/%s/%s", accessType, tenantID, dateDir, filename)
}
