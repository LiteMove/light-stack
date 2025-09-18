package service

import (
	"errors"
	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/repository"
	"github.com/LiteMove/light-stack/internal/utils"
	"github.com/LiteMove/light-stack/pkg/jwt"
	"github.com/LiteMove/light-stack/pkg/logger"
	"strings"
)

// AuthService 认证服务接口
type AuthService interface {
	// 用户登录
	Login(req *LoginRequest) (*LoginResponse, error)
	// 用户注册
	Register(req *RegisterRequest) (*model.UserProfile, error)
	// 刷新token
	RefreshToken(tokenString string) (*TokenResponse, error)
	// 验证token
	ValidateToken(tokenString string) (*jwt.Claims, error)
	// 修改密码
	ChangePassword(userID uint, oldPassword, newPassword string) error
	// 获取用户信息
	GetUserProfile(userID uint) (*model.UserProfile, error)
	// 更新用户信息
	UpdateUserProfile(userID uint, req *UpdateProfileRequest) (*model.UserProfile, error)
	// 为用户分配角色
	AssignUserRoles(userID uint, roleIDs []uint) error
	// 获取用户角色
	GetUserRoles(userID uint) ([]*model.Role, error)
}

// RoleService 角色服务接口
type RoleService interface {
	// 创建角色
	Create(req *CreateRoleRequest) (*model.RoleProfile, error)
	// 更新角色
	Update(id uint, req *UpdateRoleRequest) (*model.RoleProfile, error)
	// 删除角色
	Delete(id uint) error
	// 获取角色信息
	GetByID(id uint) (*model.RoleProfile, error)
	// 获取角色列表
	GetList(tenantID uint64, page, pageSize int, status int) ([]*model.Role, int64, error)
	// 为用户分配角色
	AssignRolesToUser(userID uint, roleIDs []uint) error
	// 移除用户角色
	RemoveUserRoles(userID uint, roleIDs []uint) error
}

// authService 认证服务实现
type authService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

