package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestUnifiedAuth_NoAuth 测试无需认证的情况
func TestUnifiedAuth_NoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", UnifiedAuth(NoAuth), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUnifiedAuth_RequireLogin_NoToken 测试需要登录但无Token的情况
func TestUnifiedAuth_RequireLogin_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", UnifiedAuth(RequireLogin), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestUnifiedAuth_RequireLogin_InvalidToken 测试需要登录但Token无效的情况
func TestUnifiedAuth_RequireLogin_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", UnifiedAuth(RequireLogin), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestUnifiedAuth_SuperAdmin_NotSuperAdmin 测试需要超管但不是超管的情况
func TestUnifiedAuth_SuperAdmin_NotSuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		// 模拟普通用户
		c.Set("userId", uint64(1))
		c.Set("is_super_admin", false)
		c.Next()
	}, UnifiedAuth(RequireSuperAdmin), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestUnifiedAuth_CustomErrorMessage 测试自定义错误消息
func TestUnifiedAuth_CustomErrorMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	customConfig := AuthConfig{
		RequireAuth:  true,
		ErrorMessage: "自定义错误消息",
	}

	router := gin.New()
	router.GET("/test", UnifiedAuth(customConfig), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "自定义错误消息")
}

// TestSimpleAuth 测试简单权限中间件
func TestSimpleAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", SimpleAuth("system:user:create"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestSimpleRoleAuth 测试简单角色权限中间件
func TestSimpleRoleAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", SimpleRoleAuth("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestAdminAuth 测试管理员权限中间件
func TestAdminAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/test", AdminAuth(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestValidateJWT_MissingToken 测试JWT验证-缺少Token
func TestValidateJWT_MissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	userID, err := validateJWT(c)

	assert.Equal(t, uint64(0), userID)
	assert.Error(t, err)
}

// TestValidateJWT_InvalidFormat 测试JWT验证-格式无效
func TestValidateJWT_InvalidFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/test", nil)
	c.Request.Header.Set("Authorization", "Invalid format")

	userID, err := validateJWT(c)

	assert.Equal(t, uint64(0), userID)
	assert.Error(t, err)
}

// TestIsSuperAdmin 测试超管检查
func TestIsSuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	// 测试不是超管
	c.Set("is_super_admin", false)
	assert.False(t, isSuperAdmin(c))

	// 测试是超管
	c.Set("is_super_admin", true)
	assert.True(t, isSuperAdmin(c))

	// 测试未设置
	c.Set("is_super_admin", nil)
	assert.False(t, isSuperAdmin(c))
}

// TestAuthConfig_WithCustomMessage 测试配置链式调用-自定义消息
func TestAuthConfig_WithCustomMessage(t *testing.T) {
	config := RequireLogin.WithCustomMessage("自定义消息")
	assert.Equal(t, "自定义消息", config.ErrorMessage)
	assert.True(t, config.RequireAuth)
}

// TestAuthConfig_WithLogic 测试配置链式调用-逻辑设置
func TestAuthConfig_WithLogic(t *testing.T) {
	config := RequireLogin.WithLogic(AuthLogicAND)
	assert.Equal(t, AuthLogicAND, config.Logic)
}

// TestAuthConfig_WithPermissions 测试配置链式调用-权限设置
func TestAuthConfig_WithPermissions(t *testing.T) {
	config := RequireLogin.WithPermissions("system:user:create", "system:user:update")
	assert.Contains(t, config.Permissions, "system:user:create")
	assert.Contains(t, config.Permissions, "system:user:update")
}

// TestAuthConfig_WithRoles 测试配置链式调用-角色设置
func TestAuthConfig_WithRoles(t *testing.T) {
	config := RequireLogin.WithRoles("admin", "user")
	assert.Contains(t, config.Roles, "admin")
	assert.Contains(t, config.Roles, "user")
}
