import request from '@/utils/request'
import type {
  {{.ClassName}},
  PageParams,
  PageResult
} from '@/api/types'

// {{.FunctionName}}相关API接口

// 获取{{.FunctionName}}列表
export function get{{pluralize .ClassName}}(params?: PageParams & {
{{- if .HasQuery }}
{{- range .QueryFields }}
{{- if eq .QueryType "BETWEEN" }}
  {{generateJSField .ColumnName}}_start?: string
  {{generateJSField .ColumnName}}_end?: string
{{- else }}
  {{generateJSField .ColumnName}}?: {{if eq .GoType "string"}}string{{else if contains .GoType "int"}}number{{else if eq .GoType "bool"}}boolean{{else}}any{{end}}
{{- end }}
{{- end }}
{{- else }}
  keyword?: string
{{- end }}
}) {
  return request<PageResult<{{.ClassName}}>>({
    url: '/v1/{{toLower .ModuleName}}/{{toLower (pluralize .BusinessName)}}',
    method: 'get',
    params
  })
}

// 获取{{.FunctionName}}详情
export function get{{.ClassName}}({{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}: {{if contains .GoType "int"}}number{{else}}string{{end}}{{break}}{{end}}{{end}}) {
  return request<{{.ClassName}}>({
    url: '/v1/{{toLower .ModuleName}}/{{toLower (pluralize .BusinessName)}}/' + {{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}{{break}}{{end}}{{end}},
    method: 'get'
  })
}

// 创建{{.FunctionName}}
export function create{{.ClassName}}(data: Partial<{{.ClassName}}>) {
  return request<{{.ClassName}}>({
    url: '/v1/{{toLower .ModuleName}}/{{toLower (pluralize .BusinessName)}}',
    method: 'post',
    data
  })
}

// 更新{{.FunctionName}}
export function update{{.ClassName}}({{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}: {{if contains .GoType "int"}}number{{else}}string{{end}}, {{break}}{{end}}{{end}}data: Partial<{{.ClassName}}>) {
  return request<{{.ClassName}}>({
    url: '/v1/{{toLower .ModuleName}}/{{toLower (pluralize .BusinessName)}}/' + {{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}{{break}}{{end}}{{end}},
    method: 'put',
    data
  })
}

// 删除{{.FunctionName}}
export function delete{{.ClassName}}({{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}: {{if contains .GoType "int"}}number{{else}}string{{end}}{{break}}{{end}}{{end}}) {
  return request<void>({
    url: '/v1/{{toLower .ModuleName}}/{{toLower (pluralize .BusinessName)}}/' + {{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}{{break}}{{end}}{{end}},
    method: 'delete'
  })
}

{{- if .HasQuery }}
// 批量删除{{.FunctionName}}
export function batchDelete{{pluralize .ClassName}}({{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}s: {{if contains .GoType "int"}}number{{else}}string{{end}}[]{{break}}{{end}}{{end}}) {
  return request<void>({
    url: '/v1/{{toLower .ModuleName}}/{{toLower (pluralize .BusinessName)}}/batch',
    method: 'delete',
    data: {
      {{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}s{{break}}{{end}}{{end}}
    }
  })
}

{{- end }}
// 导出{{.FunctionName}}
export function export{{pluralize .ClassName}}(params?: {
{{- if .HasQuery }}
{{- range .QueryFields }}
{{- if eq .QueryType "BETWEEN" }}
  {{generateJSField .ColumnName}}_start?: string
  {{generateJSField .ColumnName}}_end?: string
{{- else }}
  {{generateJSField .ColumnName}}?: {{if eq .GoType "string"}}string{{else if contains .GoType "int"}}number{{else if eq .GoType "bool"}}boolean{{else}}any{{end}}
{{- end }}
{{- end }}
{{- else }}
  keyword?: string
{{- end }}
}) {
  return request({
    url: '/v1/{{toLower .ModuleName}}/{{toLower (pluralize .BusinessName)}}/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}

// 获取活跃的{{.FunctionName}}列表（用于下拉选择等）
export function getActive{{pluralize .ClassName}}() {
  return request<{{.ClassName}}[]>({
    url: '/v1/{{toLower .ModuleName}}/{{toLower (pluralize .BusinessName)}}/active',
    method: 'get'
  })
}