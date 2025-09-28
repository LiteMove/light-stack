package service

import (
	"fmt"
	repository3 "github.com/LiteMove/light-stack/internal/modules/files/repository"
	repository2 "github.com/LiteMove/light-stack/internal/modules/system/repository"
	"runtime"
	"time"

	"github.com/LiteMove/light-stack/internal/model"
)

// DashboardService 仪表盘服务接口
type DashboardService interface {
	GetDashboardStats(userID, tenantID uint64) (interface{}, error)
	GetSystemInfo() (*SystemInfo, error)
}

// dashboardService 仪表盘服务实现
type dashboardService struct {
	userRepo   repository2.UserRepository
	tenantRepo repository2.TenantRepository
	fileRepo   *repository3.FileRepository
}

// NewDashboardService 创建仪表盘服务
func NewDashboardService(userRepo repository2.UserRepository, tenantRepo repository2.TenantRepository, fileRepo *repository3.FileRepository) DashboardService {
	return &dashboardService{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		fileRepo:   fileRepo,
	}
}

// SuperAdminStats 超管统计数据
type SuperAdminStats struct {
	UserCount   int64 `json:"userCount"`   // 用户总数
	TenantCount int64 `json:"tenantCount"` // 租户总数
	FileCount   int64 `json:"fileCount"`   // 文件总数
	LogCount    int64 `json:"logCount"`    // 日志总数
}

// TenantAdminStats 租户管理员统计数据
type TenantAdminStats struct {
	UserCount    int64      `json:"userCount"`    // 当前租户用户数
	FileCount    int64      `json:"fileCount"`    // 当前租户文件数
	LogCount     int64      `json:"logCount"`     // 当前租户日志数
	TenantName   string     `json:"tenantName"`   // 租户名称
	ExpiredAt    *time.Time `json:"expiredAt"`    // 到期时间
	Status       int        `json:"status"`       // 租户状态
	StorageUsed  int64      `json:"storageUsed"`  // 存储使用量（字节）
	StorageLimit int64      `json:"storageLimit"` // 存储限制（字节）
}

// UserWelcomeInfo 普通用户欢迎信息
type UserWelcomeInfo struct {
	Username    string `json:"username"`    // 用户名
	LastLogin   string `json:"lastLogin"`   // 最后登录时间
	TenantName  string `json:"tenantName"`  // 所属租户
	WelcomeText string `json:"welcomeText"` // 欢迎文字
}

// SystemInfo 系统信息
type SystemInfo struct {
	Version    string `json:"version"`    // 系统版本
	GoVersion  string `json:"goVersion"`  // Go版本
	StartTime  string `json:"startTime"`  // 启动时间
	Uptime     string `json:"uptime"`     // 运行时间
	OSInfo     string `json:"osInfo"`     // 操作系统信息
	ServerTime string `json:"serverTime"` // 服务器时间
}

var startTime = time.Now()

// GetDashboardStats 获取仪表盘统计数据
func (s *dashboardService) GetDashboardStats(userID, tenantID uint64) (interface{}, error) {
	// 从数据库获取用户信息（包含角色）
	user, err := s.userRepo.GetByIDWithRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 从数据库获取租户信息
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("获取租户信息失败: %w", err)
	}

	// 判断用户角色
	if user.IsSuperAdmin() {
		// 超级管理员
		return s.getSuperAdminStats()
	} else if user.IsAdmin() {
		// 租户管理员
		return s.getTenantAdminStats(tenant)
	} else {
		// 普通用户
		return s.getUserWelcomeInfo(user, tenant)
	}
}

// getSuperAdminStats 获取超管统计数据
func (s *dashboardService) getSuperAdminStats() (*SuperAdminStats, error) {
	// 获取用户总数
	userCount, err := s.userRepo.GetTotalCount()
	if err != nil {
		return nil, fmt.Errorf("获取用户总数失败: %w", err)
	}

	// 获取租户总数
	tenantCount, err := s.tenantRepo.GetTotalCount()
	if err != nil {
		return nil, fmt.Errorf("获取租户总数失败: %w", err)
	}

	// 获取文件总数
	fileCount, err := s.fileRepo.GetTotalCount()
	if err != nil {
		return nil, fmt.Errorf("获取文件总数失败: %w", err)
	}

	return &SuperAdminStats{
		UserCount:   userCount,
		TenantCount: tenantCount,
		FileCount:   fileCount,
		LogCount:    0, // TODO: 实现日志统计
	}, nil
}

// getTenantAdminStats 获取租户管理员统计数据
func (s *dashboardService) getTenantAdminStats(tenant *model.Tenant) (*TenantAdminStats, error) {
	// 获取当前租户用户数
	userCount, err := s.userRepo.GetCountByTenantID(tenant.ID)
	if err != nil {
		return nil, fmt.Errorf("获取租户用户数失败: %w", err)
	}

	// 获取当前租户文件数
	fileCount, err := s.fileRepo.GetCountByTenantID(tenant.ID)
	if err != nil {
		return nil, fmt.Errorf("获取租户文件数失败: %w", err)
	}

	return &TenantAdminStats{
		UserCount:    userCount,
		FileCount:    fileCount,
		LogCount:     0, // TODO: 实现日志统计
		TenantName:   tenant.Name,
		ExpiredAt:    tenant.ExpiredAt,
		Status:       tenant.Status,
		StorageUsed:  0, // TODO: 计算存储使用量
		StorageLimit: 0, // TODO: 从配置获取存储限制
	}, nil
}

// getUserWelcomeInfo 获取普通用户欢迎信息
func (s *dashboardService) getUserWelcomeInfo(user *model.User, tenant *model.Tenant) (*UserWelcomeInfo, error) {
	lastLogin := ""
	if user.LastLoginAt != nil {
		lastLogin = user.LastLoginAt.Format("2006-01-02 15:04:05")
	}

	welcomeText := fmt.Sprintf("欢迎回来，%s！今天也要加油工作哦～", user.Username)
	if time.Now().Hour() < 12 {
		welcomeText = fmt.Sprintf("早上好，%s！新的一天开始了，祝您工作愉快！", user.Username)
	} else if time.Now().Hour() < 18 {
		welcomeText = fmt.Sprintf("下午好，%s！今天的工作进展如何呢？", user.Username)
	} else {
		welcomeText = fmt.Sprintf("晚上好，%s！辛苦了一天，注意休息哦～", user.Username)
	}

	return &UserWelcomeInfo{
		Username:    user.Username,
		LastLogin:   lastLogin,
		TenantName:  tenant.Name,
		WelcomeText: welcomeText,
	}, nil
}

// GetSystemInfo 获取系统信息
func (s *dashboardService) GetSystemInfo() (*SystemInfo, error) {
	uptime := time.Since(startTime)
	uptimeStr := fmt.Sprintf("%.0f天%.0f小时%.0f分钟",
		uptime.Hours()/24,
		float64(int(uptime.Hours())%24),
		uptime.Minutes()-float64(int(uptime.Hours())*60))

	return &SystemInfo{
		Version:    "1.0.0", // TODO: 从配置或构建信息获取
		GoVersion:  runtime.Version(),
		StartTime:  startTime.Format("2006-01-02 15:04:05"),
		Uptime:     uptimeStr,
		OSInfo:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		ServerTime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}
