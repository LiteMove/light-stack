import request from '@/utils/request'

// 字典类型相关接口
export interface DictType {
  id: number
  name: string
  type: string
  description: string
  status: number
  createdAt: string
  updatedAt: string
}

export interface DictData {
  id: number
  dictType: string
  label: string
  value: string
  sortOrder: number
  cssClass: string
  listClass: string
  isDefault: boolean
  status: number
  remark: string
  createdAt: string
  updatedAt: string
}

export interface DictOption {
  label: string
  value: string
  cssClass?: string
  listClass?: string
  isDefault: boolean
}

// === 字典类型相关接口 ===

// 创建字典类型
export function createDictType(data: {
  name: string
  type: string
  description?: string
  status: number
}) {
  return request({
    url: '/v1/super-admin/dict/types',
    method: 'post',
    data
  })
}

// 获取字典类型列表
export function getDictTypeList(params: {
  page: number
  pageSize: number
  status?: number
  name?: string
}) {
  return request({
    url: '/v1/super-admin/dict/types',
    method: 'get',
    params
  })
}

// 获取字典类型详情
export function getDictType(id: number) {
  return request({
    url: `/v1/super-admin/dict/types/${id}`,
    method: 'get'
  })
}

// 更新字典类型
export function updateDictType(id: number, data: {
  name: string
  type: string
  description?: string
  status: number
}) {
  return request({
    url: `/v1/super-admin/dict/types/${id}`,
    method: 'put',
    data
  })
}

// 删除字典类型
export function deleteDictType(id: number) {
  return request({
    url: `/v1/super-admin/dict/types/${id}`,
    method: 'delete'
  })
}

// === 字典数据相关接口 ===

// 创建字典数据
export function createDictData(data: {
  dictType: string
  label: string
  value: string
  sortOrder?: number
  cssClass?: string
  listClass?: string
  isDefault?: boolean
  status: number
  remark?: string
}) {
  return request({
    url: '/v1/super-admin/dict/data',
    method: 'post',
    data
  })
}

// 获取字典数据列表
export function getDictDataList(dictType: string, params: {
  page: number
  pageSize: number
  status?: number
  label?: string
}) {
  return request({
    url: `/v1/super-admin/dict/data/type/${dictType}`,
    method: 'get',
    params
  })
}

// 获取字典数据详情
export function getDictData(id: number) {
  return request({
    url: `/v1/super-admin/dict/data/${id}`,
    method: 'get'
  })
}

// 更新字典数据
export function updateDictData(id: number, data: {
  dictType: string
  label: string
  value: string
  sortOrder?: number
  cssClass?: string
  listClass?: string
  isDefault?: boolean
  status: number
  remark?: string
}) {
  return request({
    url: `/v1/super-admin/dict/data/${id}`,
    method: 'put',
    data
  })
}

// 删除字典数据
export function deleteDictData(id: number) {
  return request({
    url: `/v1/super-admin/dict/data/${id}`,
    method: 'delete'
  })
}

// 批量更新字典数据状态
export function batchUpdateDictDataStatus(data: {
  ids: number[]
  status: number
}) {
  return request({
    url: '/v1/super-admin/dict/data/batch/status',
    method: 'put',
    data
  })
}

// 批量删除字典数据
export function batchDeleteDictData(data: {
  ids: number[]
}) {
  return request({
    url: '/v1/super-admin/dict/data/batch',
    method: 'delete',
    data
  })
}

// === 字典选项接口（供前端下拉框使用） ===

// 获取字典选项
export function getDictOptions(dictType: string): Promise<{ data: DictOption[] }> {
  return request({
    url: `/v1/dict/options/${dictType}`,
    method: 'get'
  })
}