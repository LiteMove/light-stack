package controller

import (
	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"
	"github.com/gin-gonic/gin"
)

// DashboardController 仪表盘控制器
type DashboardController struct {
	dashboardService service.DashboardService
}

// NewDashboardController 创建仪表盘控制器
func NewDashboardController(dashboardService service.DashboardService) *DashboardController {
	return &DashboardController{
		dashboardService: dashboardService,
	}
}

// GetDashboardStats 获取仪表盘统计数据
func (c *DashboardController) GetDashboardStats(ctx *gin.Context) {
	// 获取用户ID
	userID := ctx.GetUint64("userId")
	if userID == 0 {
		response.Unauthorized(ctx, "未登录")
		return
	}

	// 获取租户ID（从租户中间件设置）
	tenantID := ctx.GetUint64("tenant_id")
	if tenantID == 0 {
		response.BadRequest(ctx, "租户信息不存在")
		return
	}

	// 调用服务获取统计数据
	stats, err := c.dashboardService.GetDashboardStats(userID, tenantID)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, stats)
}

// GetSystemInfo 获取系统信息（超管专用）
func (c *DashboardController) GetSystemInfo(ctx *gin.Context) {
	// 调用服务获取系统信息
	systemInfo, err := c.dashboardService.GetSystemInfo()
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Success(ctx, systemInfo)
}
