import { http } from '@/utils/request'
import type { ApiResponse } from './types'

// 个人中心相关类型定义
export interface ProfileInfo {
  id: number
  username: string
  nickname: string
  email: string | null
  phone: string | null
  avatar: string
  status: number
  tenantId: number
  roles: Array<{
    id: number
    name: string
    code: string
    description?: string
  }>
  createdAt: string
  updatedAt: string
}

export interface UpdateProfileData {
  nickname: string
  email: string
  phone?: string
}

export interface ChangePasswordData {
  oldPassword: string
  newPassword: string
}

export interface FileStorageConfig {
  type: 'local' | 'oss'
  defaultPublic: boolean
  maxFileSize: number
  allowedTypes: string[]
  // 本地存储配置
  localAccessDomain?: string
  // OSS配置 - 使用自定义域名直接访问
  ossProvider?: 'aliyun' | 'tencent' | 'aws' | 'qiniu' | 'upyun'
  ossEndpoint?: string
  ossRegion?: string
  ossBucket?: string
  ossAccessKey?: string
  ossSecretKey?: string
  ossCustomDomain?: string
}

export interface TenantConfig {
  fileStorage: FileStorageConfig
  systemName?: string  // 系统名称
  logo?: string       // 系统Logo URL
  description?: string // 系统描述
  copyright?: string   // 版权信息
}

// 个人中心API
export const profileApi = {
  // 获取个人信息
  getProfile(): Promise<ApiResponse<ProfileInfo>> {
    return http.get('/v1/user/profile')
  },

  // 更新个人信息
  updateProfile(data: UpdateProfileData): Promise<ApiResponse> {
    return http.put('/v1/user/profile', data)
  },

  // 修改密码
  changePassword(data: ChangePasswordData): Promise<ApiResponse> {
    return http.put('/v1/user/password', data)
  },

  // 获取所在租户配置（仅租户管理员）
  getTenantConfig(): Promise<ApiResponse<TenantConfig>> {
    return http.get('/v1/user/tenant-config')
  },

  // 更新所在租户配置（仅租户管理员）
  updateTenantConfig(data: TenantConfig): Promise<ApiResponse> {
    return http.put('/v1/user/tenant-config', data)
  }
}