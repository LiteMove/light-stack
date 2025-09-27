package system

import (
	"github.com/LiteMove/light-stack/internal/modules/system/routes"
	"github.com/gin-gonic/gin"
)

type Module struct {
	name    string
	version string
}

func NewModule() *Module {
	return &Module{
		name:    "system",
		version: "1.0.0",
	}
}

func (m *Module) InitModule() error {
	// 初始化模块特定的配置、服务等
	return nil
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	routes.RegisterSystemRoutes(api)
}

func (m *Module) GetModuleName() string {
	return m.name
}

func (m *Module) GetModuleVersion() string {
	return m.version
}
