package service

import (
	"errors"
	"strings"

	"github.com/LiteMove/light-stack/internal/modules/auth/model"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/shared/utils"
	"github.com/LiteMove/light-stack/pkg/jwt"
	"github.com/LiteMove/light-stack/pkg/logger"
	"github.com/LiteMove/light-stack/pkg/permission"
)

// AuthService 认证服务接口
type AuthService interface {
	// 用户登录
	Login(tenantID uint64, req *LoginRequest) (*TokenResponse, error)
	// 用户注册
	Register(tenantID uint64, req *RegisterRequest) (*model.UserProfile, error)
	// 刷新token
	RefreshToken(tokenString string) (*TokenResponse, error)
	// 验证token
	ValidateToken(tokenString string) (*jwt.Claims, error)
	// 修改密码
	ChangePassword(userID uint64, oldPassword, newPassword string) error
	// 获取用户信息
	GetUserProfile(userID uint64) (*model.UserProfile, error)
	// 更新用户信息
	UpdateUserProfile(userID uint64, req *UpdateProfileRequest) (*model.UserProfile, error)
	// 为用户分配角色
	AssignUserRoles(userID uint64, roleIDs []uint64) error
	// 获取用户角色
	GetUserRoles(userID uint64) ([]*model.Role, error)
}

// RoleService 角色服务接口

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	menuRepo repository.MenuRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, menuRepo repository.MenuRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		menuRepo: menuRepo,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Email    string   `json:"email" validate:"email,max=100"`
	Password string   `json:"password" validate:"required,min=6"`
	Nickname string   `json:"nickname" validate:"max=100"`
	Phone    string   `json:"phone" validate:"max=20"`
	RoleIDs  []uint64 `json:"roleIds"` // 分配的角色ID列表
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" validate:"max=100"`
	Avatar   string `json:"avatar" validate:"max=255"`
	Phone    string `json:"phone" validate:"max=20"`
}

// TokenResponse token响应
type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

