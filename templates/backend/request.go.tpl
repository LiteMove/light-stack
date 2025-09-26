package request

import (
	"time"
)

// {{.ClassName}}CreateRequest 创建{{.FunctionName}}请求
type {{.ClassName}}CreateRequest struct {
{{- range .FormFields }}
{{- if and (not .IsPk) (not .IsIncrement) .IsInsert }}
	{{.GoField}} {{.GoType}} `json:"{{generateJSField .ColumnName}}" binding:"{{- if .IsRequired }}required{{- end }}"` // {{.ColumnComment}}
{{- end }}
{{- end }}
}

// {{.ClassName}}UpdateRequest 更新{{.FunctionName}}请求
type {{.ClassName}}UpdateRequest struct {
{{- range .FormFields }}
{{- if and (not .IsPk) (not .IsIncrement) .IsEdit }}
	{{.GoField}} {{.GoType}} `json:"{{generateJSField .ColumnName}}" binding:"{{- if .IsRequired }}required{{- end }}"` // {{.ColumnComment}}
{{- end }}
{{- end }}
}

{{- if .HasQuery }}

// {{.ClassName}}QueryRequest {{.FunctionName}}查询请求
type {{.ClassName}}QueryRequest struct {
{{- range .QueryFields }}
	{{.GoField}} {{- if eq .QueryType "BETWEEN" }}Range [2]{{.GoType}}{{- else }}{{.GoType}}{{- end }} `json:"{{generateJSField .ColumnName}}" form:"{{generateJSField .ColumnName}}"` // {{.ColumnComment}}
{{- end }}

	// 分页参数
	Page     int `json:"page" form:"page" binding:"min=1"`           // 页码
	PageSize int `json:"pageSize" form:"pageSize" binding:"min=1,max=100"` // 每页数量
}
{{- end }}

// {{.ClassName}}Response {{.FunctionName}}响应
type {{.ClassName}}Response struct {
{{- range .Fields }}
	{{.GoField}} {{.GoType}} `json:"{{generateJSField .ColumnName}}"` // {{.ColumnComment}}
{{- end }}
}

// {{.ClassName}}ListResponse {{.FunctionName}}列表响应
type {{.ClassName}}ListResponse struct {
	Total int64                      `json:"total"` // 总数
	List  []{{.ClassName}}Response   `json:"list"`  // 列表数据
}