// roleService 角色服务实现
type roleService struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) AuthService {
	return &authService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleRepo repository.RoleRepository, userRepo repository.UserRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	TenantID uint64 `json:"tenant_id"` // 租户ID，0表示系统租户
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	TenantID uint64 `json:"tenant_id"` // 租户ID，0表示系统租户
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"email,max=100"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" validate:"max=100"`
	Phone    string `json:"phone" validate:"max=20"`
	RoleIDs  []uint `json:"role_ids"` // 分配的角色ID列表
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" validate:"max=100"`
	Avatar   string `json:"avatar" validate:"max=255"`
	Phone    string `json:"phone" validate:"max=20"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	TenantID    uint64 `json:"tenant_id"`
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Code        string `json:"code" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"max=255"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
	Status      int    `json:"status" validate:"oneof=1 2"`
	SortOrder   int    `json:"sort_order"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	User  model.UserProfile `json:"user"`
	Token TokenResponse     `json:"token"`
}

// TokenResponse token响应
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Login 用户登录
func (s *authService) Login(req *LoginRequest) (*LoginResponse, error) {
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
		user, err = s.userRepo.GetByEmail(req.TenantID, req.Username)
		if err == nil {
			// 加载角色信息
			user, err = s.userRepo.GetByIDWithRoles(user.ID)
		}
	} else {
		user, err = s.userRepo.GetByUsernameWithRoles(req.TenantID, req.Username)
	}

	if err != nil {
		logger.WithField("username", req.Username).Warn("Login attempt with invalid username")
		return nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if !user.IsActive() {
		logger.WithFields(map[string]interface{}{
			"user_id": user.ID,
			"status":  user.Status,
		}).Warn("Login attempt with inactive user")
		return nil, errors.New("账户已被禁用")
	}

	// 检查用户是否被锁定
	if user.IsLocked() {
		logger.WithField("user_id", user.ID).Warn("Login attempt with locked user")
		return nil, errors.New("账户已被锁定")
	}

	// 验证密码
	if !utils.VerifyPassword(user.Password, req.Password) {
		// 记录登录失败
		s.userRepo.RecordLoginFailure(user.ID)
		logger.WithField("user_id", user.ID).Warn("Login attempt with wrong password")
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT token，使用主要角色
	var primaryRole string = "user"
	if len(user.Roles) > 0 {
		primaryRole = user.Roles[0].Code // 使用第一个角色作为主要角色
		// 如果有超级管理员角色，优先使用
		for _, role := range user.Roles {
			if role.Code == "super_admin" {
				primaryRole = role.Code
				break
			}
		}
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, primaryRole)
	if err != nil {
		logger.WithField("user_id", user.ID).Error("Failed to generate token:", err)
		return nil, errors.New("登录失败")
	}

	// 更新最后登录信息
	if err := s.userRepo.UpdateLoginInfo(user.ID, ""); err != nil {
		logger.WithField("user_id", user.ID).Warn("Failed to update login info:", err)
	}

	logger.WithField("user_id", user.ID).Info("User logged in successfully")

	return &LoginResponse{
		User: user.ToProfile(),
		Token: TokenResponse{
			AccessToken: token,
			TokenType:   "Bearer",
			ExpiresIn:   3600, // 1小时
		},
	}, nil
}

// Register 用户注册
func (s *authService) Register(req *RegisterRequest) (*model.UserProfile, error) {
	// 参数验证
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// 检查用户名是否已存在
	exists, err := s.userRepo.UsernameExists(req.TenantID, req.Username)
	if err != nil {
		logger.Error("Failed to check username existence:", err)
		return nil, errors.New("注册失败")
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		exists, err = s.userRepo.EmailExists(req.TenantID, req.Email)
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
		TenantID: req.TenantID,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Status:   1, // 启用状态
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Error("Failed to create user:", err)
		return nil, errors.New("注册失败")
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		if err := s.roleRepo.AssignRolesToUser(user.ID, req.RoleIDs); err != nil {
			logger.WithField("user_id", user.ID).Error("Failed to assign roles:", err)
			// 注册已成功，角色分配失败只记录警告
		}
	} else {
		// 如果没有指定角色，分配默认用户角色
		userRole, err := s.roleRepo.GetByCode("user")
		if err == nil {
			s.roleRepo.AssignRolesToUser(user.ID, []uint{userRole.ID})
		}
	}

	logger.WithField("user_id", user.ID).Info("User registered successfully")

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

	// 生成新token，使用当前的主要角色
	var primaryRole string = "user"
	if len(user.Roles) > 0 {
		primaryRole = user.Roles[0].Code
		for _, role := range user.Roles {
			if role.Code == "super_admin" {
				primaryRole = role.Code
				break
			}
		}
	}

	newToken, err := jwt.GenerateToken(user.ID, user.Username, primaryRole)
	if err != nil {
		logger.WithField("user_id", user.ID).Error("Failed to refresh token:", err)
		return nil, errors.New("刷新token失败")
	}

	return &TokenResponse{
		AccessToken: newToken,
		TokenType:   "Bearer",
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
func (s *authService) ChangePassword(userID uint, oldPassword, newPassword string) error {
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
		logger.WithField("user_id", userID).Error("Failed to hash new password:", err)
		return errors.New("密码修改失败")
	}

	// 更新密码
	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		logger.WithField("user_id", userID).Error("Failed to update password:", err)
		return errors.New("密码修改失败")
	}

	logger.WithField("user_id", userID).Info("Password changed successfully")
	return nil
}

// GetUserProfile 获取用户信息
func (s *authService) GetUserProfile(userID uint) (*model.UserProfile, error) {
	user, err := s.userRepo.GetByIDWithRoles(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	profile := user.ToProfile()
	return &profile, nil
}

// UpdateUserProfile 更新用户信息
func (s *authService) UpdateUserProfile(userID uint, req *UpdateProfileRequest) (*model.UserProfile, error) {
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
		user.Phone = req.Phone
	}

	// 保存更新
	if err := s.userRepo.Update(user); err != nil {
		logger.WithField("user_id", userID).Error("Failed to update user profile:", err)
		return nil, errors.New("更新失败")
	}

	logger.WithField("user_id", userID).Info("User profile updated successfully")

	// 重新获取用户信息（包含角色）
	user, _ = s.userRepo.GetByIDWithRoles(userID)
	profile := user.ToProfile()
	return &profile, nil
}

// AssignUserRoles 为用户分配角色
func (s *authService) AssignUserRoles(userID uint, roleIDs []uint) error {
	return s.roleRepo.UpdateUserRoles(userID, roleIDs)
}

// GetUserRoles 获取用户角色
func (s *authService) GetUserRoles(userID uint) ([]*model.Role, error) {
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

// 角色服务实现

// Create 创建角色
func (s *roleService) Create(req *CreateRoleRequest) (*model.RoleProfile, error) {
	// 检查角色编码是否已存在
	exists, err := s.roleRepo.CodeExists(req.TenantID, req.Code)
	if err != nil {
		logger.Error("Failed to check role code existence:", err)
		return nil, errors.New("创建失败")
	}
	if exists {
		return nil, errors.New("角色编码已存在")
	}

	// 创建角色
	role := &model.Role{
		TenantID:    req.TenantID,
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      1, // 默认启用
		SortOrder:   req.SortOrder,
	}

	if err := s.roleRepo.Create(role); err != nil {
		logger.Error("Failed to create role:", err)
		return nil, errors.New("创建失败")
	}

	logger.WithField("role_id", role.ID).Info("Role created successfully")
	profile := role.ToProfile()
	return &profile, nil
}

// Update 更新角色
func (s *roleService) Update(id uint, req *UpdateRoleRequest) (*model.RoleProfile, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	// 更新角色信息
	role.Name = req.Name
	role.Description = req.Description
	role.Status = req.Status
	role.SortOrder = req.SortOrder

	if err := s.roleRepo.Update(role); err != nil {
		logger.WithField("role_id", id).Error("Failed to update role:", err)
		return nil, errors.New("更新失败")
	}

	logger.WithField("role_id", id).Info("Role updated successfully")
	profile := role.ToProfile()
	return &profile, nil
}

// Delete 删除角色
func (s *roleService) Delete(id uint) error {
	// 检查角色是否还有用户在使用
	count, err := s.roleRepo.GetRoleUserCount(id)
	if err != nil {
		return errors.New("删除失败")
	}
	if count > 0 {
		return errors.New("该角色还有用户在使用，无法删除")
	}

	if err := s.roleRepo.Delete(id); err != nil {
		logger.WithField("role_id", id).Error("Failed to delete role:", err)
		return errors.New("删除失败")
	}

	logger.WithField("role_id", id).Info("Role deleted successfully")
	return nil
}

// GetByID 获取角色信息
func (s *roleService) GetByID(id uint) (*model.RoleProfile, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	profile := role.ToProfile()
	return &profile, nil
}

// GetList 获取角色列表
func (s *roleService) GetList(tenantID uint64, page, pageSize int, status int) ([]*model.Role, int64, error) {
	return s.roleRepo.GetList(tenantID, page, pageSize, status)
}

// AssignRolesToUser 为用户分配角色
func (s *roleService) AssignRolesToUser(userID uint, roleIDs []uint) error {
	return s.roleRepo.UpdateUserRoles(userID, roleIDs)
}

// RemoveUserRoles 移除用户角色
func (s *roleService) RemoveUserRoles(userID uint, roleIDs []uint) error {
	return s.roleRepo.RemoveUserRoles(userID, roleIDs)
}