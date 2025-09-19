import request from '@/utils/request'

export const tenantApi = {
  // 获取租户列表
  getTenantList(params) {
    return request({
      url: '/api/v1/admin/tenants',
      method: 'get',
      params
    })
  },

  // 获取租户详情
  getTenant(id) {
    return request({
      url: `/api/v1/admin/tenants/${id}`,
      method: 'get'
    })
  },

  // 创建租户
  createTenant(data) {
    return request({
      url: '/api/v1/admin/tenants',
      method: 'post',
      data
    })
  },

  // 更新租户
  updateTenant(id, data) {
    return request({
      url: `/api/v1/admin/tenants/${id}`,
      method: 'put',
      data
    })
  },

  // 删除租户
  deleteTenant(id) {
    return request({
      url: `/api/v1/admin/tenants/${id}`,
      method: 'delete'
    })
  },

  // 更新租户状态
  updateTenantStatus(id, data) {
    return request({
      url: `/api/v1/admin/tenants/${id}/status`,
      method: 'put',
      data
    })
  },

  // 检查域名是否可用
  checkDomain(domain) {
    return request({
      url: '/api/v1/admin/tenants/check-domain',
      method: 'get',
      params: { domain }
    })
  },

  // 检查租户名称是否可用
  checkName(name) {
    return request({
      url: '/api/v1/admin/tenants/check-name',
      method: 'get',
      params: { name }
    })
  },

  // 获取租户统计信息
  getTenantStats(id) {
    return request({
      url: `/api/v1/admin/tenants/${id}/stats`,
      method: 'get'
    })
  }
}