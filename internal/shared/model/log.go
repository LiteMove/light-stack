package model

import (
	"time"
)

// OperationLog 操作日志模型
type OperationLog struct {
	TenantBaseModel
	UserID       uint64 `json:"userId" gorm:"not null;index:idx_user_id" validate:"required"`
	Username     string `json:"username" gorm:"not null;size:50" validate:"required,max=50"`
	Operation    string `json:"operation" gorm:"not null;size:50;index:idx_operation" validate:"required,max=50"`
	Method       string `json:"method" gorm:"not null;size:10" validate:"required,max=10"`
	URL          string `json:"url" gorm:"not null;size:500" validate:"required,max=500"`
	Params       string `json:"params" gorm:"type:json"`
	Result       string `json:"result" gorm:"type:text"`
	ErrorMessage string `json:"error_message" gorm:"type:text"`
	IP           string `json:"ip" gorm:"not null;size:45" validate:"required,max=45"`
	UserAgent    string `json:"userAgent" gorm:"size:500" validate:"max=500"`
	Duration     int    `json:"duration"` // 执行时长（毫秒）
	Status       int    `json:"status" gorm:"not null;index:idx_status" validate:"required,oneof=1 2"`

	// 关联关系 - User 在 system 模块中
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// LoginLog 登录日志模型
type LoginLog struct {
	ID        uint64    `json:"id" gorm:"primarykey"`
	TenantID  uint64    `json:"tenantId" gorm:"not null;default:0;index:idx_tenant_id"`
	UserID    *uint64   `json:"userId" gorm:"index:idx_user_id"`
	Username  string    `json:"username" gorm:"not null;size:50;index:idx_username" validate:"required,max=50"`
	IP        string    `json:"ip" gorm:"not null;size:45" validate:"required,max=45"`
	UserAgent string    `json:"userAgent" gorm:"size:500" validate:"max=500"`
	Location  string    `json:"location" gorm:"size:100" validate:"max=100"`
	Browser   string    `json:"browser" gorm:"size:100" validate:"max=100"`
	OS        string    `json:"os" gorm:"size:100" validate:"max=100"`
	Status    int       `json:"status" gorm:"not null;index:idx_status" validate:"required,oneof=1 2"`
	Message   string    `json:"message" gorm:"size:255" validate:"max=255"`
	LoginTime time.Time `json:"loginTime" gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_login_time"`

	// 关联关系 - User 在 system 模块中
}

// TableName 指定表名
func (LoginLog) TableName() string {
	return "login_logs"
}

// OperationLogProfile 操作日志资料（简化版本）
type OperationLogProfile struct {
	ID           uint64    `json:"id"`
	TenantID     uint64    `json:"tenantId"`
	UserID       uint64    `json:"userId"`
	Username     string    `json:"username"`
	Operation    string    `json:"operation"`
	Method       string    `json:"method"`
	URL          string    `json:"url"`
	Params       string    `json:"params"`
	Result       string    `json:"result"`
	ErrorMessage string    `json:"errorMessage"`
	IP           string    `json:"ip"`
	UserAgent    string    `json:"userAgent"`
	Duration     int       `json:"duration"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ToProfile 转换为操作日志资料
func (ol *OperationLog) ToProfile() OperationLogProfile {
	return OperationLogProfile{
		ID:           ol.ID,
		TenantID:     ol.TenantID,
		UserID:       ol.UserID,
		Username:     ol.Username,
		Operation:    ol.Operation,
		Method:       ol.Method,
		URL:          ol.URL,
		Params:       ol.Params,
		Result:       ol.Result,
		ErrorMessage: ol.ErrorMessage,
		IP:           ol.IP,
		UserAgent:    ol.UserAgent,
		Duration:     ol.Duration,
		Status:       ol.Status,
		CreatedAt:    ol.CreatedAt,
	}
}

// LoginLogProfile 登录日志资料（简化版本）
type LoginLogProfile struct {
	ID        uint64    `json:"id"`
	TenantID  uint64    `json:"tenant_id"`
	UserID    *uint64   `json:"userId"`
	Username  string    `json:"username"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Location  string    `json:"location"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	Status    int       `json:"status"`
	Message   string    `json:"message"`
	LoginTime time.Time `json:"login_time"`
}

// ToProfile 转换为登录日志资料
func (ll *LoginLog) ToProfile() LoginLogProfile {
	return LoginLogProfile{
		ID:        ll.ID,
		TenantID:  ll.TenantID,
		UserID:    ll.UserID,
		Username:  ll.Username,
		IP:        ll.IP,
		UserAgent: ll.UserAgent,
		Location:  ll.Location,
		Browser:   ll.Browser,
		OS:        ll.OS,
		Status:    ll.Status,
		Message:   ll.Message,
		LoginTime: ll.LoginTime,
	}
}

// IsSuccess 检查操作是否成功
func (ol *OperationLog) IsSuccess() bool {
	return ol.Status == 1
}

// IsSuccess 检查登录是否成功
func (ll *LoginLog) IsSuccess() bool {
	return ll.Status == 1
}
