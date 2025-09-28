import { request } from '@/utils/request'
import type {
  TableInfo,
  TableColumn,
  GenTableConfig,
  GenHistory,
  SystemMenu,
  GenerateRequest,
  GenerateResponse
} from '@/types/gen'

// 获取数据库表列表
export function getTableList() {
  return request<TableInfo[]>({
    url: '/v1/gen/tables',
    method: 'get'
  })
}

// 获取表字段信息
export function getTableColumns(tableName: string) {
  return request<TableColumn[]>({
    url: `/v1/gen/tables/${tableName}/columns`,
    method: 'get'
  })
}

// 保存生成配置
export function saveGenConfig(data: Partial<GenTableConfig>) {
  return request<GenTableConfig>({
    url: '/v1/gen/config',
    method: 'post',
    data
  })
}

// 获取配置详情
export function getGenConfig(id: number) {
  return request<GenTableConfig>({
    url: `/v1/gen/config/${id}`,
    method: 'get'
  })
}

// 更新配置
export function updateGenConfig(id: number, data: Partial<GenTableConfig>) {
  return request<GenTableConfig>({
    url: `/v1/gen/config/${id}`,
    method: 'put',
    data
  })
}

// 删除配置
export function deleteGenConfig(id: number) {
  return request<void>({
    url: `/v1/gen/config/${id}`,
    method: 'delete'
  })
}

// 获取配置列表
export function getGenConfigList(params?: {
  page?: number
  pageSize?: number
  tableName?: string
  businessName?: string
}) {
  return request<{
    total: number
    list: GenTableConfig[]
  }>({
    url: '/v1/gen/config',
    method: 'get',
    params
  })
}

// 生成代码
export function generateCode(data: GenerateRequest) {
  return request<GenerateResponse>({
    url: '/v1/gen/generate',
    method: 'post',
    data
  })
}

// 预览代码
export function previewCode(configId: number) {
  return request<{
    files: { [filename: string]: string }
  }>({
    url: `/v1/gen/preview/${configId}`,
    method: 'get'
  })
}

// 下载代码包
export function downloadCode(historyId: string) {
  return request({
    url: `/v1/gen/download/${historyId}`,
    method: 'get',
    responseType: 'blob'
  })
}

// 获取系统现有菜单树
export function getSystemMenus() {
  return request<SystemMenu[]>({
    url: '/v1/gen/menus/tree',
    method: 'get'
  })
}

// 获取生成历史记录
export function getGenHistory(params?: {
  page?: number
  pageSize?: number
  tableName?: string
  status?: string
}) {
  return request<{
    total: number
    list: GenHistory[]
  }>({
    url: '/v1/gen/history',
    method: 'get',
    params
  })
}