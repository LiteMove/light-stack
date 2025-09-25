package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiteMove/light-stack/internal/container"
	"github.com/LiteMove/light-stack/internal/routes"
	"github.com/LiteMove/light-stack/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestSuite 集成测试套件
type TestSuite struct {
	Engine    *gin.Engine
	Container container.Container
	DB        interface{} // 使用interface{}避免gorm依赖
}

// NewTestSuite 创建集成测试套件
// 这个函数负责初始化整个测试环境，包括数据库连接、DI容器和路由设置
func NewTestSuite() (*TestSuite, error) {
	// 设置gin为测试模式，减少日志输出
	gin.SetMode(gin.TestMode)

	// 初始化数据库连接
	// 在实际测试中，这里应该连接到测试数据库
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("failed to initialize database")
	}

	// 创建DI容器并注册所有服务
	// 容器会自动处理所有依赖注入，包括Repository、Service、Controller
	c := container.NewContainerCompat(db)

	// 创建gin引擎
	engine := gin.New()

	// 注册所有路由模块
	// 这会加载公开接口、认证、用户、管理员、超管、静态文件等所有模块
	routes.RegisterRoutesV2(engine)

	return &TestSuite{
		Engine:    engine,
		Container: c,
		DB:        db,
	}, nil
}

// TestResponse 测试响应结构
type TestResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// MakeRequest 发送HTTP请求的辅助方法
// 这个方法封装了HTTP请求的创建和发送逻辑，简化测试代码
func (ts *TestSuite) MakeRequest(method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var bodyReader *bytes.Buffer

	// 如果有请求体，将其转换为JSON
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal request body: %v", err))
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	} else {
		bodyReader = bytes.NewBuffer(nil)
	}

	// 创建HTTP请求
	req := httptest.NewRequest(method, url, bodyReader)

	// 设置默认的Content-Type
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 添加自定义header
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 创建响应记录器并发送请求
	recorder := httptest.NewRecorder()
	ts.Engine.ServeHTTP(recorder, req)

	return recorder
}

// ParseResponse 解析响应JSON的辅助方法
func (ts *TestSuite) ParseResponse(recorder *httptest.ResponseRecorder) *TestResponse {
	var response TestResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		// 如果无法解析为标准响应格式，返回原始内容
		return &TestResponse{
			Code:    recorder.Code,
			Message: recorder.Body.String(),
			Data:    nil,
		}
	}
	return &response
}

// TestBasicRoutes 测试基础路由功能
func TestBasicRoutes(t *testing.T) {
	// 创建测试套件
	suite, err := NewTestSuite()
	if err != nil {
		t.Fatalf("Failed to create test suite: %v", err)
	}

	// 测试健康检查接口
	t.Run("Health Check", func(t *testing.T) {
		recorder := suite.MakeRequest("GET", "/api/health", nil, nil)
		assert.Equal(t, http.StatusOK, recorder.Code)

		response := suite.ParseResponse(recorder)
		assert.Contains(t, response.Message, "ok")
	})

	// 测试ping接口
	t.Run("Ping", func(t *testing.T) {
		recorder := suite.MakeRequest("GET", "/api/v1/ping", nil, nil)
		assert.Equal(t, http.StatusOK, recorder.Code)

		response := suite.ParseResponse(recorder)
		assert.Contains(t, response.Message, "pong")
	})
}

// TestAuthRoutes 测试认证相关路由
func TestAuthRoutes(t *testing.T) {
	suite, err := NewTestSuite()
	if err != nil {
		t.Fatalf("Failed to create test suite: %v", err)
	}

	// 测试登录接口（无需实际验证，只测试路由可达性）
	t.Run("Login Route Exists", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": "test",
			"password": "test",
		}

		recorder := suite.MakeRequest("POST", "/api/v1/auth/login", loginData, nil)
		// 这里不验证具体业务逻辑，只验证路由是否正确响应
		// 通常会返回400或401，而不是404（路由不存在）
		assert.NotEqual(t, http.StatusNotFound, recorder.Code)
	})

	// 测试注册接口
	t.Run("Register Route Exists", func(t *testing.T) {
		registerData := map[string]interface{}{
			"username": "testuser",
			"password": "testpass",
			"email":    "test@example.com",
		}

		recorder := suite.MakeRequest("POST", "/api/v1/auth/register", registerData, nil)
		assert.NotEqual(t, http.StatusNotFound, recorder.Code)
	})
}

// TestUserRoutes 测试用户相关路由
func TestUserRoutes(t *testing.T) {
	suite, err := NewTestSuite()
	if err != nil {
		t.Fatalf("Failed to create test suite: %v", err)
	}

	// 测试用户资料接口（需要认证，预期返回401）
	t.Run("User Profile Requires Auth", func(t *testing.T) {
		recorder := suite.MakeRequest("GET", "/api/v1/user/profile", nil, nil)
		// 没有认证token，应该返回401
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	// 测试文件上传接口（需要认证）
	t.Run("File Upload Requires Auth", func(t *testing.T) {
		recorder := suite.MakeRequest("POST", "/api/v1/files/upload", nil, nil)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}

// TestAdminRoutes 测试管理员路由
func TestAdminRoutes(t *testing.T) {
	suite, err := NewTestSuite()
	if err != nil {
		t.Fatalf("Failed to create test suite: %v", err)
	}

	// 测试管理员用户列表接口（需要认证和权限）
	t.Run("Admin Users List Requires Auth", func(t *testing.T) {
		recorder := suite.MakeRequest("GET", "/api/v1/admin/users", nil, nil)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}

// TestSuperAdminRoutes 测试超管路由
func TestSuperAdminRoutes(t *testing.T) {
	suite, err := NewTestSuite()
	if err != nil {
		t.Fatalf("Failed to create test suite: %v", err)
	}

	// 测试超管租户管理接口（需要超管权限）
	t.Run("SuperAdmin Tenants Requires Auth", func(t *testing.T) {
		recorder := suite.MakeRequest("GET", "/api/v1/super-admin/tenants", nil, nil)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
}

// TestRoutesCoverage 测试路由覆盖率
// 这个测试确保所有重要的路由都已正确注册
func TestRoutesCoverage(t *testing.T) {
	suite, err := NewTestSuite()
	if err != nil {
		t.Fatalf("Failed to create test suite: %v", err)
	}

	// 定义需要测试的关键路由
	criticalRoutes := []struct {
		method string
		path   string
		desc   string
	}{
		{"GET", "/api/health", "健康检查"},
		{"GET", "/api/v1/ping", "连通性测试"},
		{"POST", "/api/v1/auth/login", "用户登录"},
		{"POST", "/api/v1/auth/register", "用户注册"},
		{"GET", "/api/v1/user/profile", "用户资料"},
		{"POST", "/api/v1/files/upload", "文件上传"},
		{"GET", "/api/v1/admin/users", "管理员用户列表"},
		{"GET", "/api/v1/super-admin/tenants", "超管租户列表"},
	}

	// 遍历测试每个关键路由
	for _, route := range criticalRoutes {
		t.Run(fmt.Sprintf("%s %s - %s", route.method, route.path, route.desc), func(t *testing.T) {
			recorder := suite.MakeRequest(route.method, route.path, nil, nil)
			// 确保路由存在（不返回404）
			assert.NotEqual(t, http.StatusNotFound, recorder.Code,
				"Route %s %s should exist", route.method, route.path)
		})
	}
}
