package service

import (
	"errors"
	"fmt"
	repository2 "github.com/LiteMove/light-stack/internal/modules/system/repository"

	systemModel "github.com/LiteMove/light-stack/internal/modules/system/model"
	"github.com/LiteMove/light-stack/internal/shared/utils"
)

// ProfileService 个人中心服务接口
type ProfileService interface {
	// 个人信息操作
	GetProfile(userID uint64) (*systemModel.UserProfile, error)
	UpdateProfile(userID uint64, nickname, email, phone, avatar string) error
	ChangePassword(userID uint64, oldPassword, newPassword string) error

	// 租户配置操作（仅租户管理员）
	IsTenantAdmin(userID, tenantID uint64) (bool, error)
	GetTenantConfig(tenantID uint64) (*systemModel.TenantConfig, error)
	UpdateTenantConfig(tenantID uint64, config *systemModel.TenantConfig) error
}

// profileService 个人中心服务实现
type profileService struct {
	userRepo   repository2.UserRepository
	roleRepo   repository2.RoleRepository
	tenantRepo repository2.TenantRepository
}

// NewProfileService 创建个人中心服务
func NewProfileService(userRepo repository2.UserRepository, roleRepo repository2.RoleRepository, tenantRepo repository2.TenantRepository) ProfileService {
	return &profileService{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		tenantRepo: tenantRepo,
	}
}

// GetProfile 获取个人信息
func (s *profileService) GetProfile(userID uint64) (*systemModel.UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 获取用户角色
	roles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户角色失败: %w", err)
	}

	profile := &systemModel.UserProfile{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Status:   user.Status,
		TenantID: user.TenantID,
		Roles:    make([]systemModel.RoleProfile, 0, len(roles)),
	}

	// 转换角色信息
	for _, role := range roles {
		profile.Roles = append(profile.Roles, role.ToProfile())
	}

	return profile, nil
}

// UpdateProfile 更新个人信息
func (s *profileService) UpdateProfile(userID uint64, nickname, email, phone, avatar string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 检查邮箱是否已被其他用户使用
	if email != "" && (user.Email == nil || *user.Email != email) {
		exists, err := s.userRepo.EmailExists(user.TenantID, email)
		if err != nil {
			return fmt.Errorf("检查邮箱是否存在失败: %w", err)
		}
		if exists {
			return errors.New("邮箱已被其他用户使用")
		}
	}

	// 更新用户信息
	user.Nickname = nickname
	if email != "" {
		user.Email = &email
	} else {
		user.Email = nil
	}
	if phone != "" {
		user.Phone = &phone
	} else {
		user.Phone = nil
	}
	// 更新头像
	user.Avatar = avatar

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("更新用户信息失败: %w", err)
	}

	return nil
}

// ChangePassword 修改密码
func (s *profileService) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 验证旧密码
	if !utils.VerifyPassword(user.Password, oldPassword) {
		return errors.New("原密码不正确")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// IsTenantAdmin 检查用户是否为租户管理员
func (s *profileService) IsTenantAdmin(userID, tenantID uint64) (bool, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 获取用户角色
	roles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return false, fmt.Errorf("获取用户角色失败: %w", err)
	}

	// 检查是否为超级管理员（超级管理员可以管理任何租户）
	for _, role := range roles {
		if role.Code == "super_admin" || role.Name == "super_admin" || role.Name == "超级管理员" {
			return true, nil
		}
	}

	// 检查用户是否属于该租户
	if user.TenantID != tenantID {
		return false, nil
	}

	// 检查是否有租户管理员角色
	for _, role := range roles {
		// 假设角色名称为 "tenant_admin" 或者角色类型为管理员类型
		if role.Name == "tenant_admin" || role.Name == "租户管理员" || role.Code == "tenant_admin" {
			return true, nil
		}
	}

	return false, nil
}

// GetTenantConfig 获取租户配置
func (s *profileService) GetTenantConfig(tenantID uint64) (*systemModel.TenantConfig, error) {
	// 获取租户信息
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("获取租户信息失败: %w", err)
	}

	// 解析配置
	config, err := tenant.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("解析租户配置失败: %w", err)
	}

	// 确保配置有默认值
	if config.FileStorage.Type == "" {
		config.FileStorage.Type = "local"
	}
	if config.FileStorage.MaxFileSize == 0 {
		config.FileStorage.MaxFileSize = 50 * 1024 * 1024 // 50MB
	}
	if len(config.FileStorage.AllowedTypes) == 0 {
		config.FileStorage.AllowedTypes = []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".txt"}
	}
	if config.FileStorage.OSSProvider == "" {
		config.FileStorage.OSSProvider = "aliyun"
	}

	return config, nil
}

// UpdateTenantConfig 更新租户配置
func (s *profileService) UpdateTenantConfig(tenantID uint64, config *systemModel.TenantConfig) error {
	// 获取租户信息
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return fmt.Errorf("获取租户信息失败: %w", err)
	}

	// 验证配置
	if err := s.validateTenantConfig(config); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	// 设置配置
	if err := tenant.SetConfig(config); err != nil {
		return fmt.Errorf("设置租户配置失败: %w", err)
	}

	// 更新租户
	if err := s.tenantRepo.Update(tenant); err != nil {
		return fmt.Errorf("更新租户失败: %w", err)
	}

	return nil
}

// validateTenantConfig 验证租户配置
func (s *profileService) validateTenantConfig(config *systemModel.TenantConfig) error {
	// 验证文件存储配置
	fileStorage := &config.FileStorage

	// 验证存储类型
	if fileStorage.Type != "local" && fileStorage.Type != "oss" {
		return errors.New("不支持的存储类型")
	}

	// 验证文件大小限制
	if fileStorage.MaxFileSize <= 0 {
		return errors.New("文件大小限制必须大于0")
	}

	// 如果是OSS存储，验证OSS配置
	if fileStorage.Type == "oss" {
		if fileStorage.OSSProvider == "" {
			return errors.New("OSS提供商不能为空")
		}
		if fileStorage.OSSBucket == "" {
			return errors.New("OSS存储桶不能为空")
		}
		if fileStorage.OSSAccessKey == "" {
			return errors.New("OSS访问密钥不能为空")
		}
		if fileStorage.OSSSecretKey == "" {
			return errors.New("OSS密钥不能为空")
		}
	}

	return nil
}
