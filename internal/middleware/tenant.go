package middleware

import (
	"net"
	"strconv"
	"strings"

	"github.com/LiteMove/light-stack/internal/service"
	"github.com/LiteMove/light-stack/pkg/response"

	"github.com/gin-gonic/gin"
)

// TenantMiddleware 租户中间件 - 根据请求域名判断租户
func TenantMiddleware(tenantService service.TenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断是否登录，获取用户信息，如果有超级管理员的身份，则获取请求头中的X-Tenant-Id
		if isSuperAdmin := c.GetBool("is_super_admin"); isSuperAdmin {
			if tenantID := c.GetHeader("X-Tenant-Id"); tenantID != "" {
				// 转换为uint64
				tenantIDUint, err := strconv.ParseUint(tenantID, 10, 64)
				if err != nil {
					response.BadRequest(c, "无效的X-Tenant-Id: "+err.Error())
					c.Abort()
					return
				}
				c.Set("tenant_id", tenantIDUint)
				c.Next()
				return
			}
		}

		// 获取请求的Host
		host := c.Request.Host

		// 解析Host，去掉端口号
		if strings.Contains(host, ":") {
			if hostWithoutPort, _, err := net.SplitHostPort(host); err == nil {
				host = hostWithoutPort
			}
		}

		// 系统管理域名，不需要租户验证
		if host == "localhost" || host == "127.0.0.1" {
			c.Set("tenant_id", uint64(1)) // 系统租户ID为1
			c.Set("tenant_domain", "system")
			c.Next()
			return
		}

		// 根据域名获取租户信息
		tenant, err := tenantService.ValidateTenant(host)
		if err != nil {
			response.BadRequest(c, "无效的租户域名: "+err.Error())
			c.Abort()
			return
		}

		// 将租户信息存储到上下文中
		c.Set("tenant_id", tenant.ID)
		c.Set("tenant_domain", tenant.Domain)
		c.Set("tenant", tenant)

		c.Next()
	}
}

// RequireTenantMiddleware 要求租户信息的中间件
func RequireTenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否已设置租户信息
		tenantID, exists := c.Get("tenant_id")
		if !exists {
			response.BadRequest(c, "缺少租户信息")
			c.Abort()
			return
		}

		// 检查租户ID是否有效
		if tenantID == nil {
			response.BadRequest(c, "无效的租户信息")
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetTenantIDFromContext 从上下文获取租户ID的辅助函数
func GetTenantIDFromContext(c *gin.Context) (uint64, bool) {
	tenantID := c.GetUint64("tenant_id")
	if tenantID != 0 {
		return tenantID, true
	}

	return 0, false
}

// GetTenantFromContext 从上下文获取租户信息的辅助函数
func GetTenantFromContext(c *gin.Context) (interface{}, bool) {
	tenant, exists := c.Get("tenant")
	return tenant, exists
}
