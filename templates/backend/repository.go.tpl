package repository

import (
	"fmt"

	"github.com/LiteMove/light-stack/internal/model"
	"gorm.io/gorm"
)

// {{.ClassName}}Repository {{.FunctionName}}仓储
type {{.ClassName}}Repository struct {
	db *gorm.DB
}

// New{{.ClassName}}Repository 创建{{.FunctionName}}仓储
func New{{.ClassName}}Repository(db *gorm.DB) *{{.ClassName}}Repository {
	return &{{.ClassName}}Repository{
		db: db,
	}
}

// Create 创建{{.FunctionName}}
func (r *{{.ClassName}}Repository) Create({{uncapitalize .BusinessName}} *model.{{.ClassName}}) error {
	if err := r.db.Create({{uncapitalize .BusinessName}}).Error; err != nil {
		return fmt.Errorf("创建{{.FunctionName}}失败: %v", err)
	}
	return nil
}

// GetByID 根据ID获取{{.FunctionName}}
func (r *{{.ClassName}}Repository) GetByID(id uint64) (*model.{{.ClassName}}, error) {
	var {{uncapitalize .BusinessName}} model.{{.ClassName}}

	err := r.db.Where("id = ?", id).First(&{{uncapitalize .BusinessName}}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("{{.FunctionName}}不存在")
		}
		return nil, fmt.Errorf("查询{{.FunctionName}}失败: %v", err)
	}

	return &{{uncapitalize .BusinessName}}, nil
}

// Update 更新{{.FunctionName}}
func (r *{{.ClassName}}Repository) Update({{uncapitalize .BusinessName}} *model.{{.ClassName}}) error {
	if err := r.db.Save({{uncapitalize .BusinessName}}).Error; err != nil {
		return fmt.Errorf("更新{{.FunctionName}}失败: %v", err)
	}
	return nil
}

// Delete 删除{{.FunctionName}}
func (r *{{.ClassName}}Repository) Delete(id uint64) error {
	if err := r.db.Where("id = ?", id).Delete(&model.{{.ClassName}}{}).Error; err != nil {
		return fmt.Errorf("删除{{.FunctionName}}失败: %v", err)
	}
	return nil
}

// GetList 获取{{.FunctionName}}列表
func (r *{{.ClassName}}Repository) GetList({{- if .HasQuery }}query *model.{{.ClassName}}Query{{- else }}page, pageSize int{{- end }}) ([]*model.{{.ClassName}}, int64, error) {
	var {{pluralize (uncapitalize .BusinessName)}} []*model.{{.ClassName}}
	var total int64

	db := r.db.Model(&model.{{.ClassName}}{})

{{- if .HasQuery }}
	// 应用查询条件
	if query != nil {
		db = query.Apply(db)

		// 计算总数
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, fmt.Errorf("查询{{.FunctionName}}总数失败: %v", err)
		}

		// 分页查询
		if query.Page > 0 && query.PageSize > 0 {
			offset := (query.Page - 1) * query.PageSize
			db = db.Offset(offset).Limit(query.PageSize)
		}
	} else {
		// 无条件查询总数
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, fmt.Errorf("查询{{.FunctionName}}总数失败: %v", err)
		}
	}
{{- else }}
	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询{{.FunctionName}}总数失败: %v", err)
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}
{{- end }}

	// 查询数据
	if err := db.Find(&{{pluralize (uncapitalize .BusinessName)}}).Error; err != nil {
		return nil, 0, fmt.Errorf("查询{{.FunctionName}}列表失败: %v", err)
	}

	return {{pluralize (uncapitalize .BusinessName)}}, total, nil
}

// GetAll 获取所有{{.FunctionName}}
func (r *{{.ClassName}}Repository) GetAll() ([]*model.{{.ClassName}}, error) {
	var {{pluralize (uncapitalize .BusinessName)}} []*model.{{.ClassName}}

	if err := r.db.Find(&{{pluralize (uncapitalize .BusinessName)}}).Error; err != nil {
		return nil, fmt.Errorf("查询所有{{.FunctionName}}失败: %v", err)
	}

	return {{pluralize (uncapitalize .BusinessName)}}, nil
}