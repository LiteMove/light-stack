import request from '@/utils/request'

// 获取仪表盘统计数据
export function getDashboardStats() {
  return request({
    url: '/v1/dashboard/stats',
    method: 'get'
  })
}

// 获取系统信息（超管专用）
export function getSystemInfo() {
  return request({
    url: '/v1/super-admin/system/info',
    method: 'get'
  })
}