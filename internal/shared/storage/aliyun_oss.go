package storage

import (
	"fmt"
	"io"
	"strings"

	systemModel "github.com/LiteMove/light-stack/internal/modules/system/model"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliyunOSSProvider 阿里云OSS存储提供者
type AliyunOSSProvider struct {
	config *systemModel.FileStorageConfig
	client *oss.Client
	bucket *oss.Bucket
}

// NewAliyunOSSProvider 创建阿里云OSS存储提供者
func NewAliyunOSSProvider(config *systemModel.FileStorageConfig) (*AliyunOSSProvider, error) {
	// 创建OSS客户端
	client, err := oss.New(config.OSSEndpoint, config.OSSAccessKey, config.OSSSecretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create OSS client: %w", err)
	}

	// 获取Bucket
	bucket, err := client.Bucket(config.OSSBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get OSS bucket: %w", err)
	}

	return &AliyunOSSProvider{
		config: config,
		client: client,
		bucket: bucket,
	}, nil
}

// Upload 上传文件到阿里云OSS
func (p *AliyunOSSProvider) Upload(file io.Reader, path string, isPublic bool) (string, error) {
	objectKey := p.generateObjectKey(path)

	// 设置上传选项
	options := []oss.Option{}

	// 如果是公开文件，设置为公开读
	if isPublic {
		options = append(options, oss.ACL(oss.ACLPublicRead))
	} else {
		options = append(options, oss.ACL(oss.ACLPrivate))
	}

	// 上传文件
	err := p.bucket.PutObject(objectKey, file, options...)
	if err != nil {
		return "", fmt.Errorf("failed to upload to OSS: %w", err)
	}

	// 生成访问URL
	accessURL := p.GetURL(path, isPublic)
	return accessURL, nil
}

// GetURL 获取文件访问URL
func (p *AliyunOSSProvider) GetURL(path string, isPublic bool) string {
	objectKey := p.generateObjectKey(path)

	if isPublic {
		// 公开文件直接返回自定义域名URL
		if p.config.OSSCustomDomain != "" {
			// 使用OSS自定义域名，不需要填写访问域名
			return fmt.Sprintf("https://%s/%s", p.config.OSSCustomDomain, objectKey)
		} else {
			// 使用默认OSS域名
			return fmt.Sprintf("https://%s.%s/%s", p.config.OSSBucket, p.config.OSSEndpoint, objectKey)
		}
	} else {
		// 私有文件生成临时访问URL（1小时有效）
		signedURL, err := p.bucket.SignURL(objectKey, oss.HTTPGet, 3600)
		if err != nil {
			// 如果生成签名URL失败，返回需要认证的URL
			return fmt.Sprintf("/api/v1/file/download/%s", strings.ReplaceAll(path, "/", "%2F"))
		}
		return signedURL
	}
}

// Delete 删除OSS文件
func (p *AliyunOSSProvider) Delete(path string) error {
	objectKey := p.generateObjectKey(path)

	err := p.bucket.DeleteObject(objectKey)
	if err != nil {
		return fmt.Errorf("failed to delete OSS object %s: %w", objectKey, err)
	}

	return nil
}

// GetFullPath 获取OSS对象的完整路径（用于兼容接口）
func (p *AliyunOSSProvider) GetFullPath(path string, isPublic bool) string {
	return p.generateObjectKey(path)
}

// generateObjectKey 生成OSS对象键
func (p *AliyunOSSProvider) generateObjectKey(path string) string {
	// OSS对象键不需要区分public/private，通过ACL控制访问权限
	// 但为了保持路径结构一致性，仍然包含这个信息
	return path
}

// GetSignedURL 生成临时访问URL（用于私有文件下载）
func (p *AliyunOSSProvider) GetSignedURL(path string, expireSeconds int64) (string, error) {
	objectKey := p.generateObjectKey(path)

	signedURL, err := p.bucket.SignURL(objectKey, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return signedURL, nil
}

// IsObjectExists 检查对象是否存在
func (p *AliyunOSSProvider) IsObjectExists(path string) (bool, error) {
	objectKey := p.generateObjectKey(path)

	exists, err := p.bucket.IsObjectExist(objectKey)
	if err != nil {
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return exists, nil
}
