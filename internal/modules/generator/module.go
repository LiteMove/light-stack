package generator

import (
	"github.com/LiteMove/light-stack/internal/modules/generator/routes"
	"github.com/gin-gonic/gin"
)

type Module struct {
	name    string
	version string
}

func NewModule() *Module {
	return &Module{
		name:    "generator",
		version: "1.0.0",
	}
}

func (m *Module) InitModule() error {
	return nil
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	routes.RegisterGeneratorRoutes(api)
}

func (m *Module) GetModuleName() string {
	return m.name
}

func (m *Module) GetModuleVersion() string {
	return m.version
}
