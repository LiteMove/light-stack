/**
 * {{.FunctionName}}相关类型定义
 *
 * 业务名称: {{.BusinessName}}
 * 表名: {{.TableName}}
 * 作者: {{.Author}}
 * 创建日期: {{.Date}}
 */

// {{.FunctionName}}实体类型
export interface {{.ClassName}} {
{{- range .Fields }}
  {{generateJSField .ColumnName}}: {{if eq .GoType "string"}}string{{else if eq .GoType "bool"}}boolean{{else if contains .GoType "int"}}number{{else if eq .GoType "time.Time"}}string{{else if contains .GoType "float"}}number{{else}}any{{end}}{{if not .IsRequired}} | null{{end}} // {{.ColumnComment}}
{{- end }}
  createdAt?: string // 创建时间
  updatedAt?: string // 更新时间
}

{{- if .HasQuery }}
// {{.FunctionName}}查询参数类型
export interface {{.ClassName}}Query {
{{- range .QueryFields }}
{{- if eq .QueryType "BETWEEN" }}
  {{generateJSField .ColumnName}}_start?: string // {{.ColumnComment}}开始{{if contains .ColumnName "time"}}时间{{else}}值{{end}}
  {{generateJSField .ColumnName}}_end?: string   // {{.ColumnComment}}结束{{if contains .ColumnName "time"}}时间{{else}}值{{end}}
{{- else }}
  {{generateJSField .ColumnName}}?: {{if eq .GoType "string"}}string{{else if contains .GoType "int"}}number{{else if eq .GoType "bool"}}boolean{{else}}any{{end}} // {{.ColumnComment}}
{{- end }}
{{- end }}
  page?: number        // 页码
  pageSize?: number    // 每页数量
}
{{- end }}

// {{.FunctionName}}表单数据类型
export interface {{.ClassName}}FormData {
{{- range .FormFields }}
  {{generateJSField .ColumnName}}: {{if eq .GoType "string"}}string{{else if eq .GoType "bool"}}boolean{{else if contains .GoType "int"}}number{{else if eq .GoType "time.Time"}}string{{else if contains .GoType "float"}}number{{else}}any{{end}}{{if not .IsRequired}} | null{{end}} // {{.ColumnComment}}
{{- end }}
}

// {{.FunctionName}}创建请求类型
export interface Create{{.ClassName}}Request {
{{- range .FormFields }}
{{- if and (not .IsPk) .IsInsert }}
  {{generateJSField .ColumnName}}{{if not .IsRequired}}?{{end}}: {{if eq .GoType "string"}}string{{else if eq .GoType "bool"}}boolean{{else if contains .GoType "int"}}number{{else if eq .GoType "time.Time"}}string{{else if contains .GoType "float"}}number{{else}}any{{end}}{{if not .IsRequired}} | null{{end}} // {{.ColumnComment}}
{{- end }}
{{- end }}
}

// {{.FunctionName}}更新请求类型
export interface Update{{.ClassName}}Request {
{{- range .FormFields }}
{{- if and (not .IsPk) .IsEdit }}
  {{generateJSField .ColumnName}}{{if not .IsRequired}}?{{end}}: {{if eq .GoType "string"}}string{{else if eq .GoType "bool"}}boolean{{else if contains .GoType "int"}}number{{else if eq .GoType "time.Time"}}string{{else if contains .GoType "float"}}number{{else}}any{{end}}{{if not .IsRequired}} | null{{end}} // {{.ColumnComment}}
{{- end }}
{{- end }}
}

{{- if .HasQuery }}
// 批量删除{{.FunctionName}}请求类型
export interface BatchDelete{{pluralize .ClassName}}Request {
  {{range .Fields}}{{if .IsPk}}{{generateJSField .ColumnName}}s: {{if contains .GoType "int"}}number{{else}}string{{end}}[] // {{.ColumnComment}}列表{{break}}{{end}}{{end}}
}
{{- end }}

{{- $hasSelectField := false }}
{{- range .Fields }}
{{- if eq .HtmlType "select" }}
{{- if not $hasSelectField }}
{{- $hasSelectField = true }}
// {{$.FunctionName}}状态枚举
export enum {{$.ClassName}}Status {
  // TODO: 根据实际业务需求定义状态
  ACTIVE = 1,   // 激活
  INACTIVE = 0, // 未激活
}

// {{$.FunctionName}}状态标签映射
export const {{toUpper $.BusinessName}}_STATUS_MAP = {
  [{{$.ClassName}}Status.ACTIVE]: {
    label: '激活',
    type: 'success'
  },
  [{{$.ClassName}}Status.INACTIVE]: {
    label: '未激活',
    type: 'info'
  }
} as const
{{- end }}
{{- break }}
{{- end }}
{{- end }}

{{- $hasDict := false }}
{{- range .Fields }}
{{- if ne .DictType "" }}
{{- if not $hasDict }}
{{- $hasDict = true }}
// {{$.FunctionName}}字典类型
export interface {{$.ClassName}}Dict {
  {{.DictType}}: Array<{
    label: string
    value: string | number
    type?: string
  }>
}
{{- end }}
{{- end }}
{{- end }}

// 导出所有{{.FunctionName}}相关类型
export type {
  {{.ClassName}},
{{- if .HasQuery }}
  {{.ClassName}}Query,
{{- end }}
  {{.ClassName}}FormData,
  Create{{.ClassName}}Request,
  Update{{.ClassName}}Request
{{- if .HasQuery }},
  BatchDelete{{pluralize .ClassName}}Request
{{- end }}
{{- if $hasDict }},
  {{.ClassName}}Dict
{{- end }}
}

{{- if $hasSelectField }}
export { {{$.ClassName}}Status, {{toUpper $.BusinessName}}_STATUS_MAP }
{{- end }}