// Login 用户登录
func (s *authService) Login(tenantID uint64, req *LoginRequest) (*TokenResponse, error) {
	// 参数验证
	if strings.TrimSpace(req.Username) == "" {
		return nil, errors.New("用户名不能为空")
	}
	if strings.TrimSpace(req.Password) == "" {
		return nil, errors.New("密码不能为空")
	}

	// 获取用户信息（包含角色）
	var user *model.User
	var err error
	// 支持用户名或邮箱登录
	if strings.Contains(req.Username, "@") {
		user, err = s.userRepo.GetByEmail(tenantID, req.Username)
		if err == nil {
			// 加载角色信息
			user, err = s.userRepo.GetByIDWithRoles(user.ID)
		}
	} else {
		user, err = s.userRepo.GetByUsernameWithRoles(tenantID, req.Username)
	}

	if err != nil {
		logger.WithField("username", req.Username).Warn("Login attempt with invalid username")
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if !user.IsActive() {
		logger.WithFields(map[string]interface{}{
			"userId": user.ID,
			"status": user.Status,
		}).Warn("Login attempt with inactive user")
		return nil, errors.New("账户已被禁用")
	}

	// 检查用户是否被锁定
	if user.IsLocked() {
		logger.WithField("userId", user.ID).Warn("Login attempt with locked user")
		return nil, errors.New("账户已被锁定")
	}

	// 验证密码
	if !utils.VerifyPassword(user.Password, req.Password) {
		// 记录登录失败
		s.userRepo.RecordLoginFailure(user.ID)
		logger.WithField("userId", user.ID).Warn("Login attempt with wrong password")
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT token，使用主要角色
	var userRoles []string
	for _, role := range user.Roles {
		userRoles = append(userRoles, role.Code)
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, userRoles)
	if err != nil {
		logger.WithField("userId", user.ID).Error("Failed to generate token:", err)
		return nil, errors.New("登录失败")
	}

	// 更新最后登录信息
	if err := s.userRepo.UpdateLoginInfo(user.ID, ""); err != nil {
		logger.WithField("userId", user.ID).Warn("Failed to update login info:", err)
	}

	// 加载用户权限和角色到缓存
	if err := permission.LoadUserData(user.ID, s.menuRepo, s.roleRepo); err != nil {
		logger.WithField("userId", user.ID).Warn("Failed to load permissions and roles:", err)
	}

	logger.WithField("userId", user.ID).Info("User logged in successfully")

	return &TokenResponse{
		AccessToken: token,
		ExpiresIn:   3600, // 1小时
	}, nil
}

// Register 用户注册
func (s *authService) Register(tenantID uint64, req *RegisterRequest) (*model.UserProfile, error) {
	// 参数验证
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	exists, err := s.userRepo.UsernameExists(tenantID, req.Username)
	if err != nil {
		logger.Error("Failed to check username existence:", err)
		return nil, errors.New("注册失败")
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		exists, err = s.userRepo.EmailExists(tenantID, req.Email)
		if err != nil {
			logger.Error("Failed to check email existence:", err)
			return nil, errors.New("注册失败")
		}
		if exists {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password:", err)
		return nil, errors.New("注册失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Status:   1, // 启用状态
	}

	// 处理可选字段
	if req.Email != "" {
		user.Email = &req.Email
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	user.TenantID = tenantID

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("Failed to create user:", err)
		return nil, errors.New("注册失败")
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		if err := s.roleRepo.AssignRolesToUser(user.ID, req.RoleIDs); err != nil {
			logger.WithField("userId", user.ID).Error("Failed to assign roles:", err)
			// 注册已成功，角色分配失败只记录警告
		}
	} else {
		// 如果没有指定角色，分配默认用户角色
		userRole, err := s.roleRepo.GetByCode("user")
		if err == nil {
			s.roleRepo.AssignRolesToUser(user.ID, []uint64{userRole.ID})
		}
	}

	logger.WithField("userId", user.ID).Info("User registered successfully")

	// 重新获取用户信息（包含角色）
	user, _ = s.userRepo.GetByIDWithRoles(user.ID)
	profile := user.ToProfile()
	return &profile, nil
}

// RefreshToken 刷新token
func (s *authService) RefreshToken(tokenString string) (*TokenResponse, error) {
	// 解析原token
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, errors.New("无效的token")
	}

	// 检查用户是否仍然有效
	user, err := s.userRepo.GetByIDWithRoles(claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if !user.IsActive() {
		return nil, errors.New("账户已被禁用")
	}

	// 生成新token
	var userRoles []string
	for _, role := range user.Roles {
		userRoles = append(userRoles, role.Code)
	}

	newToken, err := jwt.GenerateToken(user.ID, user.Username, userRoles)
	if err != nil {
		logger.WithField("userId", user.ID).Error("Failed to refresh token:", err)
		return nil, errors.New("刷新token失败")
	}

	return &TokenResponse{
		AccessToken: newToken,
		ExpiresIn:   3600,
	}, nil
}

// ValidateToken 验证token
func (s *authService) ValidateToken(tokenString string) (*jwt.Claims, error) {
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 验证用户是否仍然有效
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if !user.IsActive() {
		return nil, errors.New("账户已被禁用")
	}

	return claims, nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(userID uint64, oldPassword, newPassword string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !utils.VerifyPassword(user.Password, oldPassword) {
		return errors.New("原密码错误")
	}

	// 验证新密码强度
	if err := utils.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		logger.WithField("userId", userID).Error("Failed to hash new password:", err)
		return errors.New("密码修改失败")
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		logger.WithField("userId", userID).Error("Failed to update password:", err)
		return errors.New("密码修改失败")
	}

	logger.WithField("userId", userID).Info("Password changed successfully")
	return nil
}

// GetUserProfile 获取用户信息（不包含菜单和权限）
func (s *authService) GetUserProfile(userID uint64) (*model.UserProfile, error) {
	user, err := s.userRepo.GetByIDWithRoles(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	profile := user.ToProfile()

	// 获取用户权限（仅返回权限码数组）
	permissions, err := s.menuRepo.GetUserPermissions(userID)
	if err != nil {
		logger.WithField("userId", userID).Warn("Failed to get user permissions:", err)
	} else {
		profile.Permissions = permissions
	}

	return &profile, nil
}

// UpdateUserProfile 更新用户信息
func (s *authService) UpdateUserProfile(userID uint64, req *UpdateProfileRequest) (*model.UserProfile, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新用户信息
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	// 保存更新
	if err := s.userRepo.Update(user); err != nil {
		logger.WithField("userId", userID).Error("Failed to update user profile:", err)
		return nil, errors.New("更新失败")
	}

	logger.WithField("userId", userID).Info("User profile updated successfully")

	// 重新获取用户信息（包含角色）
	user, _ = s.userRepo.GetByIDWithRoles(userID)
	profile := user.ToProfile()
	return &profile, nil
}

// AssignUserRoles 为用户分配角色
func (s *authService) AssignUserRoles(userID uint64, roleIDs []uint64) error {
	return s.roleRepo.UpdateUserRoles(userID, roleIDs)
}

// GetUserRoles 获取用户角色
func (s *authService) GetUserRoles(userID uint64) ([]*model.Role, error) {
	return s.roleRepo.GetUserRoles(userID)
}

// validateRegisterRequest 验证注册请求
func (s *authService) validateRegisterRequest(req *RegisterRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return errors.New("用户名不能为空")
	}
	if len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("用户名长度必须在3-50字符之间")
	}
	if strings.TrimSpace(req.Password) == "" {
		return errors.New("密码不能为空")
	}

	// 验证密码强度
	if err := utils.ValidatePasswordStrength(req.Password); err != nil {
		return err
	}

	return nil
}
