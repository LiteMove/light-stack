import request from '@/utils/request'
import type { ApiResponse, PageResponse, PageParams } from './types'

// 租户信息
interface Tenant {
  id: number
  name: string
  domain: string
  description?: string
  status: number
  isSystem: boolean
  createdAt: string
  updatedAt: string
}

// 租户创建/更新请求参数
interface TenantFormData {
  name: string
  domain: string
  description?: string
  status?: number
}

// 租户状态更新参数
interface TenantStatusData {
  status: number
}

// 租户列表查询参数
interface TenantQueryParams extends PageParams {
  name?: string
  status?: number
}

// 域名检查响应
interface DomainCheckResponse {
  available: boolean
  message?: string
}

// 名称检查响应
interface NameCheckResponse {
  available: boolean
  message?: string
}

// 租户统计信息
interface TenantStats {
  userCount: number
  roleCount: number
  menuCount: number
  apiCount: number
}

export const tenantApi = {
  // 根据域名获取租户展示信息（公开接口）
  getTenantInfo(domain?: string): Promise<ApiResponse<{
    id: number
    name: string
    systemName: string
    logo: string
    description: string
    copyright: string
  }>> {
    return request({
      url: '/tenant/info',
      method: 'get',
      params: domain ? { domain } : {}
    })
  },

  // 获取租户列表
  getTenantList(params: TenantQueryParams): Promise<ApiResponse<PageResponse<Tenant>>> {
    return request({
      url: '/v1/admin/tenants',
      method: 'get',
      params
    })
  },
    // 获取租户下拉列表
    getTenantSelectList(): Promise<ApiResponse<Tenant[]>> {
      return request({
        url: '/v1/admin/tenants/list',
        method: 'get'
      })
    },

  // 获取租户详情
  getTenant(id: number): Promise<ApiResponse<Tenant>> {
    return request({
      url: `/v1/admin/tenants/${id}`,
      method: 'get'
    })
  },

  // 创建租户
  createTenant(data: TenantFormData): Promise<ApiResponse<Tenant>> {
    return request({
      url: '/v1/admin/tenants',
      method: 'post',
      data
    })
  },

  // 更新租户
  updateTenant(id: number, data: TenantFormData): Promise<ApiResponse<Tenant>> {
    return request({
      url: `/v1/admin/tenants/${id}`,
      method: 'put',
      data
    })
  },

  // 删除租户
  deleteTenant(id: number): Promise<ApiResponse<void>> {
    return request({
      url: `/v1/admin/tenants/${id}`,
      method: 'delete'
    })
  },

  // 更新租户状态
  updateTenantStatus(id: number, data: TenantStatusData): Promise<ApiResponse<void>> {
    return request({
      url: `/v1/admin/tenants/${id}/status`,
      method: 'put',
      data
    })
  },

  // 检查域名是否可用
  checkDomain(domain: string): Promise<ApiResponse<DomainCheckResponse>> {
    return request({
      url: '/v1/admin/tenants/check-domain',
      method: 'get',
      params: { domain }
    })
  },

  // 检查租户名称是否可用
  checkName(name: string): Promise<ApiResponse<NameCheckResponse>> {
    return request({
      url: '/v1/admin/tenants/check-name',
      method: 'get',
      params: { name }
    })
  },

  // 获取租户统计信息
  getTenantStats(id: number): Promise<ApiResponse<TenantStats>> {
    return request({
      url: `/v1/admin/tenants/${id}/stats`,
      method: 'get'
    })
  }
}