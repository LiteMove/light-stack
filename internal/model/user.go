package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	TenantBaseModel
	Username      string     `json:"username" gorm:"not null;size:50;uniqueIndex:uk_tenant_username" validate:"required,min=3,max=50"`
	Password      string     `json:"-" gorm:"not null;size:255" validate:"required,min=6"`
	Nickname      string     `json:"nickname" gorm:"size:100" validate:"max=100"`
	Email         *string    `json:"email" gorm:"size:100;uniqueIndex:uk_tenant_email" validate:"omitempty,email,max=100"`
	Phone         *string    `json:"phone" gorm:"size:20;uniqueIndex:uk_tenant_phone" validate:"omitempty,max=20"`
	Avatar        string     `json:"avatar" gorm:"size:255"`
	Status        int        `json:"status" gorm:"not null;default:1;index" validate:"required,oneof=1 2 3"`
	IsSystem      bool       `json:"is_system" gorm:"not null;default:false;index"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	LastLoginIP   string     `json:"last_login_ip" gorm:"size:45"`
	LoginFailures int        `json:"login_failures" gorm:"not null;default:0"`
	LockedUntil   *time.Time `json:"locked_until"`

	// 关联关系
	Roles  []Role  `json:"roles,omitempty" gorm:"many2many:user_roles;"`
	Tenant *Tenant `json:"tenant,omitempty" gorm:"foreignKey:TenantID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserProfile 用户资料（不包含敏感信息）
type UserProfile struct {
	ID          uint64        `json:"id"`
	TenantID    uint64        `json:"tenant_id"`
	Username    string        `json:"username"`
	Nickname    string        `json:"nickname"`
	Email       *string       `json:"email"`
	Phone       *string       `json:"phone"`
	Avatar      string        `json:"avatar"`
	Status      int           `json:"status"`
	IsSystem    bool          `json:"is_system"`
	LastLoginAt *time.Time    `json:"last_login_at"`
	LastLoginIP string        `json:"last_login_ip"`
	Roles       []RoleProfile `json:"roles,omitempty"`
	RoleCodes   []string      `json:"role_codes,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// ToProfile 转换为用户资料
func (u *User) ToProfile() UserProfile {
	profile := UserProfile{
		ID:          u.ID,
		TenantID:    u.TenantID,
		Username:    u.Username,
		Nickname:    u.Nickname,
		Email:       u.Email,
		Phone:       u.Phone,
		Avatar:      u.Avatar,
		Status:      u.Status,
		IsSystem:    u.IsSystem,
		LastLoginAt: u.LastLoginAt,
		LastLoginIP: u.LastLoginIP,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}

	// 转换角色信息
	if len(u.Roles) > 0 {
		roles := make([]RoleProfile, 0, len(u.Roles))
		roleCodes := make([]string, 0, len(u.Roles))

		for _, role := range u.Roles {
			roles = append(roles, role.ToProfile())
			roleCodes = append(roleCodes, role.Code)
		}

		profile.Roles = roles
		profile.RoleCodes = roleCodes
	}

	return profile
}

// HasRole 检查用户是否拥有指定角色
func (u *User) HasRole(roleCode string) bool {
	for _, role := range u.Roles {
		if role.Code == roleCode {
			return true
		}
	}
	return false
}

// HasAnyRole 检查用户是否拥有任一指定角色
func (u *User) HasAnyRole(roleCodes ...string) bool {
	for _, code := range roleCodes {
		if u.HasRole(code) {
			return true
		}
	}
	return false
}

// IsActive 检查用户是否为活跃状态
func (u *User) IsActive() bool {
	return u.Status == 1
}

// IsLocked 检查用户是否被锁定
func (u *User) IsLocked() bool {
	if u.Status == 3 {
		return true
	}
	if u.LockedUntil != nil && u.LockedUntil.After(time.Now()) {
		return true
	}
	return false
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Status == 0 {
		u.Status = 1 // 默认启用
	}
	return nil
}
