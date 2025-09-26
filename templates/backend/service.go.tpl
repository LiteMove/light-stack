package service

import (
	"fmt"

	"github.com/LiteMove/light-stack/internal/model"
	"github.com/LiteMove/light-stack/internal/repository"
)

// {{.ClassName}}Service {{.FunctionName}}服务
type {{.ClassName}}Service struct {
	repo *repository.{{.ClassName}}Repository
}

// New{{.ClassName}}Service 创建{{.FunctionName}}服务
func New{{.ClassName}}Service(repo *repository.{{.ClassName}}Repository) *{{.ClassName}}Service {
	return &{{.ClassName}}Service{
		repo: repo,
	}
}

// Create 创建{{.FunctionName}}
func (s *{{.ClassName}}Service) Create({{uncapitalize .BusinessName}} *model.{{.ClassName}}) error {
	if err := s.validate{{.ClassName}}({{uncapitalize .BusinessName}}); err != nil {
		return err
	}

	return s.repo.Create({{uncapitalize .BusinessName}})
}

// GetByID 根据ID获取{{.FunctionName}}
func (s *{{.ClassName}}Service) GetByID(id uint64) (*model.{{.ClassName}}, error) {
	if id == 0 {
		return nil, fmt.Errorf("ID不能为空")
	}

	return s.repo.GetByID(id)
}

// Update 更新{{.FunctionName}}
func (s *{{.ClassName}}Service) Update({{uncapitalize .BusinessName}} *model.{{.ClassName}}) error {
	if {{uncapitalize .BusinessName}}.ID == 0 {
		return fmt.Errorf("ID不能为空")
	}

	if err := s.validate{{.ClassName}}({{uncapitalize .BusinessName}}); err != nil {
		return err
	}

	return s.repo.Update({{uncapitalize .BusinessName}})
}

// Delete 删除{{.FunctionName}}
func (s *{{.ClassName}}Service) Delete(id uint64) error {
	if id == 0 {
		return fmt.Errorf("ID不能为空")
	}

	// 检查是否存在
	{{uncapitalize .BusinessName}}, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("{{.FunctionName}}不存在: %v", err)
	}

	if {{uncapitalize .BusinessName}} == nil {
		return fmt.Errorf("{{.FunctionName}}不存在")
	}

	return s.repo.Delete(id)
}

// GetList 获取{{.FunctionName}}列表
func (s *{{.ClassName}}Service) GetList({{- if .HasQuery }}query *model.{{.ClassName}}Query{{- else }}page, pageSize int{{- end }}) ([]*model.{{.ClassName}}, int64, error) {
	return s.repo.GetList({{- if .HasQuery }}query{{- else }}page, pageSize{{- end }})
}

// validate{{.ClassName}} 验证{{.FunctionName}}数据
func (s *{{.ClassName}}Service) validate{{.ClassName}}({{uncapitalize .BusinessName}} *model.{{.ClassName}}) error {
{{- range .Fields }}
{{- if and .IsRequired (ne .GoType "time.Time") (not .IsPk) }}
	{{- if eq .GoType "string" }}
	if {{uncapitalize $.BusinessName}}.{{.GoField}} == "" {
		return fmt.Errorf("{{.ColumnComment}}不能为空")
	}
	{{- else if or (eq .GoType "int") (eq .GoType "int64") (eq .GoType "uint") (eq .GoType "uint64") }}
	if {{uncapitalize $.BusinessName}}.{{.GoField}} == 0 {
		return fmt.Errorf("{{.ColumnComment}}不能为空")
	}
	{{- end }}
{{- end }}
{{- end }}

	return nil
}