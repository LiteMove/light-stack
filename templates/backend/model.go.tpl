package model

import (
	"time"
{{- if .HasQuery }}
	"gorm.io/gorm"
{{- end }}
)

// {{.ClassName}} {{.FunctionName}}
type {{.ClassName}} struct {
{{- range .Fields }}
	{{.GoField}} {{.GoType}} `json:"{{generateJSField .ColumnName}}" gorm:"{{- if .IsPk }}primaryKey;{{- end }}{{- if .IsIncrement }}autoIncrement;{{- end }}column:{{.ColumnName}};{{- if not .IsRequired }}not null;{{- end }}comment:{{.ColumnComment}}"`{{- if .ColumnComment }} // {{.ColumnComment}}{{- end }}
{{- end }}
}

{{- if .HasQuery }}

// {{.ClassName}}Query {{.FunctionName}}查询参数
type {{.ClassName}}Query struct {
{{- range .QueryFields }}
	{{.GoField}} {{- if eq .QueryType "BETWEEN" }}Range [2]{{.GoType}}{{- else }}{{.GoType}}{{- end }} `json:"{{generateJSField .ColumnName}}" form:"{{generateJSField .ColumnName}}"`{{- if .ColumnComment }} // {{.ColumnComment}}{{- end }}
{{- end }}

	// 分页参数
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// Apply 应用查询条件
func (q *{{.ClassName}}Query) Apply(db *gorm.DB) *gorm.DB {
	if db == nil {
		return nil
	}

{{- range .QueryFields }}
	{{- if eq .QueryType "LIKE" }}
	if q.{{.GoField}} != "" {
		db = db.Where("{{.ColumnName}} LIKE ?", "%"+q.{{.GoField}}+"%")
	}
	{{- else if eq .QueryType "BETWEEN" }}
	if !q.{{.GoField}}Range[0].IsZero() && !q.{{.GoField}}Range[1].IsZero() {
		db = db.Where("{{.ColumnName}} BETWEEN ? AND ?", q.{{.GoField}}Range[0], q.{{.GoField}}Range[1])
	}
	{{- else if eq .QueryType "EQ" }}
	if q.{{.GoField}} != {{getDefaultValue .}} {
		db = db.Where("{{.ColumnName}} = ?", q.{{.GoField}})
	}
	{{- else if eq .QueryType "NE" }}
	if q.{{.GoField}} != {{getDefaultValue .}} {
		db = db.Where("{{.ColumnName}} != ?", q.{{.GoField}})
	}
	{{- else if eq .QueryType "GT" }}
	if q.{{.GoField}} != {{getDefaultValue .}} {
		db = db.Where("{{.ColumnName}} > ?", q.{{.GoField}})
	}
	{{- else if eq .QueryType "LT" }}
	if q.{{.GoField}} != {{getDefaultValue .}} {
		db = db.Where("{{.ColumnName}} < ?", q.{{.GoField}})
	}
	{{- end }}
{{- end }}

	return db
}
{{- end }}

// TableName 指定表名
func ({{.ClassName}}) TableName() string {
	return "{{.TableName}}"
}