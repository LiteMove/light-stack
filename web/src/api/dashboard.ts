import request from '@/utils/request'

// 获取仪表盘统计数据
export function getDashboardStats() {
  return request({
    url: '/v1/analytics/dashboard',
    method: 'get'
  })
}

// 获取系统信息
export function getSystemInfo() {
  return request({
    url: '/v1/analytics/stats',
    method: 'get'
  })
}